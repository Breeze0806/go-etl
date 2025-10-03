DROP SCHEMA IF EXISTS source CASCADE;
DROP SCHEMA IF EXISTS destination CASCADE;

CREATE SCHEMA source;
CREATE SCHEMA destination;

SET SCHEMA source;

DROP TABLE IF EXISTS type_table;

CREATE TABLE type_table
(
    t_primary         INT AUTO_INCREMENT NOT NULL,
    t_bit             BIT              NULL,
    t_tinyint         TINYINT          NULL,
    t_smallint        SMALLINT         NULL,
    t_NUMBER       NUMBER              NULL,
    t_BYTE             BYTE              NULL,
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
    t_longvarbinary   longvarbinary(255)   NULL,
    t_blob            BLOB             NULL,
    t_CLOB            CLOB             NULL,
    t_image           IMAGE             NULL,
    t_date            DATE             NULL,
    t_time            TIME             NULL,
    t_datetime        DATETIME         NULL,
    t_timestamp       TIMESTAMP        NULL ,
    t_中文列          LONGVARCHAR(256)     NULL,
    t_dec DEC(20)      NULL,
    PRIMARY KEY (t_primary)
);

INSERT INTO type_table (
    t_tinyint, t_smallint, t_NUMBER, t_BYTE, t_bigint, t_bit,
    t_float, t_double, t_decimal,
    t_char, t_varchar, t_tinytext, t_text,
    t_date, t_time, t_datetime, t_timestamp,
    t_中文列, t_dec
) VALUES (
             1, 1, 1, 1, 0, 0,
             0.0, 0.0, 0.00,
             'char', 'test_guid', 'tinytext', 'text',
             DATE '1900-01-01',
             TIME '02:00:00',
             SYSDATE, CURRENT_TIMESTAMP,
             '中文列', 9999999999999999999
         );

UPDATE type_table SET
                      t_binary = HEXTORAW('1234567890ABCDEF000000000000000000000000000000000000000000000000000000000000000000000000000000000000'),
                      t_varbinary = HEXTORAW('1234567890ABCDEF'),
                      t_longvarbinary = HEXTORAW('1234567890ABCDEF'),
                      t_blob = HEXTORAW('1234567890ABCDEF'),
                      t_CLOB = HEXTORAW('1234567890ABCDEF'),
                      t_image = HEXTORAW('121345673FBCDE')

WHERE t_primary = (SELECT MAX(t_primary) FROM type_table);
INSERT INTO type_table(t_tinyint) VALUES (NULL);
SET SCHEMA destination;
DROP TABLE IF EXISTS type_table;
CREATE TABLE type_table
(
    t_primary         INT AUTO_INCREMENT NOT NULL,
    t_bit             BIT              NULL,
    t_tinyint         TINYINT          NULL,
    t_smallint        SMALLINT         NULL,
    t_NUMBER       NUMBER               NULL,
    t_BYTE             BYTE              NULL,
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
    t_longvarbinary   longvarbinary(255)   NULL,
    t_blob            BLOB             NULL,
    t_CLOB            CLOB             NULL,
    t_image           IMAGE             NULL,
    t_date            DATE             NULL,
    t_time            TIME             NULL,
    t_datetime        DATETIME         NULL,
    t_timestamp       TIMESTAMP        NULL,
    t_中文列          LONGVARCHAR(256)     NULL,
    t_dec DEC(20)      NULL,
    PRIMARY KEY (t_primary)
);