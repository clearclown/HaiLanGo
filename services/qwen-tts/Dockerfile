# Qwen-TTS Service Dockerfile
# GPU-based Text-to-Speech service supporting 10 languages
# Requires NVIDIA GPU with CUDA support

FROM docker.io/nvidia/cuda:12.1-devel-ubuntu22.04

# Prevent interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

# Install system dependencies
RUN apt-get update && apt-get install -y \
    python3.12 \
    python3.12-venv \
    python3.12-dev \
    python3-pip \
    git \
    ffmpeg \
    libsndfile1 \
    && rm -rf /var/lib/apt/lists/*

# Set Python 3.12 as default
RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.12 1
RUN update-alternatives --install /usr/bin/python python /usr/bin/python3.12 1

# Create working directory
WORKDIR /app

# Install Python dependencies
RUN pip install --no-cache-dir --upgrade pip

# Install qwen-tts and dependencies
RUN pip install --no-cache-dir \
    qwen-tts \
    fastapi \
    uvicorn \
    python-multipart \
    soundfile \
    numpy

# Install flash-attention for better performance (optional but recommended)
RUN MAX_JOBS=4 pip install --no-cache-dir flash-attn --no-build-isolation || true

# Copy the server script
COPY server.py /app/server.py

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1

# Run the server
CMD ["python", "server.py"]
