# Distillery - Homebrew Tap by 0xGurg

Homebrew formulae for various packages.

## Packages

### ᜀᜎᜀᜎ (alaala)

Semantic memory system for AI assistants.

**Documentation:** https://github.com/0xGurg/alaala

## Installation

```bash
brew tap 0xGurg/distillery
brew install alaala
```

## Usage

```bash
# Initialize project
alaala init

# Start MCP server
alaala serve

# Show version
alaala version
```

## Uninstallation

```bash
brew uninstall alaala
brew untap 0xGurg/distillery
```

To also remove data and configuration:
```bash
rm -rf ~/.alaala
```

## Automatic Updates

Formulae in this tap are automatically updated by GoReleaser when new releases are published to their respective repositories.
