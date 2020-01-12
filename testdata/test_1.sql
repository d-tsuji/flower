-- 疎通確認(実行待ちのタスクが返却されること)
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_1', 1, -1, 'test_1', 1, 0, 0, '{}');
