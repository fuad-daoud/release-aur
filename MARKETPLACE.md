# Publishing to GitHub Marketplace

## Prerequisites Checklist

- [x] Public repository
- [x] Single `action.yml` at root
- [x] No workflow files in root (all in `.github/workflows/`)
- [x] Unique action name: "Release to AUR"
- [x] Description present
- [x] Author specified
- [x] Branding configured (icon + color)

## Publishing Steps

### 1. Accept GitHub Marketplace Developer Agreement

First time publishing? You'll need to accept the agreement:
- Go to: https://github.com/marketplace/actions
- Click "Publish an action"
- Accept the Developer Agreement

### 2. Create a Release

1. Navigate to: https://github.com/fuad-daoud/release-aur/releases/new

2. Fill in release details:
   ```
   Tag: v1.0.0
   Title: Release v1.0.0
   Description: 
   Initial release of Release to AUR action
   
   Features:
   - Generate PKGBUILD files for AUR
   - Support for multiple architectures (x86_64, aarch64)
   - Automatic version management
   - AUR package validation
   ```

3. **Important**: Check the box:
   ```
   ☑ Publish this Action to the GitHub Marketplace
   ```

4. Select categories:
   - **Primary Category**: Continuous integration
   - **Secondary Category**: Publishing

5. Click **"Publish release"**

### 3. Verify Publication

After publishing, your action will be available at:
- Marketplace: https://github.com/marketplace/actions/release-to-aur
- Direct usage: `fuad-daoud/release-aur@v1`

## Version Management

### Major Version Tags

GitHub automatically creates/updates major version tags (v1, v2, etc.) so users can reference:
```yaml
uses: fuad-daoud/release-aur@v1  # Always latest v1.x.x
uses: fuad-daoud/release-aur@v1.0.0  # Specific version
```

Our release workflow handles this automatically.

### Publishing Updates

For subsequent releases:

1. Create new tag: `v1.1.0`, `v1.2.0`, etc.
2. The release workflow will:
   - Build and test
   - Create GitHub release
   - Update major version tag (v1)
3. Users on `@v1` automatically get updates

## Unpublishing

To remove from Marketplace:

1. Go to each release
2. Click "Edit"
3. Uncheck "Publish this Action to the GitHub Marketplace"
4. Click "Update release"

## Marketplace Badge

Add to README.md:
```markdown
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Release%20to%20AUR-blue.svg?colorA=24292e&colorB=0366d6&style=flat&longCache=true&logo=github)](https://github.com/marketplace/actions/release-to-aur)
```

## Categories

Recommended categories for this action:
- **Primary**: Continuous integration
- **Secondary**: Publishing

Alternative categories:
- Deployment
- Utilities

## Branding

Current branding in `action.yml`:
```yaml
branding:
  icon: 'package'
  color: 'blue'
```

Available icons: https://feathericons.com/
Available colors: white, yellow, blue, green, orange, red, purple, gray-dark

## Verification Badge

To get the "Verified Creator" badge:
1. Become a GitHub Partner
2. Email: partnerships@github.com
3. Request verified creator badge

## Monitoring

After publication, monitor:
- Usage statistics (GitHub Insights)
- Issues from users
- Feature requests
- Security advisories

## Best Practices

1. **Semantic Versioning**: Use semver (v1.0.0, v1.1.0, v2.0.0)
2. **Changelog**: Maintain CHANGELOG.md
3. **Breaking Changes**: Bump major version
4. **Security**: Keep dependencies updated
5. **Documentation**: Keep README current
6. **Examples**: Provide working examples

## Troubleshooting

### "Name already taken"
- Action name must be unique across all of GitHub Marketplace
- Change `name:` in action.yml

### "Repository contains workflow files"
- Workflow files must be in `.github/workflows/`
- No `.yml` files in root except `action.yml`

### "Missing required fields"
- Ensure action.yml has: name, description, author, branding

### "Not seeing publish checkbox"
- Accept Marketplace Developer Agreement first
- Must be repository owner or org admin
