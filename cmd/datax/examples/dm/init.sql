-- 删除数据库（如果存在）
DROP SCHEMA IF EXISTS source CASCADE;
DROP SCHEMA IF EXISTS destination CASCADE;

-- 创建数据库
CREATE SCHEMA source;
CREATE SCHEMA destination;

-- 使用source模式
SET SCHEMA source;

-- 删除表（如果存在）
DROP TABLE IF EXISTS type_table;

-- 创建表
CREATE TABLE type_table
(
    t_primary         INT AUTO_INCREMENT NOT NULL,
    t_bit             BIT              NULL,
    t_tinyint         TINYINT          NULL,
    t_smallint        SMALLINT         NULL,
    t_mediumint       INT              NULL,
    t_int             INT              NULL,
    t_bigint          BIGINT           NULL,
    t_float           REAL             NULL,
    t_double          DOUBLE PRECISION NULL,
    t_decimal         DECIMAL(10, 2)   NULL,
    t_char            CHAR(50)         NULL,
    t_varchar         VARCHAR(60)      NULL,
    t_tinytext        VARCHAR(255)     NULL,
    t_text            TEXT             NULL,
    t_binary          BINARY(70)       NULL,
    t_varbinary       VARBINARY(80)    NULL,
    t_tinyblob        VARBINARY(255)   NULL,
    t_blob            BLOB             NULL,
    t_mediumblob      BLOB             NULL,
    t_longblob        BLOB             NULL,
    t_date            DATE             NULL,
    t_time            TIME             NULL,
    t_datetime        DATETIME         NULL,
    t_timestamp       TIMESTAMP        NULL DEFAULT CURRENT_TIMESTAMP,
    t_中文列          VARCHAR(256)     NULL,
    t_unsigned_bigint DECIMAL(20)      NULL,
    PRIMARY KEY (t_primary)
);

-- 插入基本数据（不包含二进制字段）
INSERT INTO type_table (
    t_tinyint, t_smallint, t_mediumint, t_int, t_bigint, t_bit,
    t_float, t_double, t_decimal,
    t_char, t_varchar, t_tinytext, t_text,
    t_date, t_time, t_datetime, t_timestamp,
    t_中文列, t_unsigned_bigint
) VALUES (
             1, 1, 1, 1, 0, 0,
             0.0, 0.0, 0.00,
             'char', 'test_guid', 'tinytext', 'text',
             DATE '1900-01-01',
             TIME '02:00:00',
             SYSDATE, CURRENT_TIMESTAMP,
             '中文列', 9999999999999999999
         );

-- 然后更新二进制字段 - 使用达梦支持的语法
UPDATE type_table SET
                      t_binary = HEXTORAW('1234567890ABCDEF000000000000000000000000000000000000000000000000000000000000000000000000000000000000'),
                      t_varbinary = HEXTORAW('1234567890ABCDEF'),
                      t_tinyblob = HEXTORAW('1234567890ABCDEF'),
                      t_blob = HEXTORAW('1234567890ABCDEF'),
                      t_mediumblob = HEXTORAW('1234567890ABCDEF'),
                      t_longblob = HEXTORAW('121345673FBCDE')
WHERE t_primary = (SELECT MAX(t_primary) FROM type_table);

-- 使用destination模式
SET SCHEMA destination;

-- 创建目标表
CREATE TABLE type_table
(
    t_primary         INT AUTO_INCREMENT NOT NULL,
    t_bit             BIT              NULL,
    t_tinyint         TINYINT          NULL,
    t_smallint        SMALLINT         NULL,
    t_mediumint       INT              NULL,
    t_int             INT              NULL,
    t_bigint          BIGINT           NULL,
    t_float           REAL             NULL,
    t_double          DOUBLE PRECISION NULL,
    t_decimal         DECIMAL(10, 2)   NULL,
    t_char            CHAR(50)         NULL,
    t_varchar         VARCHAR(60)      NULL,
    t_tinytext        VARCHAR(255)     NULL,
    t_text            TEXT             NULL,
    t_binary          BINARY(70)       NULL,
    t_varbinary       VARBINARY(80)    NULL,
    t_tinyblob        VARBINARY(255)   NULL,
    t_blob            BLOB             NULL,
    t_mediumblob      BLOB             NULL,
    t_longblob        BLOB             NULL,
    t_date            DATE             NULL,
    t_time            TIME             NULL,
    t_datetime        DATETIME         NULL,
    t_timestamp       TIMESTAMP        NULL DEFAULT CURRENT_TIMESTAMP,
    t_中文列          VARCHAR(256)     NULL,
    t_unsigned_bigint DECIMAL(20)      NULL,
    PRIMARY KEY (t_primary)
);