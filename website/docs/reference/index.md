---
sidebar_position: 1
---

# Reference

Technical reference documentation for skillshare.

## What are you looking for?

| Topic | Read |
|-------|------|
| Environment variables that affect skillshare | [Environment Variables](./environment-variables.md) |
| Where skillshare stores config, skills, logs, cache | [File Structure](./file-structure.md) |
| Config file format and options | [Configuration](/docs/targets/configuration) |
| All CLI commands | [Commands](/docs/commands) |

## Quick Reference

### Key Paths (Unix)

| Path | Purpose |
|------|---------|
| `~/.config/skillshare/config.yaml` | Configuration file |
| `~/.config/skillshare/skills/` | Source directory (your skills) |
| `~/.local/share/skillshare/backups/` | Backup directory |
| `~/.local/share/skillshare/trash/` | Soft-deleted skills |
| `~/.local/state/skillshare/logs/` | Operation and audit logs |
| `~/.cache/skillshare/ui/` | Downloaded web dashboard |

### Environment Variables

| Variable | Purpose |
|----------|---------|
| `SKILLSHARE_CONFIG` | Override config path |
| `GITHUB_TOKEN` | GitHub API authentication |

## See Also

- [Configuration](/docs/targets/configuration) — Config file details
- [Commands](/docs/commands) — All commands
