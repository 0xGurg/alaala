#!/bin/bash

set -e

# ᜀᜎᜀᜎ (alaala) Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/0xGurg/alaala/main/install.sh | bash

REPO="0xGurg/alaala"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="alaala"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_error() {
    echo -e "${RED}Error: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_info() {
    echo -e "${YELLOW}$1${NC}"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            echo "darwin"
            ;;
        Linux*)
            echo "linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "windows"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
}

# Detect architecture
detect_arch() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Check if running as root (for Linux)
check_sudo() {
    if [ "$(id -u)" -ne 0 ]; then
        if [ -w "$INSTALL_DIR" ]; then
            return 0
        else
            print_info "Installation requires sudo access to write to $INSTALL_DIR"
            return 1
        fi
    fi
    return 0
}

# Main installation
main() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  ᜀᜎᜀᜎ (alaala) Installer"
    echo "  Semantic memory for AI assistants"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)
    
    print_info "Detected system: $OS ($ARCH)"
    
    # Construct download URL
    ARCHIVE_NAME="${BINARY_NAME}_${OS}_${ARCH}"
    
    if [ "$OS" = "windows" ]; then
        ARCHIVE_FILE="${ARCHIVE_NAME}.zip"
    else
        ARCHIVE_FILE="${ARCHIVE_NAME}.tar.gz"
    fi
    
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ARCHIVE_FILE}"
    
    print_info "Downloading from: $DOWNLOAD_URL"
    echo ""
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT
    
    # Download archive
    if command -v curl > /dev/null 2>&1; then
        curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARCHIVE_FILE"
    elif command -v wget > /dev/null 2>&1; then
        wget -q "$DOWNLOAD_URL" -O "$TMP_DIR/$ARCHIVE_FILE"
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    if [ ! -f "$TMP_DIR/$ARCHIVE_FILE" ]; then
        print_error "Download failed"
        exit 1
    fi
    
    print_success "Downloaded successfully"
    
    # Extract archive
    print_info "Extracting archive..."
    cd "$TMP_DIR"
    
    if [ "$OS" = "windows" ]; then
        unzip -q "$ARCHIVE_FILE"
    else
        tar -xzf "$ARCHIVE_FILE"
    fi
    
    if [ ! -f "$BINARY_NAME" ]; then
        print_error "Binary not found in archive"
        exit 1
    fi
    
    print_success "Extracted successfully"
    
    # Install binary
    print_info "Installing to $INSTALL_DIR..."
    
    # Check if we need sudo
    SUDO=""
    if [ "$OS" != "windows" ] && ! check_sudo; then
        SUDO="sudo"
    fi
    
    if [ -n "$SUDO" ]; then
        print_info "Requesting sudo access..."
    fi
    
    $SUDO mkdir -p "$INSTALL_DIR"
    $SUDO mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    $SUDO chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    print_success "Installed successfully!"
    echo ""
    
    # Verify installation
    if command -v alaala > /dev/null 2>&1; then
        VERSION=$(alaala version 2>&1 | head -n 1)
        print_success "✓ alaala installed: $VERSION"
    else
        print_error "Installation completed but alaala is not in PATH"
        print_info "Add $INSTALL_DIR to your PATH or restart your shell"
        exit 1
    fi
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    print_success "Installation complete!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Quick start:"
    echo "  1. Get an API key from OpenRouter (free tier available):"
    echo "     https://openrouter.ai"
    echo ""
    echo "  2. Initialize your project:"
    echo "     cd /path/to/your/project"
    echo "     alaala init"
    echo ""
    echo "  3. Configure for Cursor:"
    echo "     Add to Cursor Settings → MCP:"
    echo ""
    echo '     {
       "mcpServers": {
         "alaala": {
           "command": "'$INSTALL_DIR/$BINARY_NAME'",
           "args": ["serve"],
           "env": {
             "OPENROUTER_API_KEY": "your-key-here"
           }
         }
       }
     }'
    echo ""
    echo "Documentation: https://github.com/$REPO"
    echo "Support: https://github.com/$REPO/issues"
    echo ""
}

main "$@"

