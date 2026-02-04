# Contributing to Outlier

Thank you for your interest in contributing to Outlier! This document provides guidelines and instructions for contributing.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Environment details** (OS, Go version, etc.)
- **Code samples** or test cases if applicable
- **Screenshots** if relevant

Use the bug report issue template when available.

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear title and description**
- **Use case** - why is this enhancement valuable?
- **Proposed solution** - how would you implement it?
- **Alternatives considered**
- **Additional context** or examples

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Follow the development setup** instructions below
3. **Make your changes** following our coding standards
4. **Add tests** for any new functionality
5. **Update documentation** (README, CHANGELOG, etc.)
6. **Run all tests** and ensure they pass
7. **Commit with clear messages** following our commit guidelines
8. **Submit a pull request**

## Development Setup

### Prerequisites

- Go 1.25.5 or later
- Git
- Make
- Docker (optional, for container testing)

### Initial Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/outlier-go
cd outlier-go

# Add upstream remote
git remote add upstream https://github.com/wingnut128/outlier-go

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test
```

### Development Workflow

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes
# ... edit files ...

# Run tests frequently
make test

# Run linter
make lint

# Check test coverage
make coverage

# Build and test the binary
make build
./bin/outlier --version

# Commit your changes
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name

# Open a pull request on GitHub
```

## Coding Standards

### Go Style Guide

- Follow the official [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting (enforced by CI)
- Use `golangci-lint` for linting (enforced by CI)
- Keep functions focused and small
- Prefer clarity over cleverness

### Code Organization

```
outlier-go/
├── cmd/           # Command-line applications
├── internal/      # Private application code
├── pkg/           # Public library code
├── docs/          # Documentation
└── examples/      # Example files
```

### Naming Conventions

- **Files:** `snake_case.go` or `camelCase.go`
- **Packages:** Short, lowercase, no underscores
- **Functions:** `PascalCase` (exported), `camelCase` (unexported)
- **Constants:** `PascalCase` or `SCREAMING_SNAKE_CASE`
- **Variables:** `camelCase`

### Testing

- Write tests for all new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Name tests clearly: `TestFunctionName_Scenario`

Example:
```go
func TestCalculatePercentile_EmptySlice(t *testing.T) {
    values := []float64{}
    _, err := CalculatePercentile(values, 50.0)
    if err == nil {
        t.Error("expected error for empty slice")
    }
}
```

### Documentation

- Add godoc comments for all exported functions, types, and packages
- Update README.md for user-facing changes
- Update CHANGELOG.md for all changes
- Include examples in godoc where helpful

Example:
```go
// CalculatePercentile calculates the percentile value using linear interpolation.
// The percentile should be between 0 and 100.
// Returns an error if the values slice is empty or percentile is out of range.
//
// Example:
//   values := []float64{1, 2, 3, 4, 5}
//   result, err := CalculatePercentile(values, 95.0)
func CalculatePercentile(values []float64, percentile float64) (float64, error) {
    // ...
}
```

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat:** New feature
- **fix:** Bug fix
- **docs:** Documentation only
- **style:** Code style changes (formatting, etc.)
- **refactor:** Code refactoring
- **perf:** Performance improvements
- **test:** Adding or updating tests
- **chore:** Maintenance tasks

### Examples

```bash
feat(calculator): add support for weighted percentiles

Implement weighted percentile calculation using linear interpolation
with custom weights for each value.

Closes #123

---

fix(server): correct CORS header handling

The CORS middleware was not properly handling preflight requests.
Updated to use gin-contrib/cors with proper configuration.

Fixes #456
```

## Pull Request Process

1. **Update CHANGELOG.md** under the `[Unreleased]` section
2. **Ensure all tests pass** (`make test`)
3. **Ensure linter passes** (`make lint`)
4. **Update documentation** as needed
5. **Keep PRs focused** - one feature/fix per PR
6. **Write clear PR descriptions** including:
   - What changed and why
   - How to test the changes
   - Related issues
   - Breaking changes (if any)

### PR Review Process

- Maintainers will review your PR within 3-5 business days
- Address review comments by pushing new commits
- Once approved, a maintainer will merge your PR
- PRs must pass all CI checks before merging

## Testing

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/calculator -v

# With coverage
make coverage

# Stress tests
make stress

# Benchmarks
make bench
```

### Writing Tests

- Place tests in `*_test.go` files
- Use `testing.T` for unit tests
- Use `testing.B` for benchmarks
- Use table-driven tests for multiple scenarios
- Test edge cases and error conditions

## Performance

- Run benchmarks for performance-critical changes: `make bench`
- Avoid unnecessary allocations
- Profile before optimizing
- Document performance expectations in tests

## Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR:** Incompatible API changes
- **MINOR:** Backwards-compatible functionality
- **PATCH:** Backwards-compatible bug fixes

## Release Process

(For maintainers)

1. Update CHANGELOG.md with version and date
2. Update version in documentation
3. Create and push tag: `git tag v1.x.x && git push origin v1.x.x`
4. Create GitHub release with changelog
5. Announce release

## Getting Help

- **Questions:** Open a GitHub Discussion
- **Bugs:** Open a GitHub Issue
- **Security:** See [SECURITY.md](SECURITY.md)

## Recognition

Contributors are recognized in:
- Git commit history
- CHANGELOG.md (for significant contributions)
- GitHub contributors page

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Outlier! 🎉
