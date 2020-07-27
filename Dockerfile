FROM centos:latest

ENV MYPATH /usr/local
WORKDIR $MYPATH/play
ENV PATH=$PATH:/usr/local/play
# 一定要注意时区啊亲, 这个是个大坑, 如果是UTC时区, 那么和我们的时间范围对不上, 就会出错!
ENV TZ=Asia/Shanghai
# all thing copy to WORKDIR
COPY ./mysqlConfig.json ./
COPY ./timingGetData ./
COPY ./lucky ./
COPY ./botMsg ./
COPY runInDocker.sh ./

ENTRYPOINT ["/bin/bash", "runInDocker.sh"]
