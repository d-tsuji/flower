DROP TABLE IF EXISTS kr_task_status;

CREATE TABLE kr_task_status
(
	job_flow_id character varying(256),
	task_id character varying(256),
	job_exec_seq numeric,
	job_depend_exec_seq numeric,
	wait_mode character varying(256),
	status character varying(128),
	response_body text,
	priority numeric,
	create_ts timestamp,
	start_ts timestamp,
	finish_ts timestamp
);

select *
from kr_task_status;

select job_flow_id, task_id, job_exec_seq, job_depend_exec_seq,
	wait_mode, status, response_body, create_ts, finish_ts
from kr_task_status;
