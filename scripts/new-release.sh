#!/usr/bin/env bash
# =============================================================================
# scripts/new-release.sh — 创建新版本目录
#
# 使用：
#   bash scripts/new-release.sh v1.0.0
# =============================================================================

set -euo pipefail

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TEMPLATE_DIR="$PROJECT_ROOT/releases/template"
TARGET_DIR="$PROJECT_ROOT/releases/$VERSION"

if [[ -d "$TARGET_DIR" ]]; then
    echo "❌ 版本目录已存在：$TARGET_DIR"
    exit 1
fi

echo "📦 创建版本目录：$TARGET_DIR"
cp -r "$TEMPLATE_DIR" "$TARGET_DIR"

# 替换模板中的版本占位符
find "$TARGET_DIR" -name "*.md" -exec sed -i.bak "s/{version}/$VERSION/g" {} \;
find "$TARGET_DIR" -name "*.bak" -delete

# 创建子目录
mkdir -p "$TARGET_DIR/sql/ddl"
mkdir -p "$TARGET_DIR/sql/dml"
mkdir -p "$TARGET_DIR/sql/rollback"
mkdir -p "$TARGET_DIR/configs/diff"
mkdir -p "$TARGET_DIR/scripts/pre-deploy"
mkdir -p "$TARGET_DIR/scripts/post-deploy"
mkdir -p "$TARGET_DIR/scripts/rollback"

echo "✅ 版本目录已创建：releases/$VERSION/"
echo ""
echo "下一步："
echo "  1. 填写 releases/$VERSION/RELEASE_NOTES.md"
echo "  2. 添加 SQL 变更到 releases/$VERSION/sql/"
echo "  3. 添加配置变更到 releases/$VERSION/configs/"
echo "  4. 添加运维脚本到 releases/$VERSION/scripts/"
echo "  5. 完成 releases/$VERSION/checklist.md"
