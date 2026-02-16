---
sidebar_position: 1
---

# Commands

Complete reference for all skillshare commands.

## Overview

| Category | Commands |
|----------|----------|
| **Core** | `init`, `install`, `uninstall`, `list`, `search`, `sync`, `status` |
| **Skill Management** | `new`, `update`, `upgrade` |
| **Target Management** | `target`, `diff` |
| **Sync Operations** | `collect`, `backup`, `restore`, `push`, `pull` |
| **Utilities** | `audit`, `log`, `doctor`, `ui` |

---

## Core Commands

| Command | Description |
|---------|-------------|
| [init](./init.md) | First-time setup |
| [install](./install.md) | Add a skill from a repo or path |
| [uninstall](./uninstall.md) | Remove a skill |
| [list](./list.md) | List all skills |
| [search](./search.md) | Search for skills |
| [sync](./sync.md) | Push skills to all targets |
| [status](./status.md) | Show sync state |

## Skill Management

| Command | Description |
|---------|-------------|
| [new](./new.md) | Create a new skill |
| [update](./update.md) | Update a skill or tracked repo |
| [upgrade](./upgrade.md) | Upgrade CLI or built-in skill |

## Target Management

| Command | Description |
|---------|-------------|
| [target](./target.md) | Manage targets |
| [diff](./diff.md) | Show differences between source and targets |

## Sync Operations

| Command | Description |
|---------|-------------|
| [collect](./collect.md) | Collect skills from target to source |
| [backup](./backup.md) | Create backup of targets |
| [restore](./restore.md) | Restore targets from backup |
| [push](./push.md) | Push to git remote |
| [pull](./pull.md) | Pull from git remote and sync |

## Utilities

| Command | Description |
|---------|-------------|
| [audit](./audit.md) | Scan skills for security threats |
| [log](./log.md) | View operations and audit logs |
| [doctor](./doctor.md) | Diagnose issues |
| [ui](./ui.md) | Launch web dashboard |

---

## Common Flags

Most commands support:

| Flag | Description |
|------|-------------|
| `--dry-run`, `-n` | Preview without making changes |
| `--help`, `-h` | Show help |

---

## Quick Reference

```bash
# Setup
skillshare init
skillshare init --remote git@github.com:you/skills.git

# Install skills
skillshare install anthropics/skills/skills/pdf
skillshare install github.com/team/skills --track

# Create skill
skillshare new my-skill

# Sync
skillshare sync
skillshare sync --dry-run

# Cross-machine
skillshare push -m "Add skill"
skillshare pull

# Status
skillshare status
skillshare list
skillshare diff

# Maintenance
skillshare update --all
skillshare audit
skillshare log
skillshare doctor
skillshare backup

# Web UI
skillshare ui
skillshare ui -p          # Project mode
```

---

## Related

- [Quick Reference](/docs/getting-started/quick-reference) — Command cheat sheet
- [Workflows](/docs/workflows) — Common usage patterns
