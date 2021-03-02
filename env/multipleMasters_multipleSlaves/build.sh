#!/bin/bash

docker-compose build
docker-compose up -d

docker exec -it mysql_master_1 /bin/bash
mysql -uroot -p123456
