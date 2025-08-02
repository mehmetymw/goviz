#!/bin/bash

# GoViz Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
GOVIZ_VERSION="${GOVIZ_VERSION:-v0.1.0}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
GITHUB_REPO="mehmetymw/goviz"

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     OS=linux;;
        Darwin*)    OS=darwin;;
        CYGWIN*|MINGW*|MSYS*) OS=windows;;
        *)          OS="unknown";;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   ARCH=amd64;;
        arm64|aarch64)  ARCH=arm64;;
        armv7l)         ARCH=arm;;
        i386|i686)      ARCH=386;;
        *)              ARCH="unknown";;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download binary
download_binary() {
    local url="$1"
    local output="$2"
    
    print_status "Downloading GoViz from $url"
    
    if command_exists curl; then
        curl -fsSL "$url" -o "$output"
    elif command_exists wget; then
        wget -q "$url" -O "$output"
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Get latest version from GitHub
get_latest_version() {
    if [ "$GOVIZ_VERSION" = "latest" ]; then
        print_status "Fetching latest version..."
        if command_exists curl; then
            GOVIZ_VERSION=$(curl -fsSL "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
        elif command_exists wget; then
            GOVIZ_VERSION=$(wget -qO- "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
        else
            print_warning "Cannot fetch latest version. Using 'v0.1.0'"
            GOVIZ_VERSION="v0.1.0"
        fi
    fi
    
    print_status "Installing GoViz $GOVIZ_VERSION"
}

# Check prerequisites
check_prerequisites() {
    # Check if Go is available for building from source as fallback
    if ! command_exists go; then
        print_warning "Go is not installed. Will attempt to download pre-built binary."
    fi
    
    # Check for Graphviz (optional but recommended)
    if ! command_exists dot; then
        print_warning "Graphviz not found. Install it for PNG/SVG generation:"
        case "$OS" in
            linux)
                print_warning "  Ubuntu/Debian: sudo apt-get install graphviz"
                print_warning "  CentOS/RHEL: sudo yum install graphviz"
                ;;
            darwin)
                print_warning "  macOS: brew install graphviz"
                ;;
            windows)
                print_warning "  Windows: Download from https://graphviz.org/download/"
                ;;
        esac
    fi
}

# Install from pre-built binary
install_binary() {
    local binary_name="goviz"
    if [ "$OS" = "windows" ]; then
        binary_name="goviz.exe"
    fi
    
    local filename="goviz_v${GOVIZ_VERSION#v}_${OS}_${ARCH}"
    if [ "$OS" = "windows" ]; then
        filename="${filename}.zip"
    else
        filename="${filename}.tar.gz"
    fi
    
    local download_url="https://github.com/$GITHUB_REPO/releases/download/$GOVIZ_VERSION/$filename"
    local temp_dir=$(mktemp -d)
    local temp_file="$temp_dir/$filename"
    
    print_status "Downloading $filename"
    
    # Download the archive
    if ! download_binary "$download_url" "$temp_file"; then
        print_error "Failed to download binary. Falling back to building from source."
        build_from_source
        return
    fi
    
    # Extract the archive
    print_status "Extracting binary..."
    cd "$temp_dir"
    
    if [ "$OS" = "windows" ]; then
        if command_exists unzip; then
            unzip -q "$temp_file"
        else
            print_error "unzip not available. Cannot extract Windows binary."
            exit 1
        fi
    else
        tar -xzf "$temp_file"
    fi
    
    # Install the binary
    install_goviz_binary "$temp_dir/$binary_name"
    
    # Cleanup
    rm -rf "$temp_dir"
}

# Build from source as fallback
build_from_source() {
    if ! command_exists go; then
        print_error "Go is required to build from source. Please install Go first."
        exit 1
    fi
    
    print_status "Building GoViz from source..."
    
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # Clone repository
    print_status "Cloning repository..."
    if command_exists git; then
        git clone "https://github.com/$GITHUB_REPO.git" .
        if [ "$GOVIZ_VERSION" != "latest" ]; then
            git checkout "$GOVIZ_VERSION"
        fi
    else
        # Download source archive
        local source_url="https://github.com/$GITHUB_REPO/archive/main.tar.gz"
        download_binary "$source_url" "source.tar.gz"
        tar -xzf "source.tar.gz" --strip-components=1
    fi
    
    # Build
    print_status "Building binary..."
    go build -ldflags="-w -s" -o goviz .
    
    # Install
    install_goviz_binary "./goviz"
    
    # Cleanup
    cd - >/dev/null
    rm -rf "$temp_dir"
}

# Install the goviz binary
install_goviz_binary() {
    local binary_path="$1"
    
    # Check if install directory exists and is writable
    if [ ! -d "$INSTALL_DIR" ]; then
        print_error "Install directory $INSTALL_DIR does not exist."
        exit 1
    fi
    
    if [ ! -w "$INSTALL_DIR" ]; then
        print_warning "Install directory $INSTALL_DIR is not writable. Trying with sudo..."
        sudo cp "$binary_path" "$INSTALL_DIR/goviz"
        sudo chmod +x "$INSTALL_DIR/goviz"
    else
        cp "$binary_path" "$INSTALL_DIR/goviz"
        chmod +x "$INSTALL_DIR/goviz"
    fi
    
    print_success "GoViz installed to $INSTALL_DIR/goviz"
}

# Check if PATH includes install directory
check_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        print_warning "$INSTALL_DIR is not in your PATH."
        print_warning "Add it to your shell profile (.bashrc, .zshrc, etc.):"
        print_warning "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bashrc"
        print_warning "  source ~/.bashrc"
    fi
}

# Verify installation
verify_installation() {
    if command_exists goviz; then
        local version=$(goviz --version 2>/dev/null || echo "unknown")
        print_success "GoViz installed successfully!"
        print_status "Version: $version"
        print_status ""
        print_status "Quick start:"
        print_status "  goviz generate          # Generate dependency tree"
        print_status "  goviz analyze           # Comprehensive analysis"
        print_status "  goviz doctor             # Health assessment"
        print_status "  goviz --help             # Show all commands"
    else
        print_error "Installation verification failed. GoViz command not found."
        check_path
        exit 1
    fi
}

# Main installation function
main() {
    echo ""
    echo "┌─────────────────────────────────────────────┐"
    echo "│          GoViz Installer v0.1.0            │"
    echo "│   Go Dependency Analysis & Visualization   │"
    echo "└─────────────────────────────────────────────┘"
    echo ""
    
    # Detect system
    detect_os
    detect_arch
    
    print_status "Detected OS: $OS"
    print_status "Detected Architecture: $ARCH"
    
    # Validate system
    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        print_error "Unsupported OS ($OS) or architecture ($ARCH)"
        print_error "Supported: linux/darwin/windows on amd64/arm64"
        exit 1
    fi
    
    # Check prerequisites
    check_prerequisites
    
    # Get version
    get_latest_version
    
    # Install
    print_status "Installing to $INSTALL_DIR"
    
    # Try binary installation first, fallback to source
    install_binary
    
    # Verify installation
    verify_installation
    
    echo ""
    print_success "🎉 GoViz installation completed!"
    echo ""
}

# Handle different installation modes
case "${1:-}" in
    --help|-h)
        echo "GoViz Installation Script"
        echo ""
        echo "Usage:"
        echo "  curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | bash"
        echo ""
        echo "Environment Variables:"
        echo "  GOVIZ_VERSION    Version to install (default: latest)"
        echo "  INSTALL_DIR      Installation directory (default: /usr/local/bin)"
        echo ""
        echo "Examples:"
        echo "  # Install latest version"
        echo "  curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | bash"
        echo ""
        echo "  # Install specific version"
        echo "  curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | GOVIZ_VERSION=v0.1.0 bash"
        echo ""
        echo "  # Install to custom directory"
        echo "  curl -fsSL https://raw.githubusercontent.com/mehmetymw/goviz/main/install.sh | INSTALL_DIR=~/bin bash"
        exit 0
        ;;
    --uninstall)
        if [ -f "$INSTALL_DIR/goviz" ]; then
            rm -f "$INSTALL_DIR/goviz"
            print_success "GoViz uninstalled successfully!"
        else
            print_error "GoViz not found in $INSTALL_DIR"
        fi
        exit 0
        ;;
    *)
        main
        ;;
esac