#!/bin/bash

docker-compose down
docker-compose build
docker-compose up -d

until docker exec mysql_master sh -c 'mysql -u root -p123456'
do
    echo "Waiting for mysql_master database connection..."
    sleep 4
done

login=`docker exec mysql_master sh -c 'mysql -u root -p123456 -e "SHOW MASTER STATUS;"'`
log_file=`echo $login | awk '{print $6}'`
log_pos=`echo $login | awk '{print $7}'`

change_master="stop slave; change master to master_host='192.168.0.254',master_port=3308,master_user='root',master_password='123456',master_log_file='$log_file',master_log_pos=$log_pos; start slave;"
mysql_login='mysql -u root -p123456 -e "'
mysql_login+="$change_master"
mysql_login+='"'
docker exec mysql_slave sh -c "$mysql_login"

docker exec mysql_slave sh -c "mysql -u root -p123456 -e 'SHOW SLAVE STATUS \G'"