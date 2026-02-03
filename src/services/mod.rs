//! External service integrations

pub mod ocr;
pub mod tts;

pub use ocr::{MockOcrProvider, OcrError, OcrProvider, OcrResult};
pub use tts::{AudioFormat, MockTtsProvider, TtsError, TtsProvider, TtsRequest, TtsResponse};
