---
sidebar_position: 4
---

# list

List all installed skills in the source directory.

```bash
skillshare list              # Compact view
skillshare list --verbose    # Detailed view
```

![list demo](/img/list-demo.png)

## Example Output

### Compact View

```
Installed skills
─────────────────────────────────────────
  → my-skill              local
  → commit-commands       github.com/user/skills
  → _team-skills:review   tracked: _team-skills

Tracked repositories
─────────────────────────────────────────
  ✓ _team-skills          3 skills, up-to-date
```

### Verbose View

```bash
skillshare list --verbose
```

```
Installed skills
─────────────────────────────────────────
  my-skill
    Source:      (local - no metadata)

  commit-commands
    Source:      github.com/user/skills
    Type:        github
    Installed:   2026-01-15

  _team-skills:review
    Tracked repo: _team-skills
    Source:      github.com/team/skills
    Type:        github
    Installed:   2026-01-10

Tracked repositories
─────────────────────────────────────────
  ✓ _team-skills          3 skills, up-to-date
  ! _other-repo           5 skills, has changes
```

## Global vs Project

Skillshare operates at two levels. The `list` command shows skills from the active level:

```
┌─────────────────────────────────────────────────────────────┐
│  GLOBAL (machine-wide)          PROJECT (repository-scoped) │
│                                                             │
│  ~/.config/skillshare/          <project>/.skillshare/      │
│  └── skills/                    └── skills/                 │
│      ├── my-skill/                  ├── local-skill/        │
│      ├── commit-cmds/               └── remote-skill/       │
│      └── _team-repo/                                        │
│                                                             │
│  ┌───────────────┐              ┌───────────────┐           │
│  │ list / list -g│              │ list -p       │           │
│  └───────┬───────┘              └───────┬───────┘           │
│          ▼                              ▼                   │
│  Installed skills               Installed skills (project)  │
│  ─────────────────              ──────────────────────────  │
│  → my-skill  local              → local-skill   local       │
│  → commit..  github.com/...     → remote-skill  github.com/ │
└─────────────────────────────────────────────────────────────┘
```

| | Global | Project |
|---|---|---|
| **Source** | `~/.config/skillshare/skills/` | `.skillshare/skills/` |
| **Flag** | `-g` or default | `-p` or auto-detected |
| **Scope** | All projects on machine | Single repository |
| **Shared via** | `push` / `pull` | git commit |

### Auto-Detection

When you run `skillshare list` without flags, skillshare automatically detects the mode:

```
skillshare list
    │
    ├── .skillshare/config.yaml exists?
    │       ├── YES → Project mode
    │       └── NO  → Global mode
```

```bash
cd my-project/            # Has .skillshare/config.yaml
skillshare list           # → Installed skills (project)

cd ~
skillshare list           # → Installed skills (global)
```

Use `-p` or `-g` to override auto-detection:

```bash
skillshare list -g        # Force global, even inside a project
skillshare list -p        # Force project, even without auto-detection
```

## Project Mode

```bash
skillshare list          # Auto-detected if .skillshare/ exists
skillshare list -p       # Explicit project mode
```

### Example Output

```
Installed skills (project)
─────────────────────────────────────────
  → my-skill            local
  → pdf                 anthropic/skills/pdf
  → review              github.com/team/tools

→ 3 skill(s): 2 remote, 1 local
```

Project list uses the same visual format as global list, with `(project)` label in the header. Skills are categorized as `local` (no metadata) or by source URL (remote).

## Options

| Flag | Description |
|------|-------------|
| `--verbose, -v` | Show detailed information (source, type, install date) |
| `--project, -p` | List project skills |
| `--help, -h` | Show help |

## Understanding the Output

### Skill Sources

| Label | Meaning |
|-------|---------|
| `local` | Created locally, no metadata |
| `github.com/...` | Installed from GitHub |
| `tracked: <repo>` | Part of a tracked repository |

### Repository Status

| Icon | Meaning |
|------|---------|
| `✓` | Up-to-date, no local changes |
| `!` | Has uncommitted changes |

## Related

- [install](/docs/commands/install) — Install skills
- [uninstall](/docs/commands/uninstall) — Remove skills
- [status](/docs/commands/status) — Show sync status
