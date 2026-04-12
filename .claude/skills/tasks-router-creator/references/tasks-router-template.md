# tasks-router SKILL.md Template

Use this template when generating a tasks-router skill. Replace all `{{PLACEHOLDER}}` values with project-specific content.

---

```markdown
---
name: tasks-router
description: Orchestrate complex tasks by analyzing instructions, selecting appropriate skills and subagents, creating execution plans, and implementing with user approval. Use for ANY multi-step task, feature request, bug fix, refactoring, or project work that benefits from structured planning and skill coordination.
---

# Tasks Router

Analyze instructions, plan work using available skills and subagents, then execute after approval.

## Phase 1: Instruction Analysis

When receiving a user instruction:

1. **Classify intent** - Determine the primary goal:
   - Feature development
   - Bug fix
   - Refactoring
   - Documentation
   - Architecture/design
   - Code review / PR
   - Investigation / research
   - Other

2. **Assess complexity** - Determine if routing is needed:
   - **Simple** (single skill suffices): Route directly to that skill
   - **Complex** (multiple skills/steps needed): Proceed to Phase 2

## Phase 2: Skill Discovery

Scan `.claude/skills/` to build available skill registry:

{{SKILL_REGISTRY}}

## Phase 3: Task Planning

Decompose the instruction into ordered tasks. For each task, assign:

- **Skill**: Which skill handles this task
- **Execution mode**: `skill-invoke` (Skill tool) | `subagent` (Task tool) | `manual` (direct implementation)
- **Dependencies**: Which tasks must complete first
- **Parallel group**: Tasks that can run concurrently

### Routing Rules

{{ROUTING_RULES}}

### Plan Template

Present the plan to the user in this format:

~~~
## Execution Plan

**Goal**: [summary of what will be accomplished]
**Estimated tasks**: N

| # | Task | Skill | Mode | Depends on |
|---|------|-------|------|------------|
| 1 | [description] | [skill-name] | [mode] | - |
| 2 | [description] | [skill-name] | [mode] | #1 |
| 3 | [description] | [skill-name] | [mode] | #1 |
| 4 | [description] | [skill-name] | [mode] | #2, #3 |

**Parallel opportunities**: Tasks #2 and #3 can run concurrently.

Proceed with this plan?
~~~

## Phase 4: Execution

After user approval:

1. **Create task tracking** - Use TaskCreate for each task in the plan
2. **Set dependencies** - Use TaskUpdate to set blockedBy relationships
3. **Execute sequentially or in parallel** following the subagent-driven-development pattern:

### Sequential Tasks (dependent)

For each task in dependency order:

```
Task tool (general-purpose):
  description: "Task N: [task name]"
  prompt: |
    You are executing Task N of the approved plan.

    **Goal**: [task description]
    **Skill to use**: Invoke the [skill-name] skill via the Skill tool
    **Context**: [relevant context from prior tasks]
    **Working directory**: [project root]

    Steps:
    1. Invoke the assigned skill
    2. Follow the skill's guidance to complete the task
    3. Verify the output meets requirements
    4. Report: what was done, files changed, any issues

    IMPORTANT: Use the Skill tool to invoke skills. Do not attempt to replicate skill behavior manually.
```

### Parallel Tasks (independent)

Dispatch multiple Task tools simultaneously for independent tasks:

```
# Dispatch all independent tasks in a single message
Task tool: "Task 2: [name]" (run_in_background: true)
Task tool: "Task 3: [name]" (run_in_background: true)

# Wait for completion, then review
```

### After Each Task

- Mark task as completed in TaskUpdate
- Review output for quality
- If issues found, dispatch fix subagent before proceeding

## Phase 5: Completion

After all tasks complete:

1. Review overall result against original instruction
2. Report summary to user:
   - What was accomplished
   - Skills used
   - Files changed
   - Any remaining items or recommendations

## Routing Quick Reference

{{ROUTING_QUICK_REFERENCE}}

## Error Handling

- **Skill not found**: Report missing skill, suggest alternatives or manual implementation
- **Subagent failure**: Retry once with more specific instructions, then report to user
- **Ambiguous intent**: Ask user to clarify before planning
- **Circular dependency**: Restructure plan to break the cycle
```

---

## Template Variables

### {{SKILL_REGISTRY}}

Replace with a table of discovered skills:

```markdown
| Skill | Description | Best For |
|-------|-------------|----------|
| brainstorming | Explore intent and design before implementation | Feature scoping, requirements |
| software-architecture | Clean architecture and DDD guidance | Code structure, design decisions |
| tdd-workflow | Test-driven development with 80%+ coverage | New features, bug fixes |
| subagent-driven-development | Dispatch subagents per task | Multi-task execution |
| architect-review | Architecture review and patterns | Design review, scalability |
| ... | ... | ... |
```

### {{ROUTING_RULES}}

Replace with project-specific routing patterns:

```markdown
#### Feature Development
1. brainstorming (skill-invoke) - Explore requirements
2. architect-review (skill-invoke) - Validate architecture
3. tdd-workflow (subagent) - Implement with tests
4. create-pr (skill-invoke) - Create pull request

#### Bug Fix
1. tdd-workflow (subagent) - Write failing test, then fix
2. commit (skill-invoke) - Commit the fix

#### Refactoring
1. architect-review (skill-invoke) - Review current structure
2. software-architecture (skill-invoke) - Plan refactoring approach
3. subagent-driven-development (subagent) - Execute refactoring tasks

#### Documentation
1. technical-writer (skill-invoke) - Create or update documentation

#### Code Review
1. review-pr (skill-invoke) - Review pull request
```

### {{ROUTING_QUICK_REFERENCE}}

Replace with a compact lookup table:

```markdown
| Keyword/Intent | Primary Skill | Supporting Skills |
|---------------|---------------|-------------------|
| "add feature", "implement", "create" | brainstorming | architect-review, tdd-workflow |
| "fix bug", "broken", "error" | tdd-workflow | subagent-driven-development |
| "refactor", "clean up", "restructure" | architect-review | software-architecture |
| "document", "readme", "guide" | technical-writer | - |
| "review PR", "check PR" | review-pr | - |
| "create PR", "submit PR" | create-pr | - |
| "diagram", "visualize", "architecture" | mermaid-diagrams | architect-review |
| "commit" | commit | - |
```
