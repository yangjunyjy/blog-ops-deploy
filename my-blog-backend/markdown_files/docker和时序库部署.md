### 1.1、docker部署流程

node1节点和node2节点安装

~~~bash
tar -zxvf docker-25.0.5.tgz
sudo cp docker/* /usr/bin/
cd /etc/systemd/system/ 
#新建文件
sudo touch docker.service
# 为 docker.service 文件添加执行权限
sudo chmod +x /etc/systemd/system/docker.service
~~~

~~~bash
sudo cat > /etc/systemd/system/docker.service << EOF
#进入编辑模式
#docker.service文件编写(不同服务器需修改部署ip)
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network-online.target firewalld.service
Wants=network-online.target

[Service]
Type=notify
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
#修改为服务器部署ip
ExecStart=/usr/bin/dockerd --selinux-enabled=false --insecure-registry=node1
ExecReload=/bin/kill -s HUP $MAINPID
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
# Uncomment TasksMax if your systemd version supports it.
# Only systemd 226 and above support this version.
#TasksMax=infinity
TimeoutStartSec=0
# set delegate yes so that systemd does not reset the cgroups of docker containers
Delegate=yes
# kill only the docker process, not all processes in the cgroup
KillMode=process
# restart the docker process if it exits prematurely
Restart=on-failure
StartLimitBurst=3
StartLimitInterval=60s

[Install]
WantedBy=multi-user.target
EOF
~~~

docker启动

~~~bash
sudo systemctl daemon-reload
# 启动，设置开机启动，查看docker服务状态
sudo systemctl start docker
sudo systemctl enable docker.service
sudo systemctl status docker
~~~

修改docker文件存储路径(看哪个路径空间大，建立docker数据存储路径)

```bash
mkdir -p /data/oneops/docker_data
#没有就创建
mkdir /etc/docker
sudo cat /etc/docker/daemon.json << EOF #没有就创建 
#添加内容如下
{
"data-root": "/home/oneops/docker_data"
}
#重启docker
EOF
sudo systemctl restart docker
```

赋予普通用户docker语句执行权限(如果登录是root用户则不用)

```bash
#查看是否存在docker用户组
sudo cat /etc/group |grep docker
#添加docker用户组
sudo groupadd docker
#查看docker.sock所属的用户组
ll /var/run/docker.sock
#更改所属用户组
sudo chgrp docker /var/run/docker.sock
sudo chgrp docker /bin/docker*
#将普通用户添加至用户组
sudo gpasswd -a oneops docker
#刷新docker用户组(普通用户无需加sudo)
sudo newgrp docker
#验证是否添加成功
id
```

docker容器间通信问题解决办法(如果遇到问题才需要执行，通过宿主机IP)

```bash
#防火墙添加受信任的IP段
firewall-cmd --permanent --zone=public --add-rich-rule='rule family=ipv4 source address=172.17.0.0/16 accept'
#重新加载防火墙
firewall-cmd --reload
```

### 1.2 时序库部署

#载入镜像

```shell
docker load < timescaledb-ha.tar
```

#为docker服务打标签

```shell
docker tag 424 telepg:1.0
```

#启动命令(注意替换密码和dockerid)

```shell
docker run -p 18922:5432 -e POSTGRES_PASSWORD='S7uV4w@XxN' -v pgdata:/var/lib/postgresql/data -d telepg:1.0
#拷贝数据库表结构文件至容器
docker cp cmdb_telepg.sql 48ebe040a783:/var/lib/postgresql/data/
#进入容器(root用户权限进入)
docker exec -it --user root 48ebe040a783 bash
#对文件进行赋权
chown postgres:postgres /var/lib/postgresql/data/cmdb_telepg.sql
#切换至postgres用户
su - postgres
#进入控制台
psql
#建库语句并生成toolkit插件(需同样切换至postgres用户，psql进入控制台)
create user cmdb_telepg with password 'S7uV4w@XxN';
create database cmdb_telepg owner cmdb_telepg;
grant all privileges on database cmdb_telepg to cmdb_telepg;
ALTER USER "cmdb_telepg" with superuser;
#切换到cmdb_telepg用户名下
\c - cmdb_telepg
CREATE EXTENSION timescaledb_toolkit;
#导入表结构
exit
#推出数据库，至root权限下
exit
cd /var/lib/postgresql/data/
psql -U cmdb_telepg -d cmdb_telepg -a -f cmdb_telepg.sql
#切换到cmdb_telepg数据库下
psql
\c - cmdb_telepg
#切换用户
exit
psql "host =127.0.0.1 port =5432 user =cmdb_telepg password='S7uV4w@XxN' dbname =cmdb_telepg"
```