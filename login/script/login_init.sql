DELIMITER //

DROP PROCEDURE IF EXISTS DROP_TABLES_LIKE //
CREATE PROCEDURE DROP_TABLES_LIKE(pattern VARCHAR(255))
BEGIN
    SELECT @str_sql := CONCAT('DROP TABLE IF EXISTS ', GROUP_CONCAT(`table_name`))
    FROM `information_schema`.`tables`
    WHERE `table_schema` = database() AND `table_name` LIKE pattern;
    IF @str_sql IS NOT NULL THEN
        PREPARE stmt FROM @str_sql; EXECUTE stmt; DROP PREPARE stmt;
    END IF;
END //

DROP PROCEDURE IF EXISTS ENSURE_SQL_TABLE //
CREATE PROCEDURE ENSURE_SQL_TABLE(pattern VARCHAR(255))
BEGIN
    SET @currdate = CURRENT_DATE();
    SET @nextdate = DATE_ADD(@currdate, INTERVAL 1 MONTH);
    SET @basetable = CONCAT('t', SUBSTRING(pattern, 2));
    SET @currtable = CONCAT(@basetable, '_', DATE_FORMAT(@currdate,'%Y%m'));
    SET @nexttable = CONCAT(@basetable, '_', DATE_FORMAT(@nextdate,'%Y%m'));
    SET @currsql = CONCAT('CREATE TABLE IF NOT EXISTS ', @currtable, ' LIKE ', pattern);
    SET @nextsql = CONCAT('CREATE TABLE IF NOT EXISTS ', @nexttable, ' LIKE ', pattern);
    PREPARE currstmt FROM @currsql; EXECUTE currstmt; DROP PREPARE currstmt;
    PREPARE nextstmt FROM @nextsql; EXECUTE nextstmt; DROP PREPARE nextstmt;
END //

DROP PROCEDURE IF EXISTS ALTER_TABLE_LIKE //
CREATE PROCEDURE ALTER_TABLE_LIKE(pattern VARCHAR(255), command VARCHAR(255))
BEGIN
    SET @currdate = CURRENT_DATE();
    SET @nextdate = DATE_ADD(@currdate, INTERVAL 1 MONTH);
    SET @basetable = CONCAT('t', SUBSTRING(pattern, 2));
    SET @currtable = CONCAT(@basetable, '_', DATE_FORMAT(@currdate,'%Y%m'));
    SET @nexttable = CONCAT(@basetable, '_', DATE_FORMAT(@nextdate,'%Y%m'));
    SET @currsql = CONCAT('CREATE TABLE IF NOT EXISTS ', @currtable, ' LIKE ', pattern);
    SET @nextsql = CONCAT('CREATE TABLE IF NOT EXISTS ', @nexttable, ' LIKE ', pattern);
    PREPARE currstmt FROM @currsql; EXECUTE currstmt; DROP PREPARE currstmt;
    PREPARE nextstmt FROM @nextsql; EXECUTE nextstmt; DROP PREPARE nextstmt;
    SET @basesql = CONCAT('ALTER TABLE ', pattern, ' ', command);
    SET @currsql = CONCAT('ALTER TABLE ', @currtable, ' ', command);
    SET @nextsql = CONCAT('ALTER TABLE ', @nexttable, ' ', command);
    PREPARE basestmt FROM @basesql; EXECUTE basestmt; DROP PREPARE basestmt;
    PREPARE currstmt FROM @currsql; EXECUTE currstmt; DROP PREPARE currstmt;
    PREPARE nextstmt FROM @nextsql; EXECUTE nextstmt; DROP PREPARE nextstmt;
END //


DROP TABLE IF EXISTS `user` //
CREATE TABLE `user` (
    `AccID` BIGINT(16) UNSIGNED NOT NULL COMMENT '账号ID',
    `Name` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '账号名',
    `SdkType` SMALLINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT '当前渠道类型ID',
    `DevID` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '设备唯一标识',
    PRIMARY KEY (`AccID`),
    UNIQUE KEY `Name`(`Name`),
    INDEX `SdkType` (`SdkType`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='账号注册表' //


DROP PROCEDURE IF EXISTS DROP_TABLES_LIKE //
DROP PROCEDURE IF EXISTS ENSURE_SQL_TABLE //
DROP PROCEDURE IF EXISTS ALTER_TABLE_LIKE //

DELIMITER ;

