drop database if exists gongde;

drop user if exists 'gongde'@'%';
-- 支持emoji：需要mysql数据库参数： character_set_server=utf8mb4
create database gongde default character set utf8mb4 collate utf8mb4_unicode_ci;
use gongde;
create user 'gongde'@'%' identified by 'gongde2022';
grant all privileges on gongde.* to 'gongde'@'%';
flush privileges;