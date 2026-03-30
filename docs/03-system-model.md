## End-to-End Flow (Refined)
- SDK client sends a request to the Frontend service.
- Frontend validates and authorizes the request, then routes it to the History service.
  - History service:
    - Interprets the request as workflow commands
    - Appends corresponding events (e.g., workflow start) to workflow history
    - Updates mutable state derived from those events
    - Atomically persists events, state, and generated tasks
    - Generates workflow, activity, and timer tasks
    - These tasks are sent to the Matching service.
  - Matching service:
    - Queues tasks by task queue and partition
    - Workers poll for tasks
    - Matching serves tasks to polling workers (load balanced)
  - Workers:
    - Execute the assigned task (workflow/activity or internal task)
    - Report completion back to History
    - History:
      - Appends completion events
      - Updates state
      - May generate further tasks

***End-to-End Flow in my own words***
SDK client sends request to FrontEnd (FE) service. The FE service authorizes and validates the request.
It then routes it to History service. The history service creates and appends a WorkflowExecutionStarted event
to the workflow history. Also, this service accumulates workflow events and state transitions before atomically
committing events, mutable state, and tasks together to persistent storage. It then generates workflow, activity,
and timer tasks that are consumed by the Matching service. The Matching service receives workflow/activity tasks
produced by History, places the tasks on a queue, and serves them to polling SDK workers and internal worker service.
The worker service executes internal background tasks produced by History via task queues.
It then reports task completion to History.


## Core Invariants
- Workflow history is the source of truth for all workflow executions.
- Workflow events, derived state updates (mutable state), and generated tasks are committed atomically and durably.
- Task execution (Workers) is decoupled from task production (History)
- Task queues are owned by Matching, not by History or Worker.
- Matching does not own workflow state; it only routes tasks.
- Workers do not own persistent workflow state; completion is recorded through History.

## Failure / Boundary Notes
- Workflow execution can always be reconstructed from persisted history; mutable state is derived and recoverable.
- If History loses durability (events), workflow state cannot be reliably recovered, leading to permanent inconsistency
unless replication or backups exist.
- If Matching is unavailable, task dispatch pauses, but tasks remain durable in History and can be retried when Matching
recovers.
- If Workers are unavailable, tasks accumulate in queues, but no state is lost; execution resumes when workers return.
