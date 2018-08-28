#!/bin/sh
useradd -s /bin/bash test
#echo "test:123456" | chpasswd

service rsyslog start

echo "Starting SSH daemon..."
mkdir -p /var/run/sshd
exec /openssh2/dist/sbin/sshd -f /openssh2/dist/etc/sshd_config -D

# rm -fr /softs/bash-4.3.30-active-syslog/
