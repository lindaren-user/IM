CREATE TABLE `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT  '用户ID',
    `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名，唯一',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
    `status` TINYINT DEFAULT 1 COMMENT '状态：0禁用 1启用',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

# TODO: 枚举类型， innodb索引, INDEX(),保留关键字
# seq:即使 created_at 时间精度到毫秒，也不能保证顺序（特别是在并发写入时）。 支持断点续传/拉取历史消息
CREATE TABLE `private_messages` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID，自增主键',
    `sender_id` BIGINT UNSIGNED NOT NULL COMMENT '发送者用户ID',
    `receiver_id` BIGINT UNSIGNED NOT NULL COMMENT '接收者用户ID（单聊对方ID）',
    `content_type` ENUM('text', 'image', 'file', 'system', 'audio') NOT NULL COMMENT '消息内容类型：文本，图片，语音等',
    `content` TEXT NOT NULL COMMENT '消息内容，文本或序列化后的多媒体数据',
    `is_read` ENUM('0', '1') DEFAULT '0' NOT NULL COMMENT '是否已读，0=未读，1=已读',
    `status` ENUM('0', '1', '2') DEFAULT '0' NOT NULL COMMENT '消息状态：0=正常，1=撤回，2=删除',
    `seq` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '会话内消息序号，用于消息排序',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '消息创建时间',
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX (`receiver_id`, `seq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='单聊消息表';

CREATE TABLE `group_messages` (
      `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID，自增主键',
      `group_id` BIGINT UNSIGNED NOT NULL COMMENT '群组ID',
      `sender_id` BIGINT UNSIGNED NOT NULL COMMENT '发送者用户ID',
      `content_type` ENUM('text', 'image', 'file', 'system', 'audio') NOT NULL COMMENT '消息内容类型：文本，图片，语音等',
      `content` TEXT NOT NULL COMMENT '消息内容',
      `status` ENUM('0', '1', '2') DEFAULT '0' NOT NULL COMMENT '消息状态：0=正常，1=撤回，2=删除',
      `seq` BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '群聊内消息序号',
      `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '消息创建时间',
      FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
      FOREIGN KEY (group_id) REFERENCES `groups`(id) ON DELETE CASCADE,
      INDEX (`group_id`, `seq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群聊消息表';

CREATE TABLE `groups` (
      `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '群聊ID，自增主键',
      `name` VARCHAR(100) NOT NULL COMMENT '群名称',
      `avatar` VARCHAR(255) DEFAULT NULL COMMENT '群头像 URL',
      `description` TEXT DEFAULT NULL COMMENT '群简介或公告',
      `owner_id` BIGINT UNSIGNED NOT NULL COMMENT '群主用户ID',
      `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
      `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
      FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
      INDEX (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群聊';

CREATE TABLE `friendships` (
       `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键，自增',
       `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
       `friend_id` BIGINT UNSIGNED NOT NULL COMMENT '好友用户ID',
       `remark` VARCHAR(100) DEFAULT NULL COMMENT '好友备注',
       `status` ENUM('pending', 'accepted', 'blocked') NOT NULL DEFAULT 'accepted' COMMENT '状态：pending=待确认，accepted=已成为好友，blocked=拉黑',
       `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加好友时间',
       FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
       FOREIGN KEY (`friend_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
       UNIQUE KEY `uniq_friend_pair` (`user_id`, `friend_id`),
       INDEX (`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友关系表';

CREATE TABLE `group_members` (
     `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
     `group_id` BIGINT UNSIGNED NOT NULL COMMENT '群聊ID',
     `user_id` BIGINT UNSIGNED NOT NULL COMMENT '成员用户ID',
     `role` ENUM('owner', 'admin', 'member') NOT NULL DEFAULT 'member' COMMENT '成员角色：群主，管理员，普通成员',
     `join_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
     `mute_until` DATETIME DEFAULT NULL COMMENT '禁言截止时间，为 NULL 表示未禁言',
     FOREIGN KEY (`group_id`) REFERENCES `groups`(`id`) ON DELETE CASCADE,
     FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
     UNIQUE KEY `uniq_group_user` (`group_id`, `user_id`),
     INDEX (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='群聊成员表';
