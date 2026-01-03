-- Charset/Collation is aligned with my.cnf (utf8mb4 / utf8mb4_ja_0900_as_cs)

CREATE TABLE IF NOT EXISTS `users` (
  `uid` CHAR(28) NOT NULL COMMENT 'ユーザーID',
  `nickname` VARCHAR(20) NOT NULL COMMENT 'ニックネーム',
  `email` VARCHAR(255) NOT NULL COMMENT 'メールアドレス',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `uk_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_ja_0900_as_cs;

CREATE TABLE IF NOT EXISTS `todo_statuses` (
  `status` CHAR(2) NOT NULL COMMENT 'ステータス',
  PRIMARY KEY (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_ja_0900_as_cs;

CREATE TABLE IF NOT EXISTS `todos` (
  `id` CHAR(36) NOT NULL COMMENT 'TodoID',
  `owner` CHAR(28) NOT NULL COMMENT '所有ユーザー',
  `status` CHAR(2) NOT NULL COMMENT 'ステータス',
  `title` VARCHAR(30) NOT NULL COMMENT 'タイトル',
  `content` TEXT NOT NULL COMMENT '内容',
  `due_datetime` DATETIME NULL COMMENT '期限日時',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`id`),
  KEY `idx_todos_owner` (`owner`),
  KEY `idx_todos_status` (`status`),
  CONSTRAINT `fk_todos_owner` FOREIGN KEY (`owner`) REFERENCES `users` (`uid`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_todos_status` FOREIGN KEY (`status`) REFERENCES `todo_statuses` (`status`)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_ja_0900_as_cs;

CREATE TABLE IF NOT EXISTS `goodlucks` (
  `user` CHAR(28) NOT NULL COMMENT 'ユーザー',
  `todo` CHAR(36) NOT NULL COMMENT 'Todo',
  PRIMARY KEY (`user`, `todo`),
  KEY `idx_goodlucks_todo` (`todo`),
  CONSTRAINT `fk_goodlucks_user` FOREIGN KEY (`user`) REFERENCES `users` (`uid`)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_goodlucks_todo` FOREIGN KEY (`todo`) REFERENCES `todos` (`id`)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_ja_0900_as_cs;


