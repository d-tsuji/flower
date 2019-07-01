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

---------------------------------------------

DROP TABLE IF EXISTS ms_task;

CREATE TABLE ms_task
(
    task_id character varying(256),
    exec_order numeric,
    endpoint character varying(256),
    method character varying(256),
    extend_parameters character varying(256),
    wait_mode character varying(256)
);

truncate ms_task;

insert into ms_task
    (task_id, exec_order, endpoint, method, extend_parameters, wait_mode)
values('sample.a.id', 1, 'http://httpbin.org/ip', 'GET', '{}', 'N');
insert into ms_task
    (task_id, exec_order, endpoint, method, extend_parameters, wait_mode)
values('sample.a.id', 2, 'http://httpbin.org/user-agent', 'GET', '{}', 'N');
insert into ms_task
    (task_id, exec_order, endpoint, method, extend_parameters, wait_mode)
values('sample.a.id', 3, 'http://codeforces.com/api/user.rating?handle=tutuz', 'GET', '{}', 'N');
insert into ms_task
    (task_id, exec_order, endpoint, method, extend_parameters, wait_mode)
values('sample.a.id', 4, 'http://httpbin.org/post', 'POST', '{param: hoge}', 'N');

insert into ms_task
    (task_id, exec_order, endpoint, method, extend_parameters, wait_mode)
values('sample.b.id', 1, 'http://codeforces.com/api/user.rating?handle=tanzaku', 'GET', '{}', 'N');
insert into ms_task
    (task_id, exec_order, endpoint, method, extend_parameters, wait_mode)
values('sample.b.id', 2, 'http://codeforces.com/api/user.rating?handle=chokudai', 'GET', '{}', 'N');
