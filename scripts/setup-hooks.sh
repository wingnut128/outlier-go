#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up Git hooks for Outlier...${NC}"

# Check if we're in a git repository
if [ ! -d .git ]; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Create pre-commit hook
echo -e "${YELLOW}Installing pre-commit hook...${NC}"
cat > .git/hooks/pre-commit << 'HOOK_EOF'
#!/bin/bash

# Pre-commit hook for Outlier Go project
# This hook runs before each commit to ensure code quality

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Running pre-commit checks...${NC}\n"

# Check if required tools are available
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}Error: $1 is not installed${NC}"
        echo "Please install $1 and try again"
        exit 1
    fi
}

# Check for Go
check_tool go

# 1. Auto-format code
echo -e "${YELLOW}[1/6] Auto-formatting code...${NC}"
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)
if [ -n "$STAGED_GO_FILES" ]; then
    echo "$STAGED_GO_FILES" | xargs gofmt -w
    echo "$STAGED_GO_FILES" | xargs git add
    echo -e "${GREEN}✓ Code formatted and staged${NC}\n"
else
    echo -e "${GREEN}✓ No Go files to format${NC}\n"
fi

# 2. Run linter
echo -e "${YELLOW}[2/6] Running linter...${NC}"
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run --timeout=3m; then
        echo -e "${RED}✗ Linter found issues${NC}"
        echo -e "${YELLOW}Some issues may be auto-fixable with: golangci-lint run --fix${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Linter passed${NC}\n"
else
    echo -e "${YELLOW}⚠ golangci-lint not installed, skipping...${NC}"
    echo -e "${YELLOW}Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest${NC}\n"
fi

# 3. Go vet
echo -e "${YELLOW}[3/6] Running go vet...${NC}"
if ! go vet ./...; then
    echo -e "${RED}✗ go vet failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ go vet passed${NC}\n"

# 4. Run tests
echo -e "${YELLOW}[4/6] Running tests...${NC}"
if ! go test -short ./...; then
    echo -e "${RED}✗ Tests failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Tests passed${NC}\n"

# 5. Check for common issues
echo -e "${YELLOW}[5/6] Checking for common issues...${NC}"

# Check for TODO/FIXME in staged files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)
if [ -n "$STAGED_GO_FILES" ]; then
    # Check for debugging statements
    if echo "$STAGED_GO_FILES" | xargs grep -n "fmt.Println\|log.Println" 2>/dev/null | grep -v "_test.go"; then
        echo -e "${YELLOW}Warning: Found debug print statements${NC}"
        echo "Consider removing them before committing"
        read -p "Continue anyway? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi

    # Check for missing error handling
    if echo "$STAGED_GO_FILES" | xargs grep -n "err :=" 2>/dev/null | grep -v "if err"; then
        echo -e "${YELLOW}Warning: Possible unhandled errors detected${NC}"
        read -p "Continue anyway? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
fi
echo -e "${GREEN}✓ Common checks passed${NC}\n"

# 6. Build check
echo -e "${YELLOW}[6/6] Checking if project builds...${NC}"
if ! go build -o /tmp/outlier-build-test ./cmd/outlier > /dev/null 2>&1; then
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi
rm -f /tmp/outlier-build-test
echo -e "${GREEN}✓ Build successful${NC}\n"

echo -e "${GREEN}All pre-commit checks passed! ✓${NC}"
echo -e "${BLUE}Proceeding with commit...${NC}\n"
HOOK_EOF

chmod +x .git/hooks/pre-commit
echo -e "${GREEN}✓ Pre-commit hook installed${NC}"

# Create pre-push hook
echo -e "${YELLOW}Installing pre-push hook...${NC}"
cat > .git/hooks/pre-push << 'HOOK_EOF'
#!/bin/bash

# Pre-push hook for Outlier Go project
# This hook runs before pushing to ensure comprehensive quality

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Running pre-push checks...${NC}\n"

# 1. Run all tests (including non-short)
echo -e "${YELLOW}[1/3] Running full test suite...${NC}"
if ! go test ./...; then
    echo -e "${RED}✗ Tests failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ All tests passed${NC}\n"

# 2. Run linter if available
echo -e "${YELLOW}[2/3] Running linter...${NC}"
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run --timeout=3m; then
        echo -e "${RED}✗ Linter found issues${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Linter passed${NC}\n"
else
    echo -e "${YELLOW}⚠ golangci-lint not found, skipping...${NC}\n"
fi

# 3. Check test coverage
echo -e "${YELLOW}[3/3] Checking test coverage...${NC}"
COVERAGE=$(go test -coverprofile=/tmp/coverage.out ./... 2>&1 | grep "coverage:" | awk '{print $2}' | sed 's/%//' | head -1)
rm -f /tmp/coverage.out

if [ -n "$COVERAGE" ]; then
    COVERAGE_INT=${COVERAGE%.*}
    if [ "$COVERAGE_INT" -lt 70 ]; then
        echo -e "${YELLOW}⚠ Warning: Test coverage is ${COVERAGE}% (target: 70%+)${NC}"
        read -p "Continue anyway? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        echo -e "${GREEN}✓ Test coverage: ${COVERAGE}%${NC}\n"
    fi
fi

echo -e "${GREEN}All pre-push checks passed! ✓${NC}"
echo -e "${BLUE}Proceeding with push...${NC}\n"
HOOK_EOF

chmod +x .git/hooks/pre-push
echo -e "${GREEN}✓ Pre-push hook installed${NC}"

# Create commit-msg hook for conventional commits
echo -e "${YELLOW}Installing commit-msg hook...${NC}"
cat > .git/hooks/commit-msg << 'HOOK_EOF'
#!/bin/bash

# Commit message hook for Outlier Go project
# Validates commit messages follow Conventional Commits format

# Colors
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

COMMIT_MSG_FILE=$1
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Skip merge commits
if echo "$COMMIT_MSG" | head -1 | grep -q "^Merge"; then
    exit 0
fi

# Conventional Commits pattern
PATTERN="^(feat|fix|docs|style|refactor|perf|test|chore|build|ci|revert)(\([a-z0-9_-]+\))?: .{1,72}"

if ! echo "$COMMIT_MSG" | head -1 | grep -Eq "$PATTERN"; then
    echo -e "${RED}Error: Commit message doesn't follow Conventional Commits format${NC}"
    echo ""
    echo -e "${YELLOW}Format: <type>(<scope>): <subject>${NC}"
    echo ""
    echo "Types: feat, fix, docs, style, refactor, perf, test, chore, build, ci, revert"
    echo ""
    echo "Examples:"
    echo "  feat(calculator): add weighted percentile support"
    echo "  fix(server): correct CORS header handling"
    echo "  docs: update README with new examples"
    echo ""
    echo "Your commit message:"
    echo "$COMMIT_MSG" | head -1
    exit 1
fi

# Check subject line length
SUBJECT_LENGTH=$(echo "$COMMIT_MSG" | head -1 | wc -c)
if [ "$SUBJECT_LENGTH" -gt 72 ]; then
    echo -e "${YELLOW}Warning: Commit subject is longer than 72 characters${NC}"
    echo "Consider shortening it or adding details to the body"
fi

exit 0
HOOK_EOF

chmod +x .git/hooks/commit-msg
echo -e "${GREEN}✓ Commit-msg hook installed${NC}"

echo ""
echo -e "${GREEN}===========================================${NC}"
echo -e "${GREEN}Git hooks successfully installed!${NC}"
echo -e "${GREEN}===========================================${NC}"
echo ""
echo -e "Installed hooks:"
echo -e "  ${BLUE}pre-commit${NC}  - Auto-format, lint, vet, test, build"
echo -e "  ${BLUE}pre-push${NC}    - Full tests, linter, coverage check"
echo -e "  ${BLUE}commit-msg${NC}  - Conventional Commits validation"
echo ""
echo -e "To skip hooks temporarily:"
echo -e "  ${YELLOW}git commit --no-verify${NC}"
echo -e "  ${YELLOW}git push --no-verify${NC}"
echo ""
echo -e "${GREEN}Happy coding! 🚀${NC}"
