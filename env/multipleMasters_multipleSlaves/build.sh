#!/bin/bash

docker-compose down
docker-compose build
docker-compose up -d

until docker exec mysql_master_1 sh -c 'mysql -u root -p123456'
do
    echo "Waiting for mysql_master database connection..."
    sleep 4
done
until docker exec mysql_master_2 sh -c 'mysql -u root -p123456'
do
    echo "Waiting for mysql_master database connection..."
    sleep 4
done

login_1=`docker exec mysql_master_1 sh -c 'mysql -u root -p123456 -e "SHOW MASTER STATUS;"'`
log_file_1=`echo $login_1 | awk '{print $6}'`
log_pos_1=`echo $login_1 | awk '{print $7}'`

login_2=`docker exec mysql_master_2 sh -c 'mysql -u root -p123456 -e "SHOW MASTER STATUS;"'`
log_file_2=`echo $login_2 | awk '{print $6}'`
log_pos_2=`echo $login_2 | awk '{print $7}'`

#mysql_master_1 change master
change_master_1="stop slave; change master to master_host='192.168.0.251',master_port=3309,master_user='root',master_password='123456',master_log_file='$log_file_2',master_log_pos=$log_pos_2; start slave;"
mysql_master_1_login='mysql -u root -p123456 -e "'
mysql_master_1_login+="$change_master_1"
mysql_master_1_login+='"'
docker exec mysql_master_1 sh -c "$mysql_master_1_login"

#mysql_master_2 change master
change_master_2="stop slave; change master to master_host='192.168.0.251',master_port=3308,master_user='root',master_password='123456',master_log_file='$log_file_1',master_log_pos=$log_pos_1; start slave;"
mysql_master_2_login='mysql -u root -p123456 -e "'
mysql_master_2_login+="$change_master_2"
mysql_master_2_login+='"'
docker exec mysql_master_2 sh -c "$mysql_master_2_login"

#mysql_slave_1 change master
change_slave_1="stop slave; change master to master_host='192.168.0.251',master_port=3308,master_user='root',master_password='123456',master_log_file='$log_file_1',master_log_pos=$log_pos_1; start slave;"
mysql_slave_1_login='mysql -u root -p123456 -e "'
mysql_slave_1_login+="$change_slave_1"
mysql_slave_1_login+='"'
docker exec mysql_slave_1 sh -c "$mysql_slave_1_login"

#mysql_slave_2 change master
change_slave_2="stop slave; change master to master_host='192.168.0.251',master_port=3309,master_user='root',master_password='123456',master_log_file='$log_file_2',master_log_pos=$log_pos_2; start slave;"
mysql_slave_2_login='mysql -u root -p123456 -e "'
mysql_slave_2_login+="$change_slave_2"
mysql_slave_2_login+='"'
docker exec mysql_slave_2 sh -c "$mysql_slave_2_login"

docker exec mysql_master_1 sh -c "mysql -u root -p123456 -e 'SHOW SLAVE STATUS \G'"
docker exec mysql_master_2 sh -c "mysql -u root -p123456 -e 'SHOW SLAVE STATUS \G'"
docker exec mysql_slave_1 sh -c "mysql -u root -p123456 -e 'SHOW SLAVE STATUS \G'"
docker exec mysql_slave_2 sh -c "mysql -u root -p123456 -e 'SHOW SLAVE STATUS \G'"