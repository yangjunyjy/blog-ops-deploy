-- ==================== 堡垒机表结构迁移 ====================

-- 1. 主机组表
CREATE TABLE IF NOT EXISTS `host_groups` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `name` VARCHAR(100) NOT NULL COMMENT '主机组名称',
    `desc` VARCHAR(255) COMMENT '描述',
    `sort` INT DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态(0:禁用,1:启用)',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_status` (`status`),
    KEY `idx_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主机组表';

-- 2. 用户组表
CREATE TABLE IF NOT EXISTS `user_groups` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `name` VARCHAR(100) NOT NULL COMMENT '用户组名称',
    `desc` VARCHAR(255) COMMENT '描述',
    `sort` INT DEFAULT 0 COMMENT '排序',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态(0:禁用,1:启用)',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_status` (`status`),
    KEY `idx_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户组表';

-- 3. 主机账号表
CREATE TABLE IF NOT EXISTS `host_accounts` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `name` VARCHAR(50) NOT NULL COMMENT '账号名称',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) COMMENT '密码（加密存储）',
    `secret_key` TEXT COMMENT '私钥内容',
    `type` TINYINT(1) NOT NULL DEFAULT 2 COMMENT '账号类型(1:root,2:普通)',
    `host_id` BIGINT UNSIGNED NOT NULL COMMENT '关联主机ID',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态(0:禁用,1:启用)',
    `remark` VARCHAR(255) COMMENT '备注',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    KEY `idx_host_id` (`host_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主机账号表';

-- 4. 主机组关联表（主机组 <-> 主机）
CREATE TABLE IF NOT EXISTS `host_group_relations` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `host_group_id` BIGINT UNSIGNED NOT NULL COMMENT '主机组ID',
    `host_id` BIGINT UNSIGNED NOT NULL COMMENT '主机ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    KEY `idx_host_group_id` (`host_group_id`),
    KEY `idx_host_id` (`host_id`),
    UNIQUE KEY `uk_group_host` (`host_group_id`, `host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主机组关联表';

-- 5. 用户组关联表（用户组 <-> 用户）
CREATE TABLE IF NOT EXISTS `user_group_relations` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `user_group_id` BIGINT UNSIGNED NOT NULL COMMENT '用户组ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    KEY `idx_user_group_id` (`user_group_id`),
    KEY `idx_user_id` (`user_id`),
    UNIQUE KEY `uk_group_user` (`user_group_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户组关联表';

-- 6. 主机用户权限关联表（用户组 <-> 主机组）
CREATE TABLE IF NOT EXISTS `host_user_permissions` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `user_group_id` BIGINT UNSIGNED NOT NULL COMMENT '用户组ID',
    `host_group_id` BIGINT UNSIGNED NOT NULL COMMENT '主机组ID',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    KEY `idx_user_group_id` (`user_group_id`),
    KEY `idx_host_group_id` (`host_group_id`),
    UNIQUE KEY `uk_user_host` (`user_group_id`, `host_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主机用户权限关联表';

-- 7. 用户审计表
CREATE TABLE IF NOT EXISTS `audit_logs` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `user_name` VARCHAR(50) NOT NULL COMMENT '用户名',
    `host_id` BIGINT UNSIGNED NOT NULL COMMENT '主机ID',
    `host_name` VARCHAR(100) NOT NULL COMMENT '主机名称',
    `host_address` VARCHAR(100) NOT NULL COMMENT '主机地址',
    `session_id` VARCHAR(100) NOT NULL COMMENT '会话ID',
    `action` TINYINT(1) NOT NULL COMMENT '操作类型(1:登录,2:执行命令,3:文件上传,4:文件下载,5:会话管理)',
    `command` TEXT COMMENT '执行的命令',
    `status` TINYINT(1) NOT NULL COMMENT '状态(1:成功,2:失败,3:警告)',
    `risk_level` TINYINT(1) NOT NULL COMMENT '风险等级(1:低,2:中,3:高,4:严重)',
    `client_ip` VARCHAR(50) COMMENT '客户端IP',
    `client_agent` VARCHAR(255) COMMENT '客户端User-Agent',
    `error_message` TEXT COMMENT '错误信息',
    `duration` BIGINT COMMENT '执行时长(毫秒)',
    `start_time` DATETIME NOT NULL COMMENT '开始时间',
    `end_time` DATETIME COMMENT '结束时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    KEY `idx_user_id` (`user_id`),
    KEY `idx_host_id` (`host_id`),
    KEY `idx_session_id` (`session_id`),
    KEY `idx_user_time` (`user_id`, `start_time`),
    KEY `idx_risk_level` (`risk_level`),
    KEY `idx_status` (`status`),
    KEY `idx_start_time` (`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户审计表';

-- 8. 批量任务表
CREATE TABLE IF NOT EXISTS `batch_tasks` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `name` VARCHAR(100) NOT NULL COMMENT '任务名称',
    `type` TINYINT(1) NOT NULL COMMENT '任务类型(1:命令,2:文件上传,3:文件下载,4:脚本)',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态(1:待执行,2:执行中,3:成功,4:失败,5:已取消)',
    `command` TEXT COMMENT '执行的命令或脚本内容',
    `source_path` VARCHAR(500) COMMENT '源文件路径（文件上传/下载）',
    `target_path` VARCHAR(500) COMMENT '目标路径',
    `script_type` VARCHAR(20) COMMENT '脚本类型(如:bash,python,shell)',
    `timeout` INT DEFAULT 300 COMMENT '超时时间(秒)',
    `success_count` INT DEFAULT 0 COMMENT '成功数量',
    `failed_count` INT DEFAULT 0 COMMENT '失败数量',
    `total_hosts` INT DEFAULT 0 COMMENT '总主机数',
    `progress` DECIMAL(5,2) DEFAULT 0 COMMENT '进度(0-100)',
    `remark` VARCHAR(255) COMMENT '备注',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `started_at` DATETIME COMMENT '开始时间',
    `finished_at` DATETIME COMMENT '完成时间',
    KEY `idx_status` (`status`),
    KEY `idx_created_by` (`created_by`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='批量任务表';

-- 9. 任务主机关联表
CREATE TABLE IF NOT EXISTS `task_host_relations` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `task_id` BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    `host_id` BIGINT UNSIGNED NOT NULL COMMENT '主机ID',
    `host_name` VARCHAR(100) NOT NULL COMMENT '主机名称',
    `host_addr` VARCHAR(100) COMMENT '主机地址',
    `account_id` BIGINT UNSIGNED COMMENT '使用的账号ID',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '执行状态',
    `output` TEXT COMMENT '执行输出',
    `error` TEXT COMMENT '错误信息',
    `duration` BIGINT COMMENT '执行时长(毫秒)',
    `started_at` DATETIME COMMENT '开始执行时间',
    `finished_at` DATETIME COMMENT '完成时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    KEY `idx_task_id` (`task_id`),
    KEY `idx_host_id` (`host_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务主机关联表';

-- 10. 定时计划表
CREATE TABLE IF NOT EXISTS `schedule_plans` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `name` VARCHAR(100) NOT NULL COMMENT '计划名称',
    `type` TINYINT(1) NOT NULL COMMENT '计划类型(1:一次,2:每天,3:每周,4:每月,5:Cron)',
    `status` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '状态(1:激活,2:暂停,3:已过期)',
    `cron_expression` VARCHAR(100) COMMENT 'Cron表达式',
    `execute_date` DATE COMMENT '执行日期(一次性执行)',
    `execute_time` TIME COMMENT '执行时间',
    `week_days` VARCHAR(20) COMMENT '周几(1-7,逗号分隔)',
    `month_day` INT COMMENT '每月几号(1-31)',
    `task_type` TINYINT(1) NOT NULL COMMENT '关联任务类型(1:命令,2:文件上传,3:文件下载,4:脚本)',
    `command` TEXT COMMENT '命令或脚本内容',
    `source_path` VARCHAR(500) COMMENT '源路径',
    `target_path` VARCHAR(500) COMMENT '目标路径',
    `script_type` VARCHAR(20) COMMENT '脚本类型',
    `timeout` INT DEFAULT 300 COMMENT '超时时间(秒)',
    `remark` VARCHAR(255) COMMENT '备注',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `last_executed_at` DATETIME COMMENT '最后执行时间',
    `next_executed_at` DATETIME COMMENT '下次执行时间',
    KEY `idx_status` (`status`),
    KEY `idx_created_by` (`created_by`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_next_executed_at` (`next_executed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='定时计划表';

-- 11. 定时计划主机关联表
CREATE TABLE IF NOT EXISTS `schedule_host_relations` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `schedule_plan_id` BIGINT UNSIGNED NOT NULL COMMENT '定时计划ID',
    `host_id` BIGINT UNSIGNED NOT NULL COMMENT '主机ID',
    `account_id` BIGINT UNSIGNED COMMENT '使用的账号ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    KEY `idx_schedule_plan_id` (`schedule_plan_id`),
    KEY `idx_host_id` (`host_id`),
    UNIQUE KEY `uk_schedule_host` (`schedule_plan_id`, `host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='定时计划主机关联表';

-- ==================== 初始化示例数据 ====================

-- 插入默认主机组
INSERT INTO `host_groups` (`id`, `name`, `desc`, `sort`, `status`, `created_by`) VALUES
(1, '生产环境组', '生产服务器组', 1, 1, 1),
(2, '测试环境组', '测试服务器组', 2, 1, 1),
(3, '开发环境组', '开发服务器组', 3, 1, 1)
ON DUPLICATE KEY UPDATE `id`=VALUES(`id`);

-- 插入默认用户组
INSERT INTO `user_groups` (`id`, `name`, `desc`, `sort`, `status`, `created_by`) VALUES
(1, '运维组', '运维人员组', 1, 1, 1),
(2, '开发组', '开发人员组', 2, 1, 1),
(3, '测试组', '测试人员组', 3, 1, 1)
ON DUPLICATE KEY UPDATE `id`=VALUES(`id`);

-- 插入默认权限配置
INSERT INTO `host_user_permissions` (`user_group_id`, `host_group_id`, `created_by`) VALUES
(1, 1, 1),  -- 运维组可以访问生产环境
(1, 2, 1),  -- 运维组可以访问测试环境
(2, 2, 1),  -- 开发组可以访问测试环境
(2, 3, 1),  -- 开发组可以访问开发环境
(3, 2, 1)   -- 测试组可以访问测试环境
ON DUPLICATE KEY UPDATE `created_by`=VALUES(`created_by`);
