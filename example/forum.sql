--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: gtfo; Type: SCHEMA; Schema: -; Owner: gtfo
--

CREATE SCHEMA gtfo;


ALTER SCHEMA gtfo OWNER TO gtfo;

SET search_path = gtfo, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account; Type: TABLE; Schema: gtfo; Owner: gtfo; Tablespace: 
--

CREATE TABLE account (
    id bigint NOT NULL,
    handle text NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    password text,
    email text NOT NULL
);


ALTER TABLE gtfo.account OWNER TO gtfo;

--
-- Name: account_id_seq; Type: SEQUENCE; Schema: gtfo; Owner: gtfo
--

CREATE SEQUENCE account_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE gtfo.account_id_seq OWNER TO gtfo;

--
-- Name: account_id_seq; Type: SEQUENCE OWNED BY; Schema: gtfo; Owner: gtfo
--

ALTER SEQUENCE account_id_seq OWNED BY account.id;


--
-- Name: entry; Type: TABLE; Schema: gtfo; Owner: gtfo; Tablespace: 
--

CREATE TABLE entry (
    id bigint NOT NULL,
    title text NOT NULL,
    body text NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    author_id bigint NOT NULL,
    forum boolean DEFAULT false NOT NULL,
    url boolean DEFAULT false NOT NULL
);


ALTER TABLE gtfo.entry OWNER TO gtfo;

--
-- Name: entry_closures; Type: TABLE; Schema: gtfo; Owner: gtfo; Tablespace: 
--

CREATE TABLE entry_closures (
    ancestor bigint NOT NULL,
    descendant bigint NOT NULL,
    depth integer NOT NULL
);


ALTER TABLE gtfo.entry_closures OWNER TO gtfo;

--
-- Name: entry_id_seq; Type: SEQUENCE; Schema: gtfo; Owner: gtfo
--

CREATE SEQUENCE entry_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE gtfo.entry_id_seq OWNER TO gtfo;

--
-- Name: entry_id_seq; Type: SEQUENCE OWNED BY; Schema: gtfo; Owner: gtfo
--

ALTER SEQUENCE entry_id_seq OWNED BY entry.id;


--
-- Name: vote; Type: TABLE; Schema: gtfo; Owner: gtfo; Tablespace: 
--

CREATE TABLE vote (
    entry_id integer NOT NULL,
    user_id integer NOT NULL,
    upvote boolean DEFAULT false NOT NULL,
    downvote boolean DEFAULT false NOT NULL,
    created time with time zone DEFAULT now() NOT NULL
);


ALTER TABLE gtfo.vote OWNER TO gtfo;

--
-- Name: id; Type: DEFAULT; Schema: gtfo; Owner: gtfo
--

ALTER TABLE ONLY account ALTER COLUMN id SET DEFAULT nextval('account_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: gtfo; Owner: gtfo
--

ALTER TABLE ONLY entry ALTER COLUMN id SET DEFAULT nextval('entry_id_seq'::regclass);


--
-- Name: PRIMARY; Type: CONSTRAINT; Schema: gtfo; Owner: gtfo; Tablespace: 
--

ALTER TABLE ONLY entry_closures
    ADD CONSTRAINT "PRIMARY" PRIMARY KEY (ancestor, descendant);


--
-- Name: account_pkey; Type: CONSTRAINT; Schema: gtfo; Owner: gtfo; Tablespace: 
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id);


--
-- Name: entry_pkey; Type: CONSTRAINT; Schema: gtfo; Owner: gtfo; Tablespace: 
--

ALTER TABLE ONLY entry
    ADD CONSTRAINT entry_pkey PRIMARY KEY (id);


--
-- Name: unique_email; Type: CONSTRAINT; Schema: gtfo; Owner: gtfo; Tablespace: 
--

ALTER TABLE ONLY account
    ADD CONSTRAINT unique_email UNIQUE (email);


--
-- Name: unique_handle; Type: CONSTRAINT; Schema: gtfo; Owner: gtfo; Tablespace: 
--

ALTER TABLE ONLY account
    ADD CONSTRAINT unique_handle UNIQUE (handle);


--
-- Name: vote_pkey; Type: CONSTRAINT; Schema: gtfo; Owner: gtfo; Tablespace: 
--

ALTER TABLE ONLY vote
    ADD CONSTRAINT vote_pkey PRIMARY KEY (entry_id, user_id);


--
-- Name: entry_closures_depth_idx; Type: INDEX; Schema: gtfo; Owner: gtfo; Tablespace: 
--

CREATE INDEX entry_closures_depth_idx ON entry_closures USING btree (depth);


--
-- Name: entry_author_id_fkey; Type: FK CONSTRAINT; Schema: gtfo; Owner: gtfo
--

ALTER TABLE ONLY entry
    ADD CONSTRAINT entry_author_id_fkey FOREIGN KEY (author_id) REFERENCES account(id);


--
-- Name: entry_closures_ancestor_fkey; Type: FK CONSTRAINT; Schema: gtfo; Owner: gtfo
--

ALTER TABLE ONLY entry_closures
    ADD CONSTRAINT entry_closures_ancestor_fkey FOREIGN KEY (ancestor) REFERENCES entry(id);


--
-- Name: entry_closures_descendant_fkey; Type: FK CONSTRAINT; Schema: gtfo; Owner: gtfo
--

ALTER TABLE ONLY entry_closures
    ADD CONSTRAINT entry_closures_descendant_fkey FOREIGN KEY (descendant) REFERENCES entry(id);


--
-- Name: vote_entry_id_fkey; Type: FK CONSTRAINT; Schema: gtfo; Owner: gtfo
--

ALTER TABLE ONLY vote
    ADD CONSTRAINT vote_entry_id_fkey FOREIGN KEY (entry_id) REFERENCES entry(id);


--
-- Name: vote_user_id_fkey; Type: FK CONSTRAINT; Schema: gtfo; Owner: gtfo
--

ALTER TABLE ONLY vote
    ADD CONSTRAINT vote_user_id_fkey FOREIGN KEY (user_id) REFERENCES account(id);


--
-- Name: gtfo; Type: ACL; Schema: -; Owner: gtfo
--

REVOKE ALL ON SCHEMA gtfo FROM PUBLIC;
REVOKE ALL ON SCHEMA gtfo FROM gtfo;
GRANT ALL ON SCHEMA gtfo TO gtfo;


--
-- Name: account; Type: ACL; Schema: gtfo; Owner: gtfo
--

REVOKE ALL ON TABLE account FROM PUBLIC;
REVOKE ALL ON TABLE account FROM gtfo;
GRANT ALL ON TABLE account TO gtfo;


--
-- Name: entry; Type: ACL; Schema: gtfo; Owner: gtfo
--

REVOKE ALL ON TABLE entry FROM PUBLIC;
REVOKE ALL ON TABLE entry FROM gtfo;
GRANT ALL ON TABLE entry TO gtfo;


--
-- Name: entry_closures; Type: ACL; Schema: gtfo; Owner: gtfo
--

REVOKE ALL ON TABLE entry_closures FROM PUBLIC;
REVOKE ALL ON TABLE entry_closures FROM gtfo;
GRANT ALL ON TABLE entry_closures TO gtfo;


--
-- Name: vote; Type: ACL; Schema: gtfo; Owner: gtfo
--

REVOKE ALL ON TABLE vote FROM PUBLIC;
REVOKE ALL ON TABLE vote FROM gtfo;
GRANT ALL ON TABLE vote TO gtfo;


--
-- PostgreSQL database dump complete
--

