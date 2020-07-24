#!/usr/bin/env bash
docker stop playlottery && docker rm playlottery;
docker rmi leo2n/playlottery:test;
docker build -t leo2n/playlottery:test . ;
docker run -d --name playlottery --network=playLottery leo2n/playlottery:test ;