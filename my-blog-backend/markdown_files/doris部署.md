# doris部署手册

（本次部署使用doris版本为2.1.9，三台机器ip使用node1；node2；node3进行替代，采用本地部署方式，部署节点为1FE3BE，包在官网即可下载）

doris官方下载地址https://doris.apache.org/download

### 一.上传包至服务器并解压

```
#需要上传的包有
jdk-8u161-linux-x64.tar.gz
mysql-5.7.22-linux-glibc2.12-x86_64.tar.gz
apache-doris-2.1.9-bin-x64.tar.gz

#解压
tar -xf 包名
```

### 二.创建doris-meta以及storage文件夹

```
#三台服务器都创建（默认可不创建）
mkdir -p /data/storage
mkdir -p /data/doris-meta
```

### 三.检查机器配置情况

```
# 关闭防火墙
systemctl stop firewalld
systemctl disable firewalld
systemctl status firewalld
#永久关闭selinux
sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
 # 临时关闭selinux
setenforce 0 
```

### 四.安装jdk

```
#上传jdk的tar包至目标路径并解压，修改用户环境变量
（jdk包大致都可以，本次选取的是自己打的包，可以根据需要进行修改）
tar  -zxvf jdk包（tar包）
cd
vi .bash_profile
#添加至末尾
#第一行目录根据jdk位置查看
export JAVA_HOME=/opt/jdk1.8.0_161
export PATH=$JAVA_HOME/bin:$PATH
export CLASSPATH=.:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar
#退出保存
:wq
source .bash_profile

```
### 五.修改fe配置文件

```
#解压后进入fe目录的conf文件夹修改fe.conf文件
cd /data/doris-2.1.9/fe/conf
vi fe.conf

# priority_networks = 10.10.10.0/24;192.168.0.0/16
修改为priority_networks = node1ip




#启动fe
./fe.start.sh --daemon
```
### 六.使用mysql客户端进行doris配置

```
#解压mysql文件夹，进入bin目录下连接doris，doris默认root用户，没有密码
./mysql -h node1 -P9030 -uroot

#修改密码
set password = password('password');

#添加be节点
ALTER SYSTEM ADD backend "node1:9050";
ALTER SYSTEM ADD backend "node2:9050";
ALTER SYSTEM ADD backend "node3:9050";
ALTER SYSTEM add BROKER broker1 "node1:8000";
ALTER SYSTEM add BROKER broker1 "node2:8000";
ALTER SYSTEM add BROKER broker1 "node3:8000";

```

### 七.修改be配置并启动be

```
#部署be的节点执行
sysctl -w vm.max_map_count=2000000
ulimit -n 65536

#进入be的conf文件夹修改be.conf
cd /data/doris-2.1.9/be/conf/
vi be.conf

# priority_networks = 10.10.10.0/24;192.168.0.0/16
按照所在机器分别修改为本机ip
priority_networks = node1
priority_networks = node2
priority_networks = node3


#启动be与broker
sh /data/doris-2.1.9/be/bin/start_be.sh --daemon
sh /data/doris-2.1.9/apache_hdfs_broker/bin/start_broker.sh --daemon

```


### 八.访问页面并验证
前端页面登录地址
http://node1:8030/

~~~
#点击System下的backends
出现三个be节点，并且节点Alive一栏都为true即可
~~~



最后进行建库建表，没问题后即为搭建完成

