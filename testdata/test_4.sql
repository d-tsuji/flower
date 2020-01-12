-- 複数のタスクが同時に取得できること
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_4_1', 1, -1, 'test_4_1', 1, 0, 0, '{}');
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_4_2', 1, -1, 'test_4_2', 1, 0, 0, '{}');

