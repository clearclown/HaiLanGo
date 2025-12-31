package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// AzureVisionClient はAzure Computer Vision APIクライアント
type AzureVisionClient struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
}

// NewAzureVisionClient は新しいAzure Computer Vision APIクライアントを作成する
func NewAzureVisionClient(endpoint, apiKey string) *AzureVisionClient {
	return &AzureVisionClient{
		endpoint:   strings.TrimSuffix(endpoint, "/"),
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

// azureReadResponse はAzure Read APIのレスポンス構造
type azureReadResponse struct {
	Status       string `json:"status"`
	AnalyzeResult struct {
		ReadResults []struct {
			Page   int `json:"page"`
			Angle  float64 `json:"angle"`
			Width  float64 `json:"width"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Lines  []struct {
				Text        string    `json:"text"`
				BoundingBox []float64 `json:"boundingBox"`
				Words       []struct {
					Text        string    `json:"text"`
					BoundingBox []float64 `json:"boundingBox"`
					Confidence  float64   `json:"confidence"`
				} `json:"words"`
			} `json:"lines"`
		} `json:"readResults"`
	} `json:"analyzeResult"`
}

// ProcessImage は画像データをOCR処理する
func (a *AzureVisionClient) ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	// Azure Computer Vision Read API (非同期版) を使用
	// Step 1: 画像をサブミットしてOperation-Locationを取得
	operationURL, err := a.submitReadRequest(ctx, imageData, languages)
	if err != nil {
		return nil, fmt.Errorf("failed to submit read request: %w", err)
	}

	// Step 2: 結果をポーリングで取得
	result, err := a.pollForResult(ctx, operationURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get read result: %w", err)
	}

	return result, nil
}

// submitReadRequest は画像をAzure Read APIにサブミットする
func (a *AzureVisionClient) submitReadRequest(ctx context.Context, imageData []byte, languages []string) (string, error) {
	// Read API エンドポイント
	url := fmt.Sprintf("%s/vision/v3.2/read/analyze", a.endpoint)

	// 言語パラメータを追加
	if len(languages) > 0 {
		url += "?language=" + languages[0]
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Ocp-Apim-Subscription-Key", a.apiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Azure Vision API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Operation-Location ヘッダーから操作URLを取得
	operationURL := resp.Header.Get("Operation-Location")
	if operationURL == "" {
		return "", fmt.Errorf("Operation-Location header not found in response")
	}

	return operationURL, nil
}

// pollForResult は結果が準備できるまでポーリングする
func (a *AzureVisionClient) pollForResult(ctx context.Context, operationURL string) (*OCRResult, error) {
	maxRetries := 30
	retryInterval := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		// コンテキストのキャンセルをチェック
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "GET", operationURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Ocp-Apim-Subscription-Key", a.apiKey)

		resp, err := a.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Azure Vision API error (status %d): %s", resp.StatusCode, string(body))
		}

		var readResp azureReadResponse
		if err := json.Unmarshal(body, &readResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		switch readResp.Status {
		case "succeeded":
			return a.parseReadResult(&readResp), nil
		case "failed":
			return nil, fmt.Errorf("OCR processing failed")
		case "running", "notStarted":
			// まだ処理中、待機して再試行
			time.Sleep(retryInterval)
			continue
		default:
			return nil, fmt.Errorf("unknown status: %s", readResp.Status)
		}
	}

	return nil, fmt.Errorf("timeout waiting for OCR result")
}

// parseReadResult はAzure Read APIのレスポンスをOCRResultに変換する
func (a *AzureVisionClient) parseReadResult(readResp *azureReadResponse) *OCRResult {
	var fullText strings.Builder
	var pages []PageOCRResult
	var totalConfidence float64
	var wordCount int

	for _, readResult := range readResp.AnalyzeResult.ReadResults {
		var pageText strings.Builder

		for _, line := range readResult.Lines {
			pageText.WriteString(line.Text)
			pageText.WriteString("\n")

			for _, word := range line.Words {
				totalConfidence += word.Confidence
				wordCount++
			}
		}

		pageTextStr := strings.TrimSpace(pageText.String())
		fullText.WriteString(pageTextStr)
		fullText.WriteString("\n")

		pages = append(pages, PageOCRResult{
			PageNumber: readResult.Page,
			Text:       pageTextStr,
			Confidence: totalConfidence / float64(wordCount),
		})
	}

	avgConfidence := 0.0
	if wordCount > 0 {
		avgConfidence = totalConfidence / float64(wordCount)
	}

	return &OCRResult{
		Text:             strings.TrimSpace(fullText.String()),
		DetectedLanguage: "", // Azure Read APIはdetected languageを返さない
		Confidence:       avgConfidence,
		Pages:            pages,
	}
}
