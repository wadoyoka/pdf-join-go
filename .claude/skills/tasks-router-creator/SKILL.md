---
name: tasks-router-creator
description: Create a tasks-router skill for any project. tasks-router analyzes user instructions, discovers available skills in .claude/skills/, plans work using skills and subagents, then executes after user approval. Use when asked to "create tasks-router", "set up task routing", "add tasks-router to project", or when a project needs orchestration of its existing skills and subagents.
---

# Tasks Router Creator

Create a `tasks-router` skill tailored to the target project. The generated skill orchestrates work by discovering available skills, planning task decomposition, and executing via subagents.

## Prerequisites

- Target project must have `.claude/skills/` directory with at least 2 skills
- `brainstorming` skill must be available (used during creation process)

## Creation Process

### Step 1: Discover Project Skills

Scan the target project's `.claude/skills/` directory:

```
1. Glob .claude/skills/*/SKILL.md
2. Read each SKILL.md's frontmatter (name + description only)
3. Build a skill registry: { name, description, path }
```

### Step 2: Analyze Project Context

Use the `brainstorming` skill to understand:

- What types of tasks does this project handle?
- Which skills complement each other? (e.g., brainstorming -> software-architecture -> tdd-workflow)
- What common workflows exist? (e.g., feature development, bug fix, refactoring)
- Are there skills that should always be used together?

### Step 3: Design Routing Rules

Based on Step 2, define routing patterns:

**Routing pattern structure:**

```yaml
pattern:
  trigger: "keyword or intent pattern"
  workflow:
    - step: 1
      skill: "skill-name"
      mode: "invoke"  # invoke = use Skill tool, reference = read for guidance
    - step: 2
      skill: "another-skill"
      mode: "subagent"  # dispatch as Task subagent
  parallel_candidates:
    - "independent tasks that can run concurrently"
```

Common routing patterns to consider:

| Intent | Typical Workflow |
|--------|-----------------|
| New feature | brainstorming -> architect-review -> tdd-workflow -> subagent-driven-development |
| Bug fix | tdd-workflow -> subagent-driven-development |
| Refactoring | architect-review -> software-architecture -> subagent-driven-development |
| PR/Review | review-pr or create-pr |
| Documentation | technical-writer |
| Architecture | brainstorming -> architect-review -> mermaid-diagrams |

### Step 4: Generate tasks-router Skill

1. Read the template at `references/tasks-router-template.md`
2. Customize with:
   - The discovered skill registry from Step 1
   - The routing patterns from Step 3
   - Project-specific workflow sequences
3. Write to target project's `.claude/skills/tasks-router/SKILL.md`

### Step 5: Verify

- Confirm the generated SKILL.md has valid frontmatter
- Confirm all referenced skills exist in the project
- Test with a sample instruction to validate routing logic

## Output

The generated `tasks-router` skill will:

1. **Analyze** user instructions to identify intent and required capabilities
2. **Discover** available skills by scanning `.claude/skills/`
3. **Plan** a task breakdown with skill/subagent assignments
4. **Present** the plan to the user for approval (with TaskCreate for tracking)
5. **Execute** approved plan using `subagent-driven-development` patterns

## Key Design Principles

- **Skill-first**: Always check if an existing skill handles the task before creating custom logic
- **Plan before execute**: Never skip the planning + approval step
- **Subagent isolation**: Each task gets a fresh subagent to prevent context pollution
- **Progressive disclosure**: Show high-level plan first, detail on demand
- **Fail fast**: If a required skill is missing, report immediately rather than proceeding without it
