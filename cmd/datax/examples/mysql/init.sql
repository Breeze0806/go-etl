DROP DATABASE if exists source;
CREATE DATABASE source DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
DROP DATABASE if exists destination;
CREATE DATABASE destination DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
use source;

drop table if exists `type_table`;

CREATE TABLE `type_table` (
	`t_primary` INT NOT NULL AUTO_INCREMENT,
	`t_bit` BIT NULL,
	`t_tinyint` TINYINT NULL,
	`t_smallint` SMALLINT NULL,
	`t_mediumint` MEDIUMINT NULL,
	`t_int` INT NULL,
	`t_bigint` BIGINT NULL,
	`t_float` FLOAT NULL,
	`t_double` DOUBLE NULL,
	`t_decimal` DECIMAL(10,2) NULL,
	`t_char` CHAR(50) NULL,
	`t_varchar` VARCHAR(60) NULL,
	`t_tinytext` TINYTEXT NULL,
	`t_text` TEXT NULL,
	`t_mediumtext` MEDIUMTEXT NULL,
	`t_longtext` LONGTEXT NULL,
	`t_binary` BINARY(70) NULL,
	`t_varbinary` VARBINARY(80) NULL,
	`t_tinyblob` TINYBLOB NULL,
	`t_blob` BLOB NULL,
	`t_mediumblob` MEDIUMBLOB NULL,
	`t_longblob` LONGBLOB NULL,
	`t_date` DATE NULL,
	`t_time` TIME NULL,
	`t_datetime` DATETIME NULL,
	`t_timestamp` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `t_中文列` VARCHAR(256) NULL,
  `t_unsigned_bigint`  BIGINT Unsigned NULL,
	PRIMARY KEY (`t_primary`)
)
COMMENT='类型表'
COLLATE='utf8_general_ci'
AUTO_INCREMENT = 1
ENGINE=InnoDB;


INSERT INTO type_table (
    t_tinyint, t_smallint, t_mediumint, t_int, t_bigint, t_bit,
    t_float, t_double, t_decimal,
    t_char, t_varchar, t_tinytext, t_text, t_mediumtext, t_longtext,
    t_binary, t_varbinary, t_tinyblob, t_blob, t_mediumblob, t_longblob,
    t_date, t_time, t_datetime, t_timestamp,
    `t_中文列`,`t_unsigned_bigint`
) VALUES (
    1, 1, 1, 1, 0, b'0', 0, 0, 0.0000,
    'char', uuid(), 'tinytext', 'text', 'mediumtext', 'longtext',
     _binary 0x1234567890ABCDEF000000000000000000000000000000000000000000000000000000000000000000000000000000000000,
     _binary 0x1234567890ABCDEF,
     _binary 0x1234567890ABCDEF,
     _binary 0x1234567890ABCDEF,
     _binary 0x1234567890ABCDEF,
     _binary 0x121345673FBCDE,
     '1900-01-01', '26:00:00', now(), '2016-11-21 14:51:53',
     '中文列',18446744073709551615
);

use destination;
CREATE TABLE `type_table` (
	`t_primary` INT NOT NULL AUTO_INCREMENT,
	`t_bit` BIT NULL,
	`t_tinyint` TINYINT NULL,
	`t_smallint` SMALLINT NULL,
	`t_mediumint` MEDIUMINT NULL,
	`t_int` INT NULL,
	`t_bigint` BIGINT NULL,
	`t_float` FLOAT NULL,
	`t_double` DOUBLE NULL,
	`t_decimal` DECIMAL(10,2) NULL,
	`t_char` CHAR(50) NULL,
	`t_varchar` VARCHAR(60) NULL,
	`t_tinytext` TINYTEXT NULL,
	`t_text` TEXT NULL,
	`t_mediumtext` MEDIUMTEXT NULL,
	`t_longtext` LONGTEXT NULL,
	`t_binary` BINARY(70) NULL,
	`t_varbinary` VARBINARY(80) NULL,
	`t_tinyblob` TINYBLOB NULL,
	`t_blob` BLOB NULL,
	`t_mediumblob` MEDIUMBLOB NULL,
	`t_longblob` LONGBLOB NULL,
	`t_date` DATE NULL,
	`t_time` TIME NULL,
	`t_datetime` DATETIME NULL,
	`t_timestamp` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `t_中文列` VARCHAR(256) NULL,
  `t_unsigned_bigint`  BIGINT Unsigned NULL,
	PRIMARY KEY (`t_primary`)
)
COMMENT='类型表'
COLLATE='utf8_general_ci'
AUTO_INCREMENT = 1
ENGINE=InnoDB;