# Release Process

## Automated Release

Releases are automated via GitHub Actions and GoReleaser.

### Creating a Release

1. Ensure all changes are merged to `main`
2. Update version in relevant files (if needed)
3. Create and push a tag:

```bash
# Create a new version tag
git tag -a v0.2.0 -m "Release v0.2.0"

# Push the tag
git push origin v0.2.0
```

4. GitHub Actions will automatically:
   - Build binaries for all platforms
   - Run tests
   - Generate checksums
   - Create GitHub Release
   - Upload assets
   - Generate changelog

### Version Scheme

We use Semantic Versioning (SemVer):
- **MAJOR.MINOR.PATCH** (e.g., v0.2.0)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Platform Support

Binaries are built for:
- **macOS**: darwin_amd64, darwin_arm64
- **Linux**: linux_amd64, linux_arm64
- **Windows**: windows_amd64, windows_arm64

### Release Artifacts

Each release includes:
- Cross-platform binaries (tar.gz/zip)
- SHA256 checksums
- LICENSE
- README.md
- QUICKSTART.md

### Manual Release (if needed)

If automated release fails:

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Test release locally
goreleaser release --snapshot --clean

# Create actual release (requires GITHUB_TOKEN)
export GITHUB_TOKEN="your-token"
goreleaser release --clean
```

## Pre-Release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated (if manual)
- [ ] Version bumped appropriately
- [ ] No uncommitted changes

## Post-Release

- [ ] Verify release on GitHub
- [ ] Test binary downloads
- [ ] Update Homebrew formula (future)
- [ ] Announce release (if applicable)

## Release Cadence

- **Patch releases**: Bug fixes as needed
- **Minor releases**: New features monthly
- **Major releases**: Breaking changes quarterly

## Troubleshooting

### GoReleaser Fails

Check the following:
1. All tags follow the `v*` pattern (e.g., v0.1.0, not 0.1.0)
2. Git is clean with no uncommitted changes
3. GITHUB_TOKEN has proper permissions
4. All build targets compile successfully

### Workflow Not Triggering

Ensure:
1. Tag was pushed to GitHub (`git push origin v0.1.0`)
2. Tag follows the `v*` pattern
3. Workflows are enabled in repository settings

### Binary Download Issues

If users report download issues:
1. Verify checksums match
2. Check file permissions on release assets
3. Test download links manually

## Version History

See [Releases](https://github.com/0xGurg/alaala/releases) for full version history.

