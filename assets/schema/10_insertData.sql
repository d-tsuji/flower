-- Test record
INSERT INTO ms_task_definition(task_id, task_seq, program, task_priority) VALUES ('sample', 1, 'EchoRandomTimeSleep', 0);
INSERT INTO ms_task_definition(task_id, task_seq, program, task_priority, param1_key, param1_value) VALUES ('sample', 2, 'EchoParamTimeSleep', 0, 'SLEEP_TIME_SECOND', '3');
INSERT INTO ms_task_definition(task_id, task_seq, program, task_priority, param1_key, param1_value, param2_key, param2_value)
VALUES ('sample', 3, 'HTTPPostRequest', 0, 'URL', 'https://postman-echo.com/post', 'BODY', '{"sample": "test"}');
