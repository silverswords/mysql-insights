### mysql-insights

In-depth study of mysql, mainly including master-slave, innodb, isolation level, benchmark, split database, split table, etc.

### Table of Contents

 - [master-slave](https://github.com/silverswords/mysql-insights#master-slave)
 - [benchmark](https://github.com/silverswords/mysql-insights#benchmark)

### [master-slave](https://github.com/silverswords/mysql-insights/tree/master/env)

    下载 env 文件，选择自己想要配置的主从形式，在 docker-compose.yml 文件存在的目录下运行 docker-compose up -d 一键启动 mysql 集群

## [benchmark](https://github.com/silverswords/mysql-insights/blob/master/benchmark)

#### 使用工具
[sysbench](https://github.com/akopytov/sysbench)
#### 脚本位置
    一般位于 /usr/share/sysbench 目录下，如果不在，使用 find / -name oltp* 查找即可
#### 工具使用
 - 数据准备

sysbench ./oltp_read_write.lua --mysql-host=192.168.0.252 --mysql-port=3306  --mysql-user=root --mysql-password=123456  --mysql-db=sakura  --tables=1 --table-size=100000 --threads=10 --events=100000 prepare
 - 开始测试

sysbench ./oltp_read_write.lua --mysql-host=192.168.0.252 --mysql-port=3006  --mysql-user=root --mysql-password=123456  --mysql-db=sakura  --tables=1 --table-size=100000 --threads=10 --events=100000 run
 - 清楚数据

sysbench ./oltp_read_write.lua --mysql-host=192.168.0.252 --mysql-port=3306  --mysql-user=root --mysql-password=123456  --mysql-db=sakura  --tables=1 --table-size=100000 --threads=10 --events=100000 cleanup
#### 参数详解
tables: 测试表格数量
table-size: 每个表格插入数据行数
threads: 并发数量，即 client 连接数量
events: 最大请求数量，即最大事务数量


