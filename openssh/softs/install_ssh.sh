OPENSSH=/openssh2
mkdir -p /openssh2/dist/
cd ${OPENSSH}
wget http://zlib.net/zlib-1.2.11.tar.gz
tar xvfz zlib-1.2.11.tar.gz
cd zlib-1.2.11
./configure --prefix=${OPENSSH}/dist/ && make && make install
cd ${OPENSSH}
wget https://www.openssl.org/source/openssl-1.0.1e.tar.gz --no-check-certificate
tar xvfz openssl-1.0.1e.tar.gz
cd openssl-1.0.1e
./config --prefix=${OPENSSH}/dist/ && make && make install
cd ${OPENSSH}
wget https://ftp.eu.openbsd.org/pub/OpenBSD/OpenSSH/portable/openssh-6.2p1.tar.gz --no-check-certificate
tar xvfz openssh-6.2p1.tar.gz
patch openssh-6.2p1/auth-passwd.c /softs/auth-password.patch 
cd openssh-6.2p1
# patch auth-password.c auth-passwd.patch
./configure --prefix=${OPENSSH}/dist/ --with-zlib=${OPENSSH}/dist --with-ssl-dir=${OPENSSH}/dist/ && make && make install
# mv -f ${OPENSSH}/dist/sbin/sshd /usr/sbin/sshd