#!/bin/bash

# GoViz Release Builder
# Builds binaries for multiple platforms and architectures

set -e

# Configuration
VERSION="${VERSION:-v0.1.0}"
OUTPUT_DIR="dist"
BINARY_NAME="goviz"

# Platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Clean and create output directory
prepare_build() {
    print_status "Preparing build environment..."
    rm -rf "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"
}

# Build for a specific platform
build_platform() {
    local platform=$1
    local goos=$(echo $platform | cut -d'/' -f1)
    local goarch=$(echo $platform | cut -d'/' -f2)
    
    local output_name="$BINARY_NAME"
    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    local build_dir="$OUTPUT_DIR/${BINARY_NAME}_${VERSION#v}_${goos}_${goarch}"
    mkdir -p "$build_dir"
    
    print_status "Building for $goos/$goarch..."
    
    # Build binary
    GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 go build \
        -ldflags="-w -s -X main.version=$VERSION" \
        -o "$build_dir/$output_name" \
        .
    
    # Copy additional files
    cp README.md "$build_dir/"
    cp LICENSE "$build_dir/" 2>/dev/null || echo "# License file not found" > "$build_dir/LICENSE"
    
    # Create archive
    cd "$OUTPUT_DIR"
    if [ "$goos" = "windows" ]; then
        zip -r "${BINARY_NAME}_${VERSION#v}_${goos}_${goarch}.zip" "${BINARY_NAME}_${VERSION#v}_${goos}_${goarch}/"
    else
        tar -czf "${BINARY_NAME}_${VERSION#v}_${goos}_${goarch}.tar.gz" "${BINARY_NAME}_${VERSION#v}_${goos}_${goarch}/"
    fi
    cd - > /dev/null
    
    print_success "Built $goos/$goarch"
}

# Generate checksums
generate_checksums() {
    print_status "Generating checksums..."
    cd "$OUTPUT_DIR"
    
    # Create checksums file
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum *.tar.gz *.zip 2>/dev/null > checksums.txt || true
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 *.tar.gz *.zip 2>/dev/null > checksums.txt || true
    else
        echo "No SHA256 tool found, skipping checksums"
    fi
    
    cd - > /dev/null
}

# Create GitHub release notes
create_release_notes() {
    local notes_file="$OUTPUT_DIR/release-notes.md"
    
    cat > "$notes_file" << EOF
# GoViz $VERSION

## 🚀 Features

- Comprehensive Go dependency analysis and visualization
- Multiple output formats: DOT, PNG, SVG, JSON, YAML, ASCII tree
- License compliance checking and analysis
- Dependency health assessment with recommendations
- CI/CD integration support

## 📦 Installation

### Quick Install (Recommended)
\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | bash
\`\`\`

### Manual Installation
Download the appropriate binary for your platform from the assets below.

## 🎯 Usage Examples

\`\`\`bash
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
\`\`\`

## 📊 Supported Platforms

- Linux (amd64, arm64, arm)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)

## 🔗 Links

- [Documentation](https://github.com/mehmetymw/goviz#readme)
- [Installation Guide](https://github.com/mehmetymw/goviz#installation)
- [Usage Examples](https://github.com/mehmetymw/goviz#usage)

---

**Full Changelog**: https://github.com/mehmetymw/goviz/compare/v0.9.0...$VERSION
EOF

    print_success "Release notes created: $notes_file"
}

# Main build process
main() {
    echo ""
    echo "┌─────────────────────────────────────────────┐"
    echo "│            GoViz Release Builder            │"
    echo "│              Version: $VERSION               │"
    echo "└─────────────────────────────────────────────┘"
    echo ""
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo "Warning: Not in a git repository"
    fi
    
    # Prepare build environment
    prepare_build
    
    # Build for all platforms
    for platform in "${PLATFORMS[@]}"; do
        build_platform "$platform"
    done
    
    # Generate checksums
    generate_checksums
    
    # Create release notes
    create_release_notes
    
    echo ""
    print_success "🎉 Build completed successfully!"
    print_status "Output directory: $OUTPUT_DIR"
    print_status "Built platforms: ${#PLATFORMS[@]}"
    echo ""
    print_status "Files created:"
    ls -la "$OUTPUT_DIR/"
    echo ""
    print_status "Next steps:"
    echo "  1. Test binaries on target platforms"
    echo "  2. Create GitHub release with: gh release create $VERSION dist/*"
    echo "  3. Update documentation if needed"
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "GoViz Release Builder"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --clean        Clean dist directory only"
        echo ""
        echo "Environment Variables:"
        echo "  VERSION        Version to build (default: v0.1.0)"
        echo ""
        echo "Examples:"
        echo "  # Build version v0.1.0"
        echo "  VERSION=v0.1.0 $0"
        echo ""
        echo "  # Clean build artifacts"
        echo "  $0 --clean"
        exit 0
        ;;
    --clean)
        print_status "Cleaning build artifacts..."
        rm -rf "$OUTPUT_DIR"
        print_success "Clean completed"
        exit 0
        ;;
    *)
        main
        ;;
esac