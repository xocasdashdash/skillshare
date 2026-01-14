<p align="center">
  <img src=".github/assets/logo.png" alt="skillshare" width="200">
</p>

<h1 align="center">skillshare</h1>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="go.mod"><img src="https://img.shields.io/github/go-mod/go-version/runkids/skillshare" alt="Go Version"></a>
  <a href="https://github.com/runkids/skillshare/releases"><img src="https://img.shields.io/github/v/release/runkids/skillshare" alt="Release"></a>
</p>

<p align="center">
  Share skills across AI CLI tools (Claude Code, Codex CLI, Cursor, Gemini CLI, OpenCode).
</p>

```
┌─────────────────────────────────────────────────────────────┐
│                  ~/.config/skillshare/skills/               │
│         my-skill/   another-skill/   shared-util/           │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
              ▼               ▼               ▼
       ┌───────────┐   ┌───────────┐   ┌───────────┐
       │  Claude   │   │   Codex   │   │  Gemini   │
       │  skills/  │   │  skills/  │   │  skills/  │
       └───────────┘   └───────────┘   └───────────┘
```

## Table of Contents

- [Why skillshare?](#why-skillshare)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Sync Modes](#sync-modes)
- [Backup & Restore](#backup--restore)
- [Configuration](#configuration)
- [Use Cases](#use-cases)
- [Sync Across Machines](#sync-across-machines)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

## Why skillshare?

Each AI CLI tool has its own skills directory, and keeping them in sync manually is tedious.

| Manual Management | skillshare |
|-------------------|------------|
| Copy to each CLI directory | Set up once, auto-sync |
| Update requires multiple copies | Change once, update everywhere |
| Easy to miss or have inconsistent versions | Always in sync |
| No change tracking | Git version control support |
| Skills lost when CLI is removed | Safe backup mechanism |

## Installation

### Homebrew (Recommended)

```bash
brew install runkids/tap/skillshare
```

### macOS

```bash
# Apple Silicon (M1/M2/M3/M4)
curl -sL https://github.com/runkids/skillshare/releases/latest/download/skillshare_darwin_arm64.tar.gz | tar xz
sudo mv skillshare /usr/local/bin/

# Intel
curl -sL https://github.com/runkids/skillshare/releases/latest/download/skillshare_darwin_amd64.tar.gz | tar xz
sudo mv skillshare /usr/local/bin/
```

### Linux

```bash
# x86_64
curl -sL https://github.com/runkids/skillshare/releases/latest/download/skillshare_linux_amd64.tar.gz | tar xz
sudo mv skillshare /usr/local/bin/

# ARM64
curl -sL https://github.com/runkids/skillshare/releases/latest/download/skillshare_linux_arm64.tar.gz | tar xz
sudo mv skillshare /usr/local/bin/
```

### Windows

Download from [Releases](https://github.com/runkids/skillshare/releases) and add to PATH.

### Verify Installation

```bash
skillshare version
```

### Uninstall

```bash
# Homebrew
brew uninstall skillshare

# Manual (curl install)
sudo rm /usr/local/bin/skillshare

# Config and data (optional)
rm -rf ~/.config/skillshare
```

## Quick Start

```bash
# 1. Initialize (auto-detects CLIs, prompts for git setup)
skillshare init

# 2. Check detected targets
skillshare status

# 3. Sync skills to all targets
skillshare sync

# 4. (Optional) Push to remote for backup
cd ~/.config/skillshare/skills
git remote add origin git@github.com:you/my-skills.git
git push -u origin main
```

## Usage

### Initialize

```bash
# Interactive mode - choose source from existing skills directories
skillshare init

# Or specify custom source
skillshare init --source ~/my-skills
```

This will:
- Detect installed CLI tools
- Optionally copy skills from existing directories
- Initialize git for version control (recommended for recovery)
- Create config at `~/.config/skillshare/config.yaml`

### Sync

```bash
# Preview what will happen
skillshare sync --dry-run

# Actually sync
skillshare sync

# Force sync (override conflicts)
skillshare sync --force
```

### Status

```bash
skillshare status
```

Shows source directory, skill count, and each target's sync state.

### Diff

```bash
# Show differences for all targets
skillshare diff

# Show differences for specific target
skillshare diff claude
```

### Pull

```bash
# Pull local skills from specific target to source
skillshare pull claude

# Pull from all targets
skillshare pull --all

# Preview what will be pulled
skillshare pull --dry-run

# Force overwrite existing skills in source
skillshare pull --force
```

Copies skills created in target directories back to source. Useful when you create skills directly in a CLI's skills directory (e.g., `~/.claude/skills/`) and want to share them with other CLIs.

**Typical workflow:**
1. Create a skill in `~/.claude/skills/my-new-skill/`
2. Run `skillshare pull claude` to copy it to source
3. Run `skillshare sync` to distribute to all targets

### Doctor

```bash
skillshare doctor
```

Diagnoses config, source directory, symlink support, and target health.

### Manage Targets

```bash
# List all targets
skillshare target list

# Show target info
skillshare target claude

# Change target sync mode
skillshare target claude --mode merge
skillshare target claude --mode symlink

# Add custom target
skillshare target add myapp ~/.myapp/skills

# Unlink target (restore skills and remove from config)
skillshare target remove myapp

# Unlink all targets
skillshare target remove --all
```

#### Target Validation

When adding a target, skillshare validates:

| Check | Behavior |
|-------|----------|
| Name format | Must start with a letter, only letters/numbers/underscores/hyphens |
| Reserved names | Cannot use `add`, `remove`, `list`, `help`, `all` |
| Path ending | Warns if path doesn't end with `skills` |
| Directory exists | Warns if target or parent directory doesn't exist |

```bash
# These will be rejected:
skillshare target add 123test ~/path    # Name starts with number
skillshare target add add ~/path        # Reserved name

# This will show a warning and ask for confirmation:
skillshare target add myapp ~/Documents  # Path doesn't look like skills dir
```

## Sync Modes

| Mode | Behavior | Backup Value |
|------|----------|--------------|
| `symlink` | Entire directory becomes a symlink to source. All targets share exactly the same skills. | Low (all synced) |
| `merge` (recommended) | Each skill is symlinked individually. Local skills in target are preserved. | High (local skills) |

**Recommendation**: Use `merge` mode for safety. It preserves local skills and is safer to manage.

Change mode per target:

```bash
skillshare target claude --mode merge
skillshare sync
```

### ⚠️ Important: Symlink Safety

When using symlinks, **deleting files through a target directory deletes the source**:

```bash
# DANGER: This deletes your SOURCE skills!
rm -rf ~/.codex/skills/*        # ❌ Deletes source!
rm -rf ~/.codex/skills/my-skill # ❌ Deletes from source!

# SAFE: Use skillshare commands instead
skillshare target remove codex  # ✅ Safely unlinks
```

**Safe practices:**
- Always use `skillshare target remove` to unlink targets
- Never manually delete files inside symlinked directories
- Keep your source in a git repo for recovery
- Use `merge` mode if you're unsure

## Backup & Restore

### Automatic Backups

Backups are created automatically before:
- `skillshare sync` - backs up targets before syncing
- `skillshare target remove` - backs up before unlinking

Backup location: `~/.config/skillshare/backups/<timestamp>/<target>/`

### Manual Backup

```bash
# Backup all targets
skillshare backup

# Backup specific target
skillshare backup claude
```

### List Backups

```bash
skillshare backup --list
```

Output:
```
All backups (1.2 MB total)
─────────────────────────────────────────
  2026-01-14_21-22-18  claude, codex   0.6 MB  ~/.config/skillshare/backups/...
  2026-01-14_21-21-55  claude          0.6 MB  ~/.config/skillshare/backups/...
```

### Cleanup Old Backups

```bash
skillshare backup --cleanup
```

Default cleanup policy:
- Keep backups for 30 days
- Keep maximum 10 backups
- Keep maximum 500 MB total

### Restore from Backup

```bash
# Restore from latest backup
skillshare restore claude

# Restore from specific backup
skillshare restore claude --from 2026-01-14_21-22-18

# Force overwrite existing files
skillshare restore claude --force
```

### When is Backup Useful?

| Mode | Backup Value | Reason |
|------|--------------|--------|
| `symlink` | Low | All targets point to source, no local data |
| `merge` | High | Local skills may exist in target |
| First sync | High | Target has original skills before linking |

In `symlink` mode, backups are mostly skipped because targets are just symlinks with no local data.

## Configuration

Config file: `~/.config/skillshare/config.yaml`

```yaml
source: ~/.config/skillshare/skills
mode: merge   # default mode for all targets
targets:
  claude:
    path: ~/.claude/skills
  codex:
    path: ~/.codex/skills
    mode: symlink   # override: use full directory symlink
  cursor:
    path: ~/.cursor/skills
  gemini:
    path: ~/.gemini/antigravity/skills
  opencode:
    path: ~/.config/opencode/skills
ignore:
  - "**/.DS_Store"
  - "**/.git/**"
```

## Use Cases

### Individual Developer

- Sync skills across multiple machines
- Unified management of all AI CLI tools

### Team Collaboration

- Share team standard skills (coding standards, review guidelines)
- Quick onboarding for new members (clone + sync)
- Review skill changes via Git PR

### Open Source Community

- Share quality skills with the community
- Fork and customize others' skill libraries

## Sync Across Machines

Use git to sync your skills across multiple machines:

### Initial Setup (Machine A)

```bash
# Initialize skillshare (git is auto-initialized)
skillshare init

# Push to remote
cd ~/.config/skillshare/skills
git remote add origin git@github.com:you/my-skills.git
git push -u origin main
```

### Clone to Another Machine (Machine B)

```bash
# Clone your skills repo
git clone git@github.com:you/my-skills.git ~/.config/skillshare/skills

# Initialize skillshare with the cloned source
skillshare init --source ~/.config/skillshare/skills

# Sync to all targets
skillshare sync
```

### Daily Workflow

```bash
# Machine A - add/update skills
cd ~/.config/skillshare/skills
# ... edit skills ...
git add . && git commit -m "Update skills" && git push

# Machine B - pull and sync
cd ~/.config/skillshare/skills
git pull
skillshare sync  # New skills are automatically symlinked
```

## FAQ

### What happens if I modify a skill in the target directory?

Since targets use symlinks, modifying a skill in any target directory actually modifies the source. The change is immediately visible in all other targets.

### How can I keep a CLI-specific skill?

Use **merge mode**. In merge mode, you can create skills directly in the target directory - they won't be overwritten by sync because skillshare only creates symlinks for skills that exist in source.

### What's the difference between merge and symlink mode?

- **Merge mode** (recommended): Individual symlinks per skill. Local skills are preserved. Safer to manage.
- **Symlink mode**: Entire directory is a symlink. All targets are identical. ⚠️ Deleting through target deletes source!

### What happens if I delete a skill in the target directory?

**⚠️ DANGER**: Since targets are symlinks, deleting files in a target directory **deletes the source files**!

```bash
rm -rf ~/.codex/skills/my-skill  # This deletes from SOURCE!
```

To safely remove targets, always use:
```bash
skillshare target remove codex
```

### Are backups created automatically?

Yes. Backups are created automatically before `sync` and `target remove` operations. However, in `symlink` mode, backups are skipped because there's no local data to backup (everything is synced).

### How do I restore from a backup?

```bash
# See available backups
skillshare backup --list

# Restore latest backup
skillshare restore claude

# Restore specific backup
skillshare restore claude --from 2026-01-14_21-22-18
```

### Does skillshare automatically initialize git?

Yes! During `skillshare init`, you'll be prompted to initialize git in the source directory. This is highly recommended because:
- Git protects against accidental deletion (`git reflog` can recover almost anything)
- Enables syncing across multiple machines
- Tracks all changes to your skills

If you skip git initialization, you'll see a warning that deleted skills cannot be recovered.

### Can I use skillshare with a private git repo?

Yes! The source directory (`~/.config/skillshare/skills`) is just a regular directory. You can use any git hosting service (GitHub, GitLab, Bitbucket, self-hosted) with any visibility setting.

## Contributing

Contributions are welcome! Here's how you can help:

### Report Issues

Found a bug or have a feature request? [Open an issue](https://github.com/runkids/skillshare/issues/new).

### Submit Pull Requests

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Clone the repo
git clone https://github.com/runkids/skillshare.git
cd skillshare

# Build and test
./build.sh

# Or manually
go build -o bin/skillshare ./cmd/skillshare
go test ./...
```

## License

MIT
