CREATE SCHEMA SOURCE;
CREATE SCHEMA dest;

create Table test.source.mytable (
    t_bit bit,
    t_int8 tinyint,
    t_int16    smallint,
    t_int32    int,
    t_int64    bigint,
    t_float32  float,
    t_float64  real,
    t_numeric  numeric(20,6),
    t_date     date,
    t_datetimeoffset datetimeoffset,
    t_datetime2 datetime2,
    t_smalldatetime smalldatetime,
    t_datetime datetime,
    t_time time,
    t_char char(8),
    t_varchar varchar(100),
    t_text text,
    t_nchar nchar(8),
    t_nvarchar nvarchar(100),
    t_ntext ntext,
    t_binary binary(100),
    t_varbinary varbinary(max)
);

create Table test.dest.mytable (
    t_bit bit,
    t_int8 tinyint,
    t_int16    smallint,
    t_int32    int,
    t_int64    bigint,
    t_float32  float,
    t_float64  real,
    t_numeric  numeric(20,6),
    t_date     date,
    t_datetimeoffset datetimeoffset,
    t_datetime2 datetime2,
    t_smalldatetime smalldatetime,
    t_datetime datetime,
    t_time time,
    t_char char(8),
    t_varchar varchar(100),
    t_text text,
    t_nchar nchar(8),
    t_nvarchar nvarchar(100),
    t_ntext ntext,
    t_binary binary(100),
    t_varbinary varbinary(max)
);

insert into test.source.mytable values
(1,127,32767,2147483647,9223372036854775807,
12345.7892,1234567890232.1334,1234567890232.1334,
'2022-09-10','2022-09-10 21:13:13',
'2022-09-10 21:13:13','2022-09-10 21:13',
'2022-09-10 21:13:13','21:13:13'
,'123','2345','abc123',
N'123',N'中文',N'中文123',
0x1230,0x1230);

insert into test.source.mytable values
(NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,
NULL,NULL,NULL,NULL,NULL,NULL,
NULL,NULL,NULL,NULL,NULL,NULL,
NULL,NULL);
