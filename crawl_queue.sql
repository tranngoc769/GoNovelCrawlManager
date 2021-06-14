/*
 Navicat Premium Data Transfer

 Source Server         : LOCAL
 Source Server Type    : MySQL
 Source Server Version : 100414
 Source Host           : localhost:3306
 Source Schema         : novelcrawler

 Target Server Type    : MySQL
 Target Server Version : 100414
 File Encoding         : 65001

 Date: 15/06/2021 02:05:37
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for crawl_queue
-- ----------------------------
DROP TABLE IF EXISTS `crawl_queue`;
CREATE TABLE `crawl_queue`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `date` datetime(0) NULL DEFAULT NULL,
  `is_delete` int(1) NULL DEFAULT NULL,
  `source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of crawl_queue
-- ----------------------------
INSERT INTO `crawl_queue` VALUES (8, 'https://', '2021-06-14 12:12:00', 0, 'wika');
INSERT INTO `crawl_queue` VALUES (9, 'https://instastatistics.com/#!/ngoctrinh89', '2021-06-15 01:58:54', 0, 'Trang web 1');

SET FOREIGN_KEY_CHECKS = 1;
