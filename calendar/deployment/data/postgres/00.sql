CREATE TABLE public."event" (
	id serial NOT NULL,
	user_id int4 NULL,
	title varchar(100) NULL,
	description text NULL,
	"time" timestamptz NULL,
	duration int8 NULL,
	notice_period int8 NULL,
	notice_time timestamptz NULL,
	CONSTRAINT event_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_event_time ON public.event USING btree ("time");
CREATE INDEX idx_event_notice_time ON public.event USING btree (notice_time);
CREATE INDEX idx_event_user_id ON public.event USING btree (user_id);

-- Permissions

ALTER TABLE public."event" OWNER TO postgres;
GRANT ALL ON TABLE public."event" TO postgres;





