#!/usr/bin/env bash
# =============================================================================
# scripts/install.sh — Harness Engineering 初始化脚本（Claude Code 专用）
#
# 说明：
#   本项目使用 Claude Code 作为 AI 编程工具。
#   .claude/ 目录是 Claude Code 的原生配置目录，无需额外渲染。
#
#   本脚本负责：
#   1. 验证 .claude/ 目录结构完整性
#   2. 验证 .service-matrix/dependencies.yaml 格式
#   3. 初始化本地配置（.harness/local.yaml）
#
# 使用：
#   bash scripts/install.sh
#   bash scripts/install.sh --validate-only   # 只验证，不修改
# =============================================================================

set -euo pipefail

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

VALIDATE_ONLY=false
while [[ $# -gt 0 ]]; do
    case $1 in
        --validate-only) VALIDATE_ONLY=true; shift ;;
        -h|--help)
            echo "Usage: $0 [--validate-only]"
            exit 0
            ;;
        *) echo -e "${RED}Unknown option: $1${NC}"; exit 1 ;;
    esac
done

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_ok() { echo -e "${GREEN}  ✅${NC} $1"; }
log_fail() { echo -e "${RED}  ❌${NC} $1"; }

ERRORS=0

echo ""
echo "============================================"
echo "  Harness Engineering — 初始化验证"
echo "  (Claude Code 专用)"
echo "============================================"
echo ""

# ============================================================
# Step 1: 验证 Claude Code 目录结构
# ============================================================
log_info "Step 1: 验证 .claude/ 目录结构..."

check_file() {
    local file="$1"
    local desc="$2"
    if [[ -f "$PROJECT_ROOT/$file" ]]; then
        log_ok "$desc ($file)"
    else
        log_fail "$desc 缺失：$file"
        ERRORS=$((ERRORS + 1))
    fi
}

check_file "CLAUDE.md" "项目规范文件"
check_file ".claude/commands/requirement-new.md" "需求新建命令"
check_file ".claude/commands/agentic-code-review.md" "代码审查命令"
check_file ".claude/commands/knowledge-extract-experience.md" "经验提取命令"
check_file ".claude/commands/service-deps.md" "服务依赖命令"
check_file ".claude/skills/managing-requirement-lifecycle.md" "需求生命周期 Skill"
check_file ".claude/skills/code-review-report.md" "代码审查 Skill"
check_file ".claude/skills/self-refinement.md" "Self-Refinement Skill"
check_file ".claude/skills/service-dependency-analyzer.md" "服务依赖分析 Skill"
check_file ".claude/skills/traceability-gate-checker.md" "追溯链校验 Skill"
check_file ".claude/agents/Definition/requirement-quality-reviewer.md" "需求评审 Agent"
check_file ".claude/agents/DetailDesign/detail-design-quality-reviewer.md" "设计评审 Agent"
check_file ".claude/agents/Implementation/code-review-preparer.md" "代码审查准备 Agent"

# ============================================================
# Step 2: 验证核心工程制品
# ============================================================
echo ""
log_info "Step 2: 验证核心工程制品..."

check_file "context/harness-framework/main-process-numbering.md" "五阶段+四门禁（唯一真相源）"
check_file "context/team/INDEX.md" "团队级知识库"
check_file "context/team/git-convention.md" "Git 规范"
check_file "context/team/error-code.md" "错误码规范"
check_file "context/team/logging.md" "日志规范"
check_file ".service-matrix/dependencies.yaml" "服务矩阵"
check_file "requirements/INDEX.md" "需求目录"
check_file "releases/INDEX.md" "版本管理目录"
check_file "releases/_template/RELEASE_NOTES.md" "版本说明模板"
check_file "releases/_template/checklist.md" "发布检查清单模板"

# ============================================================
# Step 3: 初始化本地配置（非 validate-only 模式）
# ============================================================
if [[ "$VALIDATE_ONLY" == "false" ]]; then
    echo ""
    log_info "Step 3: 初始化本地配置..."

    # 创建 .harness/local.yaml（本地配置，不提交）
    HARNESS_LOCAL="$PROJECT_ROOT/.harness/local.yaml"
    if [[ ! -f "$HARNESS_LOCAL" ]]; then
        mkdir -p "$(dirname "$HARNESS_LOCAL")"
        cat > "$HARNESS_LOCAL" << EOF
# Harness Engineering 本地配置（不提交到 git）
# 此文件已加入 .gitignore

# 当前活跃团队（覆盖 .service-matrix/dependencies.yaml 中的 default_team）
# active_team: layout

# 本地路径配置（覆盖占位符）
# business_repo: /path/to/your/business/repo
# idl_repo: /path/to/your/idl/repo
EOF
        log_ok "已创建本地配置：.harness/local.yaml"
    else
        log_ok "本地配置已存在：.harness/local.yaml"
    fi
fi

# ============================================================
# 结果汇总
# ============================================================
echo ""
echo "============================================"
if [[ $ERRORS -eq 0 ]]; then
    echo -e "${GREEN}  ✅ 验证通过！所有文件完整。${NC}"
    echo ""
    echo "下一步："
    echo "  1. 填写 .service-matrix/dependencies.yaml 中的服务信息"
    echo "  2. 更新 .harness/local.yaml 中的本地路径"
    echo "  3. 在 Claude Code 中输入 /requirement:new 开始第一个需求"
else
    echo -e "${RED}  ❌ 验证失败！发现 $ERRORS 个问题。${NC}"
    echo ""
    echo "请检查上述缺失文件，然后重新运行此脚本。"
    exit 1
fi
echo "============================================"
echo ""
