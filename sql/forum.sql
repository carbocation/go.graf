--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: askbitcoin; Type: SCHEMA; Schema: -; Owner: projectuser
--

CREATE SCHEMA askbitcoin;


ALTER SCHEMA askbitcoin OWNER TO projectuser;

SET search_path = askbitcoin, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account; Type: TABLE; Schema: askbitcoin; Owner: projectuser; Tablespace: 
--

CREATE TABLE account (
    id bigint NOT NULL,
    handle text NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    password text
);


ALTER TABLE askbitcoin.account OWNER TO projectuser;

--
-- Name: account_id_seq; Type: SEQUENCE; Schema: askbitcoin; Owner: projectuser
--

CREATE SEQUENCE account_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE askbitcoin.account_id_seq OWNER TO projectuser;

--
-- Name: account_id_seq; Type: SEQUENCE OWNED BY; Schema: askbitcoin; Owner: projectuser
--

ALTER SEQUENCE account_id_seq OWNED BY account.id;


--
-- Name: entry; Type: TABLE; Schema: askbitcoin; Owner: projectuser; Tablespace: 
--

CREATE TABLE entry (
    id bigint NOT NULL,
    title text NOT NULL,
    body text NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    author_id bigint NOT NULL
);


ALTER TABLE askbitcoin.entry OWNER TO projectuser;

--
-- Name: entry_closures; Type: TABLE; Schema: askbitcoin; Owner: projectuser; Tablespace: 
--

CREATE TABLE entry_closures (
    ancestor bigint NOT NULL,
    descendant bigint NOT NULL,
    depth integer NOT NULL
);


ALTER TABLE askbitcoin.entry_closures OWNER TO projectuser;

--
-- Name: entry_id_seq; Type: SEQUENCE; Schema: askbitcoin; Owner: projectuser
--

CREATE SEQUENCE entry_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE askbitcoin.entry_id_seq OWNER TO projectuser;

--
-- Name: entry_id_seq; Type: SEQUENCE OWNED BY; Schema: askbitcoin; Owner: projectuser
--

ALTER SEQUENCE entry_id_seq OWNED BY entry.id;


--
-- Name: id; Type: DEFAULT; Schema: askbitcoin; Owner: projectuser
--

ALTER TABLE ONLY account ALTER COLUMN id SET DEFAULT nextval('account_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: askbitcoin; Owner: projectuser
--

ALTER TABLE ONLY entry ALTER COLUMN id SET DEFAULT nextval('entry_id_seq'::regclass);


--
-- Name: PRIMARY; Type: CONSTRAINT; Schema: askbitcoin; Owner: projectuser; Tablespace: 
--

ALTER TABLE ONLY entry_closures
    ADD CONSTRAINT "PRIMARY" PRIMARY KEY (ancestor, descendant);


--
-- Name: account_pkey; Type: CONSTRAINT; Schema: askbitcoin; Owner: projectuser; Tablespace: 
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id);


--
-- Name: entry_pkey; Type: CONSTRAINT; Schema: askbitcoin; Owner: projectuser; Tablespace: 
--

ALTER TABLE ONLY entry
    ADD CONSTRAINT entry_pkey PRIMARY KEY (id);


--
-- Name: entry_author_id_fkey; Type: FK CONSTRAINT; Schema: askbitcoin; Owner: projectuser
--

ALTER TABLE ONLY entry
    ADD CONSTRAINT entry_author_id_fkey FOREIGN KEY (author_id) REFERENCES account(id);


--
-- Name: entry_closures_ancestor_fkey; Type: FK CONSTRAINT; Schema: askbitcoin; Owner: projectuser
--

ALTER TABLE ONLY entry_closures
    ADD CONSTRAINT entry_closures_ancestor_fkey FOREIGN KEY (ancestor) REFERENCES entry(id);


--
-- Name: entry_closures_descendant_fkey; Type: FK CONSTRAINT; Schema: askbitcoin; Owner: projectuser
--

ALTER TABLE ONLY entry_closures
    ADD CONSTRAINT entry_closures_descendant_fkey FOREIGN KEY (descendant) REFERENCES entry(id);


--
-- Name: askbitcoin; Type: ACL; Schema: -; Owner: projectuser
--

REVOKE ALL ON SCHEMA askbitcoin FROM PUBLIC;
REVOKE ALL ON SCHEMA askbitcoin FROM projectuser;
GRANT ALL ON SCHEMA askbitcoin TO projectuser;


--
-- PostgreSQL database dump complete
--

