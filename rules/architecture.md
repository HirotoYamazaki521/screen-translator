# アーキテクチャ制約

詳細な設計・計画は [README.md](../README.md) を参照。

## データフロー

```
スクリーンキャプチャ → 差分検出 → Claude Vision API → Wails EventsEmit → フロントエンド
```

## パッケージ責務

| パッケージ | 責務 |
| --------- | ---- |
| `internal/config` | `~/.screen-translator/config.json` の読み書き |
| `internal/capture` | `kbinani/screenshot` を使った画面キャプチャ |
| `internal/translate` | Claude Vision API 呼び出し |
| `app.go` | Wailsバインディング・キャプチャループ・差分検出 |

## 重要な制約

- **APIキー**: 環境変数 `ANTHROPIC_API_KEY` のみ。config.json に保存しない
- **差分検出**: 前回と同じ画像ならAPIを叩かない（16x16ダウンスケール→MD5で判定）
- **クロスプラットフォーム**: ファイルパスは `os.UserHomeDir()` で構築する
- **タイムアウト**: Claude API は30秒。エラー時リトライなし（次インターバルで再試行）
- **キャプチャループ**: goroutineで動作。`context.CancelFunc` で停止制御

## 設定ファイル構造

```json
{
  "display_index": 0,
  "interval_seconds": 5,
  "capture_region": null,
  "language_from": "en",
  "language_to": "ja",
  "font_size": 18
}
```

`capture_region` が `null` の場合は全画面キャプチャ。指定する場合:

```json
{"x": 0, "y": 0, "width": 1920, "height": 1080}
```

## Wailsイベント名

```text
"translation:updated"  →  {text: string, timestamp: string}
"translation:status"   →  {status: "idle"|"capturing"|"translating"|"error", message?: string}
```
