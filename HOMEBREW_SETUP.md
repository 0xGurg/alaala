# Homebrew Setup - Simplified Approach

## Using Main Repository

The Homebrew formula is stored directly in this repository under `Formula/alaala.rb`. No separate tap repository needed!

## How It Works

1. **Formula Location:** `Formula/alaala.rb` in the main repo
2. **Installation:** Users install directly from the formula URL
3. **Auto-Updates:** GoReleaser automatically updates the formula in this repo on each release

## User Installation

Users install with:
```bash
brew install https://raw.githubusercontent.com/0xGurg/alaala/main/Formula/alaala.rb
```

No need to `brew tap` first!

## Setup Required

### Add GitHub Token to Secrets (Optional)

If you want GoReleaser to auto-update the formula:

1. The workflow already has access to `GITHUB_TOKEN`
2. No additional setup needed since we're using the same repo!

GoReleaser will automatically commit formula updates to the `Formula/` directory on each release.

## Testing

After the v0.1.0 release completes (GitHub Actions finishes):

```bash
# Install
brew install https://raw.githubusercontent.com/0xGurg/alaala/main/Formula/alaala.rb

# Verify
alaala version

# Uninstall
brew uninstall alaala
```

## Updating the Formula

### Automatic (via GoReleaser)

When you create a new release:
```bash
git tag -a v0.1.1 -m "Release v0.1.1"
git push origin v0.1.1
```

GoReleaser will:
1. Build binaries
2. Update `Formula/alaala.rb` with new version and SHA256s
3. Commit the changes to main branch
4. Create GitHub release

### Manual (if needed)

If you need to manually update the formula:

1. Get the SHA256 checksums from the release
2. Update `Formula/alaala.rb`:
   - Change `version "x.x.x"`
   - Update URLs to new version
   - Update SHA256 values
3. Commit and push

## Benefits of This Approach

- ✅ No separate repository to maintain
- ✅ Formula lives with the code
- ✅ Simpler workflow
- ✅ Still get all Homebrew benefits
- ✅ Users can install with one command
- ✅ Auto-updates via GoReleaser

## Alternative: Official Tap (Future)

Later, you can create a separate `homebrew-alaala` tap for:
- Shorter install command: `brew tap 0xGurg/alaala && brew install alaala`
- Multiple formulas in one tap
- Professional appearance

But for now, the direct URL approach works perfectly!
