-- postgis for geospatial queries
CREATE EXTENSION IF NOT EXISTS postgis;

-- users
CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email VARCHAR(255) UNIQUE NOT NULL,
	username VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	role VARCHAR(100) NOT NULL DEFAULT 'user',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- joints
CREATE TABLE IF NOT EXISTS joints(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	-- location data stored as raw coordinates and point in postgis
	latitude DECIMAL(9, 6) NOT NULL,
	longitude DECIMAL(9, 6) NOT NULL,
	location GEOGRAPHY(POINT) NOT NULL,
	description VARCHAR(5000),
	is_approved BOOLEAN NOT NULL DEFAULT false,
	creator_id UUID NOT NULL REFERENCES users(id),
	photo_url VARCHAR(255),
	upvotes INTEGER NOT NULL DEFAULT 0 CHECK(upvotes >= 0),
	downvotes INTEGER NOT NULL DEFAULT 0 CHECK(downvotes >= 0),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- index for joint location
CREATE INDEX IF NOT EXISTS idx_joints_location ON joints USING GIST(location);

-- votes
CREATE TABLE IF NOT EXISTS votes(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id),
	joint_id UUID NOT NULL REFERENCES joints(id),
	direction VARCHAR(10) NOT NULL CHECK(direction IN ('up', 'down')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

	-- only one user vote per joint is allowed
	UNIQUE(user_id, joint_id)
);

-- complaints
CREATE TABLE IF NOT EXISTS complaints(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	joint_id UUID NOT NULL REFERENCES joints(id),
	user_id UUID NOT NULL REFERENCES users(id),
	reason VARCHAR(255) NOT NULL,
	status VARCHAR(50) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

	-- only one user vote per joint is allowed
	UNIQUE(user_id, joint_id)
);
