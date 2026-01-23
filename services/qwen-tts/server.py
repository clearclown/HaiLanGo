"""
Qwen-TTS REST API Server
Provides a simple REST API for text-to-speech synthesis using Qwen3-TTS.

Endpoints:
  POST /synthesize - Synthesize text to speech
  GET /voices - List available voices
  GET /languages - List supported languages
  GET /health - Health check
"""

import os
import io
import base64
import logging
from typing import Optional, List
from contextlib import asynccontextmanager

import torch
import soundfile as sf
from fastapi import FastAPI, HTTPException
from fastapi.responses import Response, JSONResponse
from pydantic import BaseModel

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Global model instance
model = None
tokenizer = None

# Environment variables
MODEL_NAME = os.getenv("QWEN_TTS_MODEL", "Qwen/Qwen3-TTS-12Hz-0.6B-CustomVoice")
DEVICE = os.getenv("DEVICE", "cuda:0")
DTYPE = os.getenv("DTYPE", "bfloat16")

# Supported languages
SUPPORTED_LANGUAGES = [
    {"code": "zh", "name": "Chinese", "native_name": "中文"},
    {"code": "en", "name": "English", "native_name": "English"},
    {"code": "ja", "name": "Japanese", "native_name": "日本語"},
    {"code": "ko", "name": "Korean", "native_name": "한국어"},
    {"code": "de", "name": "German", "native_name": "Deutsch"},
    {"code": "fr", "name": "French", "native_name": "Français"},
    {"code": "ru", "name": "Russian", "native_name": "Русский"},
    {"code": "pt", "name": "Portuguese", "native_name": "Português"},
    {"code": "es", "name": "Spanish", "native_name": "Español"},
    {"code": "it", "name": "Italian", "native_name": "Italiano"},
]

# Available speakers for CustomVoice model
SPEAKERS = [
    {"id": "Vivian", "name": "Vivian", "language": "zh", "gender": "female", "description": "Bright, slightly edgy young female voice"},
    {"id": "Serena", "name": "Serena", "language": "zh", "gender": "female", "description": "Warm, gentle young female voice"},
    {"id": "Uncle_Fu", "name": "Uncle Fu", "language": "zh", "gender": "male", "description": "Seasoned male voice with a low, mellow timbre"},
    {"id": "Dylan", "name": "Dylan", "language": "zh", "gender": "male", "description": "Youthful Beijing male voice with a clear, natural timbre"},
    {"id": "Eric", "name": "Eric", "language": "zh", "gender": "male", "description": "Lively Chengdu male voice with a slightly husky brightness"},
    {"id": "Ryan", "name": "Ryan", "language": "en", "gender": "male", "description": "Dynamic male voice with strong rhythmic drive"},
    {"id": "Aiden", "name": "Aiden", "language": "en", "gender": "male", "description": "Sunny American male voice with a clear midrange"},
    {"id": "Ono_Anna", "name": "Ono Anna", "language": "ja", "gender": "female", "description": "Playful Japanese female voice with a light, nimble timbre"},
    {"id": "Sohee", "name": "Sohee", "language": "ko", "gender": "female", "description": "Warm Korean female voice with rich emotion"},
]

# Language code mapping
LANGUAGE_MAP = {
    "zh": "Chinese",
    "zh-CN": "Chinese",
    "zh-TW": "Chinese",
    "en": "English",
    "en-US": "English",
    "en-GB": "English",
    "ja": "Japanese",
    "ja-JP": "Japanese",
    "ko": "Korean",
    "ko-KR": "Korean",
    "de": "German",
    "de-DE": "German",
    "fr": "French",
    "fr-FR": "French",
    "ru": "Russian",
    "ru-RU": "Russian",
    "pt": "Portuguese",
    "pt-BR": "Portuguese",
    "pt-PT": "Portuguese",
    "es": "Spanish",
    "es-ES": "Spanish",
    "it": "Italian",
    "it-IT": "Italian",
}


class SynthesizeRequest(BaseModel):
    text: str
    language: str = "Auto"
    voice: str = "Vivian"
    instruct: Optional[str] = None
    format: str = "mp3"


class SynthesizeResponse(BaseModel):
    audio: str  # base64 encoded
    format: str
    sample_rate: int
    duration_ms: int


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Load model on startup."""
    global model

    logger.info(f"Loading Qwen-TTS model: {MODEL_NAME}")
    logger.info(f"Device: {DEVICE}, Dtype: {DTYPE}")

    try:
        from qwen_tts import Qwen3TTSModel

        # Determine dtype
        dtype = torch.bfloat16 if DTYPE == "bfloat16" else torch.float16

        # Try to use flash attention if available
        attn_impl = "flash_attention_2"
        try:
            import flash_attn
        except ImportError:
            logger.warning("Flash attention not available, using default attention")
            attn_impl = "eager"

        model = Qwen3TTSModel.from_pretrained(
            MODEL_NAME,
            device_map=DEVICE,
            dtype=dtype,
            attn_implementation=attn_impl,
        )

        logger.info("Model loaded successfully!")

    except Exception as e:
        logger.error(f"Failed to load model: {e}")
        logger.warning("Running in mock mode - audio will be empty")
        model = None

    yield

    # Cleanup
    if model is not None:
        del model
        torch.cuda.empty_cache()


app = FastAPI(
    title="Qwen-TTS API",
    description="Text-to-Speech API using Qwen3-TTS",
    version="1.0.0",
    lifespan=lifespan,
)


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {
        "status": "healthy",
        "model": MODEL_NAME,
        "model_loaded": model is not None,
        "device": DEVICE,
    }


@app.get("/languages")
async def get_languages():
    """Get list of supported languages."""
    return {"languages": SUPPORTED_LANGUAGES}


@app.get("/voices")
async def get_voices(language: Optional[str] = None):
    """Get list of available voices."""
    voices = SPEAKERS

    if language:
        # Map language code
        lang_name = LANGUAGE_MAP.get(language, language)
        # Filter by language (all speakers can speak all languages, but native is best)
        # Return all voices but mark native language
        voices = [
            {**v, "is_native": v["language"] == language or LANGUAGE_MAP.get(v["language"]) == lang_name}
            for v in SPEAKERS
        ]

    return {"voices": voices}


@app.post("/synthesize")
async def synthesize(request: SynthesizeRequest):
    """Synthesize text to speech."""

    if model is None:
        raise HTTPException(status_code=503, detail="Model not loaded")

    try:
        # Map language code
        language = LANGUAGE_MAP.get(request.language, request.language)
        if language == request.language and request.language != "Auto":
            # Try short code
            short_code = request.language.split("-")[0]
            language = LANGUAGE_MAP.get(short_code, "Auto")

        logger.info(f"Synthesizing: text='{request.text[:50]}...', language={language}, voice={request.voice}")

        # Generate audio
        if "CustomVoice" in MODEL_NAME:
            wavs, sr = model.generate_custom_voice(
                text=request.text,
                language=language,
                speaker=request.voice,
                instruct=request.instruct or "",
            )
        elif "VoiceDesign" in MODEL_NAME:
            wavs, sr = model.generate_voice_design(
                text=request.text,
                language=language,
                instruct=request.instruct or "Natural speaking voice",
            )
        else:
            # Base model - use default voice
            wavs, sr = model.generate_voice_clone(
                text=request.text,
                language=language,
                x_vector_only_mode=True,
            )

        # Convert to bytes
        audio_buffer = io.BytesIO()

        if request.format == "wav":
            sf.write(audio_buffer, wavs[0], sr, format="WAV")
        else:
            # Default to MP3-like format (actually WAV, but can be converted)
            sf.write(audio_buffer, wavs[0], sr, format="WAV")

        audio_buffer.seek(0)
        audio_bytes = audio_buffer.read()

        # Calculate duration
        duration_ms = int(len(wavs[0]) / sr * 1000)

        # Return as base64
        return SynthesizeResponse(
            audio=base64.b64encode(audio_bytes).decode("utf-8"),
            format="wav",
            sample_rate=sr,
            duration_ms=duration_ms,
        )

    except Exception as e:
        logger.error(f"Synthesis failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/synthesize/stream")
async def synthesize_stream(request: SynthesizeRequest):
    """Synthesize text to speech and return audio directly."""

    if model is None:
        raise HTTPException(status_code=503, detail="Model not loaded")

    try:
        # Map language code
        language = LANGUAGE_MAP.get(request.language, request.language)
        if language == request.language and request.language != "Auto":
            short_code = request.language.split("-")[0]
            language = LANGUAGE_MAP.get(short_code, "Auto")

        # Generate audio
        if "CustomVoice" in MODEL_NAME:
            wavs, sr = model.generate_custom_voice(
                text=request.text,
                language=language,
                speaker=request.voice,
                instruct=request.instruct or "",
            )
        elif "VoiceDesign" in MODEL_NAME:
            wavs, sr = model.generate_voice_design(
                text=request.text,
                language=language,
                instruct=request.instruct or "Natural speaking voice",
            )
        else:
            wavs, sr = model.generate_voice_clone(
                text=request.text,
                language=language,
                x_vector_only_mode=True,
            )

        # Convert to bytes
        audio_buffer = io.BytesIO()
        sf.write(audio_buffer, wavs[0], sr, format="WAV")
        audio_buffer.seek(0)

        return Response(
            content=audio_buffer.read(),
            media_type="audio/wav",
            headers={
                "Content-Disposition": "attachment; filename=speech.wav",
                "X-Sample-Rate": str(sr),
                "X-Duration-Ms": str(int(len(wavs[0]) / sr * 1000)),
            }
        )

    except Exception as e:
        logger.error(f"Synthesis failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    import uvicorn

    port = int(os.getenv("PORT", "8000"))
    host = os.getenv("HOST", "0.0.0.0")

    logger.info(f"Starting Qwen-TTS server on {host}:{port}")
    uvicorn.run(app, host=host, port=port)
