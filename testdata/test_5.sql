-- 同一task_idが実行待ちになっているときに、同時に実行されないこと
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_5_1', 1, -1, 'test_5', 1, 0, 0, '{}');
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_5_2', 1, -1, 'test_5', 1, 0, 0, '{}');

