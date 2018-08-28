FROM debian:jessie

RUN export DEBIAN_FRONTEND='noninteractive' && \
    apt-get update -qq && \
    apt-get install -qqy --no-install-recommends openssh-server rsyslog wget patch make gcc curl libc6-dev net-tools vim && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* && \
    mkdir -p /softs/bash-4.3.30-active-syslog/ && mkdir -p /home/test

COPY softs/ /softs/

RUN chmod +x /softs/bash-4.3.30-active-syslog/install-bash-syslog.sh
RUN chmod +x /softs/install_ssh.sh

WORKDIR /softs/bash-4.3.30-active-syslog/
RUN ["/bin/bash", "./install-bash-syslog.sh"]

WORKDIR /softs/
RUN ["/bin/bash", "./install_ssh.sh"]

EXPOSE 22/tcp

COPY entrypoint.sh /etc/systemd/system/

COPY Shanghai /etc/localtime

RUN chmod +x /etc/systemd/system/entrypoint.sh

ENTRYPOINT ["/etc/systemd/system/entrypoint.sh"]

# docker run -itd  -p 22:22 -v $CUR_DIR/logs/openssh/auth.log:/var/log/auth.log  -v $CUR_DIR/logs/openssh/syslog.log:/var/log/syslog.log openssh:v0.2
