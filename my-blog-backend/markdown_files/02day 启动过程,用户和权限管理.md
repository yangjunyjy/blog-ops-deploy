[TOC]

### 一.Linux启动流程

#### 1.1CentOS7启动流程

```shell
1.BIOS自检：
   BIOS是一种固化在主板上的固件程序，负责硬件初始化和启动引导。
   BIOS在启动时会进行硬件自检（POST，Power-On Self-Test），检查硬件设备是否正常。
   如果硬件检测通过，BIOS 会根据启动顺序（如硬盘、光驱、USB 等）查找可引导设备。
   
2.读取MBR引导代码到内存中执行

3.加载引导程序GRUB
  引导程序的作用是加载操作系统内核到内存中。
  GRUB是目前最常用的引导程序，它会显示一个菜单，允许用户选择启动的操作系统或内核版本。
  
4.加载内核
  当引导程序（如 GRUB）启动时，它会从磁盘加载 Linux 内核到内存中。
  内核文件通常是一个压缩的镜像文件（如 vmlinuz）。
  内核在加载过程中会初始化硬件设备，如 CPU、内存、硬盘等。
  内核会挂载真实根文件系统（通常是只读模式）:通过/sbin/init(实际为systemd)挂载/文件系统，并启动用户空间的第一个进程。
  切换根目录:内核将控制权交给systemd(PID 1)
  
5.启动用户空间(systemd 初始化)
  systemd 作为第一个进程：启动后执行以下步骤：
  初始化默认目标（target）：默认为 default.target（通常是 multi-user.target 或 graphical.target）。
  解析单元文件：加载 /usr/lib/systemd/system/ 和 /etc/systemd/system/ 中的 .service、.target 等单元文件。
  并行启动服务：根据依赖关系并行启动服务，加速启动过程。
  
6.目标(targets)和服务启动
目标（Targets）：类似于旧版的“运行级别”，但更灵活。常见目标：
    Target	描述	对应旧运行级别
    emergency.target	紧急模式（单用户，无网络）	无直接对应
    rescue.target	    救援模式（基本系统）	      单用户模式（1）
    multi-user.target	多用户命令行模式	       运行级别 3
    graphical.target	图形界面模式	            运行级别 5
    服务启动：依次启动目标依赖的服务（如网络、日志、SSH 等）。

7.用户登录
   登录管理器：如果目标是 graphical.target，启动显示管理器（如 GDM）进入图形登录界面。
   命令行登录：在 multi-user.target 下，启动 getty 服务，提供 TTY 命令行登录。
   
8.用户会话
   用户登录后，系统会启动用户会话，加载用户的配置文件（如 .bashrc、.profile 等）。
   用户可以启动应用程序，进行日常工作。
   
总结
Linux 启动流程可以概括为以下步骤：
1.BIOS/UEFI 启动：硬件初始化和自检。
2.MBR/GPT 和引导程序：加载引导程序（如 GRUB）。
3.加载内核：引导程序加载 Linux 内核。
4.启动用户空间：内核启动 init 或 systemd。
5.运行级别/目标：systemd 根据目标启动服务。
6.用户登录：启动登录管理器，用户登录。
7.用户会话：用户开始使用系统。
```

---

#### 1.2.注意

在用户登录的时候，还有一点需要特别注意

```sh
不管是什么用户，只要登录就必执行下面的内容
/etc/profile       所有用户通用的变量写在这个文件

/etc/profile.d     所有用户通用的脚本文件写在这个目录
```

在上面的2个文件或目录执行完以后，就会执行该用户自己的家目录下面的

```sh
.bashrc           如果想给某一个用户单独设置一些变量，就在这个用户的家目录下的.bashrc文件里面设置，这个变量只有这个用户才能使用
```

---



### 二./etc 目录核心文件

```
/etc/sysconfig/network-scripts/ifcfg-ens33  Centos7网卡配置文件
/etc/resolv.conf DNS服务器配置文件，优先级小于网卡配置文件
DNS就是域名解析，把域名解析成ip地址
/etc/hosts  域名和IP地址映射解析文件，一般用于局域网内部解析
windows里面对应的文件是
C:\Windows\System32\drivers\etc\HOSTS
/etc/fstab  系统自动挂载配置文件，所有设备都必须先进行挂载才能去使用，谨慎更改!!!
/etc/rc.local 存放系统开机自启的命令，是个链接文件指向的是/etc/rc.d/rc.local
/etc/profile /etc/bashrc  系统变量
/etc/issue  开机显示详细信息
/etc/motd   欢迎词，提示词，全局
/etc/redhat-release 查看系统版本
/etc/sysctl.conf 内核参数设置,重要
/etc/init.d Centos7之前的软件启动目录
```

---



### 三./usr 目录核心文件

```
/usr/local  编译安装软件的默认位置路径
/usr/src    存放源码的路径
/usr/sbin   存放开机时不需要用到的命令或脚本
/usr/share  存放帮助文档或共享文件
```

---



### 四./var 目录核心文件

```
/var/log  记录软件跟系统信息文件的目录，log日志本质就是记录计算机行为或软件行为
/var/log/messages 系统级别日志，记录Linux系统运行行为的文件，正常与非正常都记录，每周生成一个新的日志，这种机制叫做轮询
/var/log/secure  用户登录信息日志,用户所有的登录信息包含密码是否正确，也是每周生成一个新的日志
/var/log/dmesg  记录硬件信息日志
```

---



### 五./proc 目录核心文件

```
/proc 内存相关目录
/proc/meminfo 查看系统内存的，对应命令是free -h
/proc/cpuinfo 查看cpu相关信息，对应命令是lscpu
/proc/loadavg 查看负载(系统繁忙程度的综合指标),对应命令是uptime

[root@didiyun proc]# uptime
 10:59:37 up 23:07,  2 users,  load average: 0.00, 0.01, 0.05
 10:59:37 当前时间
 up 23:07 开机到现在多长时间
 2 users 几个用户登录
 load average: 0.00, 0.01, 0.05
 （负载）（平均） 1分钟  5分钟  15分钟
```

----

### 六.主机名

```shell
主机名方便我们管理员进行标识
存放在/etc/hostname文件下

1.查看主机名
cat /etc/hostname 

2.更改
2.1.临时更改
hostname 主机名
#刷新
bash

2.2.永久更改
hostnamectl set-hostname 主机名
#刷新
bash
```

<img src="F:\运维笔记\笔记\图片\hostname.PNG" style="zoom:80%;" />

---

### 七.Selinux介绍

```
安全增强式Linux是一个Linux内核的安全模块，其提供了访问控制安全策略机制，包括了强制访问控制。是美国安全局制定的规则，让Linux遵循这套安全规则，但是过于严格，我们工作的企业都是处于关闭状态，可以利用其它手段来实现安全目的。
```

#### 7.1.查看关闭SeLinux

```shell
1.查看
getenforce 
Disabled 关闭

2.临时关闭
setenforce 0
重启系统后就失效了

3.永久关闭
/etc/selinux/config Selinux配置文件
[root@didiyun ~]# vim /etc/selinux/config
SELINUX=disabled
```

***

### 八.Linux中的通配符

```
什么是通配符？
就是键盘上面的特殊符号 "~","*","?","[]" ，在Linux系统中我们可以用特殊符号替代一些功能，这些符号的用法都是Linux系统规定好的，一般在命令行中去使用。
```

#### 8.1.模糊匹配

```
"*" 代表全部，任何字符，零个或多个
注意*不能匹配隐藏文件，也就是说*是无法找到隐藏文件的

"?" 匹配单个字符，有且只有一个，精准

"[]" 匹配中括号里面的任意!!字符!!，可以是连续字符，也可以是不连续字符

"!" 取反，也可以用 "^" 代替
```



***

#### 8.2其他

```
";" 命令分隔符

"\" 转义字符，让字符还原本意

"{起始..结束}" 生成序列

"&&" and 前面执行正确才执行后面的内容
"||" or  前面执行错误才执行后面的内容
```

***

### 九.用户与用户组基本概念

```shell
为什么要有用户跟用户组？
Linux系统规定所有文件跟进程都必须拥有所有者
```

#### 9.1.用户的基础概念

```
Linux是多用户多进程的系统，用户是管理Linux系统存在的，用户分为用户名跟ID相当于我们的名字跟身份证号，名字是为了方便我们进行区分的，用户名也不能重复，系统只认id号，ID号是唯一的区别，叫UID
用户被分为三类
第一类 超级管理员 root，UID为0的就是超级管理员，有且只有一个
第二类 虚拟用户，又被称为傀儡用户，能满足文件或进程属主的要求，又不会给管理带来风险。uid1-999
第三类 普通用户是实际存在的用户，允许登录，管理员身份创建，帮助管理员管理系统的，权限限制在家目录，在一些系统目录下只能读不能写，不能进入到/root目录，uid1000起步默认是1000-60000.

用户的详细信息记录在/etc/passwd下
```

#### 9.2.用户管理基础命令

```shell
id 查看用户详细信息(uid/gid/用户组)
id 显示当前登录用户的uid等信息
id user 显示指定用户的uid等信息

whoami查看当前登录用户

useradd 新建用户

passwd 更改用户密码

userdel 删除用户
	-r 删除家目录
```

#### 9.3.用户组的基础概念

```
用户组相当于一个集合，组织，就跟人的一个家庭体系一样，Linux每个用户都必须要有一个组织，这个组织就是用户组。
用户组的名称也是为了方便我们管理员进行辨别的，用户组跟用户一样也有一个唯一标识，叫group id 简写为gid,默认情况下我们创建用户，Linux系统会自动创建一个和用户名相同的用户组，uid跟gid也相同。
用户组与用户的对应关系
用户与用户组不是一一对应的关系
一个用户可以对应多个用户组
多个用户对应一个用户组
多个用户对应多个用户组
```

#### 9.4.用户组管理基础命令

```
用户组是可以单独存在的，理解即可

groupadd 创建用户组

groupdel 删除用户组
```

***

### 十.用户用户组相关文件

#### 10.1./etc/passwd 用户信息文件

```shell
[root@didiyun ~]# vim /etc/passwd
每一行都是一个用户信息，冒号作为分隔，一共七列
root:x:0:0:root:/root:/bin/bash

第一列 用户名
第二列 密码，x代表有密码
第三列 UID
第四列 GID
第五列 用户说明
第六列 用户的家目录
第七列 shell解释器
/sbin/nologin                禁止用户登录
```

#### 10.2./etc/shadow 密码文件

```shell
[root@didiyun ~]# vim /etc/shadow
每一行都是一个用户密码信息，冒号分隔，共九列
root:$6$YmsIxPxlQzriz6/J$0YOXF0xfE.tCYIxylZX9XSlxa3ByO.N9UbppWhraYJgio8C/Hi7OgKUrZjUvGwoWtSp6a6vPHhKj9mIbt2X53/:18372:0:99999:7:::
wb:!!:18562:0:99999:7:::

第一列 用户名称
第二列 用户密码(SHA512加密)
第三列 更新密码的时间(时间戳) date -d "1970-01-01 18562 days"
第四列 两次密码修改的时间间隔，跟第三列的时间比较
第五列 密码的有效期，跟第三列时间比较
第六列 密码失效前多少天给予提示，和第五列时间比较
第七列 密码过期后的宽限时间，和第五列时间比较
第八列 密码失效时间，使用时间戳，自1970-1-1以来的总天数作为账号失效时间，这个时间戳以外的其他字段无效
第九列 保留，等待新功能加入
```

#### 10.3./etc/group 用户组文件

```shell
[root@didiyun ~]# vim /etc/group
每一行都是一个用户组信息，冒号分隔，共四列
root:x:0:

第一列 用户组名
第二列 组密码位
第三列 GID
第四列 组成员
```

### 十一.用户用户组命令进阶

#### 11.1.useradd 进阶

```shell
useradd 新建用户
常见选项
	-u 指定UID，一般不会去指定
	-c 添加说明
	-d 指定家目录
	-s 指定shell解释器
	-G 把用户加入到其他用户组(附加组)
	-M 不创建家目录
	-g 把用户加入到其他用户组(初始组)
	
useradd 新建用户的默认选项参考 /etc/default/useradd 和 /etc/login.defs

创建虚拟用户
useradd -s /sbin/nologin -M mysql
```

#### 11.2.passwd 进阶

```shell
passwd 更改密码
常见选项
	-l 暂时锁定用户
	-u 解锁用户
	--stdin 从标准输入接收密码
echo "520yjYJ," | passwd --stdin mysql
```

#### 11.3.su 切换用户

```
su - 用户名
	- 切换用户的同时环境变量一起切过去
高权限切换到低权限不需要密码
低权限切换到高权限需要输入高权限用户的密码
```

#### 11.4.gpasswd 用户组添加删除

```
gpasswd 修改已存在用户组信息
常用选项
	-a 把用户加入用户组
	-d 把用户移出用户组
语法
gpasswd -a 用户 用户组
gpasswd -d 用户 用户组
```



### 十二.useradd默认选项详解

#### 12.1./etc/default/useradd

```
[root@didiyun ~]# vim /etc/default/useradd 
# useradd defaults file
GROUP=100
HOME=/home
INACTIVE=-1
EXPIRE=
SHELL=/bin/bash
SKEL=/etc/skel
CREATE_MAIL_SPOOL =yes

GROUP 用户初始组就是GID100这个用户组，这个选项我们并没有去采用，我们的机制是每个用户都有一个默认跟自己同名的用户组
HOME 家目录位置
INACTIVE 密码过期后的宽限天数，-1 表示密码永不过期，对应/etc/shadow 第7列
EXPIRE 密码失效时间，默认为空，表示没有失效时间
SHELL 默认的命令解释器
SKEL 指定我们创建用户的时候使用那个模板文件
CREATE_MAIL_SPOOL 每一个用户都创建邮箱
```

#### 12.2.etc/skel/

```
skel 是 skeleton 的缩写,每当你新建一个用户的时候 (通过 useradd 命令),/etc/skel 目录下的文件,都会原封不动的复制到新建用户的家目录下～
```

#### 12.3./etc/login.defs

```
[root@didiyun skel]# vim /etc/login.defs
MAIL_DIR        /var/spool/mail
PASS_MAX_DAYS   99999
PASS_MIN_DAYS   0
PASS_MIN_LEN    5
PASS_WARN_AGE   7
UID_MIN         1000
UID_MAX         60000
CREATE_HOME     yes
UMASK           077
USERGROUPS_ENAB yes
ENCRYPT_METHOD SHA512

MAIL_DIR 邮箱位置
PASS_MAX_DAYS 密码有效期
PASS_MIN_DAYS 密码修改间隔时间
PASS_MIN_LEN  密码最小长度，废弃状态，PAM模块代替它
PASS_WARN_AGE 密码到期前的警告天数
UID_MIN  UID_MAX   UID的范围
CREATE_HOME 建立用户时，是否建立家目录
UMASK 建立家目录的默认权限
USERGROUPS_ENAB  删除用户时，是否删除用户组
ENCRYPT_METHOD 用户密码加密方式
```

***

###                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             十三.Linux文件权限

```sh
如果发现自己机器的时间不对，先看一下时区是否设置正确，如果时区设置没有问题，那么可以使用下面的方式来进行解决
如果机器可以联网，那么就执行下面的命令
时间同步命令
ntpdate ntp3.aliyun.com
执行完成后，同步BIOS时钟，强制把系统时间写入CMOS，命令如下：
hwclock -w
这样就可以保证机器重启后，时间也是正确的了
如果机器不能联网，那么就使用下面的命令设置时间
date -s "月份/日期/年份 小时:分钟:秒"
执行完成后，同步BIOS时钟，强制把系统时间写入CMOS，命令如下：
hwclock -w
```



### 十四、文件属性详解

```shell
-rwxrwxrwx 1 root root 25 Dec 17 10:06 hello.py

第一项
第一个字符(文件属性)
d 代表目录 
- 代表普通文件
l 代表链接符号(软连接)
b(block) 块设备，一般指存储设备 如硬盘，光盘
c 字符设备，一般指输入串行接口 中介 如键盘，鼠标
s socket，套接字，网络中不同主机上的应用进程之间进行双向通信的端点

第二到第九个字符(文件权限)
共9个字符，三个为一组(-rwx)
rwx简介：
r- read 代表可读权限
w- write 代表可写权限
x- execute 代表可执行权限
- 代表没有权限

三组分别为：
第一组：用户(代表文件的主人)
第二组：用户组(代表家庭，一个群组)
第三组：其他人(既不是主人，也不是家庭成员)
用户简写u(user)，用户组简写g(group)，其他人o(other)，代表全部a(all)包含ugo

第二项
. (点) 代表selinux标识符

第三项
硬链接数

第四项
代表文件/目录的用户(拥有者)

第五项
代表文件/目录的用户组(拥有者的家庭成员)

第六项
文件/目录大小，默认单位是bytes(字节)

第七到第九项
文件/目录时间(创建时间/最近修改时间)

第十项
文件/目录的名称
```

### 十五、权限

```shell
rwx- 数字权限
r---4
w---2
x---1
----0

常见权限
drwx-r-x-r-x 目录默认权限 (755)
drwxrwxrwx   目录最大权限 (777) 
-rw-r--r--   文件默认权限 (644)
-rwx-rwx-rwx-   文件最大权限 (777)

```



#### 15.1.Linux权限管理命令

```
chmod 更改文件或者目录权限
chown 更改文件用户
chgrp 更改文件用户组
```

#### 15.2.chmod 更改文件权限

```powershell
chmod [命令选项] [符号组合|数字权限] 文件/目录
常见选项
	-R 递归
	--reference 根据其他文件权限设置文件
	chmod --reference 其他文件  要设置权限的文件
	u=用户,g=用户组,o=其他用户,a=所有用户

1# 为用户添加x执行权限

2# 删除所有用户的r读权限

3#用户设置权限rwx,用户组设置为rx,其他人为r

4# 设置用户读写权限，用户组读权限，其他人没有任何权限

5# 设置目录权限为755，子目录下递归设置

6# 根据其他文件权限设置文件
```

#### 15.3.chown 更改文件用户和用户组

```powershell
chown [命令选项] 用户:用户组 文件名称
chown [命令选项] 用户 文件名称
chown [命令选项] :用户组 文件名称
常见选项
	-R 递归
	
1# 更改文件用户

2# 更改文件用户和用户组
```

#### 15.4.chgrp 更改文件用户组

```powershell
chgrp [命令选项] 用户组 文件/目录
用的不多，了解即可
常见选项
	-R 递归
```

***

#### 15.5.chattr,lsattr

```sh
chattr
这个命令的主要作用就是配合i参数一起使用，这样的效果让文件内容不能被更改
一般对系统的重要文件，比如/etc/passwd,/etc/shadow这些和用户账号密码的文件，就需要使用这个命令来保证系统用户的安全性
chattr +i /etc/passwd
这样就可以让/etc/passwd的内容，无法被更改
如果想取消就可以使用
chattr -i /etc/passwd
这样/etc/passwd文件的内容就又可以被更改了
lsattr命令就是查看文件是否被chattr命令设置了权限
lsattr /etc/passwd
```

***

### 十六.读写执行权限对文件/目录

#### 16.1.读写执行对于文件

```
对于文件权限
读权限(r)，代表可以查看文件中的内容。命令举例 cat less more
写权限(w)，代表可以修改文件中的内容，如果想要删除这个文件，必须要有这个文件上级目录的w写权限，文件本身的写权限只能针对文件的内容进行增删改，命令举例 vim echo(> >> 重定向)
执行权限(x)，代表文件可以被系统执行，对于文件来说执行权限x 就是最高权限。命令举例 /bin/sh
```

#### 16.2.读写执行对于目录

```
对于目录权限
读权限(r)，代表我们可以查看目录下的内容。命令举例 ls
写权限(w)，代表我们可以在目录下面新增，删除，拷贝等操作，对于目录来讲写权限就是最高权限。命令举例 cp mkdir touch rm mv
执行权限(x),代表我们可以进入目录，命令举例 cd
目录最小权限要给4+1 读跟执行权限，如果进不去目录，无法查看目录下的具体内容
```



### 十七.默认权限与特殊权限

#### 17.1.什么是默认权限

```powershell
Linux下我们创建目录或者文件，都是由umask默认来控制权限的，文件默认最高权限为666，默认系统不会添加x执行权限，文件最高权限为x 可执行权限，目录默认最高权限为777。我们创建的文件默认权限为 -rw--r--r-- = 644 ，创建的目录默认权限为 -rwxr-xr-x = 755
```

#### 17.2.umask命令

```shell
[root@weibo ~]# umask 
0022
[root@weibo ~]# umask -S
u=rwx,g=rx,o=rx

常见选项
	-S 字符显示权限

umask加减计算权限为错误，要字符计算
```

#### 17.3.特殊权限，理解即可

```
Linux默认权限为9位，但是还有三位特殊权限。分别为SUID／SGID／粘滞位
```

#### 17.4.SUID

```powershell
SUID，用大写的S进行表示，位置跟用户位的x权限重合，有x权限就显示为s，没有x权限就显示为S，数字表示为4
作用：
让普通用户，也可以执行特殊命令，例如passwd命令，SUID允许我们突破9位权限，让普通用户再执行某个二进制命令拥有跟root一样的执行权
```

#### 17.5.SGID

```powershell
SGID，用大写的S进行表示，位置跟用户组位的x权限重合，有x权限就显示为s，没有x权限就显示为S，数字表示为2
作用：
让没有在用户组的普通用户，也可以执行特殊命令，例如locate命令，SGID允许我们突破9位权限，让普通用户再执行某个二进制命令拥有跟用户组成员一样的执行权
```

#### 17.6.粘滞位

```powershell
粘滞位，用大写的T进行表示，位置跟其他用户位的x权限重合，有x权限就显示为t，没有x权限就显示为T，数字表示为1
作用：
最常见的就是/tmp 目录，所有用户都可以在里面增删改查，但是不能删除目录本身啊。
[root@weibo ~]# ll -d /tmp/
drwxrwxrwt. 10 root root 4096 Nov 21 00:17 /tmp/
```

|               | SUID                                   | SGID                                          | 粘滞位                                            |
| ------------- | -------------------------------------- | --------------------------------------------- | ------------------------------------------------- |
| 标识字符      | S                                      | S                                             | T                                                 |
| 占据位置      | 用户基本权限x位                        | 用户组基本权限x位                             | 其他用户基本权限x位                               |
| 基本权限有x位 | s                                      | s                                             | t                                                 |
| 数字表示      | 4                                      | 2                                             | 1                                                 |
| 作用举例      | 命令：passwd 对应的文件/usr/bin/passwd | 命令：locate 对应的数据库文件 /usr/bin/locate | /tmp 目录，所有用户都可增删改查，但是不能删除本身 |

***

### 十八.普通用户授权

#### 18.1.sudo基本概念

```
sudo，给普通用户赋予管理员权限，Linux对普通用户限制很严格，很多命令无法直接去使用，sudo可以把指定的命令给到指定的用户，授权赋予的权限越详细，普通用户得到的权限越小，赋予的权限越简单，普通用户权限越大。
```

#### 18.2.sudo的格式

```
在100行
[root@didiyun ~]# visudo      ==  vim /etc/sudoers
root    ALL=(ALL)       ALL
byf     ALL=(ALL)   /usr/bin/rm

第一段 root 用户名或者用户组名称，你打算给谁赋予权限，如果是组前面要加%
第二段 ALL 表示主机，被管理者IP，支撑服务叫NIS，废弃状态，直接写ALL
第三段 (ALL) 代表用户可以切换成什么身份，ALL表示所有 包含root
第四段 ALL 表示赋予用户什么命令，ALL代表所有命令
```

