# Git Hooks セットアップガイド

このプロジェクトでは、commit前に自動的にlintチェックを実行するGit hooksを設定できます。

## 目的

GHAのCIでlintが失敗するのを防ぐため、commit時にローカルで事前にチェックします。

## セットアップ方法

### 自動セットアップ（推奨）

プロジェクトルートで以下のスクリプトを実行してください：

```bash
./scripts/setup-git-hooks.sh
```

このスクリプトは以下を実行します：
1. `golangci-lint`がインストールされているか確認
2. 未インストールの場合はインストールを促す
3. pre-commitフックを`.git/hooks/`にセットアップ

### 手動セットアップ

#### 1. golangci-lintをインストール

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### 2. pre-commitフックをコピー

```bash
cp scripts/pre-commit.sample .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

## 動作確認

pre-commitフックが正しく動作しているか確認：

```bash
# 何かファイルを変更してステージング
echo "// test" >> main.go
git add main.go

# コミットを試みる（lint チェックが実行される）
git commit -m "test commit"
```

出力例：
```
Running pre-commit checks...
Checking Go files with golangci-lint...
✓ Lint check passed
```

## 使い方

### 通常のcommit

普通にコミットするだけで、自動的にlintチェックが実行されます：

```bash
git add .
git commit -m "feat: add new feature"
```

### Lintエラーがある場合

Lintエラーがある場合、コミットは中断されます：

```
✗ Lint check failed

Please fix the lint errors before committing.
You can run 'golangci-lint run' to see all errors.

To skip this check (not recommended), use: git commit --no-verify
```

エラーを修正してから再度コミットしてください：

```bash
# エラー内容を確認
golangci-lint run

# エラーを修正
vim path/to/file.go

# 再度コミット
git add .
git commit -m "feat: add new feature"
```

### Lintチェックをスキップ（非推奨）

緊急の場合など、どうしてもlintチェックをスキップしたい場合：

```bash
git commit --no-verify -m "WIP: temporary commit"
```

⚠️ **注意**: `--no-verify`を使用すると、CIでlintが失敗する可能性があります。

## トラブルシューティング

### golangci-lintが見つからない

```
Error: golangci-lint is not installed
```

**解決方法:**
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# PATHが通っているか確認
which golangci-lint

# PATHに追加が必要な場合（bashの例）
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

### pre-commitフックが実行されない

**確認事項:**
1. ファイルが存在するか: `ls -la .git/hooks/pre-commit`
2. 実行権限があるか: `chmod +x .git/hooks/pre-commit`
3. Goファイルがステージングされているか: `git diff --cached --name-only | grep '\.go$'`

### タイムアウトエラー

大きなプロジェクトでタイムアウトする場合、`.git/hooks/pre-commit`の`--timeout`値を増やしてください：

```bash
# 現在: --timeout=3m
# 変更後: --timeout=5m
golangci-lint run --timeout=5m
```

### Go 1.25のgo.mod調整について

このプロジェクトはGo 1.25を使用していますが、golangci-lintがGo 1.23でビルドされているため、pre-commitフックは一時的にgo.modを調整します。

**処理内容:**
1. commit前にgo.modをバックアップ
2. `go 1.25`を`go 1.23`に一時変更
3. lintを実行
4. go.modを元に戻す

この処理は自動で行われ、ユーザーの操作は不要です。

## チーム開発での推奨事項

### 1. READMEに記載

プロジェクトのREADMEに以下を追加することを推奨：

```markdown
## 開発環境のセットアップ

...

### Git Hooksのセットアップ

commit前に自動でlintチェックを実行するように設定できます：

\`\`\`bash
./scripts/setup-git-hooks.sh
\`\`\`

詳細は[Git Hooks セットアップガイド](docs/git-hooks.md)を参照してください。
```

### 2. オンボーディングチェックリストに追加

新メンバーのオンボーディング時に、Git hooksのセットアップを含めてください。

### 3. CI/CDとの一貫性

pre-commitフックのlint設定は、GHAのワークフローと同じ設定を使用しています。

### 4. 任意参加

Git hooksの使用は任意です。使いたくない人は無理に使う必要はありません。

## 設定ファイル

### .golangci.yml

プロジェクトルートの`.golangci.yml`でlintの設定を管理しています。
この設定は、ローカルのpre-commitフックとGHAの両方で使用されます。

## その他のhooks

将来的に追加できるhooksの例：

- **pre-push**: push前にテストを実行
- **commit-msg**: コミットメッセージの形式をチェック
- **post-merge**: merge後に依存関係を更新

必要に応じて`scripts/setup-git-hooks.sh`を拡張してください。

## 参考リンク

- [Git Hooks Documentation](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [GitHub Actions Workflow](.github/workflows/test.yml)
