set names 'utf8';
set character_set_database = 'utf8';
set character_set_server = 'utf8';

DROP DATABASE IF EXISTS `game`;
CREATE DATABASE `game` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;

USE `game`;

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT "id",
  `weixin` varchar(100) DEFAULT '' COMMENT "微信id",
  `deviceMac` varchar(100) COMMENT "设备码",
  `name` varchar(100) COMMENT "名字",
  `forbid` int(11) DEFAULT 0 COMMENT "禁止",
  `image` varchar(500) COMMENT "图像",
  `sex` int(11) DEFAULT 0 COMMENT "性别",
  `cardNum` bigint(20) DEFAULT 0 COMMENT "房间卡数",
  `lastLoginTime` bigint(20) DEFAULT 0 COMMENT "最后登陆时间",
  `lastLoginIp` varchar(20) COMMENT "最后登陆ip",
  `updateTime` bigint(20) DEFAULT 0 COMMENT "更新时间",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  `deleteTime` bigint(20)  DEFAULT 0 COMMENT "删除时间",
  `state` int(11) DEFAULT 0 COMMENT "状态",
  `location` varchar(100) COMMENT "位置",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for t_room
-- ----------------------------
DROP TABLE IF EXISTS `t_room`;
CREATE TABLE `t_room` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `roomType` bigint(20) NOT NULL COMMENT "房间类型",
  `roomConfig` varchar(5000) NOT NULL COMMENT "房间配置",
  `round` int(11) DEFAULT 0 COMMENT "盘数",
  `cost` int(11) DEFAULT 0 COMMENT "花费",
  `ownerId` bigint(20) COMMENT "房间主人id",
  `updateTime` bigint(20) DEFAULT 0 COMMENT "更新时间",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  `deleteTime` bigint(20)  DEFAULT 0 COMMENT "删除时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8;



-- ----------------------------
-- Table structure for t_news
-- ----------------------------
DROP TABLE IF EXISTS `t_news`;
CREATE TABLE `t_news` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT "id",
  `content` text  COMMENT "内容",
  `updateTime` bigint(20) DEFAULT 0 COMMENT "更新时间",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  `deleteTime` bigint(20)  DEFAULT 0 COMMENT "删除时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO t_news (content) VALUES("长沙麻将");

-- ----------------------------
-- Table structure for t_notice
-- ----------------------------
DROP TABLE IF EXISTS `t_notice`;
CREATE TABLE `t_notice` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT "id",
  `content` text  COMMENT "内容",
  `updateTime` bigint(20) DEFAULT 0 COMMENT "更新时间",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  `deleteTime` bigint(20)  DEFAULT 0 COMMENT "删除时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO t_notice (content) VALUES("跑马灯1");
INSERT INTO t_notice (content) VALUES("跑马灯2");


-- ----------------------------
-- Table structure for t_room_record
-- ----------------------------
DROP TABLE IF EXISTS `t_room_record`;
CREATE TABLE `t_room_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `roomType` bigint(20) NOT NULL COMMENT "房间类型",
  `roomId` bigint(20) NOT NULL COMMENT "房间id",
  `ownerId` bigint(20) COMMENT "房间主人id",
  `player1` bigint(20) COMMENT "玩家1",
  `player2` bigint(20) COMMENT "玩家2",
  `player3` bigint(20) COMMENT "玩家3",
  `player4` bigint(20) COMMENT "玩家4",
  `settle` TEXT COMMENT "结算",
  `updateTime` bigint(20) DEFAULT 0 COMMENT "更新时间",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  `deleteTime` bigint(20)  DEFAULT 0 COMMENT "删除时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for t_round
-- ----------------------------
DROP TABLE IF EXISTS `t_round`;
CREATE TABLE `t_round` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `roomRecordId` bigint(20) NOT NULL COMMENT "房间纪录", 
  `roomType` bigint(20) NOT NULL COMMENT "房间类型",
  `roomId` bigint(20) NOT NULL COMMENT "房间id",
  `round` int(11) COMMENT "当前盘数",
  `totalRound` int(11) COMMENT "总盘数",
  `logs` MEDIUMTEXT COMMENT "日志",
  `config` varchar(5000) COMMENT "配置",
  `settle` TEXT COMMENT "结算",
  `updateTime` bigint(20) DEFAULT 0 COMMENT "更新时间",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  `deleteTime` bigint(20)  DEFAULT 0 COMMENT "删除时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for t_card_record
-- ----------------------------
DROP TABLE IF EXISTS `t_card_record`;
CREATE TABLE `t_card_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `userId` bigint(20) NOT NULL COMMENT "用户id",
  `changeNum` bigint(20) COMMENT "卡数",
  `reason` int(11) COMMENT "改变原因",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for t_tasks
-- ----------------------------
DROP TABLE IF EXISTS `t_tasks`;
CREATE TABLE `t_tasks` (
  `id` int(11) NOT NULL COMMENT "id", 
  `reward` int(11) DEFAULT 0 COMMENT "任务奖励", 
  `targetNum` int(11) DEFAULT 0 COMMENT "任务目标次数", 
  `content` text  COMMENT "描述",
  `updateTime` bigint(20) DEFAULT 0 COMMENT "更新时间",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  `deleteTime` bigint(20)  DEFAULT 0 COMMENT "删除时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO t_tasks (id,reward,targetNum,content) VALUES(1,1,1,"每日分享任务");
INSERT INTO t_tasks (id,reward,targetNum,content) VALUES(2,1,10,"每日游戏场数任务");
INSERT INTO t_tasks (id,reward,targetNum,content) VALUES(3,1,20,"每日赢得比赛次数任务");


-- ----------------------------
-- Table structure for t_user_task
-- ----------------------------
DROP TABLE IF EXISTS `t_user_task`;
CREATE TABLE `t_user_task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `taskId` int(11) DEFAULT 0 COMMENT "任务奖励", 
  `userId` bigint(20) DEFAULT 0 COMMENT "任务奖励", 
  `reward` int(11) DEFAULT 0 COMMENT "任务奖励", 
  `finish` int(11) DEFAULT 0 COMMENT "任务奖励",  
  `targetNum` int(11) DEFAULT 0 COMMENT "任务目标次数",  
  `updateTime` bigint(20) DEFAULT 0 COMMENT "更新时间",
  `createTime` bigint(20) DEFAULT 0 COMMENT "创建时间",
  `deleteTime` bigint(20)  DEFAULT 0 COMMENT "删除时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8;

use game;

select * from t_tasks;


select * from t_user_task;

select * from t_user;

update t_user set cardNum=10 where id!=0;

select * from t_card_record;

