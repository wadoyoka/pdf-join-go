---
name: tasks-router
description: Orchestrate complex tasks by analyzing instructions, selecting appropriate skills and subagents, creating execution plans, and implementing with user approval. Use for ANY multi-step task, feature request, bug fix, refactoring, or project work that benefits from structured planning and skill coordination.
---

# Tasks Router

Analyze instructions, plan work using available skills and subagents, then execute after approval.

## Project Context

This is **nigopdf**, a Go CLI tool for PDF merging and splitting built with cobra + pdfcpu. The project follows a clean `cmd/` (thin CLI wrappers) and `internal/` (core logic) separation pattern. Tests are table-driven integration tests using real PDF operations.

## Phase 1: Instruction Analysis

When receiving a user instruction:

1. **Classify intent** - Determine the primary goal:
   - New CLI feature / subcommand
   - Bug fix
   - Refactoring
   - Documentation
   - Architecture / design
   - Code review / PR
   - Release / deployment
   - Investigation / research
   - Other

2. **Assess complexity** - Determine if routing is needed:
   - **Simple** (single skill suffices): Route directly to that skill
   - **Complex** (multiple skills/steps needed): Proceed to Phase 2

## Phase 2: Skill Discovery

Available skills in this project:

| Skill | Description | Best For |
|-------|-------------|----------|
| brainstorming | Explore intent, requirements, and design before implementation | Feature scoping, design decisions |
| golang-patterns | Idiomatic Go patterns and best practices | Go code quality, conventions |
| software-architecture | Quality-focused architecture guidance | Code structure, design patterns |
| architect-review | Architecture review for integrity and scalability | Design review, architecture decisions |
| tdd-workflow | Test-driven development with 80%+ coverage | New features, bug fixes with tests |
| subagent-driven-development | Dispatch subagents per task with quality gates | Multi-task parallel execution |
| spec | Spec Driven Development with design and task documents | Complex features with approval gates |
| commit | Create git commits | Committing changes |
| create-pr | Create GitHub pull requests | Submitting changes for review |
| review-pr | Review-only PR analysis | Code review, PR feedback |
| github | GitHub CLI interactions (issues, PRs, CI runs) | GitHub operations |
| technical-writer | Documentation, guides, API references | README, docs, usage guides |
| mermaid-diagrams | Software diagrams in Mermaid syntax | Architecture visualization |

Available agents (in `.claude/agents/`):

| Agent | Description | Best For |
|-------|-------------|----------|
| go-test-runner | Run tests, check coverage, write new tests | Test execution and validation |
| pdf-cli-developer | Implement CLI features, subcommands, bug fixes | Core feature development |
| release-manager | Release prep, versioning, GoReleaser/CI | Publishing new versions |
| go-code-reviewer | Static analysis, linting, Go best practices review | Code quality checks |

## Phase 3: Task Planning

Decompose the instruction into ordered tasks. For each task, assign:

- **Skill/Agent**: Which skill or agent handles this task
- **Execution mode**: `skill-invoke` (Skill tool) | `subagent` (Agent tool) | `manual` (direct implementation)
- **Dependencies**: Which tasks must complete first
- **Parallel group**: Tasks that can run concurrently

### Routing Rules

#### New CLI Feature (subcommand / flag)
1. brainstorming (skill-invoke) - Explore requirements and design
2. golang-patterns (skill-invoke) - Review Go patterns for the feature
3. pdf-cli-developer (subagent) - Implement in cmd/ and internal/
4. go-test-runner (subagent) - Write and run tests
5. go-code-reviewer (subagent) - Review code quality
6. commit (skill-invoke) - Commit changes

#### Bug Fix
1. go-test-runner (subagent) - Write failing test to reproduce
2. pdf-cli-developer (subagent) - Fix the bug
3. go-test-runner (subagent) - Verify all tests pass
4. commit (skill-invoke) - Commit the fix

#### Refactoring
1. architect-review (skill-invoke) - Review current structure
2. software-architecture (skill-invoke) - Plan refactoring approach
3. subagent-driven-development (skill-invoke) - Execute refactoring tasks
4. go-test-runner (subagent) - Verify no regressions
5. go-code-reviewer (subagent) - Review final code quality

#### Release
1. release-manager (subagent) - Check version, validate config, summarize changes
2. commit (skill-invoke) - Commit any pre-release changes
3. Manual - User creates tag and pushes

#### Documentation
1. technical-writer (skill-invoke) - Create or update documentation
2. commit (skill-invoke) - Commit changes

#### Architecture / Design
1. brainstorming (skill-invoke) - Explore design space
2. architect-review (skill-invoke) - Validate architecture
3. mermaid-diagrams (skill-invoke) - Visualize architecture

#### Code Review / PR
1. go-code-reviewer (subagent) - Run linters and review code
2. review-pr (skill-invoke) - Review pull request
   OR
1. create-pr (skill-invoke) - Create pull request

#### Complex Feature (Spec-Driven)
1. brainstorming (skill-invoke) - Explore requirements
2. spec (skill-invoke) - Create design doc with approval gates
3. subagent-driven-development (skill-invoke) - Execute implementation tasks
4. go-test-runner (subagent) - Verify test coverage
5. create-pr (skill-invoke) - Create pull request

### Plan Template

Present the plan to the user in this format:

~~~
## Execution Plan

**Goal**: [summary of what will be accomplished]
**Estimated tasks**: N

| # | Task | Skill/Agent | Mode | Depends on |
|---|------|-------------|------|------------|
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
Agent tool (general-purpose):
  description: "Task N: [task name]"
  prompt: |
    You are executing Task N of the approved plan.

    **Goal**: [task description]
    **Skill to use**: Invoke the [skill-name] skill via the Skill tool
    **Context**: [relevant context from prior tasks]
    **Working directory**: /Users/userName/github-projects/github.com/wadoyoka/nigopdf

    Steps:
    1. Invoke the assigned skill
    2. Follow the skill's guidance to complete the task
    3. Verify the output meets requirements
    4. Report: what was done, files changed, any issues

    IMPORTANT: Use the Skill tool to invoke skills. Do not attempt to replicate skill behavior manually.
```

### Parallel Tasks (independent)

Dispatch multiple Agent tools simultaneously for independent tasks:

```
# Dispatch all independent tasks in a single message
Agent tool: "Task 2: [name]" (run_in_background: true)
Agent tool: "Task 3: [name]" (run_in_background: true)

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
   - Skills/agents used
   - Files changed
   - Any remaining items or recommendations

## Routing Quick Reference

| Keyword/Intent | Primary Skill/Agent | Supporting Skills/Agents |
|---------------|---------------------|--------------------------|
| "add feature", "implement", "create command" | brainstorming | pdf-cli-developer, go-test-runner, golang-patterns |
| "fix bug", "broken", "error", "fail" | go-test-runner | pdf-cli-developer |
| "refactor", "clean up", "restructure" | architect-review | software-architecture, go-code-reviewer |
| "document", "readme", "guide", "usage" | technical-writer | - |
| "review PR", "check PR" | review-pr | go-code-reviewer |
| "create PR", "submit PR" | create-pr | - |
| "release", "version", "publish" | release-manager | commit |
| "test", "coverage", "verify" | go-test-runner | - |
| "lint", "quality", "review code" | go-code-reviewer | golang-patterns |
| "diagram", "visualize", "architecture" | mermaid-diagrams | architect-review |
| "commit" | commit | - |
| "design", "plan", "spec" | brainstorming | spec, architect-review |
| "merge PDF", "split PDF", "CLI" | pdf-cli-developer | golang-patterns |

## Error Handling

- **Skill not found**: Report missing skill, suggest alternatives or manual implementation
- **Subagent failure**: Retry once with more specific instructions, then report to user
- **Ambiguous intent**: Ask user to clarify before planning
- **Circular dependency**: Restructure plan to break the cycle
