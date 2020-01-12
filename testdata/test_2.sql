-- 実行待ちの次のタスクを取得できること
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_2', 1, -1, 'test_2', 1, 3, 0, '{}');
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_2', 2, 1,  'test_2', 2, 0, 0, '{}');
