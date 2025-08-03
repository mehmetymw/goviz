# GoViz - Go Dependency Analysis & Visualization

<div align="center">

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge\&logo=go\&logoColor=white)
![CLI](https://img.shields.io/badge/CLI-Tool-brightgreen?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)

**Analyze, visualize, and audit your Go dependencies for security, license compliance, and health.**

</div>

---

## 🔍 Features

* **Dependency Tree Generation** (ASCII, PNG, SVG, Graphviz DOT)
* **Security Analysis** (vulnerability scanning)
* **License Compliance** (detection, risk checking)
* **Health Monitoring** (last update, score, maintenance status)
* **CI/CD Friendly** (JSON/YAML outputs, non-zero exit codes)

---

## 🚀 Quick Start

### Installation

```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | bash
```

### Usage

```bash
goviz generate --format tree         # ASCII tree in terminal
goviz generate --format png -o out.png  # Visual diagram
goviz doctor                         # Health score + update info
goviz licenses                       # License analysis
goviz analyze --format json          # Full report in JSON
```

---

## 🎬 Demos

### ASCII Tree View

![Tree Demo](https://github.com/mehmetymw/goviz/blob/main/demos/output-demo.gif)

### Security Analysis

![Security Demo](https://github.com/mehmetymw/goviz/blob/main/demos/goviz-security-demo.gif)

---

## 📦 CI/CD Integration

### GitHub Actions

```yaml
- name: Install GoViz
  run: |
    git clone https://github.com/mehmetymw/goviz.git
    cd goviz && go build -o goviz . && sudo mv goviz /usr/local/bin/

- run: goviz doctor --format json --output health.json
- run: goviz licenses --format json --output licenses.json
```

---

**Built for the Go community – helping developers govern dependencies in the age of AI.**

---

İstersen `contributing`, `architecture`, veya `json output` gibi bölümleri de ayrıca "extended" versiyon olarak README'nin sonuna veya ayrı bir dosyaya alabiliriz. İlgileniyorsan onu da çıkarabilirim.
