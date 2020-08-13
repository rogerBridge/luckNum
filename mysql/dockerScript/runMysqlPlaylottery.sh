#!/usr/bin/env bash
docker stop mysqlplaylottery && docker rm mysqlplaylottery;
docker network create playLottery;
docker rmi leo2n/mysqlplaylottery:test;
docker build -t leo2n/mysqlplaylottery:test .;
docker run -d --name mysqlplaylottery \
  --restart=always -p 3300:3306 \
  -v $HOME/docker_container/mysqlplaylottery/conf.d:/etc/mysql/conf.d \
  -v $HOME/docker_container/mysqlplaylottery/data:/var/lib/mysql \
  -v $PWD/initScripts:/docker-entrypoint-initdb.d \
  --network=playLottery \
  --network-alias=playlotterymysql \
  leo2n/mysqlplaylottery:test;