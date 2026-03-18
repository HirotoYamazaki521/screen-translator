# Goコーディング規約

## バージョン

- Go 1.22 以上を使用する

## パッケージ構成

```
internal/capture/    # スクリーンキャプチャのロジック
internal/translate/  # Claude API呼び出しのロジック
internal/config/     # 設定の読み書き
```

- `internal/` — 外部に公開しないパッケージ
- `pkg/` — 汎用ユーティリティ（必要になったら追加）
- `cmd/` — エントリポイント（Wailsの場合はmain.goで代替）

## 命名規約

- Go標準の MixedCaps / mixedCaps に従う
- インターフェース名は動詞+erが望ましい（例: `Capturer`, `Translator`）
- 略語は大文字で統一（例: `APIKey`, `URLPath`）

## エラーハンドリング

```go
// 新しいエラーを作る場合
errors.New("something went wrong")

// コンテキストを付加する場合
fmt.Errorf("capture failed: %w", err)

// panicは使わない（recoveryが必要な場面を除く）
```

- エラーは呼び出し元に返す。握りつぶさない
- ログ出力してから return するのではなく、上位でまとめてハンドリングする

## ログ

- `log/slog` パッケージを使用する（Go 1.21+標準）
- `fmt.Println` によるデバッグログは残さない

```go
slog.Info("capture started", "monitor", monitorIndex)
slog.Error("api call failed", "err", err)
```

## 設定管理

- APIキーは環境変数 `ANTHROPIC_API_KEY` から取得する
- その他の設定（インターバル・モニター選択など）はJSONファイルで管理する
- `github.com/joho/godotenv` で `.env` ファイルをサポートする

## 依存関係

- 外部パッケージは必要最小限にとどめる
- `go mod tidy` を定期的に実行してクリーンに保つ
