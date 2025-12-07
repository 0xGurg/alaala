#!/bin/bash

set -e

# ᜀᜎᜀᜎ (alaala) Uninstallation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/0xGurg/alaala/main/uninstall.sh | bash
# Or with flags: ./uninstall.sh --dry-run

BINARY_PATH="/usr/local/bin/alaala"
DATA_DIR="$HOME/.alaala"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Flags
DRY_RUN=false

# Parse arguments
for arg in "$@"; do
    case $arg in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
    esac
done

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

print_header() {
    echo -e "${BLUE}$1${NC}"
}

# Calculate directory size
get_dir_size() {
    if [ -d "$1" ]; then
        du -sh "$1" 2>/dev/null | cut -f1
    else
        echo "0B"
    fi
}

# Get file size
get_file_size() {
    if [ -f "$1" ]; then
        du -sh "$1" 2>/dev/null | cut -f1
    else
        echo "0B"
    fi
}

# Check if Weaviate container exists
check_weaviate() {
    if command -v docker > /dev/null 2>&1; then
        docker ps -a --format '{{.Names}}' | grep -q "^weaviate$"
        return $?
    fi
    return 1
}

# Show what will be removed
show_removal_plan() {
    echo ""
    print_header "The following items will be analyzed for removal:"
    echo ""
    
    # Check binary
    if [ -f "$BINARY_PATH" ]; then
        echo "  ✓ Binary: $BINARY_PATH"
    else
        echo "  ✗ Binary: Not found"
    fi
    
    # Check data directory
    if [ -d "$DATA_DIR" ]; then
        echo "  ✓ Data directory: $DATA_DIR ($(get_dir_size "$DATA_DIR"))"
        
        if [ -f "$DATA_DIR/config.yaml" ]; then
            echo "    - config.yaml"
        fi
        if [ -f "$DATA_DIR/alaala.db" ]; then
            echo "    - alaala.db ($(get_file_size "$DATA_DIR/alaala.db"))"
        fi
        if [ -f "$DATA_DIR/alaala.log" ]; then
            echo "    - alaala.log ($(get_file_size "$DATA_DIR/alaala.log"))"
        fi
        if [ -d "$DATA_DIR/weaviate-data" ]; then
            echo "    - weaviate-data/ ($(get_dir_size "$DATA_DIR/weaviate-data"))"
        fi
    else
        echo "  ✗ Data directory: Not found"
    fi
    
    # Check Weaviate container
    if check_weaviate; then
        WEAVIATE_STATUS=$(docker inspect -f '{{.State.Status}}' weaviate 2>/dev/null || echo "unknown")
        echo "  ✓ Weaviate container: weaviate (status: $WEAVIATE_STATUS)"
    else
        echo "  ✗ Weaviate container: Not found"
    fi
    
    echo ""
}

# Create backup
create_backup() {
    if [ ! -d "$DATA_DIR" ]; then
        print_info "No data to backup"
        return 0
    fi
    
    BACKUP_DIR="$HOME/.alaala-backup-$(date +%Y%m%d-%H%M%S)"
    
    if [ "$DRY_RUN" = true ]; then
        print_info "[DRY RUN] Would create backup at: $BACKUP_DIR"
        return 0
    fi
    
    print_info "Creating backup..."
    cp -r "$DATA_DIR" "$BACKUP_DIR"
    print_success "Backup created at: $BACKUP_DIR"
    echo ""
}

# Remove binary
remove_binary() {
    if [ ! -f "$BINARY_PATH" ]; then
        print_info "Binary not found, skipping"
        return 0
    fi
    
    if [ "$DRY_RUN" = true ]; then
        print_info "[DRY RUN] Would remove: $BINARY_PATH"
        return 0
    fi
    
    print_info "Removing binary..."
    
    # Check if we need sudo
    if [ -w "$(dirname "$BINARY_PATH")" ]; then
        rm "$BINARY_PATH"
    else
        print_info "Requesting sudo access to remove binary..."
        sudo rm "$BINARY_PATH"
    fi
    
    print_success "Binary removed: $BINARY_PATH"
}

# Remove data directory
remove_data() {
    if [ ! -d "$DATA_DIR" ]; then
        print_info "Data directory not found, skipping"
        return 0
    fi
    
    if [ "$DRY_RUN" = true ]; then
        print_info "[DRY RUN] Would remove: $DATA_DIR"
        return 0
    fi
    
    print_info "Removing data directory..."
    rm -rf "$DATA_DIR"
    print_success "Data directory removed: $DATA_DIR"
}

# Remove Weaviate container
remove_weaviate() {
    if ! check_weaviate; then
        print_info "Weaviate container not found, skipping"
        return 0
    fi
    
    if [ "$DRY_RUN" = true ]; then
        print_info "[DRY RUN] Would remove Weaviate container"
        return 0
    fi
    
    print_info "Removing Weaviate container..."
    docker stop weaviate 2>/dev/null || true
    docker rm weaviate 2>/dev/null || true
    print_success "Weaviate container removed"
}

# Main uninstall function
main() {
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  ᜀᜎᜀᜎ (alaala) Uninstaller"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    if [ "$DRY_RUN" = true ]; then
        print_info "DRY RUN MODE - No files will be deleted"
    fi
    
    # Show what will be removed
    show_removal_plan
    
    # Interactive menu
    echo "What would you like to remove?"
    echo ""
    echo "  [1] alaala binary only"
    echo "  [2] Binary + configuration"
    echo "  [3] Everything (binary, config, data, Weaviate)"
    echo "  [4] Create backup and remove everything"
    echo "  [c] Cancel"
    echo ""
    read -p "Choice: " -r CHOICE
    
    case $CHOICE in
        1)
            echo ""
            print_header "Removing binary only..."
            remove_binary
            ;;
        2)
            echo ""
            print_header "Removing binary and configuration..."
            remove_binary
            remove_data
            ;;
        3)
            echo ""
            print_info "This will remove EVERYTHING. Type 'yes' to confirm: "
            read -r CONFIRM
            if [ "$CONFIRM" != "yes" ]; then
                print_error "Aborted"
                exit 1
            fi
            echo ""
            print_header "Removing everything..."
            remove_binary
            remove_data
            remove_weaviate
            ;;
        4)
            echo ""
            print_header "Creating backup and removing everything..."
            create_backup
            
            print_info "Type 'yes' to confirm removal: "
            read -r CONFIRM
            if [ "$CONFIRM" != "yes" ]; then
                print_error "Aborted (backup preserved)"
                exit 1
            fi
            
            remove_binary
            remove_data
            remove_weaviate
            ;;
        [cC])
            print_info "Cancelled"
            exit 0
            ;;
        *)
            print_error "Invalid choice"
            exit 1
            ;;
    esac
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    print_success "Uninstallation complete!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    if [ "$CHOICE" = "4" ] && [ -d "$BACKUP_DIR" ]; then
        print_info "Your data was backed up to: $BACKUP_DIR"
        echo ""
    fi
    
    echo "Thank you for using alaala!"
    echo "To reinstall: curl -fsSL https://raw.githubusercontent.com/0xGurg/alaala/main/install.sh | bash"
    echo ""
}

main "$@"

