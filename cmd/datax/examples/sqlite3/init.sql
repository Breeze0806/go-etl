drop table if exists "type_table";
create table "type_table" (
    "t_integer" integer,
    "t_real" real,
    "t_text" text
);

insert into "type_table" values (1, 1.01, 123456);

drop table if exists "type_table_copy";
create table "type_table_copy" (
    "t_integer" integer,
    "t_real" real,
    "t_text" text
);
