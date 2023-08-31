package main

func (env Env) CreateSegmentsTable() error {
	query := `CREATE TABLE IF NOT EXISTS public.segments
	(
		id SERIAL,
		segment text COLLATE pg_catalog."default",
		CONSTRAINT segments_pkey PRIMARY KEY (id)
	)
	
	TABLESPACE pg_default;
	
	ALTER TABLE IF EXISTS public.segments
		OWNER to "avito-user";`

	_, err := env.DB.Exec(query)
	return err
}

func (env Env) CreateUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS public.users
	(
		id integer NOT NULL,
		segments text[] COLLATE pg_catalog."default",
		CONSTRAINT users_pkey PRIMARY KEY (id)
	)
	
	TABLESPACE pg_default;
	
	ALTER TABLE IF EXISTS public.users
		OWNER to "avito-user";`

	_, err := env.DB.Exec(query)
	return err
}