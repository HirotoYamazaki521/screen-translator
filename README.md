# screen-translator

英語画面をリアルタイムで日本語に翻訳するデスクトップアプリ。
ゲームやアプリの画面をキャプチャし、Claude Vision APIで翻訳結果をサブウィンドウに表示する。

デュアルモニター環境を想定（モニター1でゲームなどを表示、モニター2で翻訳を表示）。

---

## 技術スタック

| 役割 | 技術 |
|------|------|
| 言語 | Go 1.22+ |
| デスクトップフレームワーク | Wails v2 |
| フロントエンド | React + TypeScript + Tailwind CSS |
| スクリーンキャプチャ | `github.com/kbinani/screenshot` |
| 翻訳・OCR | Claude Vision API（`claude-sonnet-4-6`） |
| 設定管理 | JSON ローカルファイル + 環境変数 |

---

## アーキテクチャ

### データフロー

```
[モニター1: ゲームなど]
        ↓ kbinani/screenshot でキャプチャ
[スクリーンショット画像]
        ↓ 差分検出（16x16 MD5ハッシュ）
        ↓ 変化あり時のみ
[base64エンコード]
        ↓ Claude Vision API (claude-sonnet-4-6)
[翻訳テキスト]
        ↓ Wails EventsEmit
[モニター2: 翻訳アプリ画面]
```

### OCR + 翻訳方針

**Claude Vision API 一本**（OCRと翻訳を1回のAPIコールで完結）を採用。

- シンプルな実装、Win/Mac クロスプラットフォーム対応
- ゲームUIの複雑な背景・特殊フォントでも高精度
- 差分検出により不要なAPIコールを抑制（コスト削減）

#### Claude へのプロンプト

```
この画像はゲームやアプリのスクリーンショットです。
画像内に表示されている英語のテキストをすべて抽出し、自然な日本語に翻訳してください。

- 翻訳結果のテキストのみを出力する（前置きや説明は不要）
- 画像にテキストが含まれない場合は「（テキストなし）」と出力する
- UI要素（ボタン名、メニュー項目）も含めて翻訳する
- ゲーム固有の用語はカタカナで表記する（例: "Health" → "ヘルス"）
- 元のレイアウト構造（改行・段落）をできるだけ保持する
```

### 差分検出アルゴリズム

```go
func hashImage(img image.Image) string {
    // 16x16にダウンスケール → グレースケール変換 → MD5
    // タイムスタンプ等の微細変化は意図的に無視（APIコスト節約）
}
```

---

## ディレクトリ構成

```
screen-translator/
├── main.go                          # Wailsエントリポイント
├── app.go                           # Wailsバインディング（フロントエンド公開メソッド群）
├── wails.json
├── go.mod / go.sum
├── build/
│   ├── darwin/Info.plist            # NSScreenCaptureUsageDescription を記述
│   └── windows/wails.exe.manifest
├── internal/
│   ├── config/config.go             # Config struct + Load/Save
│   ├── capture/capturer.go          # Capturer interface + ScreenCapturer
│   └── translate/translator.go      # Translator interface + ClaudeTranslator
└── frontend/
    └── src/
        ├── App.tsx
        ├── components/
        │   ├── StatusBar.tsx
        │   ├── TranslationArea.tsx  # フェードインアニメーション
        │   ├── SettingsPanel.tsx
        │   ├── Toast.tsx
        │   └── Footer.tsx
        └── hooks/
            ├── useTranslation.ts    # Wailsイベントリスナー
            └── useSettings.ts       # 設定読み書き
```

---

## セットアップ

### 前提条件

- Go 1.22+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- macOS: Xcodeコマンドラインツール（`xcode-select --install`）
- `ANTHROPIC_API_KEY` 環境変数

### 開発サーバー起動

```bash
# 依存関係インストール
go mod tidy

# 開発サーバー起動（ホットリロード）
wails dev
```

### 環境変数の設定

```bash
# macOS (~/.zshenv)
export ANTHROPIC_API_KEY=your_api_key_here

# または .env ファイル（プロジェクトルート）
ANTHROPIC_API_KEY=your_api_key_here
```

---

## 設定ファイル

保存場所: `~/.screen-translator/config.json`（初回起動時に自動生成）

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

| フィールド | 説明 | デフォルト |
|-----------|------|----------|
| `display_index` | キャプチャするモニターのインデックス | `0` |
| `interval_seconds` | キャプチャ間隔（秒） | `5` |
| `capture_region` | キャプチャ範囲（null=全画面） | `null` |
| `language_from` | 翻訳元言語 | `"en"` |
| `language_to` | 翻訳先言語 | `"ja"` |
| `font_size` | 翻訳テキストのフォントサイズ（px） | `18` |

APIキーは config.json に保存しない（環境変数 `ANTHROPIC_API_KEY` のみ）。

---

## Wailsバインディング設計

### app.go の公開メソッド

```go
GetConfig() config.Config
SaveConfig(cfg config.Config) error
GetAPIKeyStatus() bool   // APIキーが設定されているかだけ返す（キー自体は返さない）
StartCapture() error
StopCapture()
IsCapturing() bool
GetDisplayCount() int
```

### フロントエンドへのイベント

| イベント名 | ペイロード | タイミング |
|-----------|-----------|----------|
| `translation:updated` | `{text: string, timestamp: string}` | 翻訳結果が更新されたとき |
| `translation:status` | `{status: "idle"\|"capturing"\|"translating"\|"error", message?: string}` | ステータスが変化したとき |

---

## ビルドと配布

### macOS

```bash
# Universal Binary（Intel + Apple Silicon両対応）
wails build -platform darwin/universal

# DMGの作成
hdiutil create -volname "Screen Translator" \
    -srcfolder build/bin/screen-translator.app \
    -ov -format UDZO \
    dist/screen-translator-mac.dmg
```

**初回起動時の注意**: 公証なしのアプリのため「開発元不明」の警告が出る。
右クリック→「開く」で回避できる旨をユーザーに伝える。

macOSのスクリーンキャプチャ権限: システム設定 → プライバシーとセキュリティ → 画面収録 でアプリを許可する。

### Windows

```bash
wails build -platform windows/amd64
```

WebView2が必要（Windows 10 21H2以降は標準搭載）。

---

## コスト感（Claude Vision API）

`claude-sonnet-4-6` の画像トークンコスト目安:
- 1080p スクリーンショット ≒ 1500〜2000トークン
- 5秒インターバル・1時間プレイ（差分あり時のみAPIコール）≒ $0.5〜1.0

差分検出によりゲームの静止シーン（メニュー待機など）では API コールをスキップするため、実際のコストはさらに低くなる。

---

## 実装フェーズ

### Phase 1 — プロジェクト骨格
ブランチ: `feature/wails-init`

- [x] `wails init` でプロジェクト生成（React + TypeScript テンプレート）
- [ ] `go.mod` に依存関係追加（`kbinani/screenshot`, `anthropic-go`）
- [ ] `wails dev` で起動確認

### Phase 2 — Goバックエンド
各モジュールを依存関係の順に実装する。

| ブランチ | 内容 |
|---------|------|
| `feature/backend-config` | `internal/config/config.go` — Config struct + Load/Save |
| `feature/backend-capture` | `internal/capture/capturer.go` — スクリーンキャプチャ + 差分検出 |
| `feature/backend-translate` | `internal/translate/translator.go` — Claude Vision API呼び出し |
| `feature/backend-app` | `app.go` — 上記3つを組み合わせてWailsに公開 |

### Phase 3 — フロントエンド
ブランチ: `feature/frontend-components`

- [ ] Tailwind CSS設定
- [ ] コンポーネント（`TranslationArea`, `StatusBar`, `SettingsPanel`, `Toast`, `Footer`）
- [ ] hooks（`useTranslation`, `useSettings`）
- [ ] フェードインアニメーション

### Phase 4 — ビルド・配布
ブランチ: `feature/build-dist`

- [ ] macOS: Universal Binary → DMG
- [ ] Windows: exe
