-- =============================================
-- RBAC权限管理系统表结构
-- =============================================

-- 1. 用户表 (sys_user)
CREATE TABLE IF NOT EXISTS `sys_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像URL',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `gender` tinyint DEFAULT '0' COMMENT '性别 0:未知 1:男 2:女',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `status` tinyint DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `is_admin` tinyint DEFAULT '0' COMMENT '是否超级管理员 0:否 1:是',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint DEFAULT NULL COMMENT '更新人ID',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `deleted` tinyint DEFAULT '0' COMMENT '删除标记 0:正常 1:删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_dept_id` (`dept_id`),
  KEY `idx_status` (`status`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统用户表';

-- 2. 角色表 (sys_role)
CREATE TABLE IF NOT EXISTS `sys_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `role_code` varchar(50) NOT NULL COMMENT '角色编码',
  `role_name` varchar(50) NOT NULL COMMENT '角色名称',
  `role_desc` varchar(200) DEFAULT NULL COMMENT '角色描述',
  `data_scope` tinyint DEFAULT '1' COMMENT '数据权限范围 1:全部数据 2:本部门及以下 3:本部门 4:仅本人',
  `status` tinyint DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `sort` int DEFAULT '0' COMMENT '排序',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint DEFAULT NULL COMMENT '更新人ID',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `deleted` tinyint DEFAULT '0' COMMENT '删除标记 0:正常 1:删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_code` (`role_code`),
  KEY `idx_status` (`status`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统角色表';

-- 3. 菜单表 (sys_menu)
CREATE TABLE IF NOT EXISTS `sys_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父菜单ID 0:顶级菜单',
  `menu_type` tinyint NOT NULL COMMENT '菜单类型 1:目录 2:菜单 3:按钮',
  `menu_name` varchar(50) NOT NULL COMMENT '菜单名称',
  `menu_code` varchar(100) NOT NULL COMMENT '菜单编码（权限标识）',
  `path` varchar(200) DEFAULT NULL COMMENT '路由路径',
  `component` varchar(200) DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(100) DEFAULT NULL COMMENT '图标',
  `sort` int DEFAULT '0' COMMENT '排序',
  `is_frame` tinyint DEFAULT '0' COMMENT '是否外链 0:否 1:是',
  `is_cache` tinyint DEFAULT '0' COMMENT '是否缓存 0:否 1:是',
  `is_visible` tinyint DEFAULT '1' COMMENT '是否显示 0:隐藏 1:显示',
  `status` tinyint DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `perms` varchar(500) DEFAULT NULL COMMENT '权限标识（多个用逗号分隔）',
  `query_params` varchar(500) DEFAULT NULL COMMENT '路由参数',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint DEFAULT NULL COMMENT '更新人ID',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `deleted` tinyint DEFAULT '0' COMMENT '删除标记 0:正常 1:删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_menu_code` (`menu_code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_menu_type` (`menu_type`),
  KEY `idx_sort` (`sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统菜单表';

-- 4. 部门表 (sys_dept)
CREATE TABLE IF NOT EXISTS `sys_dept` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '部门ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父部门ID 0:顶级部门',
  `dept_code` varchar(50) NOT NULL COMMENT '部门编码',
  `dept_name` varchar(100) NOT NULL COMMENT '部门名称',
  `leader_id` bigint DEFAULT NULL COMMENT '负责人ID',
  `leader_name` varchar(50) DEFAULT NULL COMMENT '负责人姓名',
  `phone` varchar(20) DEFAULT NULL COMMENT '联系电话',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `sort` int DEFAULT '0' COMMENT '排序',
  `status` tinyint DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint DEFAULT NULL COMMENT '更新人ID',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `deleted` tinyint DEFAULT '0' COMMENT '删除标记 0:正常 1:删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dept_code` (`dept_code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='部门表';

-- 5. 用户角色关联表 (sys_user_role)
CREATE TABLE IF NOT EXISTS `sys_user_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户角色关联表';

-- 6. 角色菜单关联表 (sys_role_menu)
CREATE TABLE IF NOT EXISTS `sys_role_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`,`menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色菜单关联表';

-- 7. 字典类型表 (sys_dict_type)
CREATE TABLE IF NOT EXISTS `sys_dict_type` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典主键',
  `dict_type` varchar(100) NOT NULL COMMENT '字典类型',
  `dict_name` varchar(100) NOT NULL COMMENT '字典名称',
  `dict_source` tinyint DEFAULT '1' COMMENT '数据来源 1:系统内置 2:用户自定义',
  `dict_scope` tinyint DEFAULT '1' COMMENT '作用范围 1:全局 2:租户',
  `is_system` tinyint DEFAULT '0' COMMENT '是否系统字典 0:否 1:是',
  `status` tinyint DEFAULT '0' COMMENT '状态 0:正常 1:停用',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint DEFAULT NULL COMMENT '更新人ID',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `deleted` tinyint DEFAULT '0' COMMENT '删除标记 0:正常 1:删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dict_type` (`dict_type`),
  KEY `idx_status` (`status`),
  KEY `idx_is_system` (`is_system`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='字典类型表';

-- 8. 字典数据表 (sys_dict_data)
CREATE TABLE IF NOT EXISTS `sys_dict_data` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典编码',
  `dict_type` varchar(100) NOT NULL COMMENT '字典类型',
  `dict_label` varchar(100) NOT NULL COMMENT '字典标签',
  `dict_value` varchar(100) NOT NULL COMMENT '字典键值',
  `dict_sort` int DEFAULT '0' COMMENT '字典排序',
  `parent_value` varchar(100) DEFAULT NULL COMMENT '父级值（支持树形字典）',
  `css_class` varchar(100) DEFAULT NULL COMMENT '样式属性',
  `list_class` varchar(100) DEFAULT NULL COMMENT '表格回显样式',
  `is_default` tinyint DEFAULT '0' COMMENT '是否默认 0:否 1:是',
  `status` tinyint DEFAULT '0' COMMENT '状态 0:正常 1:停用',
  `extend_json` json DEFAULT NULL COMMENT '扩展字段',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint DEFAULT NULL COMMENT '更新人ID',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `deleted` tinyint DEFAULT '0' COMMENT '删除标记 0:正常 1:删除',
  PRIMARY KEY (`id`),
  KEY `idx_dict_type` (`dict_type`),
  KEY `idx_dict_sort` (`dict_sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='字典数据表';

-- 9. 操作日志表 (sys_operation_log)
CREATE TABLE IF NOT EXISTS `sys_operation_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `title` varchar(100) DEFAULT NULL COMMENT '操作模块',
  `business_type` tinyint DEFAULT NULL COMMENT '业务类型 1:新增 2:修改 3:删除 4:授权 5:导出 6:导入 7:强退 8:清空数据',
  `method` varchar(100) DEFAULT NULL COMMENT '方法名称',
  `request_method` varchar(10) DEFAULT NULL COMMENT '请求方式 GET/POST/PUT/DELETE',
  `user_id` bigint DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '操作用户名',
  `dept_name` varchar(100) DEFAULT NULL COMMENT '部门名称',
  `operation_ip` varchar(50) DEFAULT NULL COMMENT '操作IP',
  `operation_location` varchar(255) DEFAULT NULL COMMENT '操作地点',
  `request_url` varchar(500) DEFAULT NULL COMMENT '请求URL',
  `request_params` text COMMENT '请求参数',
  `response_result` text COMMENT '返回结果',
  `status` tinyint DEFAULT '0' COMMENT '操作状态 0:成功 1:失败',
  `error_msg` text COMMENT '错误信息',
  `operation_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  `cost_time` bigint DEFAULT NULL COMMENT '耗时（毫秒）',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_operation_time` (`operation_time`),
  KEY `idx_status` (`status`),
  KEY `idx_business_type` (`business_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='操作日志表';

-- 10. 登录日志表 (sys_login_log)
CREATE TABLE IF NOT EXISTS `sys_login_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `username` varchar(50) DEFAULT NULL COMMENT '用户名',
  `login_ip` varchar(50) DEFAULT NULL COMMENT '登录IP',
  `login_location` varchar(255) DEFAULT NULL COMMENT '登录地点',
  `browser` varchar(100) DEFAULT NULL COMMENT '浏览器类型',
  `os` varchar(100) DEFAULT NULL COMMENT '操作系统',
  `login_status` tinyint DEFAULT '0' COMMENT '登录状态 0:成功 1:失败',
  `login_message` varchar(500) DEFAULT NULL COMMENT '登录消息',
  `login_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`),
  KEY `idx_login_time` (`login_time`),
  KEY `idx_login_status` (`login_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='登录日志表';

-- 11. 岗位表 (sys_post)
CREATE TABLE IF NOT EXISTS `sys_post` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
  `post_code` varchar(50) NOT NULL COMMENT '岗位编码',
  `post_name` varchar(50) NOT NULL COMMENT '岗位名称',
  `post_sort` int DEFAULT '0' COMMENT '显示顺序',
  `status` tinyint DEFAULT '1' COMMENT '状态 0:停用 1:正常',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_by` bigint DEFAULT NULL COMMENT '更新人ID',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `deleted` tinyint DEFAULT '0' COMMENT '删除标记 0:正常 1:删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_code` (`post_code`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='岗位信息表';

-- 12. 用户岗位关联表 (sys_user_post)
CREATE TABLE IF NOT EXISTS `sys_user_post` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `post_id` bigint NOT NULL COMMENT '岗位ID',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_post` (`user_id`,`post_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户岗位关联表';

-- =============================================
-- 初始化数据
-- =============================================

-- 初始化管理员账号 (密码: admin123)
INSERT INTO `sys_user` (`id`, `username`, `password`, `email`, `real_name`, `nickname`, `status`, `is_admin`, `create_time`)
VALUES (1, 'admin', '$2a$10$bRmmJQlWfbPO75ySdMM9fewSGex3.aNQz6jy3zE6S1Y8QI.S9PYe6', '2641070015@qq.com', '超级管理员', 'admin', 1, 1, NOW())
ON DUPLICATE KEY UPDATE `id`=`id`;

-- 初始化管理员角色
INSERT INTO `sys_role` (`id`, `role_code`, `role_name`, `role_desc`, `data_scope`, `status`, `sort`, `create_time`)
VALUES (1, 'admin', '超级管理员', '拥有所有权限', 1, 1, 1, NOW())
ON DUPLICATE KEY UPDATE `id`=`id`;

-- 初始化管理员用户角色关联
INSERT INTO `sys_user_role` (`user_id`, `role_id`, `create_time`)
VALUES (1, 1, NOW())
ON DUPLICATE KEY UPDATE `user_id`=`user_id`;
