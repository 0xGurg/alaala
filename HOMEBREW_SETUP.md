# Homebrew Tap Setup Instructions

## Using Existing Distillery Repository

Using your existing `0xGurg/distillery` repository as the Homebrew tap.

### Step 1: Push Formula to Distillery Repository

```bash
# Clone your distillery repository
cd ~/projects
git clone https://github.com/0xGurg/distillery.git
cd distillery

# Copy the prepared formula
cp -r ~/projects/alaala/homebrew-tap/* .

# Add and push
git add .
git commit -m "Add alaala Homebrew formula"
git push origin main
```

### Step 3: Create GitHub Personal Access Token

1. Go to https://github.com/settings/tokens/new
2. Token name: `alaala-homebrew-tap`
3. Select scopes:
   - ✅ `repo` (Full control of private repositories)
4. Click "Generate token"
5. **Copy the token** (starts with `ghp_`)

### Step 4: Add Secret to alaala Repository

1. Go to https://github.com/0xGurg/alaala/settings/secrets/actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Value: Paste the token from Step 3
5. Click "Add secret"

### Step 5: Update Formula with SHA256 Checksums

After the v0.1.0 release completes on GitHub Actions:

1. Download the checksums:
```bash
curl -L https://github.com/0xGurg/alaala/releases/download/v0.1.0/checksums.txt
```

2. Update the formula in distillery:
```bash
cd ~/projects/distillery

# Edit Formula/alaala.rb
# Replace PLACEHOLDER_SHA256_* with actual SHA256 values from checksums.txt

git add Formula/alaala.rb
git commit -m "Update SHA256 checksums for alaala v0.1.0"
git push origin main
```

### That's It!

Users can now install with:
```bash
brew tap 0xGurg/distillery
brew install alaala
```

### Auto-Updates on Future Releases

From now on, when you create a new release:
```bash
git tag -a v0.1.1 -m "Release v0.1.1"
git push origin v0.1.1
```

GoReleaser will automatically:
1. Build binaries
2. Update the formula in `homebrew-alaala` repo
3. Create GitHub release
4. Users get updates with `brew upgrade alaala`

## Testing Installation

After completing the steps above:

```bash
# Install
brew tap 0xGurg/distillery
brew install alaala

# Verify
alaala version

# Uninstall (for testing)
brew uninstall alaala
brew untap 0xGurg/distillery
```

## Troubleshooting

### "Formula not found"
- Make sure the tap repository is public
- Verify Formula/alaala.rb exists in the tap repo
- Try `brew update` and try again

### "Checksum mismatch"
- Update the SHA256 values in the formula
- Get correct checksums from the release

### "GoReleaser not updating formula"
- Verify HOMEBREW_TAP_GITHUB_TOKEN is set correctly
- Check GitHub Actions logs for errors
- Ensure token has `repo` scope

## File Structure

After setup, you'll have:

**Main repo (0xGurg/alaala):**
- Source code
- .goreleaser.yml (with brews section)
- homebrew-tap/ (template for reference)

**Tap repo (0xGurg/distillery):**
- Formula/alaala.rb (auto-updated by GoReleaser)
- README.md (optional)

## Benefits

- ✅ Clean command: `brew tap 0xGurg/distillery && brew install alaala`
- ✅ Professional appearance
- ✅ Standard Homebrew workflow
- ✅ Automatic formula updates
- ✅ Separate concerns (code vs distribution)
