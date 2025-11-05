#!/bin/bash

# Git hooksをセットアップするスクリプト
# チームメンバーが簡単にhooksを有効化できるようにするためのヘルパースクリプト

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
HOOKS_DIR="$PROJECT_ROOT/.git/hooks"

echo "Setting up Git hooks..."
echo ""

# golangci-lintがインストールされているか確認
if ! command -v golangci-lint &> /dev/null; then
    echo "⚠️  Warning: golangci-lint is not installed"
    echo ""
    echo "The pre-commit hook requires golangci-lint."
    echo "Please install it with:"
    echo ""
    echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    echo ""
    read -p "Do you want to install it now? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        echo "✓ golangci-lint installed successfully"
    else
        echo "Skipping golangci-lint installation"
        echo "Note: The pre-commit hook will fail until golangci-lint is installed"
    fi
    echo ""
fi

# pre-commitフックをコピー
PRE_COMMIT_HOOK="$HOOKS_DIR/pre-commit"

if [ -f "$PRE_COMMIT_HOOK" ]; then
    echo "⚠️  A pre-commit hook already exists"
    echo ""
    read -p "Do you want to overwrite it? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Skipping pre-commit hook setup"
        exit 0
    fi
fi

# pre-commitフックの内容を作成
cat > "$PRE_COMMIT_HOOK" << 'EOF'
#!/bin/bash

# pre-commit hook for running golangci-lint
# このフックはcommit前にlintチェックを実行します

set -e

echo "Running pre-commit checks..."

# golangci-lintがインストールされているか確認
if ! command -v golangci-lint &> /dev/null; then
    echo "Error: golangci-lint is not installed"
    echo "Please install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    exit 1
fi

# ステージングされたGoファイルを取得
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$STAGED_GO_FILES" ]; then
    echo "No Go files staged for commit"
    exit 0
fi

echo "Checking Go files with golangci-lint..."

# GHAと同じように一時的にgo.modを調整
GO_MOD_MODIFIED=false
if grep -q "^go 1.25" go.mod 2>/dev/null; then
    echo "Temporarily adjusting go.mod for golangci-lint compatibility..."
    cp go.mod go.mod.backup
    sed -i.tmp 's/^go 1.25/go 1.23/' go.mod
    sed -i.tmp '/^toolchain go1.25.0/d' go.mod
    rm -f go.mod.tmp
    GO_MOD_MODIFIED=true
fi

# golangci-lintを実行
# --new-from-rev=HEAD を使うとステージングされた変更のみをチェック
# ただし、全体をチェックする方が確実なので、通常のrunを使用
if golangci-lint run --timeout=3m; then
    echo "✓ Lint check passed"
    LINT_RESULT=0
else
    echo "✗ Lint check failed"
    echo ""
    echo "Please fix the lint errors before committing."
    echo "You can run 'golangci-lint run' to see all errors."
    echo ""
    echo "To skip this check (not recommended), use: git commit --no-verify"
    LINT_RESULT=1
fi

# go.modを元に戻す
if [ "$GO_MOD_MODIFIED" = true ]; then
    echo "Restoring go.mod..."
    mv go.mod.backup go.mod
fi

exit $LINT_RESULT
EOF

# 実行権限を付与
chmod +x "$PRE_COMMIT_HOOK"

echo "✓ Pre-commit hook has been set up successfully"
echo ""
echo "The hook will run golangci-lint before each commit."
echo ""
echo "To temporarily skip the hook, use: git commit --no-verify"
echo ""
echo "Done!"
