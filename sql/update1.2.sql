set names 'utf8';
set character_set_database = 'utf8';
set character_set_server = 'utf8';

USE `game`;

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
INSERT INTO t_tasks (id,reward,targetNum,content) VALUES(3,2,20,"每日赢得比赛次数任务");

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
