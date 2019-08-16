/*
 Navicat Premium Data Transfer

 Source Server         : 测试网络
 Source Server Type    : MySQL
 Source Server Version : 50726
 Source Host           : 118.31.121.239:3306
 Source Schema         : tokensky

 Target Server Type    : MySQL
 Target Server Version : 50726
 File Encoding         : 65001

 Date: 01/08/2019 15:08:16
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for rms_backend_user
-- ----------------------------
DROP TABLE IF EXISTS `rms_backend_user`;
CREATE TABLE `rms_backend_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `real_name` varchar(255) NOT NULL DEFAULT '',
  `user_name` varchar(255) NOT NULL DEFAULT '',
  `user_pwd` varchar(255) NOT NULL DEFAULT '',
  `is_super` tinyint(1) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  `mobile` varchar(16) NOT NULL DEFAULT '',
  `email` varchar(256) NOT NULL DEFAULT '',
  `avatar` varchar(256) NOT NULL DEFAULT '',
  `roles_str` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of rms_backend_user
-- ----------------------------
BEGIN;
INSERT INTO `rms_backend_user` VALUES (21, 'root', 'root', '3b22b5ff9d8bdfb216b9b302dc9e5362', 1, 1, '123', '123@123.com', '', '22');
INSERT INTO `rms_backend_user` VALUES (22, 'admin', 'admin', '65e06396c33cd60db4e6ffaaac71522a', 1, 1, '123', '123@123.com', '', '');
COMMIT;

-- ----------------------------
-- Table structure for rms_model_record
-- ----------------------------
DROP TABLE IF EXISTS `rms_model_record`;
CREATE TABLE `rms_model_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `handle` varchar(255) DEFAULT NULL,
  `model` varchar(255) DEFAULT NULL,
  `tbid` varchar(255) DEFAULT NULL,
  `old_data` text,
  `new_data` text,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rms_resource
-- ----------------------------
DROP TABLE IF EXISTS `rms_resource`;
CREATE TABLE `rms_resource` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rtype` int(11) NOT NULL DEFAULT '0',
  `name` varchar(64) NOT NULL DEFAULT '',
  `parent_id` int(11) DEFAULT NULL,
  `seq` int(11) NOT NULL DEFAULT '0',
  `icon` varchar(32) NOT NULL DEFAULT '',
  `url_for` varchar(256) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=186 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of rms_resource
-- ----------------------------
BEGIN;
INSERT INTO `rms_resource` VALUES (96, 1, '身份审核列表', 121, 1, '', 'TokenskyRealAuthController.DataGrid');
INSERT INTO `rms_resource` VALUES (98, 0, '权限管理', NULL, 1, '', '');
INSERT INTO `rms_resource` VALUES (99, 1, '权限管理-角色管理', 98, 1, '', 'AdminRoleController.DataGrid');
INSERT INTO `rms_resource` VALUES (100, 2, '权限管理-角色管理-新增编辑', 99, 1, '', 'AdminRoleController.Edit');
INSERT INTO `rms_resource` VALUES (101, 2, '权限管理-角色管理-删除', 99, 1, '', 'AdminRoleController.Delete');
INSERT INTO `rms_resource` VALUES (102, 2, '权限管理-角色管理-分配资源', 99, 1, '', 'AdminRoleController.Allocate');
INSERT INTO `rms_resource` VALUES (103, 2, '权限管理-角色管理-更新seq', 99, 1, '', 'AdminRoleController.UpdateSeq');
INSERT INTO `rms_resource` VALUES (104, 1, '权限管理-用户管理', 98, 1, '', 'AdminBackendUserController.DataGrid');
INSERT INTO `rms_resource` VALUES (105, 2, '权限管理-用户管理-新增编辑', 99, 1, '', 'AdminBackendUserController.Edit');
INSERT INTO `rms_resource` VALUES (106, 2, '权限管理-角色管理-删除', 99, 1, '', 'AdminBackendUserController.Delete');
INSERT INTO `rms_resource` VALUES (107, 1, '权限管理-资源管理(查询会获取当前所有的权限)', 98, 1, '', 'AdminResourceController.TreeGrid');
INSERT INTO `rms_resource` VALUES (108, 2, '权限管理-资源管理-编辑删除', 107, 1, '', 'AdminResourceController.Edit');
INSERT INTO `rms_resource` VALUES (109, 2, '获取所有权限(当角色已拥有权限choice为true)', 107, 1, '', 'AdminResourceController.TreeGridByRole');
INSERT INTO `rms_resource` VALUES (110, 2, '用户有权管理的菜单列表（包括区域）', 107, 1, '', 'AdminResourceController.UserMenuTree');
INSERT INTO `rms_resource` VALUES (111, 2, 'url校验', 107, 1, '', 'AdminResourceController.CheckUrlFor');
INSERT INTO `rms_resource` VALUES (112, 2, '资源管理-删除', 107, 1, '', 'AdminResourceController.Delete');
INSERT INTO `rms_resource` VALUES (113, 2, '更新资源swq', 107, 1, '', 'AdminRoleController.UpdateSeq');
INSERT INTO `rms_resource` VALUES (114, 2, '权限管理-用户管理-新增编辑', 104, 1, '', 'AdminBackendUserController.Edit');
INSERT INTO `rms_resource` VALUES (115, 2, '权限管理-用户管理-删除', 104, 1, '', 'AdminBackendUserController.Delete');
INSERT INTO `rms_resource` VALUES (116, 2, '获取所有权限(当角色已拥有权限choice为true)', 99, 1, '', 'AdminResourceController.TreeGridByRole');
INSERT INTO `rms_resource` VALUES (117, 2, '身份审核列表-审核', 96, 1, '', 'TokenskyRealAuthController.Auditing');
INSERT INTO `rms_resource` VALUES (121, 0, '用户管理', NULL, 0, '', '');
INSERT INTO `rms_resource` VALUES (122, 1, '用户管理-列表', 121, 1, '', 'TokenskyUserController.DataGrid');
INSERT INTO `rms_resource` VALUES (123, 2, '用户管理-用户列表(资产)', 122, 2, '', 'TokenskyUserBalanceController.GetBalances');
INSERT INTO `rms_resource` VALUES (124, 0, '充提币管理', NULL, 1, '', '');
INSERT INTO `rms_resource` VALUES (125, 1, '提币管理-提币审核列表', 124, 1, '', 'TokenskyUserTibiController.DataGrid');
INSERT INTO `rms_resource` VALUES (126, 2, '提币管理-提币审核列表-审核', 125, 1, '', 'TokenskyUserTibiController.Examine');
INSERT INTO `rms_resource` VALUES (127, 0, '消息管理', NULL, 1, '', '');
INSERT INTO `rms_resource` VALUES (128, 1, '消息管理-消息列表', 127, 1, '', 'TokenskyMessageController.DataGrid');
INSERT INTO `rms_resource` VALUES (129, 2, '消息管理-消息列表-新增/编辑', 128, 1, '', 'TokenskyMessageController.Edit');
INSERT INTO `rms_resource` VALUES (130, 2, '消息管理-消息列表-删除', 128, 2, '', 'TokenskyMessageController.Delete');
INSERT INTO `rms_resource` VALUES (131, 0, '财务管理', NULL, 1, '', '');
INSERT INTO `rms_resource` VALUES (132, 1, '财务管理-交易明细', 131, 2, '', 'TokenskyTransactionRecordController.DataGrid');
INSERT INTO `rms_resource` VALUES (134, 1, '提币管理-提币配置', 124, 2, '', 'TokenskyTibiConfigController.DataGrid');
INSERT INTO `rms_resource` VALUES (135, 2, '提币管理-提币配置-新增/编辑', 134, 1, '', 'TokenskyTibiConfigController.Edit');
INSERT INTO `rms_resource` VALUES (136, 0, '算力合约', NULL, 1, '', '');
INSERT INTO `rms_resource` VALUES (137, 1, '算力合约-算力合约分类列表', 136, 1, '', 'HashrateCategoryController.DataGrid');
INSERT INTO `rms_resource` VALUES (138, 2, '算力合约-算力合约分类-编辑', 137, 2, '', 'HashrateCategoryController.Edit');
INSERT INTO `rms_resource` VALUES (139, 2, '算力合约-算力合约分类-删除', 137, 1, '', 'HashrateCategoryController.Delete');
INSERT INTO `rms_resource` VALUES (140, 1, '算力合约-算力合约', 136, 1, '', 'HashrateTreatyController.DataGrid');
INSERT INTO `rms_resource` VALUES (141, 2, '算力合约-算力合约-新增/编辑', 140, 1, '', 'HashrateTreatyController.Edit');
INSERT INTO `rms_resource` VALUES (142, 2, '算力合约-算力合约-删除', 140, 2, '', 'HashrateTreatyController.Delete');
INSERT INTO `rms_resource` VALUES (143, 2, '算力合约-算力合约-合约上下架', 140, 2, '', 'HashrateTreatyController.Shelves');
INSERT INTO `rms_resource` VALUES (144, 1, '算力合约-订单', 136, 3, '', 'HashrateOrderController.DataGrid');
INSERT INTO `rms_resource` VALUES (145, 0, 'OTC', NULL, 1, '', '');
INSERT INTO `rms_resource` VALUES (147, 2, 'otc-配置', 145, 2, '', 'OtcConfController.GetConf');
INSERT INTO `rms_resource` VALUES (148, 2, 'OTC-配置-编辑', 147, 1, '', 'OtcConfController.Edit');
INSERT INTO `rms_resource` VALUES (149, 1, 'otc-申诉列表', 145, 1, '', 'OtcAppealController.DataGrid');
INSERT INTO `rms_resource` VALUES (150, 2, 'otc-申诉列表-申诉', 149, 1, '', 'OtcAppealController.Examine');
INSERT INTO `rms_resource` VALUES (151, 1, 'otc-订单', 145, 1, '', 'OtcOrderController.DataGrid');
INSERT INTO `rms_resource` VALUES (152, 1, 'otc-委托单', 145, 2, '', 'OtcEntrustOrderController.DataGrid');
INSERT INTO `rms_resource` VALUES (153, 1, '用户管理-黑名单', 121, 2, '', 'RoleBlackListController.DataGrid');
INSERT INTO `rms_resource` VALUES (154, 2, '用户管理-黑名单-新增/编辑', 153, 2, '', 'RoleBlackListController.Edit');
INSERT INTO `rms_resource` VALUES (155, 2, '用户管理-黑名单-删除', 153, 2, '', 'RoleBlackListController.Delete');
INSERT INTO `rms_resource` VALUES (156, 2, '用户管理-用户列表-设置用户等级', 122, 1, '', 'TokenskyUserController.SetLevel');
INSERT INTO `rms_resource` VALUES (157, 1, '算力合约-订单收益', 136, 1, '', 'HashrateOrderProfitController.DataGrid');
INSERT INTO `rms_resource` VALUES (158, 1, '算力合约-算力合约收益发放记录', 136, 2, '', 'HashrateSendBalanceRecordController.DataGrid');
INSERT INTO `rms_resource` VALUES (160, 1, '充币管理-充币记录', 124, 2, '', 'TokenskyUserDepositController.DataGrid');
INSERT INTO `rms_resource` VALUES (163, 2, '算力合约-算力合约收益发放记录-收益发放', 158, 1, '', 'HashrateSendBalanceRecordController.SendBalcnce');
INSERT INTO `rms_resource` VALUES (164, 0, '运营管理', NULL, 0, '', '');
INSERT INTO `rms_resource` VALUES (165, 1, 'Banner配置', 164, 0, '', 'OperationBannerController.DataGrid');
INSERT INTO `rms_resource` VALUES (166, 2, '运营-Banner-编辑', 165, 0, '', 'OperationBannerController.Edit');
INSERT INTO `rms_resource` VALUES (167, 2, '运营-banner-删除', 165, 0, '', 'OperationBannerController.Delete');
INSERT INTO `rms_resource` VALUES (168, 0, '理财', NULL, 0, '', '');
INSERT INTO `rms_resource` VALUES (169, 1, '理财分类表', 168, 1, '', 'FinancialCategoryController.DataGrid');
INSERT INTO `rms_resource` VALUES (170, 2, '理财分类-编辑', 169, 1, '', 'FinancialCategoryController.Edit');
INSERT INTO `rms_resource` VALUES (171, 2, '理财分类-删除', 169, 0, '', 'FinancialCategoryController.Delete');
INSERT INTO `rms_resource` VALUES (172, 1, '理财配置', 168, 0, '', 'FinancialProductController.DataGrid');
INSERT INTO `rms_resource` VALUES (173, 2, '理财配置-编辑', 172, 1, '', 'FinancialProductController.Edit');
INSERT INTO `rms_resource` VALUES (174, 2, '理财配置-删除', 172, 1, '', 'FinancialProductController.Delete');
INSERT INTO `rms_resource` VALUES (175, 2, '理财配置-上下架', 172, 1, '', 'FinancialProductController.TheUpper');
INSERT INTO `rms_resource` VALUES (176, 1, '理财配置-修改记录表', 172, 1, '', 'FinancialProductHistoricalRecordController.DataGrid');
INSERT INTO `rms_resource` VALUES (177, 1, '理财用户订单表', 168, 1, '', 'FinancialOrderController.DataGrid');
INSERT INTO `rms_resource` VALUES (178, 1, '理财用户订单-提币表', 177, 1, '', 'FinancialOrderWithdrawalController.DataGrid');
INSERT INTO `rms_resource` VALUES (179, 1, '理财收益表', 168, 1, '', 'FinancialProfitController.DataGrid');
INSERT INTO `rms_resource` VALUES (180, 0, '借贷', NULL, 1, '', '');
INSERT INTO `rms_resource` VALUES (181, 1, '借贷-配置表', 180, 1, '', 'BorrowConfController.DataGrid');
INSERT INTO `rms_resource` VALUES (182, 2, '借贷-配置表-新增编辑', 181, 1, '', 'BorrowConfController.Edit');
INSERT INTO `rms_resource` VALUES (183, 2, '借贷-配置表-上下架', 181, 1, '', 'BorrowConfController.TheUpper');
INSERT INTO `rms_resource` VALUES (184, 2, '用户管理-用户列表-设置是否推广', 122, 1, '', 'TokenskyUserController.SetInvitation');
INSERT INTO `rms_resource` VALUES (185, 2, '用户管理-用户列表-获取邀请连接', 122, 1, '', 'TokenskyUserController.GetAddr');
COMMIT;

-- ----------------------------
-- Table structure for rms_role
-- ----------------------------
DROP TABLE IF EXISTS `rms_role`;
CREATE TABLE `rms_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `seq` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of rms_role
-- ----------------------------
BEGIN;
INSERT INTO `rms_role` VALUES (41, '管理员', 0);
INSERT INTO `rms_role` VALUES (42, '观众', 0);
COMMIT;

-- ----------------------------
-- Table structure for rms_role_backenduser_rel
-- ----------------------------
DROP TABLE IF EXISTS `rms_role_backenduser_rel`;
CREATE TABLE `rms_role_backenduser_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `admin_role_id` int(11) NOT NULL,
  `admin_backend_user_id` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=146 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of rms_role_backenduser_rel
-- ----------------------------
BEGIN;
INSERT INTO `rms_role_backenduser_rel` VALUES (145, 22, 21, '2019-08-01 15:07:46');
COMMIT;

-- ----------------------------
-- Table structure for rms_role_resource_rel
-- ----------------------------
DROP TABLE IF EXISTS `rms_role_resource_rel`;
CREATE TABLE `rms_role_resource_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created` datetime NOT NULL,
  `admin_role_id` int(11) NOT NULL,
  `admin_resource_id` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1298 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of rms_role_resource_rel
-- ----------------------------
BEGIN;
INSERT INTO `rms_role_resource_rel` VALUES (1239, '2019-06-18 13:17:54', 42, 121);
INSERT INTO `rms_role_resource_rel` VALUES (1240, '2019-06-18 13:17:54', 42, 122);
INSERT INTO `rms_role_resource_rel` VALUES (1241, '2019-06-18 13:17:54', 42, 153);
INSERT INTO `rms_role_resource_rel` VALUES (1243, '2019-06-18 13:17:54', 42, 96);
INSERT INTO `rms_role_resource_rel` VALUES (1244, '2019-06-18 13:17:54', 42, 124);
INSERT INTO `rms_role_resource_rel` VALUES (1245, '2019-06-18 13:17:54', 42, 125);
INSERT INTO `rms_role_resource_rel` VALUES (1246, '2019-06-18 13:17:54', 42, 134);
INSERT INTO `rms_role_resource_rel` VALUES (1247, '2019-06-18 13:17:54', 42, 127);
INSERT INTO `rms_role_resource_rel` VALUES (1248, '2019-06-18 13:17:54', 42, 128);
INSERT INTO `rms_role_resource_rel` VALUES (1249, '2019-06-18 13:17:54', 42, 131);
INSERT INTO `rms_role_resource_rel` VALUES (1250, '2019-06-18 13:17:54', 42, 132);
INSERT INTO `rms_role_resource_rel` VALUES (1251, '2019-06-18 13:17:54', 42, 136);
INSERT INTO `rms_role_resource_rel` VALUES (1252, '2019-06-18 13:17:54', 42, 137);
INSERT INTO `rms_role_resource_rel` VALUES (1253, '2019-06-18 13:17:54', 42, 140);
INSERT INTO `rms_role_resource_rel` VALUES (1254, '2019-06-18 13:17:54', 42, 144);
INSERT INTO `rms_role_resource_rel` VALUES (1255, '2019-06-18 13:17:54', 42, 145);
INSERT INTO `rms_role_resource_rel` VALUES (1256, '2019-06-18 13:17:54', 42, 149);
INSERT INTO `rms_role_resource_rel` VALUES (1257, '2019-06-18 13:17:54', 42, 151);
INSERT INTO `rms_role_resource_rel` VALUES (1258, '2019-06-18 13:17:54', 42, 147);
INSERT INTO `rms_role_resource_rel` VALUES (1259, '2019-06-18 13:17:54', 42, 152);
INSERT INTO `rms_role_resource_rel` VALUES (1277, '2019-06-24 13:30:55', 41, 121);
INSERT INTO `rms_role_resource_rel` VALUES (1278, '2019-06-24 13:30:55', 41, 122);
INSERT INTO `rms_role_resource_rel` VALUES (1279, '2019-06-24 13:30:55', 41, 156);
INSERT INTO `rms_role_resource_rel` VALUES (1280, '2019-06-24 13:30:55', 41, 123);
INSERT INTO `rms_role_resource_rel` VALUES (1281, '2019-06-24 13:30:55', 41, 153);
INSERT INTO `rms_role_resource_rel` VALUES (1282, '2019-06-24 13:30:55', 41, 154);
INSERT INTO `rms_role_resource_rel` VALUES (1283, '2019-06-24 13:30:55', 41, 155);
INSERT INTO `rms_role_resource_rel` VALUES (1286, '2019-06-24 13:30:55', 41, 160);
INSERT INTO `rms_role_resource_rel` VALUES (1288, '2019-06-24 13:30:55', 41, 98);
INSERT INTO `rms_role_resource_rel` VALUES (1289, '2019-06-24 13:30:55', 41, 99);
INSERT INTO `rms_role_resource_rel` VALUES (1290, '2019-06-24 13:30:55', 41, 104);
INSERT INTO `rms_role_resource_rel` VALUES (1291, '2019-06-24 13:30:55', 41, 107);
INSERT INTO `rms_role_resource_rel` VALUES (1292, '2019-06-24 13:30:55', 41, 124);
INSERT INTO `rms_role_resource_rel` VALUES (1293, '2019-06-24 13:30:55', 41, 127);
INSERT INTO `rms_role_resource_rel` VALUES (1294, '2019-06-24 13:30:55', 41, 131);
INSERT INTO `rms_role_resource_rel` VALUES (1295, '2019-06-24 13:30:55', 41, 132);
INSERT INTO `rms_role_resource_rel` VALUES (1296, '2019-06-24 13:30:55', 41, 136);
INSERT INTO `rms_role_resource_rel` VALUES (1297, '2019-06-24 13:30:55', 41, 145);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
