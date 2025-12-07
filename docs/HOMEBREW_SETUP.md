# Homebrew Tap Setup - Using Distillery

## Quick Setup Guide

Your distillery repository (`git@github.com:0xGurg/distillery.git`) will be used as the Homebrew tap.

### Step 1: Push Formula to Distillery

```bash
# Clone your distillery repo
cd ~/projects
git clone git@github.com:0xGurg/distillery.git
cd distillery

# Copy the formula files from alaala repo
cp -r ~/projects/alaala/homebrew-tap/* .

# Commit and push
git add .
git commit -m "Add alaala Homebrew formula"
git push origin main
```

### Step 2: Create Fine-Grained Personal Access Token

1. Go to https://github.com/settings/personal-access-tokens/new
2. Token name: `alaala-homebrew-tap`
3. Expiration: 1 year (recommended) or custom
4. Repository access:
   - Select: **Only select repositories**
   - Choose: `0xGurg/distillery`
5. Permissions:
   - Repository permissions:
     - ✅ Contents: **Read and write**
6. Click "Generate token"
7. **Copy the token** (starts with `github_pat_`)

### Step 3: Add Secret to alaala Repository

1. Go to https://github.com/0xGurg/alaala/settings/secrets/actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Value: Paste the token from Step 2
5. Click "Add secret"

### Step 4: Wait for v0.1.0 Release to Complete

The GitHub Actions workflow for v0.1.0 is currently running. Once it completes:

1. Go to https://github.com/0xGurg/alaala/releases/tag/v0.1.0
2. Download `checksums.txt`
3. Find the SHA256 values for each platform:
   - `alaala_darwin_arm64.tar.gz`
   - `alaala_darwin_amd64.tar.gz`
   - `alaala_linux_arm64.tar.gz`
   - `alaala_linux_amd64.tar.gz`

### Step 5: Update SHA256 Checksums in Formula

```bash
cd ~/projects/distillery

# Edit Formula/alaala.rb
# Replace each PLACEHOLDER_SHA256_* with actual values from checksums.txt

git add Formula/alaala.rb
git commit -m "Update alaala formula with v0.1.0 checksums"
git push origin main
```

### That's It!

Users can now install with:
```bash
brew tap 0xGurg/distillery
brew install alaala
```

## Testing

After completing all steps:

```bash
# Install
brew tap 0xGurg/distillery
brew install alaala

# Verify
alaala version

# Test
alaala init
```

## Auto-Updates for Future Releases

From now on, when you create a new release:

```bash
git tag -a v0.1.1 -m "Release v0.1.1"
git push origin v0.1.1
```

GoReleaser will automatically:
1. Build binaries
2. Calculate SHA256 checksums
3. Update `Formula/alaala.rb` in the distillery repo
4. Commit with message: "brew: update formula for v0.1.1"
5. Users get updates with `brew upgrade alaala`

No manual SHA256 updates needed after the initial setup!

## Files Structure

**In distillery repo:**
```
distillery/
├── Formula/
│   └── alaala.rb        # Auto-updated by GoReleaser
└── README.md            # Generic tap description
```

**In alaala repo:**
```
alaala/
├── homebrew-tap/        # Template/reference only
│   ├── Formula/alaala.rb
│   └── README.md
└── .goreleaser.yml      # Points to distillery repo
```

## Quick Reference

**Installation command for users:**
```bash
brew tap 0xGurg/distillery && brew install alaala
```

**Uninstall:**
```bash
brew uninstall alaala
```

**Upgrade:**
```bash
brew upgrade alaala
```

## Troubleshooting

### "Formula not found"
- Ensure distillery repo is public
- Verify Formula/alaala.rb exists in distillery repo
- Try `brew update` and retry

### "Checksum mismatch"
- Update SHA256 values in Formula/alaala.rb
- Get correct values from checksums.txt in the release

### GoReleaser not updating formula
- Verify HOMEBREW_TAP_GITHUB_TOKEN is set in alaala repo secrets
- Check token has `repo` scope
- Check GitHub Actions logs for errors
