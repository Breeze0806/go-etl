CREATE DATABASE source DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
CREATE DATABASE destination DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
CREATE TABLE source.split (
	id bigint(20),
	dt date,
	str varchar(10)
);

CREATE TABLE destination.split (
	id bigint(20),
	dt date,
	str varchar(10)
);