
CREATE TABLE `member_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `sex` tinyint NOT NULL DEFAULT 1,
  `email` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `addr` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `position` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `id_card_num` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `photo` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `introducer_id` int(11)  NOT NULL DEFAULT 0,
  `note` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `add_time` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone` (`phone`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=40000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

 CREATE TABLE `member_data` (
  `user_id` int(11) NOT NULL DEFAULT 0,
  `member_level` tinyint NOT NULL DEFAULT 1,
  `password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `amount` int NOT NULL DEFAULT 0,
  `point` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 0,
  `effect_time` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;