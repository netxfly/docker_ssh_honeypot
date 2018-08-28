 #!/bin/bash
echo "start to build docker images.\n"

echo "build openssh image"
docker build -t openssh:v0.4 openssh

CUR_DIR=`pwd`

# echo "clean history logs"
rm -f $CUR_DIR/logs/openssh/auth.log
rm -f $CUR_DIR/logs/openssh/syslog

# echo "create log file"
mkdir -p $CUR_DIR/logs/openssh && touch $CUR_DIR/logs/openssh/auth.log && touch $CUR_DIR/logs/openssh/syslog

# chown -R nobody: $CUR_DIR/logs/

echo "stop all docker container"
docker stop $(docker ps -q) 

echo "start to run docker container.\n"

echo "start openssh"
docker run -itd -p 22:22 -h server02 -v $CUR_DIR/logs/openssh/auth.log:/var/log/auth.log  -v $CUR_DIR/logs/openssh/syslog:/var/log/syslog openssh:v0.4

echo "set fs.inotify parameter"
sysctl -w fs.inotify.max_user_watches=100000

echo "start sandbox log"
cd sandbox && chmod a+x ./main  && ./main
