# Gitワークフロー規約

## ブランチ命名

| 種別 | 形式 | 例 |
| ---- | ---- | -- |
| 新機能 | `feature/xxx` | `feature/capture-loop` |
| バグ修正 | `fix/xxx` | `fix/diff-detection` |

## 実装完了時のフロー

ブランチ上での実装が完了したら、以下の順で実行する。

1. 変更をコミット
2. ブランチをプッシュ（`git push -u origin <branch>`）
3. `gh pr create` で PR を自動作成する

### PRタイトル

- **日本語**でわかりやすく書く
- 例: `internal/capture: スクリーンキャプチャと差分検出を追加`

### PR作成コマンド

```bash
gh pr create \
  --title "<日本語タイトル>" \
  --body "$(cat <<'EOF'
## 概要

<!-- 何をしたか1〜2行で -->

## 変更内容

- 変更点1
- 変更点2

## 動作確認

- [ ] 確認項目1
- [ ] 確認項目2

## 備考

<!-- 補足事項があれば -->
EOF
)"
```

- `tmp/pr-description.md` は作成しない（PR本文は `gh pr create` に直接渡す）
- PR作成後、GitHub の URL を出力して確認できるようにする
