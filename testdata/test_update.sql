-- 疎通確認(実行待ちのタスクが実行中に更新されること)
INSERT INTO kr_task_stat(task_flow_id, task_exec_seq, depends_task_exec_seq, task_id, task_seq, exec_status, task_priority, parameters) VALUES ('test_flow_update_1', 1, -1, 'test_update', 1, 0, 0, '{}');
