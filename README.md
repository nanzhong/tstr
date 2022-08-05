# tstr

:construction: **NOTE** tstr is still a work in progerss. The initial minimal set of core functionality is complete, but there are still a significant number additional features being actively worked on. Before of the first release, backwards compatibility will only be a best effort attempt and is not a guaranatee. :construction:

<hr>

tstr is a test orchestration, visualization, and analysis platform.

tstr is not a testing framework, and it does not impose any patterns on how tests are written. Instead, tests are treated as black block workloads packaged up in containers that can be registered with tstr. tstr handles the scheduling, running, and collecting of results and logs allowing them to be retrieved, visualized, and analyzed from a central location.

![9D114254-B2D1-41CB-9D2F-C6336BB6D6A8](https://user-images.githubusercontent.com/224216/183086353-3262a65f-9c01-43f2-82ff-a77522cde690.jpeg)

There is a live demo instance of the tstr that tracks the main branch of this repository:
- Web UI https://demo.tstr.dev (read only)
- gRPC API grpc.demo.tstr.dev:443 (restricted access)
- json API https://json.demo.tstr.dev (restricted access)

## Features

- Test language/framework agnostic. If you can package up the test workload as a container, tstr can orchestrate it.
- Declarative API driven test configuration.
- Test scheduling via cron style syntax.
- Dynamic test runner registration and test run placement.
- Test matrix scheduling via labels.
- Web UI for visualization.
- gRPC and json APIs for integration/extension (the web UI is built this way).
  - Token based api authentication and scope based authorization.

## How It Works

### Concepts

#### Test

Tests are the fundamental building block that are represented as the combination of:
1. a black box containerized workload that returns a 0 exit code on success
2. configuration that describes how to run the workload

The user is responsible for writing, building and publishing **1**, while tstr is responsible for following **2**.

The configuration consist of:
- The description of the containerized test workload (container image, cmd, args, environment variables, timeout, etc.)
- The set of labels (k/v pairs) that can describe the test which will be used for determining where the test runs
- The schedule (optional) on which the test should be run

#### Run

Runs represent an instance of an execution of a configured test. It contains a snapshot of the configuration and labels that describe the execution and placement of the test.

Each run stores the timing (scheduled at, started at, finished at), result (pass, fail, error, unknown), logs (stdout and stderr), as well as any result data (custom extracted metadata) for the test.

#### Runner

Runners are the placement targets that execute test workloads. Each runner has a set of accept and reject labels selectors that are used to determine where tests can run by matching them against the configured labels of an instance of a test run.

### Architecture

```
                                       ┌─────────────┐
                                    ┌──┤ tstr runner ├─┐
     ┌──────────┐                   │  └─┬───────────┘ ├─┐
     │ tstr ctl ├───┐               │    └─┬───────────┘ │
     └──────────┘   │               │      └─────────────┘
                    │               │
┌───────────────┐   │               │     ┌──────────┐
│               ├───┴──────► gRPC ──┴────►│          │
│ Other clients │                         │ tstr api │
│               ├───┬──────► json ───────►│          │
└───────────────┘   │                     └───┬──▲───┘
                    │                         │  │
      ┌─────────┐   │                ┌────────▼──┴───────────┐
      │ tstr ui ├───┘                │ Database (PostgreSQL) │
      └─────────┘                    └───────────────────────┘
```

TODO

## Usage

tstr is bundled as a single binary that includes all the components.

### `tstr api`

TODO

### `tstr ui`

TODO

### `tstr run`

TODO

### `tstr ctl`

TODO
