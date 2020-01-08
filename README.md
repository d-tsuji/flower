# Flower [![Go Report Card](https://goreportcard.com/badge/github.com/d-tsuji/flower)](https://goreportcard.com/report/github.com/d-tsuji/flower) ![License MIT](https://img.shields.io/badge/license-MIT-blue.svg) [![Actions Status](https://github.com/d-tsuji/flower/workflows/build/badge.svg)](https://github.com/d-tsuji/flower/actions) [![GoDoc](https://godoc.org/github.com/d-tsuji/flower?status.svg)](https://godoc.org/github.com/d-tsuji/flower)

Flower is a workflow engine. Manages the execution of a series of tasks that make up a workflow.

- feature
  - Control task execution order
  - Easily find out which task caused the error
  - Rerun or recover from a task
  - Flow control for multiple workflows
  - Prioritize workflow

## System Overview

![System overview](/doc/images/system_overview.png "System overview")

## Tasks

Tasks that compose a workflow are defined in [DAG](https://en.wikipedia.org/wiki/Directed_acyclic_graph) as follows.

![Task structure](/doc/images/task_structure.png "Task structure")

## Getting Started

Here is how to start flower using docker-compose.

```
$ docker-compose up
Starting flower_db_1 ... done
Starting register    ... done
Starting watcher     ... done
Attaching to flower_db_1, watcher, register
db_1        |
db_1        | PostgreSQL Database directory appears to contain a database; Skipping initialization
db_1        |
db_1        | 2020-01-05 14:59:26.373 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
db_1        | 2020-01-05 14:59:26.373 UTC [1] LOG:  listening on IPv6 address "::", port 5432
db_1        | 2020-01-05 14:59:26.386 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
db_1        | 2020-01-05 14:59:26.445 UTC [12] LOG:  database system was shut down at 2020-01-05 14:59:15 UTC
db_1        | 2020-01-05 14:59:26.454 UTC [1] LOG:  database system is ready to accept connections
watcher     | 2020/01/05 14:59:27 Waiting for: tcp://db:5432
watcher     | 2020/01/05 14:59:27 Connected to tcp://db:5432
watcher     | 2020/01/05 14:59:27 [dispatcher] starting worker: 1
watcher     | 2020/01/05 14:59:27 [dispatcher] starting worker: 2
watcher     | 2020/01/05 14:59:27 [dispatcher] starting worker: 3
watcher     | 2020/01/05 14:59:27 [dispatcher] starting worker: 4
watcher     | 2020/01/05 14:59:27 [dispatcher] starting worker: 5
register    | 2020/01/05 14:59:27 Waiting for: tcp://db:5432
register    | 2020/01/05 14:59:27 Connected to tcp://db:5432
register    | 2020/01/05 14:59:27 [register] starting server on address: 0.0.0.0:8000
watcher     | 2020/01/05 14:59:32 [watcher] watching task...
```

Note: Application of *watcher* and *register* depend on starting Database. Therefore, it is controlled using [dockerize](https://github.com/jwilder/dockerize).

## Configuration

Flower consists of two main tables. **ms_task_definition** and **kr_task_stat**.

### `ms_task_definition`

`ms_task_definition` is a table that defines the tasks that make up the workflow.

| Column        | Primary Key | Data type     | Constraint |
| ------------- | :---------: | ------------- | ---------- |
| task_id       |     ✔️      | varchar(256)  | NOT NULL   |
| task_seq      |     ✔️      | numeric       | NOT NULL   |
| program       |             | varchar(256)  | NOT NULL   |
| task_priority |             | numeric       | NOT NULL   |
| param1_key    |             | varchar(1024) |            |
| param1_value  |             | varchar(1024) |            |
| param2_key    |             | varchar(1024) |            |
| param2_value  |             | varchar(1024) |            |
| param3_key    |             | varchar(1024) |            |
| param3_value  |             | varchar(1024) |            |
| param4_key    |             | varchar(1024) |            |
| param4_value  |             | varchar(1024) |            |
| param5_key    |             | varchar(1024) |            |
| param5_value  |             | varchar(1024) |            |

We have registered a series of tasks that make up a workflow in the master in advance. The following is an example of a record to be registered. The workflow called `example` consists of three tasks. Register the tasks you want to execute in a series of workflows as records. If you register a workflow, you need to register a series of tasks in `ms_task_definition`.

Actually, the Go program registered in the `program` column is executed by reflection. The tasks that make up your workflow are implemented as Go programs and registered in the master as `program`. This is very useful if you want to use the same task in different workflows.

#### Example

| task_id | task_seq | program             | task_priority | param1_key        | param1_value                  | param2_key | param2_value   | ... |
| ------- | -------- | ------------------- | ------------- | ----------------- | ----------------------------- | ---------- | -------------- | --- |
| example | 1        | EchoRandomTimeSleep | 10            |                   |                               |            |                | ... |
| example | 2        | EchoParamTimeSleep  | 10            | SLEEP_TIME_SECOND | 3                             |            |                | ... |
| example | 3        | HTTPPostRequest     | 10            | URL               | https://postman-echo.com/post | BODY       | {"id": "test"} | ... |

### `kr_task_stat`

`kr_task_stat` is a table that manages the execution of the tasks that make up the workflow. The task is registered as a DAG in `kr_task_stat`.

| Column                | Primary Key | Data type                | Constraint |
| --------------------- | :---------: | ------------------------ | ---------- |
| task_flow_id          |     ✔️      | varchar(256)             | NOT NULL   |
| task_exec_seq         |     ✔️      | numeric                  | NOT NULL   |
| depends_task_exec_seq |             | numeric                  | NOT NULL   |
| task_id               |             | varchar(256)             | NOT NULL   |
| task_seq              |             | numeric                  | NOT NULL   |
| exec_status           |             | numeric                  | NOT NULL   |
| task_priority         |             | numeric                  | NOT NULL   |
| parameters            |             | json                     | NOT NULL   |
| registered_ts         |             | timestamp with time zone |            |
| started_ts            |             | timestamp with time zone |            |
| finished_ts           |             | timestamp with time zone |            |
| suspended_ts          |             | timestamp with time zone |            |

Note: We can register a task as waiting by executing an HTTP request or a job. Currently, only the following HTTP requests are supported. With the following HTTP request, the task of the workflow registered in `ms_task_definition` is registered in `kr_task_stat` as waiting to be executed.

#### Example

The following curl command is a command to call the execution of the workflow whose task_id is `example`.

```console
$ curl -X POST localhost:8000/register/example
```

The above command registers the workflow as waiting task to be executed in `kr_task_stat`. The following records are created.

| task_flow_id                         | task_exec_seq | depends_task_exec_seq | task_id | task_seq | exec_status | task_priority | parameters                                                          |
| ------------------------------------ | ------------- | --------------------- | ------- | -------- | ----------- | ------------- | ------------------------------------------------------------------- |
| da03a7a9-31e5-11ea-8ff9-0242ac1f0003 | 1             | -1                    | example | 1        | 3           | 0             | {}                                                                  |
| da03a7a9-31e5-11ea-8ff9-0242ac1f0003 | 2             | 1                     | example | 2        | 1           | 0             | {"SLEEP_TIME_SECOND":"3"}                                           |
| da03a7a9-31e5-11ea-8ff9-0242ac1f0003 | 3             | 2                     | example | 3        | 0           | 0             | {"BODY":"{\"id\": \"test\"}","URL":"https://postman-echo.com/post"} |

`task_status` is a value indicating the task execution status as follows.

| value | status  | description                                      |
| ----- | ------- | ------------------------------------------------ |
| 0     | Wait    | The task that are waiting to be executed         |
| 1     | Running | The task in Running                              |
| 2     | Suspend | The task that has been suspended for some reason |
| 3     | Finish  | The task finished                                |
| 9     | Ignore  | The task ignored                                 |

## HTTP API

Registration of workflow execution is performed via HTTP API.

Overview of endpoints:

- [`POST /register/{task_id}`](#post-registertask_id): Registration of workflow to execute.

### `POST /register/{task_id}`

Registration of workflow to execute in `kr_task_stat`.

#### Request

```http
POST /register/{task_id}
```

##### Parameters <!-- omit in toc -->

- `task_id` (string,required): Id for registering execution of workflow. Must be registered in `ms_task_definition`.

#### Response

```json
{
  "status": "succeeded",
  "taskId": "`task_id`"
}
```

## Author

Tsuji Daishiro

## LICENSE

This software is licensed under the MIT license, see [LICENSE](https://github.com/d-tsuji/flower/blob/master/LICENSE) for more information.
