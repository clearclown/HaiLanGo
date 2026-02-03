//! OCR (Optical Character Recognition) service abstraction

use async_trait::async_trait;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum OcrError {
    #[error("Service unavailable")]
    ServiceUnavailable,
    #[error("Invalid image format")]
    InvalidFormat,
    #[error("Processing failed: {0}")]
    ProcessingFailed(String),
    #[error("Rate limit exceeded")]
    RateLimitExceeded,
}

/// Result of OCR processing
#[derive(Debug, Clone)]
pub struct OcrResult {
    pub text: String,
    pub confidence: f32,
    pub language_detected: Option<String>,
    pub bounding_boxes: Vec<BoundingBox>,
}

/// Bounding box for text region
#[derive(Debug, Clone)]
pub struct BoundingBox {
    pub x: f32,
    pub y: f32,
    pub width: f32,
    pub height: f32,
    pub text: String,
}

/// OCR Provider trait - implement for different providers
#[async_trait]
pub trait OcrProvider: Send + Sync {
    /// Extract text from an image
    async fn extract_text(&self, image_data: &[u8]) -> Result<OcrResult, OcrError>;

    /// Extract text from a PDF page
    async fn extract_text_pdf(&self, pdf_data: &[u8], page: usize) -> Result<OcrResult, OcrError>;
}

/// Mock OCR provider for development and testing
pub struct MockOcrProvider;

#[async_trait]
impl OcrProvider for MockOcrProvider {
    async fn extract_text(&self, _image_data: &[u8]) -> Result<OcrResult, OcrError> {
        Ok(OcrResult {
            text: "Mock extracted text from image".to_string(),
            confidence: 0.95,
            language_detected: Some("en".to_string()),
            bounding_boxes: vec![],
        })
    }

    async fn extract_text_pdf(&self, _pdf_data: &[u8], page: usize) -> Result<OcrResult, OcrError> {
        Ok(OcrResult {
            text: format!("Mock extracted text from PDF page {}", page),
            confidence: 0.90,
            language_detected: Some("en".to_string()),
            bounding_boxes: vec![],
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_mock_ocr_extract_text() {
        let provider = MockOcrProvider;
        let result = provider.extract_text(&[0u8; 100]).await.unwrap();

        assert!(!result.text.is_empty());
        assert!(result.confidence > 0.0);
    }

    #[tokio::test]
    async fn test_mock_ocr_extract_pdf() {
        let provider = MockOcrProvider;
        let result = provider.extract_text_pdf(&[0u8; 100], 1).await.unwrap();

        assert!(result.text.contains("page 1"));
    }

    #[tokio::test]
    async fn test_ocr_result_structure() {
        let result = OcrResult {
            text: "Test text".to_string(),
            confidence: 0.95,
            language_detected: Some("en".to_string()),
            bounding_boxes: vec![],
        };

        assert_eq!(result.text, "Test text");
        assert_eq!(result.confidence, 0.95);
    }

    #[tokio::test]
    async fn test_bounding_box_creation() {
        let bbox = BoundingBox {
            x: 10.0,
            y: 20.0,
            width: 100.0,
            height: 50.0,
            text: "hello".to_string(),
        };

        assert_eq!(bbox.x, 10.0);
        assert_eq!(bbox.text, "hello");
    }

    #[test]
    fn test_ocr_error_display() {
        let error = OcrError::ServiceUnavailable;
        assert_eq!(error.to_string(), "Service unavailable");

        let error = OcrError::InvalidFormat;
        assert_eq!(error.to_string(), "Invalid image format");
    }
}
