package translate

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const (
	model      = "claude-sonnet-4-6"
	apiTimeout = 30 * time.Second
	prompt     = `この画像はゲームやアプリのスクリーンショットです。
画像内に表示されている英語のテキストをすべて抽出し、自然な日本語に翻訳してください。

- 翻訳結果のテキストのみを出力する（前置きや説明は不要）
- 画像にテキストが含まれない場合は「（テキストなし）」と出力する
- UI要素（ボタン名、メニュー項目）も含めて翻訳する
- ゲーム固有の用語はカタカナで表記する（例: "Health" → "ヘルス"）
- 元のレイアウト構造（改行・段落）をできるだけ保持する`
)

// Translator は翻訳処理のインターフェース。
type Translator interface {
	// Translate はPNG画像バイト列を受け取り、翻訳結果テキストを返す。
	Translate(ctx context.Context, pngBytes []byte) (string, error)
}

// ClaudeTranslator は Claude Vision API を使って翻訳する。
type ClaudeTranslator struct {
	client anthropic.Client
}

// New は新しい ClaudeTranslator を返す。
// APIキーは環境変数 ANTHROPIC_API_KEY から自動取得される。
func New(apiKey string) *ClaudeTranslator {
	return &ClaudeTranslator{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
	}
}

// Translate はスクリーンショットのPNGバイト列を翻訳する。
func (t *ClaudeTranslator) Translate(ctx context.Context, pngBytes []byte) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()

	encoded := base64.StdEncoding.EncodeToString(pngBytes)

	slog.Info("translate: sending to Claude Vision API", "bytes", len(pngBytes))

	msg, err := t.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     model,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewImageBlockBase64("image/png", encoded),
				anthropic.NewTextBlock(prompt),
			),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude api call failed: %w", err)
	}

	if len(msg.Content) == 0 {
		return "", fmt.Errorf("empty response from claude")
	}

	result := msg.Content[0].Text
	slog.Info("translate: done", "chars", len(result))
	return result, nil
}
