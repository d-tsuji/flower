-- 同一task_idが実行待ちになっているときに、先に進んでいるタスクを優先的に選択すること
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_6_1', 1, -1, 'test_6', 1, 0, 0, '{}');
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_6_1', 2,  1, 'test_6', 2, 0, 0, '{}');
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_6_2', 1, -1, 'test_6', 1, 3, 0, '{}');
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_6_2', 2,  1, 'test_6', 2, 0, 0, '{}');

