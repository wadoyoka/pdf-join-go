---
allowed-tools: Bash(git add:*), Bash(git status:*), Bash(git commit:*), Bash(git diff:*), Bash(git branch:*), Bash(git log:*), Task
description: Create a git commit
---

# Smart Commit Command

This command analyzes staged changes, performs code review, and generates an appropriate commit message before creating a commit.

## Context
- Current git status: !`git status`
- Staged changes: !`git diff --cached`
- Unstaged changes: !`git diff`
- Current branch: !`git branch --show-current`
- Recent commits: !`git log --oneline -10 2>/dev/null || echo "(No commits yet)"`

## Execution Steps

### 1. Analyze Changes
- Check staged changes with `git diff --cached`
- Review changed files with `git status`
- If no staged changes exist, inform the user and exit

### 2. Detect Commit Convention
- Analyze recent commits from `git log --oneline -10` to detect the existing commit message style:
  - **Conventional Commits** (e.g. `feat: add login`, `fix(auth): resolve token issue`)
  - **Emoji prefix** (e.g. `✨ feat: add login`, `🐛 fix: resolve bug`)
  - **Plain style** (e.g. `Add login feature`, `Fix token issue`)
  - **Other patterns** (ticket prefixes like `[PROJ-123]`, etc.)
- Detect the language used in commit messages (English, Japanese, etc.)
- Use the detected convention for new commit messages
- If no prior commits exist, default to **Conventional Commits** in English

### 3. Code Review
- Launch a subagent using the Task tool for code review:
  ```
  Task(subagent_type=general-purpose, prompt="Review the following staged diff. Report any bugs, security issues, performance problems, or readability concerns. If no issues are found, report 'No issues found'.\n\n<diff>\n{staged diff output}\n</diff>")
  ```

### 4. Review Result Handling

**If no issues found:**
- Display: "Review result: No issues found"
- Ask user whether to proceed with the commit
  - **yes** → Proceed to Step 6
  - **no** → Abort

**If issues found:**
- Display review summary with all issues
- Ask user whether to auto-fix with a subagent
  - **yes** → Proceed to Step 5
  - **no** → Ask whether to commit as-is
    - **yes** → Proceed to Step 6
    - **no** → Abort

### 5. Auto-fix (if requested)
- Track fix iteration count (max 3 cycles to prevent infinite loops)
- If already at 3 cycles:
  - Display: "Fix cycle limit (3) reached. Please fix manually."
  - Abort
- Launch a subagent using the Task tool for auto-fix:
  ```
  Task(subagent_type=general-purpose, prompt="Fix the following issues found during code review. Read the target files and apply corrections.\n\nIssues:\n{review issues}\n\nDiff:\n{staged diff output}")
  ```
- After fix completes, stage the modified files:
  ```bash
  git add <fixed files>
  ```
- Increment fix iteration count
- Return to Step 3 (re-review)

### 6. Generate Commit Message & Commit
- Analyze changes and generate a commit message **following the convention detected in Step 2**
- Keep the subject line concise (under 72 characters recommended)
- Stage all changes and commit immediately:
  ```bash
  git add FILE_NAMES
  git commit -m COMMIT_MESSAGE
  ```

### 7. Handle Hooks
- **pre-commit hook** (if it modifies files):
  ```bash
  git status
  git add EDITED_FILE_NAMES
  git commit -m COMMIT_MESSAGE
  ```
- **commit-msg hook** (if format validation fails):
  - Display the validation error message
  - Ask user whether to revise the commit message
    - **yes** → Show the rejected message, ask for corrected version, then retry commit
    - **no** → Abort

## Commit Types Reference (Default Fallback)

Used when no existing convention is detected:

| Emoji | Type | Usage |
|-------|------|-------|
| ✨ | feat | New feature |
| 🐛 | fix | Bug fix |
| 🔨 | refactor | Code refactoring |
| 💄 | style | Style changes |
| ✅ | test | Add or modify tests |
| ⚰️ | remove | Remove code |
| 📚 | docs | Documentation changes |
| 🔍️ | seo | SEO improvements |
| ⚡️ | perf | Performance improvements |
| 💚 | format | Code formatting |
| 📦️ | package | Package updates |
| ⏪️ | revert | Revert changes |
| 🔀 | merge | Merge branches |
| 🚧 | wip | Work in progress |
| 🧑‍💻 | dev | Development changes |
| 🚀 | release | Release |
| 🔧 | chore | Build, tooling, config changes |
| 👷 | ci | CI/CD configuration changes |

## Important Notes

- If no staged changes exist, inform the user and exit
- Follow the repository's existing commit style and language
- If changes span multiple types, focus on the most significant change
- **Push is not performed — the user must push manually**
