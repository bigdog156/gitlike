# Rename Summary: todo-cli → gitlike

## ✅ Completed Changes

### 1. Go Module and Imports
- **go.mod**: Changed module name from `todo-cli` to `gitlike`
- **Import paths**: Updated all Go files to use `gitlike/` imports:
  - `gitlike/commands`
  - `gitlike/models`
  - `gitlike/storage`
  - `gitlike/git`
  - `gitlike/remote`

### 2. Binary Name and CLI Interface
- **Command name**: Changed from `todo` to `gitlike`
- **Binary output**: Now builds as `gitlike` instead of `todo-cli`
- **Help text**: Updated all descriptions to "GitLike CLI"
- **Usage examples**: Updated all command examples to use `gitlike`

### 3. Documentation Updates
- **INSTALL.md**: Updated all references from `todo-cli` to `gitlike`
- **HOMEBREW_PUBLISHING.md**: Updated for GitLike CLI publishing
- **Release script**: Updated to build `gitlike-*` binaries

### 4. Homebrew Publishing
- **Tap name**: Changed from `homebrew-todocli` to `homebrew-gitlike`
- **Formula class**: Changed from `TodoCli` to `Gitlike`
- **Binary names**: All release binaries now use `gitlike-` prefix
- **Installation command**: `brew install gitlike`

### 5. Release Configuration
- **GoReleaser**: Updated project name to `gitlike`
- **Binary outputs**: 
  - `gitlike-darwin-amd64`
  - `gitlike-darwin-arm64`
  - `gitlike-linux-amd64`
  - `gitlike-windows-amd64.exe`

## 🎯 User Impact

### Before (todo-cli)
```bash
# Build
go build -o todo-cli

# Usage
todo-cli branch create feature-auth
todo-cli todo add "New task"
todo-cli --help

# Homebrew
brew tap bigdog156/todocli
brew install todo-cli
```

### After (gitlike)
```bash
# Build  
go build -o gitlike

# Usage
gitlike branch create feature-auth
gitlike todo add "New task"
gitlike --help

# Homebrew
brew tap bigdog156/gitlike
brew install gitlike
```

## 📦 Release Process

### Manual Release
1. Run: `./scripts/release.sh 1.0.0`
2. Create GitHub release with generated binaries
3. Create `homebrew-gitlike` repository
4. Copy `homebrew-formula-template.rb` to `gitlike.rb` in tap

### Automated Release (GoReleaser)
1. Install: `brew install goreleaser`
2. Run: `goreleaser release --rm-dist`
3. Automatically creates release and updates tap

## ✨ Benefits of the Rename

1. **Better Branding**: "GitLike" clearly describes the Git-like workflow
2. **Shorter Commands**: `gitlike` is more concise than `todo-cli`
3. **Professional Name**: Sounds more like a serious developer tool
4. **Clear Purpose**: Name immediately conveys the Git-like task management concept
5. **Homebrew Friendly**: Easier to remember `brew install gitlike`

## 🔧 Technical Notes

- All Go import paths successfully updated
- Module dependencies remain the same
- Backward compatibility: Old data files still work (stored in `~/.tododata/`)
- All functionality preserved - just renamed interface
- Git integration still works as before
- Remote repository features unchanged

The rename is complete and the application is ready for release as "GitLike CLI"!
