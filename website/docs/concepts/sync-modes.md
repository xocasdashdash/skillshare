---
sidebar_position: 3
---

# Sync Modes

How skillshare links source to targets.

:::tip When does this matter?
Choose merge mode (default) when you want per-skill symlinks and to preserve local skills in targets. Choose symlink mode when you want the entire directory linked and don't need local target skills.
:::

## Overview

| Mode | Behavior | Use Case |
|------|----------|----------|
| `merge` | Each skill symlinked individually | **Default.** Preserves local skills. |
| `symlink` | Entire directory is one symlink | Exact copies everywhere. |
| `copy` | Each skill copied (no symlinks) | When symlinks are not desired or supported. |

---

## Merge Mode (Default)

Each skill is symlinked individually. Local skills in the target are preserved.

```
Source                          Target (claude)
─────────────────────────────────────────────────────────────
skills/                         ~/.claude/skills/
├── my-skill/        ────────►  ├── my-skill/ → (symlink)
├── another/         ────────►  ├── another/  → (symlink)
└── ...                         ├── local-only/  (preserved)
                                └── ...
```

**Advantages:**
- Keep target-specific skills (not synced)
- Mix installed and local skills
- Granular control
- Per-target include/exclude filtering

**When to use:**
- You want some skills only in specific AI CLIs
- You want to try local skills before syncing
- You want one source but different skill subsets per target

### Filter strategy in merge mode

`include` and `exclude` are evaluated per target in this order:
1. `include` keeps matching names
2. `exclude` removes from that kept set

Quick choices:
- Use `include` when the target should get only a small subset
- Use `exclude` when the target should get almost everything
- Use `include + exclude` when you need a broad subset with explicit carve-outs

Behavior when rules change:
- Previously synced source-linked entries that become filtered out are removed on next `sync`
- Existing local non-symlink folders in target are preserved

See [Target Configuration](/docs/targets/configuration#include--exclude-target-filters) for full examples.

---

## Symlink Mode

The entire target directory is a single symlink to source.

```
Source                          Target (claude)
─────────────────────────────────────────────────────────────
skills/              ────────►  ~/.claude/skills → (symlink to source)
├── my-skill/
├── another/
└── ...
```

**Advantages:**
- All targets are identical
- Simpler to manage
- No orphaned symlinks

**When to use:**
- You want all AI CLIs to have exactly the same skills
- You don't need target-specific skills

**Warning:** In symlink mode, deleting through target deletes source!
```bash
rm -rf ~/.claude/skills/my-skill  # ❌ Deletes from SOURCE
skillshare target remove claude   # ✅ Safe way to unlink
```

---

## Copy Mode

Each skill is **copied** into the target directory (no symlinks). Same filtering as merge mode (include/exclude, per-skill targets).

**When to use:**
- Symlinks are not supported or not desired (e.g. some network or shared filesystems)
- You want a true snapshot of skills in the target; changes in source require `skillshare sync` to update the target

**Behavior:**
- `sync` copies each managed skill directory into the target (flat names, e.g. `my__nested__skill`)
- Existing copies are left unchanged unless you run `sync --force` (then they are overwritten)
- Orphan copies (skills removed from source) are pruned on sync

```bash
skillshare target claude --mode copy
skillshare sync
```

---

## Changing Mode

### Per-target

```bash
# Switch to symlink mode
skillshare target claude --mode symlink
skillshare sync

# Switch back to merge mode
skillshare target claude --mode merge
skillshare sync

# Use copy mode (no symlinks)
skillshare target claude --mode copy
skillshare sync
```

### Default mode

Set in config for new targets:

```yaml
# ~/.config/skillshare/config.yaml
mode: merge  # or symlink, or copy

targets:
  claude:
    path: ~/.claude/skills
    # inherits default mode

  codex:
    path: ~/.codex/skills
    mode: symlink  # override default (merge, symlink, or copy)
```

---

## Mode Comparison

| Aspect | Merge | Symlink |
|--------|-------|---------|
| Local skills preserved | ✅ Yes | ❌ No |
| All targets identical | ❌ Can differ | ✅ Yes |
| Per-target include/exclude | ✅ Yes | ❌ Ignored |
| Orphan cleanup needed | ✅ Yes | ❌ No |
| Delete safety | ✅ Safe | ⚠️ Caution |
| Complexity | Higher | Lower |

---

## Orphaned Symlinks

In merge mode, when you remove a skill from source, the target symlinks become "orphaned" (pointing to nothing).

**Handling:** `skillshare sync` automatically prunes orphaned symlinks:

```
$ skillshare sync
✓ claude: merged (5 linked, 2 local, 0 updated, 1 pruned)
                                                  ^^^^^^^^
```

---

## See Also

- [sync](/docs/commands/sync) — Run sync to apply mode changes
- [target](/docs/commands/target) — Change a target's sync mode
- [Source & Targets](./source-and-targets.md) — The core architecture
- [Configuration](/docs/targets/configuration) — Per-target settings
