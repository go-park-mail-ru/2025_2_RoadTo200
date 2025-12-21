--
-- PostgreSQL database dump
--

\restrict rIyIofCS2IkLKDQ3CcY6nyDU4JZCGcfwg73KHoketnKstgiqxcLiORySxpVQTe2

-- Dumped from database version 17.6 (Debian 17.6-1.pgdg13+1)
-- Dumped by pg_dump version 17.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: gender_enum; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.gender_enum AS ENUM (
    'male',
    'female',
    'other'
);


ALTER TYPE public.gender_enum OWNER TO admin;

--
-- Name: gender_preference_enum; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.gender_preference_enum AS ENUM (
    'male',
    'female',
    'both'
);


ALTER TYPE public.gender_preference_enum OWNER TO admin;

--
-- Name: interest_theme_enum; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.interest_theme_enum AS ENUM (
    'workout',
    'fun',
    'party',
    'chill',
    'love',
    'relax',
    'yoga',
    'friendship',
    'culture',
    'cinema'
);


ALTER TYPE public.interest_theme_enum OWNER TO admin;

--
-- Name: notification_type_enum; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.notification_type_enum AS ENUM (
    'match',
    'super_like',
    'like',
    'message'
);


ALTER TYPE public.notification_type_enum OWNER TO admin;

--
-- Name: plan_type_enum; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.plan_type_enum AS ENUM (
    'week',
    'month',
    'quarter'
);


ALTER TYPE public.plan_type_enum OWNER TO admin;

--
-- Name: strike_reason_type; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.strike_reason_type AS ENUM (
    'spam',
    'fake_profile',
    'offensive_content',
    'harassment',
    'inappropriate_content',
    'underage',
    'copyright_violation',
    'other'
);


ALTER TYPE public.strike_reason_type OWNER TO admin;

--
-- Name: strike_status_type; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.strike_status_type AS ENUM (
    'pending',
    'approved',
    'rejected',
    'resolved'
);


ALTER TYPE public.strike_status_type OWNER TO admin;

--
-- Name: swipe_type_enum; Type: TYPE; Schema: public; Owner: admin
--

CREATE TYPE public.swipe_type_enum AS ENUM (
    'like',
    'dislike',
    'super_like'
);


ALTER TYPE public.swipe_type_enum OWNER TO admin;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: admin
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO admin;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: interest; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.interest (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    theme public.interest_theme_enum NOT NULL
);


ALTER TABLE public.interest OWNER TO admin;

--
-- Name: match; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.match (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user1_id uuid NOT NULL,
    user2_id uuid NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    matched_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone,
    CONSTRAINT match_no_self_match_check CHECK ((user1_id <> user2_id))
);


ALTER TABLE public.match OWNER TO admin;

--
-- Name: COLUMN match.expires_at; Type: COMMENT; Schema: public; Owner: admin
--

COMMENT ON COLUMN public.match.expires_at IS 'Время истечения 24-часового окна. NULL если кто-то написал сообщение (матч активен навсегда)';


--
-- Name: message; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.message (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    match_id uuid NOT NULL,
    sender_id uuid NOT NULL,
    receiver_id uuid NOT NULL,
    content text NOT NULL,
    is_read boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.message OWNER TO admin;

--
-- Name: notification; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.notification (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    type public.notification_type_enum NOT NULL,
    from_user_id uuid,
    match_id uuid,
    is_read boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.notification OWNER TO admin;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO admin;

--
-- Name: strike; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.strike (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    reporter_id uuid NOT NULL,
    target_user_id uuid NOT NULL,
    type public.strike_reason_type NOT NULL,
    reason text,
    status public.strike_status_type DEFAULT 'pending'::public.strike_status_type NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
    moderator_id uuid,
    moderator_note text,
    CONSTRAINT chk_strikes_dates CHECK ((created_at <= COALESCE(updated_at, now()))),
    CONSTRAINT chk_strikes_self_report CHECK ((reporter_id <> target_user_id)),
    CONSTRAINT strike_note_check CHECK (((length(moderator_note) >= 1) AND (length(moderator_note) <= 250))),
    CONSTRAINT strike_reason_check CHECK (((length(reason) >= 1) AND (length(reason) <= 250)))
);


ALTER TABLE public.strike OWNER TO admin;

--
-- Name: subscription; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.subscription (
    user_id uuid NOT NULL,
    plan_type public.plan_type_enum NOT NULL,
    start_date timestamp with time zone DEFAULT now() NOT NULL,
    end_date timestamp with time zone NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.subscription OWNER TO admin;

--
-- Name: swipe; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.swipe (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    swiper_user_id uuid NOT NULL,
    target_user_id uuid NOT NULL,
    swipe_type public.swipe_type_enum NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT swipe_no_self_swipe_check CHECK ((swiper_user_id <> target_user_id))
);


ALTER TABLE public.swipe OWNER TO admin;

--
-- Name: user; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public."user" (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email text NOT NULL,
    phone text,
    name text NOT NULL,
    password text NOT NULL,
    birth_date date,
    gender public.gender_enum,
    bio text,
    city text,
    artist text,
    quote text,
    is_verified boolean DEFAULT false NOT NULL,
    last_active timestamp with time zone DEFAULT now() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    is_premium boolean DEFAULT false NOT NULL,
    super_likes_count integer DEFAULT 3 NOT NULL,
    CONSTRAINT user_age_check CHECK ((birth_date <= ((now() - '18 years'::interval))::date)),
    CONSTRAINT user_artist_check CHECK (((length(artist) >= 1) AND (length(artist) <= 50))),
    CONSTRAINT user_bio_length_check CHECK ((length(TRIM(BOTH FROM bio)) < 255)),
    CONSTRAINT user_city_check CHECK (((length(city) >= 1) AND (length(city) <= 50))),
    CONSTRAINT user_email_check CHECK ((email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'::text)),
    CONSTRAINT user_name_length_check CHECK (((length(name) >= 1) AND (length(name) <= 50))),
    CONSTRAINT user_password_length_check CHECK (((length(TRIM(BOTH FROM password)) >= 8) AND (length(TRIM(BOTH FROM password)) <= 60))),
    CONSTRAINT user_phone_format_check CHECK (((phone IS NULL) OR (phone ~ '^\+?[0-9\s\-\(\)]{10,20}$'::text))),
    CONSTRAINT user_quote_check CHECK (((length(quote) >= 1) AND (length(quote) <= 50)))
);


ALTER TABLE public."user" OWNER TO admin;

--
-- Name: user_photo; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.user_photo (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    photo_url text NOT NULL,
    display_order smallint DEFAULT 0 NOT NULL,
    is_approved boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.user_photo OWNER TO admin;

--
-- Name: user_preference; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.user_preference (
    user_id uuid NOT NULL,
    show_gender public.gender_preference_enum DEFAULT 'both'::public.gender_preference_enum NOT NULL,
    age_min smallint DEFAULT 18 NOT NULL,
    age_max smallint DEFAULT 99 NOT NULL,
    max_distance integer DEFAULT 100 NOT NULL,
    global_search boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.user_preference OWNER TO admin;

--
-- Name: interest interest_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.interest
    ADD CONSTRAINT interest_pkey PRIMARY KEY (id);


--
-- Name: interest interest_user_id_theme_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.interest
    ADD CONSTRAINT interest_user_id_theme_key UNIQUE (user_id, theme);


--
-- Name: match match_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_pkey PRIMARY KEY (id);


--
-- Name: message message_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_pkey PRIMARY KEY (id);


--
-- Name: notification notification_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: strike strike_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.strike
    ADD CONSTRAINT strike_pkey PRIMARY KEY (id);


--
-- Name: strike strike_reporter_id_target_user_id_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.strike
    ADD CONSTRAINT strike_reporter_id_target_user_id_key UNIQUE (reporter_id, target_user_id);


--
-- Name: subscription subscription_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.subscription
    ADD CONSTRAINT subscription_pkey PRIMARY KEY (user_id);


--
-- Name: swipe swipe_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.swipe
    ADD CONSTRAINT swipe_pkey PRIMARY KEY (id);


--
-- Name: user user_email_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_email_key UNIQUE (email);


--
-- Name: user user_phone_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_phone_key UNIQUE (phone);


--
-- Name: user_photo user_photo_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_photo
    ADD CONSTRAINT user_photo_pkey PRIMARY KEY (id);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: user_preference user_preference_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_preference
    ADD CONSTRAINT user_preference_pkey PRIMARY KEY (user_id);


--
-- Name: idx_match_expires_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_match_expires_at ON public.match USING btree (expires_at) WHERE ((expires_at IS NOT NULL) AND (is_active = true));


--
-- Name: idx_match_user1_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_match_user1_id ON public.match USING btree (user1_id);


--
-- Name: idx_match_user2_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_match_user2_id ON public.match USING btree (user2_id);


--
-- Name: idx_message_created_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_message_created_at ON public.message USING btree (created_at DESC);


--
-- Name: idx_message_match_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_message_match_id ON public.message USING btree (match_id);


--
-- Name: idx_message_match_id_created_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_message_match_id_created_at ON public.message USING btree (match_id, created_at DESC);


--
-- Name: idx_message_receiver_unread; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_message_receiver_unread ON public.message USING btree (receiver_id, is_read) WHERE (is_read = false);


--
-- Name: idx_message_sender_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_message_sender_id ON public.message USING btree (sender_id);


--
-- Name: idx_notification_created_at; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_notification_created_at ON public.notification USING btree (created_at DESC);


--
-- Name: idx_notification_user_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_notification_user_id ON public.notification USING btree (user_id);


--
-- Name: idx_notification_user_unread; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_notification_user_unread ON public.notification USING btree (user_id, is_read) WHERE (is_read = false);


--
-- Name: idx_swipe_swiper_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_swipe_swiper_id ON public.swipe USING btree (swiper_user_id);


--
-- Name: idx_swipe_target_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_swipe_target_id ON public.swipe USING btree (target_user_id);


--
-- Name: idx_user_email; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_user_email ON public."user" USING btree (email);


--
-- Name: idx_user_phone; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_user_phone ON public."user" USING btree (phone);


--
-- Name: idx_user_photo_user_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_user_photo_user_id ON public.user_photo USING btree (user_id);


--
-- Name: user_preference update_user_preference_updated_at; Type: TRIGGER; Schema: public; Owner: admin
--

CREATE TRIGGER update_user_preference_updated_at BEFORE UPDATE ON public.user_preference FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: user update_user_updated_at; Type: TRIGGER; Schema: public; Owner: admin
--

CREATE TRIGGER update_user_updated_at BEFORE UPDATE ON public."user" FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: interest interest_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.interest
    ADD CONSTRAINT interest_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: match match_user1_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_user1_id_fkey FOREIGN KEY (user1_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: match match_user2_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_user2_id_fkey FOREIGN KEY (user2_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: message message_match_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_match_id_fkey FOREIGN KEY (match_id) REFERENCES public.match(id) ON DELETE CASCADE;


--
-- Name: message message_receiver_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_receiver_id_fkey FOREIGN KEY (receiver_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: message message_sender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: notification notification_from_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_from_user_id_fkey FOREIGN KEY (from_user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: notification notification_match_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_match_id_fkey FOREIGN KEY (match_id) REFERENCES public.match(id) ON DELETE CASCADE;


--
-- Name: notification notification_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.notification
    ADD CONSTRAINT notification_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: strike strike_moderator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.strike
    ADD CONSTRAINT strike_moderator_id_fkey FOREIGN KEY (moderator_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: strike strike_reporter_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.strike
    ADD CONSTRAINT strike_reporter_id_fkey FOREIGN KEY (reporter_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: strike strike_target_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.strike
    ADD CONSTRAINT strike_target_user_id_fkey FOREIGN KEY (target_user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: subscription subscription_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.subscription
    ADD CONSTRAINT subscription_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: swipe swipe_swiper_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.swipe
    ADD CONSTRAINT swipe_swiper_user_id_fkey FOREIGN KEY (swiper_user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: swipe swipe_target_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.swipe
    ADD CONSTRAINT swipe_target_user_id_fkey FOREIGN KEY (target_user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: user_photo user_photo_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_photo
    ADD CONSTRAINT user_photo_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- Name: user_preference user_preference_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_preference
    ADD CONSTRAINT user_preference_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict rIyIofCS2IkLKDQ3CcY6nyDU4JZCGcfwg73KHoketnKstgiqxcLiORySxpVQTe2

