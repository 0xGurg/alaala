# Homebrew Tap Setup Instructions

## Manual Steps Required

Since I cannot create GitHub repositories or add secrets through the API, you'll need to complete these steps manually:

### Step 1: Create Homebrew Tap Repository

1. Go to https://github.com/new
2. Create a new repository:
   - **Name:** `homebrew-alaala` (must start with `homebrew-`)
   - **Description:** Homebrew formulae for alaala
   - **Visibility:** Public
   - **Initialize:** No README, no .gitignore, no license

3. Clone and initialize the tap:
```bash
cd ~/projects
git clone https://github.com/0xGurg/homebrew-alaala.git
cd homebrew-alaala

# Copy the prepared files
cp -r ~/projects/alaala/homebrew-tap/* .

# Commit and push
git add .
git commit -m "Initial Homebrew formula for alaala"
git push origin main
```

### Step 2: Create GitHub Personal Access Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token" → "Generate new token (classic)"
3. Give it a name: `alaala-homebrew-tap`
4. Select scopes:
   - ✅ `repo` (Full control of private repositories)
5. Click "Generate token"
6. **Copy the token** (starts with `ghp_`)

### Step 3: Add Secret to alaala Repository

1. Go to https://github.com/0xGurg/alaala/settings/secrets/actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Value: Paste the token from Step 2
5. Click "Add secret"

### Step 4: Update SHA256 Checksums

After v0.1.0 release completes:

1. Download the checksums file:
```bash
curl -L https://github.com/0xGurg/alaala/releases/download/v0.1.0/checksums.txt
```

2. Update `Formula/alaala.rb` with actual SHA256 values for each platform
3. Commit and push to homebrew-alaala repo

### Step 5: Test Installation

```bash
# Remove existing installation if any
which alaala && sudo rm $(which alaala)

# Install via Homebrew
brew tap 0xGurg/alaala
brew install alaala

# Verify
alaala version
```

## What's Already Prepared

I've created these files in `homebrew-tap/` directory:
- ✅ `Formula/alaala.rb` - Complete Homebrew formula
- ✅ `README.md` - Tap repository README
- ✅ `.goreleaser.yml` updated with brews section (in main repo)
- ✅ Documentation updates prepared

## After Manual Steps Complete

Once you've completed steps 1-3 above, I can continue with:
- Updating documentation to use brew
- Removing install.sh and uninstall.sh
- Creating a new release to test the auto-update

## Verification

After setup, the workflow should:
1. You push a new tag (e.g., v0.1.1)
2. GitHub Actions runs GoReleaser
3. GoReleaser builds binaries
4. GoReleaser automatically updates homebrew-alaala repo
5. Users get the update with `brew upgrade alaala`

Let me know when you've completed the manual steps!

