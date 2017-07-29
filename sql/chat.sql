CREATE DATABASE IF NOT EXISTS core;
use core;

CREATE TABLE IF NOT EXISTS `message` (
  `id`        int(16) unsigned NOT NULL AUTO_INCREMENT,
  `messageid` int(16) unsigned NOT NULL,
  `source`    varchar(128)     NOT NULL,
  `target`    varchar(128)     NOT NULL,
  `type`      int(16)          NOT NULL,
  `version`   int(16)          NOT NULL,
  `issend`    bool             NOT NULL DEFAULT FALSE,
  `content`   text             NOT NULL,
  `created`   datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
