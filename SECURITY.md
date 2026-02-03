# Security Policy

## Supported Versions

We release patches for security vulnerabilities. The following versions are currently supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of Outlier seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please do NOT:

- Open a public GitHub issue for security vulnerabilities
- Post about the vulnerability in public forums
- Attempt to exploit the vulnerability on production systems

### Please DO:

**Report security vulnerabilities via GitHub Security Advisories:**

1. Go to https://github.com/wingnut128/outlier-go/security/advisories
2. Click "Report a vulnerability"
3. Provide a detailed description of the vulnerability
4. Include steps to reproduce if possible
5. Suggest a fix if you have one

**Alternatively, email security reports to:**
- Email: [Your security contact email]
- Subject: "[SECURITY] Outlier Vulnerability Report"

### What to include in your report:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if any)
- Your name/handle (for acknowledgment)

### Response Timeline:

- **Initial Response:** Within 48 hours
- **Status Update:** Within 5 business days
- **Fix Timeline:** Varies by severity
  - Critical: 7 days
  - High: 14 days
  - Medium: 30 days
  - Low: 60 days

## Security Best Practices

### For Users

When using Outlier in production:

1. **Keep Updated:** Always use the latest version
2. **Validate Input:** When accepting user-provided data files, validate file sizes and formats
3. **Rate Limiting:** Implement rate limiting on API endpoints to prevent DoS
4. **Authentication:** Add authentication/authorization for production API deployments
5. **HTTPS Only:** Always use HTTPS in production environments
6. **Environment Variables:** Never commit API keys (like HONEYCOMB_API_KEY) to version control

### For Developers

When contributing:

1. **Never hardcode secrets:** Use environment variables
2. **Validate all inputs:** Check file sizes, percentile ranges, data formats
3. **Run security scans:** Use `gosec` and CodeQL before submitting PRs
4. **Follow Go security best practices:** https://golang.org/doc/security/
5. **Keep dependencies updated:** Dependabot will create PRs automatically

## Known Security Considerations

### Input Validation

- **Large Datasets:** The API accepts up to 100MB payloads. Consider your server's memory capacity
- **Percentile Range:** Input is validated (0-100) but always verify in your application
- **File Uploads:** CSV/JSON files are parsed - validate file sources in production

### Memory Usage

- The calculator creates a copy of input data for sorting (~1.5x memory usage)
- Very large datasets (>10M values) may require significant memory

### Denial of Service (DoS)

- Consider implementing:
  - Rate limiting on API endpoints
  - Maximum request size limits
  - Timeout configurations
  - Circuit breakers for high-load scenarios

## Security Features

### Current Implementations

- ✅ Input validation (empty datasets, percentile ranges)
- ✅ No SQL injection risk (no database)
- ✅ No XSS risk (API-only, no HTML rendering)
- ✅ CORS properly configured
- ✅ Dependencies regularly updated (Dependabot)
- ✅ Automated security scanning (CodeQL, Gosec)
- ✅ Safe concurrency (no shared mutable state)

### Automated Security Scanning

This repository uses:

1. **GitHub CodeQL:** Advanced semantic code analysis
2. **Gosec:** Go security checker
3. **Dependabot:** Automated dependency updates
4. **GitHub Security Advisories:** Vulnerability tracking

## Security Updates

Security updates are released as patch versions (e.g., 1.0.1) and are documented in:
- GitHub Security Advisories
- CHANGELOG.md
- Release notes

Subscribe to security updates:
- Watch this repository for "Security alerts"
- Enable notifications for GitHub Security Advisories

## Acknowledgments

We appreciate the security research community and will acknowledge researchers who responsibly disclose vulnerabilities (unless they prefer to remain anonymous).

## Questions?

For non-security-related questions, please open a regular GitHub issue or discussion.

---

*Last updated: 2026-02-03*
