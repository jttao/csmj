set names 'utf8';
set character_set_database = 'utf8';
set character_set_server = 'utf8';


USE `game`;

ALTER TABLE t_user ADD state int(11) DEFAULT 0;


