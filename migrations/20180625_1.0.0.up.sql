CREATE TABLE users (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	name TEXT NOT NULL,
	metadata jsonb NOT NULL DEFAULT '{}',
	archived boolean NOT NULL DEFAULT false,
	archived_on timestamp,
	created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE records (
	id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
	description TEXT NOT NULL,
    category TEXT NOT NULL,
    cents INTEGER NOT NULL,
	metadata jsonb NOT NULL DEFAULT '{}',
	archived boolean NOT NULL DEFAULT false,
	archived_on timestamp,
	created_at timestamp NOT NULL DEFAULT NOW()
);