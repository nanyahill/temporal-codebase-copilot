# Core Services Overview

## Frontend Service

**Primary responsibility**

- Handles all external API requests (gRPC / HTTP) from clients and SDKs.
- Performs request validation, authorization, and routing.
- Does **not** execute workflow logic or persist workflow state.

**Entry points (APIs / handlers)**

- gRPC WorkflowService handlers exposed to clients and SDKs.
- HTTP API server (for REST-style access where applicable).
- Service bootstrap via server initialization (wiring, not business logic).

**What state it owns (or does not own)**

State it owns:

- No durable workflow state.
- Only transient request-scoped state during API handling.

State it does not own:

- Workflow execution state
- Workflow history
- Task queues
- Retry or backoff state

**What other services it talks to**

- History service (for workflow execution–related operations)
- Matching service (for task queue–related operations)
- Worker service (indirectly, via internal task routing)
- Persistence layer only indirectly (never as the system of record)

**Key abstractions / concepts**

- gRPC service boundaries
- Request validation and authorization
- Namespace routing
- Service-level interception (metrics, auth, rate limiting)
- Separation of API handling from execution logic

**Notes / open questions**

- Which requests are routed to history vs matching in which cases?
- How much validation happens here vs downstream services?

## History Service

**Primary responsibility**
- Workflow history events are the source of truth for workflow execution.
- Accumulates workflow events and state transitions before committing them.
- Transforms mutable states into workflow mutations/snapshots and persists them to the persistence layer.
- Generates workflow, activity, and timer tasks that are later consumed via Matching.
- Partitions workflow executions across shards for scalability.
- Enforces workflow state transitions / invariants
- Does **not** create task queues nor match task to task queues nor does it consume and execute workflow tasks.

**What state it owns (or does not own)**
  - Workflow history events
  - Mutable workflow state
  - Shard ownership of workflow executions

State it does not own:
  - Task queue matching
  - Workflow/activity execution
  - Client request validation (frontend)

**What other services it talks to**
- Matching service (for task queue–related operations)
- Worker service (indirectly, via internal task routing)
- Persistence layer (via ExecutionManager) as the system of record for workflow history and execution state.

**Key abstractions / concepts**

- Mutable state
- Workflow context
- Workflow snapshot
- Shards
- History Events (event sourcing)
- Workflow Mutation / transaction boundary
- Task generation (transfer/timer tasks)

**Notes / open questions**

- How are workflow states reconstructed from History events?
- How does History service interact with Matching service?

## Matching Service

**Primary responsibility**
- Manages task queues and partitions
- Receives workflow/activity tasks produced by History and dispatches them to polling workers
- Serves poll requests from workers and load balances/scales task delivery

**What state it owns (or does not own)**
- Task queue state (including pending tasks and queue metadata)
- Poller coordination state
- Task queue partition state

State it does not own:
- Worker lifecycle
- Workflow/activity task execution
- Workflow history and mutable state

**What other services it talks to**
- History service (produces workflow/activity tasks that Matching ingests and queues)
- Worker service (polls Matching for tasks; Matching coordinates task delivery)


## Worker Service

**Primary responsibility**
- Supports system level workflows such as replication, and timers
- Consumes and executes internal background tasks produced by History via queues

**What state it owns (or does not own)**
- Ephemeral task execution state / processing context

State it does not own:
- Client request validation (frontend)
- Workflow/activity code execution (SDK workers)
- Workflow history and mutable state (History)
- Poller coordination state (Matching)
- Task queues (Matching)