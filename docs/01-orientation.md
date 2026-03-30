# Core Services:

- ### service/frontend

- ### service/history

- ### service/matching

- ### service/worker

## Repo map (top-level)

    - server code
        - /service
        - /cmd
    - /common contains shared runtime infrastructure used by all server services (not SDKs).
    - API
        - /api
    - Client code
        - /client

### Where the Temporal server starts:

- /cmd/server/main.go

### Where are the entry points?

- /service
- /cmd
- /api

### What binaries/services appear to be launched?

- frontend
- history
- matching
- worker

#### Key subsystems
- persistence
- dynamic config
- auth
- logging / metrics

### Core services you’ll scope to?

- frontend, history,matching, and worker

### Confirm where frontend, history, matching live?

- ### **/service**
  #### Write 1 sentence each on what they seem to do (even if unsure):
        - frontend: handles all external API requests (gRPC) from clients and SDKs; performs validation and routing **but does not execute workflow logic**.
        - history: The History Service handles two main types of request (gRPC) relating to an individual Workflow Execution:
            1. **Requests originating from the User Application** (Start / Cancel / Query / Update / Signal / Reset etc).
            2. **Requests originating from Temporal Workers**, on completion of a Workflow Task or Activity Task.
            See [docs/architecture/history-service.md] (../../docs/architecture/history-service.md)
        - matching: Mannages task queues and task queue partitions.
            - See [docs/architecture/matching-service.md](../../docs/architecture/matching-service.md).
        - worker:  a role in Temporal service used for hosting any components responsible for performing background processing on the Temporal cluster, e.g. replicator, activity workers, etc.. Distinct from SDK workers that execute user-defined workflow/activity code.

### Out of scope for Phase 1

- SDks
- tools

### List directories you’ll ignore initially (SDKs, samples, UI, etc.)

#### Code generation & templates
  - **/*.tmpl
  - **/generated/**
   - **/*_generated.go
- **/tools/** (codegen helpers)

Why:
These are inputs to code generation, not runtime logic. Including them hurts retrieval precision and adds noise.


#### Build, CI, and repo plumbing

  - .github/
  - .buildkite/
  - .circleci/
  - Makefile
  - go.mod, go.sum
  - Dockerfile*
  - .gitignore

Why:
These don’t participate in workflow execution, routing, or persistence logic.


#### Tests (for now)

Ignore for Phase 1.
  - **/*_test.go
  - tests/
  - testdata/

Why:
Tests are useful later for:

#### validation

  - evaluation
  - edge cases
But they obscure the main execution path during initial understanding.
You’ll come back to them later if needed.

#### SDKs and client libraries
Ignore completely in Phase 1.
  - client/
  - sdk/
  - common/sdk/ (if present)
  - language-specific SDK folders

Why:
You are intentionally starting with server-side execution.
SDKs introduce conceptual overload (workers, APIs) too early.

You already planned this — good instinct.

#### UI, samples, and demos

  - samples/
  - ui/
  - web/
  - examples/

Why:
These are consumers of the system, not part of the system.

### Open questions

- How do workers and task queues map to SWF activity workers and queues?
- What is the main workflow worker called in Temporal? In SWF, it is called the decider.

### Notes
Architecture docs are used for cross-checking understanding; code is the primary source of truth.