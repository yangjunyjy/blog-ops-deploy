# doris可查pg数据库配置

### doris配置

```shell
#要连接到 PostgreSQL 数据库，您需要PostgreSQL 11.x 或更高版本PostgreSQL 数据库的 JDBC 驱动程序，您可以从 Maven 仓库下载最新或指定版本的 PostgreSQL JDBC 驱动程序。推荐使用 PostgreSQL JDBC Driver 42.5.x 及以上版本。

https://doris.apache.org/zh-CN/docs/2.0/lakehouse/database/jdbc

type	固定为 jdbc
user	数据源用户名
password	数据源密码
jdbc_url	数据源连接 URL
driver_url	数据源 JDBC 驱动程序的路径
driver_class	数据源 JDBC 驱动程序的类名

lower_case_table_names	"false"	是否以小写的形式同步 jdbc 外部数据源的库名和表名
only_specified_database	"false"	是否只同步 JDBC URL 中指定的数据源的 Database（此处的 Database 为映射到 Doris 的 Database 层级）
include_database_list	""	当 only_specified_database=true 时，指定同步多个 Database，以','分隔。Database 名称是大小写敏感的。
exclude_database_list	""	当 only_specified_database=true 时，指定不需要同步的多个 Database，以','分割。Database 名称是大小写敏感的。

connection_pool_min_size	1	定义连接池的最小连接数，用于初始化连接池并保证在启用保活机制时至少有该数量的连接处于活跃状态。
connection_pool_max_size	30	定义连接池的最大连接数，每个 Catalog 对应的每个 FE 或 BE 节点最多可持有此数量的连接。
connection_pool_max_wait_time	5000	如果连接池中没有可用连接，定义客户端等待连接的最大毫秒数。
connection_pool_max_life_time	1800000	设置连接在连接池中保持活跃的最大时长（毫秒）。超时的连接将被回收。同时，此值的一半将作为连接池的最小逐出空闲时间，达到该时间的连接将成为逐出候选对象。
connection_pool_keep_alive	false	仅在 BE 节点上有效，用于决定是否保持达到最小逐出空闲时间但未到最大生命周期的连接活跃。默认关闭，以减少不必要的资源使用。



#实操
文件名。如 mysql-connector-j-8.3.0.jar。需将 Jar 包预先存放在 FE 和 BE 部署目录下的 jdbc_drivers/ 目录下。系统会自动在这个目录下寻找。该目录的位置，也可以由 fe.conf 和 be.conf 中的 jdbc_drivers_dir 配置修改

#修改路径配置后，在doris下执行一下sql，更具自己的配置更改
CREATE CATALOG pg_cmdb PROPERTIES (
    "type"="jdbc",
    "user"="cmdb",
    "password"="密码",
    "jdbc_url" = "jdbc:postgresql://10.167.90.12:18922/cmdb_db",
    "driver_url" = "postgresql-42.5.5.jar",
    "driver_class" = "org.postgresql.Driver"
)

#检验
SHOW CATALOGS;
SHOW CREATE CATALOG pg_cmdb;

-- 查看所有CATALOG的基本信息
SELECT * FROM information_schema.catalogs;

-- 查看特定CATALOG的详细信息
SELECT * FROM information_schema.catalogs WHERE catalog_name = 'pg_cmdb';

#查库
SHOW DATABASES FROM pg_cmdb;
SHOW TABLES FROM pg_cmdb.cmdb_db;  -- 假设cmdb_db是您的schema名
#查表
SELECT * FROM pg_cmdb.[schema名].[表名];
SELECT * FROM pg_cmdb.public.designpage limit 10;


-- 3. 删除CATALOG
DROP CATALOG IF EXISTS pg_cmdb;
-- 4. 验证是否删除成功
SHOW CATALOGS;
```

