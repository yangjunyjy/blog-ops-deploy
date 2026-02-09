# 堡垒机数据库设计文档

## 设计原则

1. **不使用外键约束**：所有关联关系通过中间表实现，但不设置 FOREIGN KEY 约束
2. **RBAC 权限模型**：通过用户组 -> 主机组 的映射实现基于角色的访问控制
3. **审计可追溯**：所有操作记录到审计表，支持风险分析和行为追踪
4. **批量任务支持**：支持对多个主机执行批量任务
5. **定时任务调度**：支持多种定时策略执行自动化任务

## 表结构说明

### 1. host_groups（主机组表）
用于管理主机的分组，方便批量管理和权限控制。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| name | VARCHAR(100) | 主机组名称（唯一） |
| desc | VARCHAR(255) | 描述 |
| sort | INT | 排序 |
| status | TINYINT(1) | 状态(0:禁用,1:启用) |
| created_by | BIGINT | 创建人ID |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### 2. user_groups（用户组表）
用于管理用户的分组，实现基于组的权限管理。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| name | VARCHAR(100) | 用户组名称（唯一） |
| desc | VARCHAR(255) | 描述 |
| sort | INT | 排序 |
| status | TINYINT(1) | 状态(0:禁用,1:启用) |
| created_by | BIGINT | 创建人ID |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### 3. host_accounts（主机账号表）
管理每台主机的多个账号，支持密钥和密码两种认证方式。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| name | VARCHAR(50) | 账号名称 |
| username | VARCHAR(50) | 用户名 |
| password | VARCHAR(255) | 密码（加密存储） |
| secret_key | TEXT | 私钥内容 |
| type | TINYINT(1) | 账号类型(1:root,2:普通) |
| host_id | BIGINT | 关联主机ID |
| status | TINYINT(1) | 状态(0:禁用,1:启用) |
| remark | VARCHAR(255) | 备注 |
| created_by | BIGINT | 创建人ID |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### 4. host_group_relations（主机组关联表）
主机组与主机的多对多关系表。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| host_group_id | BIGINT | 主机组ID |
| host_id | BIGINT | 主机ID |
| created_at | DATETIME | 创建时间 |

**唯一索引**: uk_group_host (host_group_id, host_id)

### 5. user_group_relations（用户组关联表）
用户组与用户的多对多关系表。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| user_group_id | BIGINT | 用户组ID |
| user_id | BIGINT | 用户ID |
| created_at | DATETIME | 创建时间 |

**唯一索引**: uk_group_user (user_group_id, user_id)

### 6. host_user_permissions（主机用户权限关联表）
实现 RBAC：用户组对主机组拥有访问权限。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| user_group_id | BIGINT | 用户组ID |
| host_group_id | BIGINT | 主机组ID |
| created_by | BIGINT | 创建人ID |
| created_at | DATETIME | 创建时间 |

**唯一索引**: uk_user_host (user_group_id, host_group_id)

### 7. audit_logs（用户审计表）
记录用户的所有操作，用于审计分析和异常行为检测。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| user_id | BIGINT | 用户ID |
| user_name | VARCHAR(50) | 用户名 |
| host_id | BIGINT | 主机ID |
| host_name | VARCHAR(100) | 主机名称 |
| host_address | VARCHAR(100) | 主机地址 |
| session_id | VARCHAR(100) | 会话ID |
| action | TINYINT(1) | 操作类型(1:登录,2:执行命令,3:文件上传,4:文件下载,5:会话管理) |
| command | TEXT | 执行的命令 |
| status | TINYINT(1) | 状态(1:成功,2:失败,3:警告) |
| risk_level | TINYINT(1) | 风险等级(1:低,2:中,3:高,4:严重) |
| client_ip | VARCHAR(50) | 客户端IP |
| client_agent | VARCHAR(255) | 客户端User-Agent |
| error_message | TEXT | 错误信息 |
| duration | BIGINT | 执行时长(毫秒) |
| start_time | DATETIME | 开始时间 |
| end_time | DATETIME | 结束时间 |
| created_at | DATETIME | 创建时间 |

**索引**:
- idx_user_time (user_id, start_time)
- idx_host_id (host_id)
- idx_session_id (session_id)
- idx_risk_level (risk_level)

### 8. batch_tasks（批量任务表）
用于执行多个主机的批量任务。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| name | VARCHAR(100) | 任务名称 |
| type | TINYINT(1) | 任务类型(1:命令,2:文件上传,3:文件下载,4:脚本) |
| status | TINYINT(1) | 状态(1:待执行,2:执行中,3:成功,4:失败,5:已取消) |
| command | TEXT | 执行的命令或脚本内容 |
| source_path | VARCHAR(500) | 源文件路径（文件上传/下载） |
| target_path | VARCHAR(500) | 目标路径 |
| script_type | VARCHAR(20) | 脚本类型 |
| timeout | INT | 超时时间(秒) |
| success_count | INT | 成功数量 |
| failed_count | INT | 失败数量 |
| total_hosts | INT | 总主机数 |
| progress | DECIMAL(5,2) | 进度(0-100) |
| remark | VARCHAR(255) | 备注 |
| created_by | BIGINT | 创建人ID |
| created_at | DATETIME | 创建时间 |
| started_at | DATETIME | 开始时间 |
| finished_at | DATETIME | 完成时间 |

### 9. task_host_relations（任务主机关联表）
记录批量任务中每台主机的执行情况。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| task_id | BIGINT | 任务ID |
| host_id | BIGINT | 主机ID |
| host_name | VARCHAR(100) | 主机名称 |
| host_addr | VARCHAR(100) | 主机地址 |
| account_id | BIGINT | 使用的账号ID |
| status | TINYINT(1) | 执行状态 |
| output | TEXT | 执行输出 |
| error | TEXT | 错误信息 |
| duration | BIGINT | 执行时长(毫秒) |
| started_at | DATETIME | 开始执行时间 |
| finished_at | DATETIME | 完成时间 |
| created_at | DATETIME | 创建时间 |

### 10. schedule_plans（定时计划表）
用于设定定时计划，支持多种定时策略。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| name | VARCHAR(100) | 计划名称 |
| type | TINYINT(1) | 计划类型(1:一次,2:每天,3:每周,4:每月,5:Cron) |
| status | TINYINT(1) | 状态(1:激活,2:暂停,3:已过期) |
| cron_expression | VARCHAR(100) | Cron表达式 |
| execute_date | DATE | 执行日期(一次性执行) |
| execute_time | TIME | 执行时间 |
| week_days | VARCHAR(20) | 周几(1-7,逗号分隔) |
| month_day | INT | 每月几号(1-31) |
| task_type | TINYINT(1) | 关联任务类型(1:命令,2:文件上传,3:文件下载,4:脚本) |
| command | TEXT | 命令或脚本内容 |
| source_path | VARCHAR(500) | 源路径 |
| target_path | VARCHAR(500) | 目标路径 |
| script_type | VARCHAR(20) | 脚本类型 |
| timeout | INT | 超时时间(秒) |
| remark | VARCHAR(255) | 备注 |
| created_by | BIGINT | 创建人ID |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |
| last_executed_at | DATETIME | 最后执行时间 |
| next_executed_at | DATETIME | 下次执行时间 |

### 11. schedule_host_relations（定时计划主机关联表）
定时计划与主机的关联关系。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| schedule_plan_id | BIGINT | 定时计划ID |
| host_id | BIGINT | 主机ID |
| account_id | BIGINT | 使用的账号ID |
| created_at | DATETIME | 创建时间 |

**唯一索引**: uk_schedule_host (schedule_plan_id, host_id)

## 枚举类型说明

### AccountType（账号类型）
- `1` - RootAccount: root 账号
- `2` - NormalAccount: 普通账号

### AuditStatus（审计状态）
- `1` - AuditSuccess: 成功
- `2` - AuditFailed: 失败
- `3` - AuditWarning: 警告

### AuditAction（审计操作类型）
- `1` - LoginAction: 登录
- `2` - ExecuteAction: 执行命令
- `3` - FileUploadAction: 文件上传
- `4` - FileDownloadAction: 文件下载
- `5` - SessionAction: 会话管理

### RiskLevel（风险等级）
- `1` - LowRisk: 低风险
- `2` - MediumRisk: 中风险
- `3` - HighRisk: 高风险
- `4` - CriticalRisk: 严重风险

### TaskStatus（任务状态）
- `1` - TaskPending: 待执行
- `2` - TaskRunning: 执行中
- `3` - TaskSuccess: 成功
- `4` - TaskFailed: 失败
- `5` - TaskCanceled: 已取消

### TaskType（任务类型）
- `1` - CommandTask: 命令任务
- `2` - FileUploadTask: 文件上传
- `3` - FileDownloadTask: 文件下载
- `4` - ScriptTask: 脚本任务

### ScheduleType（定时计划类型）
- `1` - OnceSchedule: 执行一次
- `2` - DailySchedule: 每天
- `3` - WeeklySchedule: 每周
- `4` - MonthlySchedule: 每月
- `5` - CronSchedule: Cron表达式

### ScheduleStatus（定时计划状态）
- `1` - ScheduleActive: 激活
- `2` - SchedulePaused: 暂停
- `3` - ScheduleExpired: 已过期

## 权限模型

### RBAC 流程
1. 用户属于一个或多个用户组（user_group_relations）
2. 用户组对一个或多个主机组有权限（host_user_permissions）
3. 主机组包含多台主机（host_group_relations）

### 权限判断逻辑
```go
// 判断用户是否有权访问某台主机
func CanUserAccessHost(userID, hostID uint) bool {
    // 1. 获取用户所属的用户组
    userGroupIDs := GetUserGroupIDs(userID)

    // 2. 获取这些用户组有权限的主机组
    hostGroupIDs := GetHostGroupIDsByUserGroups(userGroupIDs)

    // 3. 获取这些主机组包含的主机
    hostIDs := GetHostIDsByHostGroups(hostGroupIDs)

    // 4. 判断目标主机是否在列表中
    return Contains(hostIDs, hostID)
}
```

## 使用示例

### 查询用户可访问的所有主机
```sql
SELECT DISTINCT rh.*
FROM remote_hosts rh
INNER JOIN host_group_relations hgr ON rh.id = hgr.host_id
INNER JOIN host_user_permissions hup ON hgr.host_group_id = hup.host_group_id
INNER JOIN user_group_relations ugr ON hup.user_group_id = ugr.user_group_id
WHERE ugr.user_id = ?
  AND hgr.host_group_id IN (
    SELECT host_group_id FROM host_group_relations WHERE host_id = ?
  )
  AND rh.status = 1
  AND hgr.host_group_id IN (
    SELECT id FROM host_groups WHERE status = 1
  )
ORDER BY hgr.host_group_id, rh.id;
```

### 查询用户审计日志
```sql
SELECT *
FROM audit_logs
WHERE user_id = ?
  AND start_time >= ?
  AND start_time <= ?
ORDER BY start_time DESC
LIMIT 20;
```

### 查询高风险操作
```sql
SELECT al.*, u.username, u.nickname
FROM audit_logs al
INNER JOIN users u ON al.user_id = u.id
WHERE al.risk_level >= 3
  AND al.created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
ORDER BY al.created_at DESC;
```

### 查询批量任务执行情况
```sql
SELECT
    bt.id,
    bt.name,
    bt.status,
    bt.total_hosts,
    bt.success_count,
    bt.failed_count,
    bt.progress,
    COUNT(CASE WHEN thr.status = 3 THEN 1 END) AS completed_hosts,
    COUNT(CASE WHEN thr.status = 4 THEN 1 END) AS failed_hosts
FROM batch_tasks bt
LEFT JOIN task_host_relations thr ON bt.id = thr.task_id
WHERE bt.id = ?
GROUP BY bt.id;
```

## 后续开发建议

1. **审计规则引擎**：基于审计日志实现智能风险检测
2. **命令过滤**：对敏感命令（如 rm、dd）进行拦截或二次确认
3. **会话录制**：支持 SSH 会话视频录制和回放
4. **文件传输审计**：记录所有文件上传下载操作
5. **任务重试机制**：批量任务失败时支持自动重试
6. **通知提醒**：任务完成或异常时发送通知
7. **多租户支持**：支持企业级多租户场景
