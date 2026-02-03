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

/// Google Cloud Vision OCR provider
pub struct GoogleVisionOcrProvider {
    api_key: String,
    endpoint: String,
}

impl GoogleVisionOcrProvider {
    /// Create a new Google Vision OCR provider
    pub fn new(api_key: &str) -> Self {
        Self {
            api_key: api_key.to_string(),
            endpoint: "https://vision.googleapis.com/v1/images:annotate".to_string(),
        }
    }

    /// Create from environment variable
    pub fn from_env() -> Option<Self> {
        std::env::var("GOOGLE_CLOUD_VISION_API_KEY")
            .ok()
            .map(|key| Self::new(&key))
    }

    /// Build the request body for OCR
    fn build_request_body(&self, image_data: &[u8]) -> serde_json::Value {
        use base64::{Engine, engine::general_purpose::STANDARD};
        let encoded = STANDARD.encode(image_data);

        serde_json::json!({
            "requests": [{
                "image": {
                    "content": encoded
                },
                "features": [{
                    "type": "TEXT_DETECTION",
                    "maxResults": 50
                }, {
                    "type": "DOCUMENT_TEXT_DETECTION"
                }],
                "imageContext": {
                    "languageHints": ["zh", "ja", "ko", "en"]
                }
            }]
        })
    }
}

#[async_trait]
impl OcrProvider for GoogleVisionOcrProvider {
    async fn extract_text(&self, image_data: &[u8]) -> Result<OcrResult, OcrError> {
        let client = reqwest::Client::new();
        let url = format!("{}?key={}", self.endpoint, self.api_key);
        let body = self.build_request_body(image_data);

        let response = client
            .post(&url)
            .json(&body)
            .send()
            .await
            .map_err(|_e| OcrError::ServiceUnavailable)?;

        if response.status() == 429 {
            return Err(OcrError::RateLimitExceeded);
        }

        if !response.status().is_success() {
            return Err(OcrError::ProcessingFailed(format!(
                "API returned status {}",
                response.status()
            )));
        }

        let result: serde_json::Value = response
            .json()
            .await
            .map_err(|e| OcrError::ProcessingFailed(e.to_string()))?;

        // Parse the response
        let text = result["responses"][0]["fullTextAnnotation"]["text"]
            .as_str()
            .unwrap_or("")
            .to_string();

        let language = result["responses"][0]["fullTextAnnotation"]["pages"][0]["property"]
            ["detectedLanguages"][0]["languageCode"]
            .as_str()
            .map(String::from);

        let confidence = result["responses"][0]["fullTextAnnotation"]["pages"][0]["confidence"]
            .as_f64()
            .unwrap_or(0.0) as f32;

        // Parse bounding boxes
        let mut bounding_boxes = Vec::new();
        if let Some(annotations) = result["responses"][0]["textAnnotations"].as_array() {
            for (i, annotation) in annotations.iter().enumerate() {
                if i == 0 {
                    continue; // Skip first annotation (full text)
                }
                if let Some(vertices) = annotation["boundingPoly"]["vertices"].as_array() {
                    let x = vertices[0]["x"].as_f64().unwrap_or(0.0) as f32;
                    let y = vertices[0]["y"].as_f64().unwrap_or(0.0) as f32;
                    let x2 = vertices[2]["x"].as_f64().unwrap_or(0.0) as f32;
                    let y2 = vertices[2]["y"].as_f64().unwrap_or(0.0) as f32;

                    bounding_boxes.push(BoundingBox {
                        x,
                        y,
                        width: x2 - x,
                        height: y2 - y,
                        text: annotation["description"].as_str().unwrap_or("").to_string(),
                    });
                }
            }
        }

        Ok(OcrResult {
            text,
            confidence,
            language_detected: language,
            bounding_boxes,
        })
    }

    async fn extract_text_pdf(&self, pdf_data: &[u8], _page: usize) -> Result<OcrResult, OcrError> {
        // For PDF, we would need to convert to image first or use Cloud Vision PDF feature
        // For now, treat it as image data (works for single-page PDFs with images)
        self.extract_text(pdf_data).await
    }
}

/// Factory function to create the appropriate OCR provider
pub fn create_ocr_provider() -> Box<dyn OcrProvider> {
    if let Some(provider) = GoogleVisionOcrProvider::from_env() {
        Box::new(provider)
    } else {
        tracing::warn!("GOOGLE_CLOUD_VISION_API_KEY not set, using mock OCR provider");
        Box::new(MockOcrProvider)
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
