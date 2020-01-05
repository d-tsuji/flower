-- タスク実行時のステータス管理テーブル
DROP TABLE IF EXISTS kr_task_stat;
CREATE TABLE IF NOT EXISTS kr_task_stat (
    task_flow_id              varchar(256) NOT NULL
,   task_exec_seq             numeric NOT NULL
,   depends_task_exec_seq     numeric NOT NULL
,   parameters                json
,   task_id                   varchar(256) NOT NULL
,   task_seq                  numeric NOT NULL
,   exec_status               numeric NOT NULL
,   task_priority             numeric NOT NULL

,   PRIMARY KEY (task_flow_id, task_exec_seq)
);

-- ワークフローを構成するタスクを管理するマスタテーブル
DROP TABLE IF EXISTS ms_task_definition;
CREATE TABLE IF NOT EXISTS ms_task_definition (
    task_id                  varchar(256) NOT NULL
,   task_seq                 numeric NOT NULL
,   program                  varchar(256) NOT NULL
,   task_priority            numeric NOT NULL
,   param1_key               varchar(1024)
,   param1_value             varchar(1024)
,   param2_key               varchar(1024)
,   param2_value             varchar(1024)
,   param3_key               varchar(1024)
,   param3_value             varchar(1024)
,   param4_key               varchar(1024)
,   param4_value             varchar(1024)
,   param5_key               varchar(1024)
,   param5_value             varchar(1024)

,   PRIMARY KEY (task_id, task_seq)
);

-- テスト用のレコード
INSERT INTO ms_task_definition(task_id, task_seq, program, task_priority, param1_key, param1_value) VALUES ('sample', 1, 'Test1', 10, 'hoge', 'huga', 'piyo', 'foo123');
INSERT INTO ms_task_definition(task_id, task_seq, program, task_priority) VALUES ('sample', 2, 'Test2', 10);
INSERT INTO ms_task_definition(task_id, task_seq, program, task_priority) VALUES ('sample', 3, 'Test3', 10);