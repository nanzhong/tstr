SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: access_token_scope; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.access_token_scope AS ENUM (
    'admin',
    'control_r',
    'control_rw'
);


--
-- Name: run_result; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.run_result AS ENUM (
    'unknown',
    'pass',
    'fail',
    'error'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: access_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.access_tokens (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying NOT NULL,
    token_hash character varying NOT NULL,
    scopes public.access_token_scope[],
    issued_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamp with time zone,
    revoked_at timestamp with time zone
);


--
-- Name: runners; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.runners (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying NOT NULL,
    accept_test_labels jsonb,
    reject_test_labels jsonb,
    registered_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    approved_at timestamp with time zone,
    revoked_at timestamp with time zone,
    last_heartbeat_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: runs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.runs (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    test_id uuid,
    test_run_config_id uuid,
    runner_id uuid,
    result public.run_result DEFAULT 'unknown'::public.run_result,
    logs jsonb,
    scheduled_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    started_at timestamp with time zone,
    finished_at timestamp with time zone
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(255) NOT NULL
);


--
-- Name: test_run_configs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.test_run_configs (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    test_id uuid,
    container_image character varying NOT NULL,
    command character varying,
    args character varying[],
    env jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: test_suites; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.test_suites (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying NOT NULL,
    labels jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    archived_at timestamp with time zone
);


--
-- Name: tests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tests (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying NOT NULL,
    labels jsonb,
    cron_schedule character varying,
    registered_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    archived_at timestamp with time zone
);


--
-- Name: access_tokens access_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.access_tokens
    ADD CONSTRAINT access_tokens_pkey PRIMARY KEY (id);


--
-- Name: runners runners_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.runners
    ADD CONSTRAINT runners_pkey PRIMARY KEY (id);


--
-- Name: runs runs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.runs
    ADD CONSTRAINT runs_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: test_run_configs test_run_configs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.test_run_configs
    ADD CONSTRAINT test_run_configs_pkey PRIMARY KEY (id);


--
-- Name: test_suites test_suites_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.test_suites
    ADD CONSTRAINT test_suites_name_key UNIQUE (name);


--
-- Name: test_suites test_suites_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.test_suites
    ADD CONSTRAINT test_suites_pkey PRIMARY KEY (id);


--
-- Name: tests tests_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tests
    ADD CONSTRAINT tests_name_key UNIQUE (name);


--
-- Name: tests tests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tests
    ADD CONSTRAINT tests_pkey PRIMARY KEY (id);


--
-- Name: access_tokens_revoked_at_expires_at_token_hash_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX access_tokens_revoked_at_expires_at_token_hash_idx ON public.access_tokens USING btree (revoked_at, expires_at, token_hash);


--
-- Name: runners_last_heartbeat_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX runners_last_heartbeat_at_idx ON public.runners USING btree (last_heartbeat_at);


--
-- Name: runs_result_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX runs_result_idx ON public.runs USING btree (result);


--
-- Name: runs_runner_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX runs_runner_id_idx ON public.runs USING btree (runner_id);


--
-- Name: runs_scheduled_at_started_at_finished_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX runs_scheduled_at_started_at_finished_at_idx ON public.runs USING btree (scheduled_at, started_at, finished_at);


--
-- Name: runs_test_id_test_run_config_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX runs_test_id_test_run_config_id_idx ON public.runs USING btree (test_id, test_run_config_id);


--
-- Name: test_run_configs_created_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX test_run_configs_created_at_idx ON public.test_run_configs USING btree (created_at);


--
-- Name: test_suites_archived_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX test_suites_archived_at_idx ON public.test_suites USING btree (archived_at);


--
-- Name: test_suites_unique_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX test_suites_unique_name ON public.test_suites USING btree (name) WHERE (archived_at IS NULL);


--
-- Name: tests_labels_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tests_labels_idx ON public.tests USING gin (labels);


--
-- Name: tests_unique_name; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX tests_unique_name ON public.tests USING btree (name) WHERE (archived_at IS NULL);


--
-- Name: runs runs_runner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.runs
    ADD CONSTRAINT runs_runner_id_fkey FOREIGN KEY (runner_id) REFERENCES public.runners(id);


--
-- Name: runs runs_test_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.runs
    ADD CONSTRAINT runs_test_id_fkey FOREIGN KEY (test_id) REFERENCES public.tests(id);


--
-- Name: runs runs_test_run_config_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.runs
    ADD CONSTRAINT runs_test_run_config_id_fkey FOREIGN KEY (test_run_config_id) REFERENCES public.test_run_configs(id);


--
-- Name: test_run_configs test_run_configs_test_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.test_run_configs
    ADD CONSTRAINT test_run_configs_test_id_fkey FOREIGN KEY (test_id) REFERENCES public.tests(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20220426134459');
