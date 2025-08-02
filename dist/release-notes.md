# GoViz v0.1.0

## 🚀 Features

- Comprehensive Go dependency analysis and visualization
- Multiple output formats: DOT, PNG, SVG, JSON, YAML, ASCII tree
- License compliance checking and analysis
- Dependency health assessment with recommendations
- CI/CD integration support

## 📦 Installation

### Quick Install (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | bash
```

### Manual Installation
Download the appropriate binary for your platform from the assets below.

## 🎯 Usage Examples

```bash
# Basic dependency analysis
goviz generate

# Health assessment
goviz doctor

# License compliance check
goviz licenses

# Generate PNG visualization
goviz generate --format png --output deps.png

# JSON output for CI/CD
goviz analyze --format json --output analysis.json
```

## 📊 Supported Platforms

- Linux (amd64, arm64, arm)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)

## 🔗 Links

- [Documentation](https://github.com/mehmetymw/goviz#readme)
- [Installation Guide](https://github.com/mehmetymw/goviz#installation)
- [Usage Examples](https://github.com/mehmetymw/goviz#usage)

---

**Full Changelog**: https://github.com/mehmetymw/goviz/compare/v0.9.0...v0.1.0
