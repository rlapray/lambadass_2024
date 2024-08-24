--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Drop databases (except postgres and template1)
--

DROP DATABASE pguser;




--
-- Drop roles
--

DROP ROLE pguser;


--
-- Roles
--

CREATE ROLE pguser;
ALTER ROLE pguser WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:oiBm7rSr1t4WKTmmI8KEQQ==$SVJCf6zeUdw/HC41tQszLjk9FnwbhysuMSm2JvJGvoc=:WgPuHaQhBThrv8oUh+j0T19oeZdCxOYmIubO1AShKE0=';

--
-- User Configurations
--








--
-- Databases
--

--
-- Database "template1" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2 (Debian 16.2-1.pgdg120+2)
-- Dumped by pg_dump version 16.2 (Debian 16.2-1.pgdg120+2)

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

UPDATE pg_catalog.pg_database SET datistemplate = false WHERE datname = 'template1';
DROP DATABASE template1;
--
-- Name: template1; Type: DATABASE; Schema: -; Owner: pguser
--

CREATE DATABASE template1 WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE template1 OWNER TO pguser;

\connect template1

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
-- Name: DATABASE template1; Type: COMMENT; Schema: -; Owner: pguser
--

COMMENT ON DATABASE template1 IS 'default template for new databases';


--
-- Name: template1; Type: DATABASE PROPERTIES; Schema: -; Owner: pguser
--

ALTER DATABASE template1 IS_TEMPLATE = true;


\connect template1

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
-- Name: DATABASE template1; Type: ACL; Schema: -; Owner: pguser
--

REVOKE CONNECT,TEMPORARY ON DATABASE template1 FROM PUBLIC;
GRANT CONNECT ON DATABASE template1 TO PUBLIC;


--
-- PostgreSQL database dump complete
--

--
-- Database "pguser" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2 (Debian 16.2-1.pgdg120+2)
-- Dumped by pg_dump version 16.2 (Debian 16.2-1.pgdg120+2)

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
-- Name: pguser; Type: DATABASE; Schema: -; Owner: pguser
--

CREATE DATABASE pguser WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE pguser OWNER TO pguser;

\connect pguser

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: pet; Type: TABLE; Schema: public; Owner: pguser
--

CREATE TABLE public.pet (
    id uuid NOT NULL,
    name text,
    race_id uuid
);


ALTER TABLE public.pet OWNER TO pguser;

--
-- Name: race; Type: TABLE; Schema: public; Owner: pguser
--

CREATE TABLE public.race (
    id uuid NOT NULL,
    name text
);


ALTER TABLE public.race OWNER TO pguser;

--
-- Name: test; Type: TABLE; Schema: public; Owner: pguser
--

CREATE TABLE public.test (
    name text
);


ALTER TABLE public.test OWNER TO pguser;

--
-- Data for Name: pet; Type: TABLE DATA; Schema: public; Owner: pguser
--

COPY public.pet (id, name, race_id) FROM stdin;
752cd664-4267-493e-b831-1d45879bf5a9	Pet 1	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d45879bf5a0	Pet 3	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d4587abf5b0	Pet 5	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d4587abf5b3	bite	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d45879bf5b0	Pet 4	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d4587abf5b4	aaf bleuh	018ff77d-2dd8-7d82-beb4-54fa45d88878
018f83e5-79c2-74e3-8cb0-0dc11957f43a	Pet 5	018ff77d-2dd8-7d82-beb4-54fa45d88878
018f83e6-a405-75b0-8c7f-c1e0c22e7709	Pet 6	018ff77d-2dd8-7d82-beb4-54fa45d88878
018f83e8-114f-7468-8b1e-3b2fc0b3fa77	Pet 7	018ff77d-2dd8-7d82-beb4-54fa45d88878
018fdd92-0316-711a-badc-d338e85b2626	bleuh	018ff77d-2dd8-7d82-beb4-54fa45d88878
018fde9e-76d9-7375-8fad-f76c20fb2a6e	bleuh	018ff77d-2dd8-7d82-beb4-54fa45d88878
018fde9e-85a6-740a-a8d7-0f8eb627f200	bleuh	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d4587abf5b1	Pet 16	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d4587abf5b2	Pet 17	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d4587abf5b5	bleuh	018ff77d-2dd8-7d82-beb4-54fa45d88878
752cd664-4267-493e-b831-1d4587abf5b6	bleuh	018ff77d-2dd8-7d82-beb4-54fa45d88878
\.


--
-- Data for Name: race; Type: TABLE DATA; Schema: public; Owner: pguser
--

COPY public.race (id, name) FROM stdin;
018ff77d-2dd8-7d82-beb4-54fa45d88878	Cat
\.


--
-- Data for Name: test; Type: TABLE DATA; Schema: public; Owner: pguser
--

COPY public.test (name) FROM stdin;
2024-05-09 21:47:32
2024-05-09 21:47:33
2024-05-09 21:47:34
2024-05-09 21:47:34
2024-05-09 21:47:34
\.


--
-- Name: pet pet_pkey; Type: CONSTRAINT; Schema: public; Owner: pguser
--

ALTER TABLE ONLY public.pet
    ADD CONSTRAINT pet_pkey PRIMARY KEY (id);


--
-- Name: race race_pkey; Type: CONSTRAINT; Schema: public; Owner: pguser
--

ALTER TABLE ONLY public.race
    ADD CONSTRAINT race_pkey PRIMARY KEY (id);


--
-- Name: pet fk_race; Type: FK CONSTRAINT; Schema: public; Owner: pguser
--

ALTER TABLE ONLY public.pet
    ADD CONSTRAINT fk_race FOREIGN KEY (race_id) REFERENCES public.race(id);


--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2 (Debian 16.2-1.pgdg120+2)
-- Dumped by pg_dump version 16.2 (Debian 16.2-1.pgdg120+2)

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

DROP DATABASE postgres;
--
-- Name: postgres; Type: DATABASE; Schema: -; Owner: pguser
--

CREATE DATABASE postgres WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE postgres OWNER TO pguser;

\connect postgres

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
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: pguser
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

