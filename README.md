# flower [![Go Report Card](https://goreportcard.com/badge/github.com/d-tsuji/flower)](https://goreportcard.com/report/github.com/d-tsuji/flower) ![License MIT](https://img.shields.io/badge/license-MIT-blue.svg) [![Actions Status](https://github.com/d-tsuji/flower/workflows/build/badge.svg)](https://github.com/d-tsuji/flower/actions)

Flower is a workflow engine. Manages the execution of a series of tasks that make up a workflow. It manages the status of a series of tasks to be executed, and has a mechanism to quickly find a recovery point in the event of an error. Similarly, it has a mechanism that makes recovery such as reruns easy. Supports parallel execution of tasks and flow control by worker pool.

## System Overview

![System overview](/doc/images/system_overview.png "System overview")

## Tasks

Tasks that compose a workflow are defined in DAG as follows.

![Task structure](/doc/images/task_structure.png "Task structure")

## Usage

```
$ docker-compose -d up
```

### Detail

We have registered a series of tasks that make up a workflow in the master in advance.

| task_id | task_seq | program | task_priority |
| ------- | -------- | ------- | ------------- |
| sample  | 1        | Test1   | 10            |
| sample  | 2        | Test2   | 10            |
| sample  | 3        | Test3   | 10            |

We can register a task as pending by executing an HTTP request or a job. Currently, only the following HTTP requests are supported. With the following HTTP request, the task of the workflow registered in `ms_task_definition` is registered in `kr_task_stat` as waiting to be executed.

```
$ curl -X POST -H 'Content-Type:application/json' localhost:8000/register -i -d '{"taskId": "sample"}'
```

The above command registers the task waiting to be executed in `kr_task_stat`. The following records are created. (exec_status = 0 is status to wait execution.)

| task_flow_id | task_exec_seq | depends_task_exec_seq | task_id | task_seq | exec_status | task_priority |
| ------------ | ------------- | --------------------- | ------- | -------- | ----------- | ------------- |
| xxxxxxxxx    | 1             | 0                     | sample  | 1        | 0           | 10            |
| xxxxxxxxx    | 2             | 1                     | sample  | 2        | 0           | 10            |
| xxxxxxxxx    | 3             | 2                     | sample  | 3        | 0           | 10            |

## LICENSE

This software is licensed under the MIT license, see [LICENSE](https://github.com/d-tsuji/flower/blob/master/LICENSE) for more information.
