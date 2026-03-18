# screen-translator

英語画面をリアルタイムで日本語に翻訳するデスクトップアプリ。
詳細は [README.md](README.md) を参照。

## 技術スタック

| 役割 | 技術 |
|------|------|
| 言語 | Go 1.22+ |
| デスクトップフレームワーク | Wails v2 |
| フロントエンド | React + TypeScript + Tailwind CSS |
| スクリーンキャプチャ | `github.com/kbinani/screenshot` |
| 翻訳・OCR | Claude Vision API（`claude-sonnet-4-6`） |
| 設定管理 | JSON ローカルファイル + 環境変数 |

## 開発コマンド

```bash
wails dev        # 開発サーバー起動（ホットリロード）
wails build      # ビルド（Mac/Windows）
go mod tidy      # 依存関係整理
```

## 環境変数

```
ANTHROPIC_API_KEY=your_api_key_here
```

## コーディングルール

- [rules/go.md](rules/go.md) — Goコーディング規約
- [rules/architecture.md](rules/architecture.md) — アーキテクチャ制約
- [rules/ui.md](rules/ui.md) — UI/UX設計方針
- [rules/git.md](rules/git.md) — Gitワークフロー（ブランチ命名・PR説明生成）
