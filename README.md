 # GoViz - Go Dependency Analysis & Visualization Tool

<div align="center">

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![CLI](https://img.shields.io/badge/CLI-Tool-brightgreen?style=for-the-badge)
![Production](https://img.shields.io/badge/Production-Ready-success?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)

**A comprehensive, production-ready CLI tool for Go dependency analysis, visualization, and compliance monitoring.**

*Born from the need to address AI-assisted development challenges in production environments.*

</div>

---

## 🎯 **Project Mission**

GoViz was created to solve real-world challenges in modern Go development, particularly those arising from AI-assisted coding practices. While AI tools have revolutionized development speed and accessibility, they also introduce new challenges that production teams must address.

### 🤖 **AI-Assisted Development: The Double-Edged Sword**

#### ✅ **Advantages of AI-Assisted Coding**
- **🚀 Rapid Prototyping**: Lightning-fast feature implementation
- **📚 Knowledge Democratization**: Access to best practices and patterns
- **🔄 Code Generation**: Automated boilerplate and repetitive code
- **🎯 Problem Solving**: Quick solutions to complex algorithmic challenges
- **📖 Documentation**: Auto-generated comprehensive documentation
- **🔧 Refactoring**: Intelligent code restructuring and optimization

#### ⚠️ **Critical Challenges & Risks**
- **🔒 Security Vulnerabilities**: AI may suggest insecure patterns or outdated practices
- **📊 Dependency Bloat**: Tendency to add unnecessary dependencies without proper analysis
- **🏗️ Architecture Inconsistency**: Mixing different architectural patterns within projects
- **📈 Repository Size Explosion**: Uncontrolled growth in project complexity
- **🔍 Lack of Oversight**: Missing comprehensive dependency analysis and health monitoring
- **⚖️ Compliance Issues**: Overlooked license compatibility and legal requirements
- **🎭 Hidden Technical Debt**: AI-generated code may accumulate unseen maintenance burdens

### 🎯 **How GoViz Addresses These Challenges**

GoViz was specifically designed to provide the **production-ready governance and oversight layer** that AI-assisted development often lacks:

#### 🛡️ **Security & Compliance**
```bash
goviz analyze              # Comprehensive security and dependency analysis
goviz licenses            # License compliance checking and risk assessment
```

#### 📊 **Dependency Health Monitoring**
```bash
goviz doctor              # Health assessment with update recommendations
goviz generate --format json  # CI/CD integration for automated monitoring
```

#### 🔍 **Visibility & Governance**
```bash
goviz generate --format png   # Visual dependency mapping for stakeholders
goviz analyze --format yaml   # Structured reporting for compliance teams
```

---

## 🎬 **Live Demos**

See GoViz in action with our interactive CLI demonstrations:

### 🌳 Dependency Tree, JSON, YAML Visualization
![GoViz Tree Demo](https://github.com/mehmetymw/goviz/blob/main/demos/output-demo.gif)
*Real-time ASCII tree generation showing Go module dependencies*

### 🔒 Security Vulnerability Analysis  
![GoViz Security Demo](https://github.com/mehmetymw/goviz/blob/main/demos/goviz-security-demo.gif)
*Comprehensive security scanning with JSON output for CI/CD integration*


---

## 🚀 **Features & Capabilities**

### 🔍 **Comprehensive Analysis Engine**
- **📋 Complete Dependency Parsing**: Deep analysis of `go.mod` and `go.sum` files
- **🌐 Transitive Dependency Resolution**: Maps complete dependency chains and relationships
- **⚡ Version Conflict Detection**: Identifies and reports potential compatibility issues
- **🔒 Security Framework**: Extensible foundation for vulnerability scanning integration
- **📊 Health Assessment**: Package maintenance status and update recommendations

### 🎨 **Multi-Format Visualization**
- **🌳 ASCII Tree**: Clean, terminal-friendly dependency trees
- **📈 Graphviz DOT**: Professional network diagrams with metadata
- **🖼️ PNG/SVG Export**: High-quality visual outputs for presentations
- **📄 JSON/YAML Reports**: Structured data for automation and CI/CD pipelines
- **🏷️ Enhanced Metadata**: License, security, and health information overlay

### ⚖️ **Production-Ready Compliance**
- **📜 License Analysis**: Automatic license detection and compatibility checking
- **🚨 Risk Assessment**: Identifies unknown licenses and potential legal issues
- **📋 Compliance Reporting**: Detailed breakdowns for legal and security teams
- **🔄 CI/CD Integration**: Automated compliance checking in development pipelines

### 🏥 **Dependency Health Monitoring**
- **📊 Health Scoring**: Comprehensive dependency health assessment
- **📅 Maintenance Tracking**: Last update dates and maintenance status
- **🔄 Update Recommendations**: Intelligent suggestions for dependency updates
- **⚠️ Risk Identification**: Flags for stale, abandoned, or problematic packages

---

## 📦 **Installation**

### **🚀 Quick Install** (Recommended)
```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | bash

# Install specific version
curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | GOVIZ_VERSION=v0.1.0 bash

# Install to custom directory
curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | INSTALL_DIR=~/bin bash
```

### **📦 Manual Installation**
Download pre-built binaries from [GitHub Releases](https://github.com/mehmetymw/goviz/releases):

```bash
# Linux/macOS
wget https://github.com/mehmetymw/goviz/releases/download/v0.1.0/goviz_0.1.0_linux_amd64.tar.gz
tar -xzf goviz_0.1.0_linux_amd64.tar.gz
sudo mv goviz /usr/local/bin/

# macOS via Homebrew (coming soon)
# brew install mehmetymw/tap/goviz
```

### **🔧 From Source**
```bash
git clone https://github.com/mehmetymw/goviz.git
cd goviz
go build -ldflags="-w -s" -o goviz .
sudo mv goviz /usr/local/bin/
```

### **📋 Prerequisites**
- **Optional: Graphviz** (for PNG/SVG generation):
  ```bash
  # Ubuntu/Debian
  sudo apt-get install graphviz
  
  # macOS
  brew install graphviz
  
  # CentOS/RHEL
  sudo yum install graphviz
  
  # Windows
  # Download from https://graphviz.org/download/
  ```

### **✅ Verify Installation**
```bash
goviz --version
goviz --help
```

---

## 🎮 **Usage Guide**

### **Quick Start**
```bash
# Basic dependency analysis
goviz generate

# Comprehensive health assessment  
goviz doctor

# License compliance check
goviz licenses

# Full analysis with all checks
goviz analyze
```

### **Visualization Outputs**
```bash
# Terminal-friendly tree view
goviz generate --format tree

# Professional network diagram
goviz generate --format png --output deps.png

# Scalable vector graphics
goviz generate --format svg --output deps.svg

# Graphviz DOT file for custom processing
goviz generate --format dot --output deps.dot
```

### **CI/CD Integration**
```bash
# JSON output for automated processing
goviz analyze --format json --output analysis.json

# YAML for configuration management
goviz generate --format yaml --output deps.yaml

# Exit codes for pipeline control
goviz doctor  # Returns non-zero on critical health issues
```

### **Advanced Analysis**
```bash
# Detailed health assessment with recommendations
goviz doctor --show-outdated

# License compliance with compatibility analysis
goviz licenses --check-compatibility

# Full dependency analysis with security checks
goviz analyze --format json | jq '.security_issues'
```

---

## 📊 **Command Reference**

### **🎯 Core Commands**

#### `goviz generate [path] [flags]`
**Primary visualization and dependency mapping**

```bash
Flags:
  -f, --format string   Output format: tree, ascii, dot, png, svg, json, yaml (default "tree")
  -o, --output string   Output file path
  -h, --help           Show command help
```

**Examples:**
```bash
goviz generate                                    # ASCII tree in terminal
goviz generate --format png --output graph.png   # PNG visualization
goviz generate /path/to/project --format json     # JSON report for specific project
```

#### `goviz analyze [path] [flags]`
**Comprehensive dependency analysis and reporting**

```bash
Flags:
  -f, --format string   Output format: text, json, yaml (default "text")
  -o, --output string   Output file path
  -h, --help           Show command help
```

**Analysis includes:**
- 📊 Dependency statistics and breakdown
- ⚡ Version conflict detection
- 🔒 Security issue identification
- 📜 License compatibility assessment
- 💡 Actionable recommendations

#### `goviz doctor [path] [flags]`
**Dependency health assessment and maintenance recommendations**

```bash
Flags:
  -f, --format string     Output format: text, json, yaml (default "text")
  -o, --output string     Output file path
      --show-outdated     Show detailed outdated package information (default true)
  -h, --help             Show command help
```

**Health Assessment:**
- 🎯 Overall health scoring (0-100)
- 📊 Package categorization (well-maintained, outdated, stale)
- 📅 Last update tracking
- 🔄 Update recommendations with commands
- ⚠️ Risk identification for critical packages

#### `goviz licenses [path] [flags]`
**License compliance analysis and risk assessment**

```bash
Flags:
  -f, --format string        Output format: text, json, yaml (default "text")
  -o, --output string        Output file path
      --check-compatibility  Perform license compatibility analysis (default true)
  -h, --help                Show command help
```

**Compliance Features:**
- 📜 Automatic license detection
- ⚖️ Compatibility risk assessment
- 🚨 Unknown license identification
- 📊 License distribution analysis
- 💼 Production compliance reporting

---

## 🏗️ **Architecture & Design**

### **Clean Architecture Principles**
```
goviz/
├── main.go                    # Clean entry point
├── cmd/                       # CLI interface layer
│   ├── root.go               # Root command with enhanced help
│   ├── generate.go           # Visualization command
│   ├── analyze.go            # Analysis command
│   ├── doctor.go             # Health assessment command
│   └── licenses.go           # License compliance command
├── pkg/
│   ├── parser/               # Data parsing layer
│   │   ├── modfile.go       # go.mod file parsing
│   │   └── gosum.go         # go.sum transitive dependency parsing
│   ├── graph/                # Domain logic layer
│   │   ├── dependency.go    # Core dependency graph structures
│   │   └── enhanced.go      # Advanced analysis capabilities
│   └── output/               # Presentation layer
│       ├── dot.go           # Graphviz DOT generation
│       ├── ascii.go         # Terminal tree visualization
│       ├── structured.go    # JSON/YAML structured output
│       └── png.go           # Image generation (PNG/SVG)
└── README.md                 # This documentation
```

### **Core Dependencies**
- **CLI Framework**: `github.com/spf13/cobra` - Industry-standard Go CLI framework
- **Go Module Parsing**: `golang.org/x/mod` - Official Go module manipulation
- **Graph Visualization**: `github.com/awalterschulze/gographviz` - Graphviz DOT generation
- **Structured Output**: `gopkg.in/yaml.v3` - YAML serialization
- **Terminal Enhancement**: `github.com/fatih/color` - Rich terminal output

### **Design Principles**
1. **🎯 Single Responsibility**: Each package has a clear, focused purpose
2. **🔌 Dependency Inversion**: Interfaces define contracts, implementations are swappable
3. **📦 Modularity**: Features can be enabled/disabled independently
4. **🧪 Testability**: Pure functions and dependency injection enable easy testing
5. **🔧 Extensibility**: Plugin architecture for custom analyzers and output formats

---

## 🔄 **CI/CD Integration Examples**

### **GitHub Actions Workflow**
```yaml
name: Dependency Analysis
on: [push, pull_request]

jobs:
  dependency-analysis:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.20'
      
      - name: Install GoViz
        run: |
          git clone https://github.com/mehmetymw/goviz.git
          cd goviz && go build -o goviz .
          sudo mv goviz /usr/local/bin/
      
      - name: Dependency Health Check
        run: goviz doctor --format json --output health-report.json
      
      - name: License Compliance Check
        run: goviz licenses --format json --output license-report.json
      
      - name: Upload Reports
        uses: actions/upload-artifact@v3
        with:
          name: dependency-reports
          path: |
            health-report.json
            license-report.json
```

### **Git CI Pipeline**
```yaml
dependency_analysis:
  stage: test
  image: golang:1.20
  before_script:
    - apt-get update && apt-get install -y graphviz
    - git clone https://github.com/mehmetymw/goviz.git
    - cd goviz && go build -o goviz . && mv goviz /usr/local/bin/
  script:
    - goviz analyze --format json --output analysis.json
    - goviz doctor --format yaml --output health.yaml
  artifacts:
    reports:
      junit: analysis.json
    paths:
      - analysis.json
      - health.yaml
```

### **Jenkins Pipeline**
```groovy
pipeline {
    agent any
    stages {
        stage('Dependency Analysis') {
            steps {
                script {
                    sh 'goviz analyze --format json --output analysis.json'
                    sh 'goviz doctor'  // Fails build on critical health issues
                    archiveArtifacts artifacts: 'analysis.json'
                }
            }
        }
    }
}
```

---

## 📈 **Output Format Examples**

### **JSON Report Structure**
```json
{
  "metadata": {
    "generated_at": "2024-01-15T10:30:00Z",
    "tool": "goviz",
    "version": "v0.1.0"
  },
  "module": {
    "name": "github.com/yourorg/yourproject",
    "go_version": "1.20",
    "path": "/path/to/project"
  },
  "statistics": {
    "total_dependencies": 45,
    "direct_dependencies": 12,
    "indirect_dependencies": 33,
    "unique_licenses": 8,
    "version_conflicts": 0,
    "security_issues": 1
  },
  "dependencies": [
    {
      "name": "github.com/gin-gonic/gin",
      "version": "v1.9.1",
      "direct": true,
      "license": "MIT",
      "health_score": 95,
      "last_update": "2023-12-15T00:00:00Z"
    }
  ],
  "licenses_summary": {
    "MIT": 15,
    "Apache-2.0": 12,
    "BSD-3-Clause": 8,
    "Unknown": 2
  }
}
```

### **Terminal ASCII Output**
```
🔍 Dependency Analysis Report
============================

Module: github.com/yourorg/yourproject
Go Version: 1.20

📊 Statistics:
  Total Dependencies: 45
  Direct Dependencies: 12
  Indirect Dependencies: 33
  Unique Licenses: 8

✅ No version conflicts detected

📄 License Summary:
  • MIT: 15 packages
  • Apache-2.0: 12 packages
  • BSD-3-Clause: 8 packages
  • Unknown: 2 packages

💡 Recommendations:
  • Review licenses for 2 unknown packages
  • Consider running 'go mod tidy' to clean up dependencies
  • Use 'goviz doctor' for detailed package health analysis
```

---

## 🛡️ **Security Considerations**

### **AI Development Challenges Addressed**
1. **🔍 Dependency Oversight**: AI tools often add dependencies without proper analysis
2. **📊 Bloat Prevention**: Identifies unnecessary or duplicate functionality
3. **🔒 Security Gaps**: Provides framework for vulnerability scanning integration
4. **⚖️ Compliance Blind Spots**: Ensures license compatibility in AI-generated solutions

### **Security Best Practices**
- **🔒 No Network Calls**: All analysis performed locally on existing files
- **📁 File System Safety**: Read-only operations on dependency files
- **🔍 Input Validation**: Robust parsing with error handling
- **🚫 No Code Execution**: Static analysis only, no dynamic code evaluation

---

## 🤝 **Contributing**

### **Development Setup**
```bash
# Clone the repository
git clone https://github.com/mehmetymw/goviz.git
cd goviz

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build locally
go build -o goviz .

# Run linting
golangci-lint run
```

### **Contributing Guidelines**
1. **🔧 Fork and Branch**: Create feature branches from `main`
2. **✅ Test Coverage**: Ensure new features have comprehensive tests
3. **📝 Documentation**: Update README and add inline documentation
4. **🎯 Single Responsibility**: Keep PRs focused on single features/fixes
5. **🔍 Code Review**: All changes require review before merging

### **Feature Requests & Bug Reports**
- **🐛 Bug Reports**: Use GitHub issues with reproduction steps
- **💡 Feature Requests**: Propose enhancements with use cases
- **🔒 Security Issues**: Report privately via email

---

## 🔮 **Roadmap & Future Enhancements**

### **🎯 Planned Features**
- **🔒 Security Integration**: Direct integration with Go vulnerability database
- **🌐 Real-time Health**: Live package health monitoring via API integration
- **📊 Trend Analysis**: Historical dependency health tracking
- **🔗 License Fetching**: Automatic license detection from source repositories
- **🎨 Interactive UI**: Web-based dashboard for large projects
- **📈 Metrics Export**: Prometheus/Grafana integration for monitoring

### **🔧 Technical Improvements**
- **⚡ Performance**: Parallel processing for large dependency trees
- **🧪 Testing**: Comprehensive test suite with integration tests
- **📦 Distribution**: Binary releases and package manager support
- **🔌 Plugin System**: Extensible architecture for custom analyzers

---

## 📜 **License**

MIT License - see [LICENSE](LICENSE) file for details.

---

## 🙏 **Acknowledgments**

- **Go Team**: For the excellent `golang.org/x/mod` package
- **Cobra Authors**: For the outstanding CLI framework
- **Graphviz Project**: For powerful graph visualization capabilities
- **Open Source Community**: For the ecosystem that makes projects like this possible
- **AI Community**: For highlighting the need for better governance tools in AI-assisted development

---

## 📞 **Support & Contact**

- **📧 Email**: [your-email@domain.com]
- **🐛 Issues**: [GitHub Issues](https://github.com/mehmetymw/goviz/issues)
- **💬 Discussions**: [GitHub Discussions](https://github.com/mehmetymw/goviz/discussions)
- **📱 Twitter**: [@yourhandle]

---

<div align="center">

**Built with ❤️ for the Go community**

*Empowering developers to maintain secure, healthy, and compliant codebases in the age of AI-assisted development.*

![Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?style=for-the-badge&logo=go)

</div>