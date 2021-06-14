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

 Date: 15/06/2021 02:05:47
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for novel
-- ----------------------------
DROP TABLE IF EXISTS `novel`;
CREATE TABLE `novel`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content` longtext CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `url` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `date` datetime(0) NULL DEFAULT NULL,
  `is_delete` int(1) NULL DEFAULT NULL,
  `caption` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 28 CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of novel
-- ----------------------------
INSERT INTO `novel` VALUES (21, 'sss', 'ss4s', '0000-00-00 00:00:00', 1, 'xx');
INSERT INTO `novel` VALUES (22, 'Content', 'Url', '0000-00-00 00:00:00', 1, 'xx');
INSERT INTO `novel` VALUES (23, 'Content', 'Url', '0000-00-00 00:00:00', 1, 'xx');
INSERT INTO `novel` VALUES (24, 'Content', 'Url', '0000-00-00 00:00:00', 1, 'xx');
INSERT INTO `novel` VALUES (25, 'Content', 'Url', '0000-00-00 00:00:00', 1, 'xx');
INSERT INTO `novel` VALUES (26, 'Content', 'Url', '0000-00-00 00:00:00', 0, 'xx');
INSERT INTO `novel` VALUES (27, 'Content', 'Url', '0000-00-00 00:00:00', 0, 'xx');

SET FOREIGN_KEY_CHECKS = 1;
