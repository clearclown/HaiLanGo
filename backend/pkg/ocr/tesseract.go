package ocr

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TesseractClient はTesseract OCRクライアント
type TesseractClient struct {
	tessDataPath string
}

// NewTesseractClient は新しいTesseract OCRクライアントを作成する
func NewTesseractClient() *TesseractClient {
	tessDataPath := os.Getenv("TESSDATA_PREFIX")
	if tessDataPath == "" {
		tessDataPath = "/usr/share/tesseract-ocr/4.00/tessdata"
	}
	return &TesseractClient{
		tessDataPath: tessDataPath,
	}
}

// ProcessImage は画像データをOCR処理する
func (t *TesseractClient) ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	// 一時ファイルに画像を書き込み
	tmpFile, err := os.CreateTemp("", "tesseract-input-*.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(imageData); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to write image data: %w", err)
	}
	tmpFile.Close()

	// 出力ファイルのベース名を準備
	outputBase := tmpFile.Name() + "-out"
	outputFile := outputBase + ".txt"
	defer os.Remove(outputFile)

	// tesseractコマンドを構築
	args := []string{tmpFile.Name(), outputBase}

	// 言語設定
	if len(languages) > 0 {
		// Tesseractの言語コード形式に変換（例: ja+eng）
		langStr := strings.Join(languages, "+")
		args = append(args, "-l", langStr)
	}

	// tesseractコマンドを実行
	cmd := exec.CommandContext(ctx, "tesseract", args...)

	// TESSDATA_PREFIXを設定
	if t.tessDataPath != "" {
		cmd.Env = append(os.Environ(), "TESSDATA_PREFIX="+t.tessDataPath)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("tesseract failed: %w, stderr: %s", err, stderr.String())
	}

	// 結果を読み取り
	resultData, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	text := strings.TrimSpace(string(resultData))

	// 言語検出（Tesseractは自動検出しないので、指定された言語を使用）
	detectedLanguage := ""
	if len(languages) > 0 {
		detectedLanguage = languages[0]
	}

	return &OCRResult{
		Text:             text,
		DetectedLanguage: detectedLanguage,
		Confidence:       0.85, // Tesseractは信頼度を提供しないので固定値
		Pages: []PageOCRResult{
			{
				PageNumber: 1,
				Text:       text,
				Confidence: 0.85,
			},
		},
	}, nil
}
