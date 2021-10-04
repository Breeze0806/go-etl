create schema source;
create schema destination;

create table "source"."type_table"(
    "t_bigserial" bigserial not NULL,
    "t_serial" serial,
    "t_smallserial" smallserial,
    "t_boolean" boolean,
    "t_smallint" smallint,
    "t_integer" integer,
    "t_bigint" bigint,
    "t_real" real,
    "t_double" double precision,
    "t_decimal" decimal(18,2),
    "t_numeric" numeric(18,2),
    "t_varchar" varchar(200),
    "t_char" char(200),
    "t_text" text,
    "t_timestamp" timestamp,
    "t_timestamp_tz" timestamp with time zone,
    "t_date" date,
    "t_time" time,
    "t_time_tz" time with time zone,
    PRIMARY KEY("t_bigserial")
);

INSERT into "source"."type_table"
VALUES (1,1,1,NULL
,NULL,NULL,NULL
,NULL,NULL,NULL
,NULL,NULL,NULL
,NULL,NULL,NULL
,NULL,NULL,NULL);
INSERT into "source"."type_table" VALUES
(9223372036854775807,2147483647,32767,true
,-32768, -2147483648,-9223372036854775808
,12323273.345,2.34567e30,123456789123430.67
,123456789123430.67,'中文12as;','中文12as;'
,'中文12as;','2021-10-31 16:16:16.123','2021-10-31 16:16:16.123'
,'2021-10-31','16:16:16.123','16:16:16.123');

create table "destination"."type_table"(
    "t_bigserial" bigserial not NULL,
    "t_serial" serial,
    "t_smallserial" smallserial,
    "t_boolean" boolean,
    "t_smallint" smallint,
    "t_integer" integer,
    "t_bigint" bigint,
    "t_real" real,
    "t_double" double precision,
    "t_decimal" decimal(18,2),
    "t_numeric" numeric(18,2),
    "t_varchar" varchar(200),
    "t_char" char(200),
    "t_text" text,
    "t_timestamp" timestamp,
    "t_timestamp_tz" timestamp with time zone,
    "t_date" date,
    "t_time" time,
    "t_time_tz" time with time zone,
    PRIMARY KEY("t_bigserial")
);