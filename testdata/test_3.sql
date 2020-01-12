-- 依存するタスクが実行中の場合に、実行待ちの次のタスクが取得されないこと
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_3', 1, -1, 'test_3', 1, 1, 0, '{}');
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_3', 2,  1, 'test_3', 2, 0, 0, '{}');

