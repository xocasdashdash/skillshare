#!/usr/bin/env bash
# Devcontainer post-create setup — mirrors sandbox_playground_up.sh
# for an out-of-the-box demo experience inside VS Code / Codespaces.
#
# Environment assumptions (from docker-compose.yml + devcontainer.json):
#   HOME=/tmp  |  binary at /workspace/bin/skillshare  |  PATH includes /workspace/bin
set -euo pipefail

if [ "${HOME:-}" != "/tmp" ] || [ ! -d /workspace ] || [ ! -f /workspace/go.mod ]; then
  echo "Refusing to run: expected devcontainer context (HOME=/tmp, /workspace mounted)." >&2
  exit 1
fi
cd /workspace

# ── 1. Build CLI ────────────────────────────────────────────────────
echo "▸ Building skillshare binary …"
make build

# ── 1b. Install frontend dependencies ─────────────────────────────
echo "▸ Installing UI dependencies …"
(cd /workspace/ui && pnpm install --frozen-lockfile)
echo "▸ Installing website dependencies …"
(cd /workspace/website && pnpm install --frozen-lockfile)

# ── 2. Install shortcut symlinks to PATH ─────────────────────────
# Devcontainer-specific scripts live in .devcontainer/bin/ (source-controlled,
# survives make clean). Only ephemeral symlinks are created here.
echo "▸ Installing shortcut commands …"
ln -sf /workspace/bin/skillshare /workspace/bin/ss
ln -sf /workspace/.devcontainer/dev-servers.sh /workspace/.devcontainer/bin/dev-servers

# ── 3. Global mode init ────────────────────────────────────────────
echo "▸ Initializing global mode …"
mkdir -p "$HOME/.claude/skills"
skillshare init -g --no-copy --all-targets --no-git --skill

# ── 4. Create demo content (shared with sandbox playground) ───────
echo "▸ Creating demo content …"
SKILLS="$HOME/.config/skillshare/skills"
CFG="$HOME/.config/skillshare"
DEMO="$HOME/demo-project"
/workspace/scripts/create-demo-content.sh "$SKILLS" "$CFG" "$DEMO"

# ── Done ────────────────────────────────────────────────────────────
echo ""
echo "══════════════════════════════════════════════════════════"
echo "  Devcontainer ready!"
echo "══════════════════════════════════════════════════════════"
echo ""
echo "Quick start:"
echo "  ss status               # check current state"
echo "  ss list                 # see all skills"
echo "  ui                      # global-mode dashboard (auto-started)"
echo "  ui-p                    # switch to project-mode dashboard"
echo "  docs                    # open documentation site"
echo ""
echo "Dev servers (auto-started on container open):"
echo "  dev-servers status      # check which servers are running"
echo "  dev-servers restart     # restart all dev servers"
echo "  dev-servers logs vite   # tail Vite log (api|vite|docusaurus)"
echo ""
echo "Audit playground:"
echo "  ss audit                # scan all skills, see findings"
echo "  cd ~/demo-project && ss audit  # project-level scan"
