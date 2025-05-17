# postgres_source
CREATE SCHEMA source;

CREATE TABLE source.split (
	id bigint,
	dt date,
	str varchar(10)
);

# postgres_dest
CREATE SCHEMA dest;

CREATE TABLE dest.split (
	id bigint,
	dt date,
	str varchar(10)
);