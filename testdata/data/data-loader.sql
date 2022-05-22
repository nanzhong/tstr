create table if not exists json_tests (row text);
\copy json_tests from '/home/h/code/tstr/testdata/data/tests.json';
insert into tests (select q.* from json_tests, json_populate_record(null::tests, row::json) as q) on conflict do nothing;
drop table json_tests;

create table if not exists json_test_run_configs (row text);
\copy json_test_run_configs from '/home/h/code/tstr/testdata/data/test_run_configs.json';
insert into test_run_configs (select q.* from json_test_run_configs, json_populate_record(null::test_run_configs, row::json) as q) on conflict do nothing;
drop table json_test_run_configs;


create table if not exists json_runners (row text);
\copy json_runners from '/home/h/code/tstr/testdata/data/runners.json';
insert into runners (select q.* from json_runners, json_populate_record(null::runners, row::json) as q) on conflict do nothing;
drop table json_runners;


create table if not exists json_runs (row text);
\copy json_runs from '/home/h/code/tstr/testdata/data/runs.json';
insert into runs (select q.* from json_runs, json_populate_record(null::runs, row::json) as q) on conflict do nothing;
drop table json_runs;

