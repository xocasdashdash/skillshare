---
sidebar_position: 7
---

# Project Setup

Set up project-level skills from scratch — skills scoped to a single repository, shared with your team via git.

## When to Use Project Mode

| Scenario | Use |
|----------|-----|
| Skills for a specific project | **Project mode** |
| Skills shared across all projects | Global mode |
| Team-specific skills committed to repo | **Project mode** |
| Personal skills on multiple machines | Global mode |
| Monorepo with project-specific AI context | **Project mode** |

---

## Step-by-Step Setup

### Step 1: Initialize

Run `skillshare init -p` in your project root:

```bash
cd my-project
skillshare init -p
```

```
┌─────────────────────────────────────────────────────────────────┐
│ skillshare init -p                                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. Create .skillshare/ directory                                 │
│    → .skillshare/skills/                                         │
│    → .skillshare/config.yaml                                     │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Detect AI CLI directories in project                          │
│    → Found: .claude/, .cursor/                                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Create target skill directories                               │
│    → .claude/skills/                                             │
│    → .cursor/skills/                                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. Write config.yaml                                             │
│    → targets: claude-code, cursor                                │
└─────────────────────────────────────────────────────────────────┘
```

You can also specify targets directly:

```bash
skillshare init -p --targets claude-code,cursor
```

### Step 2: Create Local Skills

Create skills manually or with `skillshare new`:

```bash
# Using skillshare new
skillshare new my-skill -p

# Or manually
mkdir -p .skillshare/skills/my-skill
cat > .skillshare/skills/my-skill/SKILL.md << 'EOF'
---
name: my-skill
description: Project-specific coding guidelines
---
# My Skill

Your skill content here...
EOF
```

### Step 3: Install Remote Skills

Install skills from GitHub into the project:

```bash
skillshare install anthropics/skills/skills/pdf -p
skillshare install github.com/team/shared-skills/review -p
```

Remote skills are:
- Installed to `.skillshare/skills/<name>/`
- Recorded in `.skillshare/config.yaml` under `skills:`
- Added to `.skillshare/.gitignore` (cloned content not committed)

### Step 4: Sync to Targets

```bash
skillshare sync
```

Creates symlinks from `.skillshare/skills/` to each target directory. Auto-detects project mode.

### Step 5: Commit to Version Control

```bash
git add .skillshare/
git commit -m "Add project-level skills"
```

**What gets committed:**
- `.skillshare/config.yaml` — targets and remote skill list
- `.skillshare/.gitignore` — ignore patterns for cloned skills
- `.skillshare/skills/<local-skills>/` — local skill content

**What's ignored:**
- Remote skill directories (re-installed from config)

---

## New Team Member Onboarding

When a new member clones the repo:

```bash
# 1. Clone the project
git clone github.com/team/my-project
cd my-project

# 2. Install remote skills from config
skillshare install -p

# 3. Sync to targets
skillshare sync
```

`skillshare install -p` (no URL) reads `.skillshare/config.yaml` and installs all listed remote skills.

---

## Custom Target Paths

Targets support both known names and custom paths:

```yaml
# .skillshare/config.yaml
targets:
  - claude-code                    # Known name → .claude/skills/
  - cursor                         # Known name → .cursor/skills/
  - name: custom-tool              # Custom path
    path: ./tools/ai/skills        # Relative to project root
  - name: another-tool
    path: ~/global/path/skills     # Absolute path with ~ expansion
```

---

## Full Config Example

```yaml
targets:
  - claude-code
  - cursor
  - name: windsurf
    path: .windsurf/skills

skills:
  - name: pdf
    source: anthropic/skills/pdf
  - name: code-review
    source: github.com/team/skills/code-review
```

---

## Coexistence with Global Mode

Project and global skills work independently:

```
Global source                       Project source
~/.config/skillshare/skills/        .skillshare/skills/
├── personal-skill/                 ├── project-skill/
└── _team-repo/                     └── remote-skill/
         │                                   │
         ▼                                   ▼
   ~/.claude/skills/                .claude/skills/
   (global targets)                 (project targets — separate)
```

- Project mode targets are **project-local** (e.g., `.claude/skills/` inside the project)
- Global mode targets are **system-wide** (e.g., `~/.claude/skills/`)
- They don't conflict — different directories, different scope

---

## Related

- [Project Skills](/docs/concepts/project-skills) — Core concepts
- [Project Workflow](/docs/workflows/project-workflow) — Daily operations
- [Team Sharing](/docs/guides/team-sharing) — Global-mode team sharing with tracked repos
- [Commands: init](/docs/commands/init) — Init command details
- [Commands: install](/docs/commands/install) — Install command details
