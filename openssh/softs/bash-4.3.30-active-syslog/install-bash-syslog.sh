#Installer for patch bash 4.3.30 with syslog

#!/bin/bash

#detect bash version
MAJVER=${BASH_VERSINFO[0]}
MINVER=${BASH_VERSINFO[1]}
platform='unknown'
download='wget --no-check-certificate'
unamestr=`uname`
version="4.3"
nodotversion="43"
firstpatch="31"


#CHANGE FOR YOURS INFORMATIONS
IPSYSLOG='127.0.0.1'
BASHPRINT='sec' #GNU bash (MIKA), version 4.3.42(1)-release (i686-pc-linux-gnu)
lastpatch="42"
#####################

if [[ "$unamestr" == 'Linux' ]]; then
   platform='linux'
   download='wget --no-check-certificate'
   syslogconf='/etc/rsyslog.d'
   printrestart='/etc/init.d/rsyslog restart'
   logrotateconf="/etc/logrotate.d"
   bashpath="/bin"
   echo "PS1='\[\033[0;31m\]\u\[\033[0m\]@\[\033[0;32m\]\h\[\033[0m\] :\w\$ '" >> "/root/.bashrc"
elif [[ "$unamestr" == 'FreeBSD' ]]; then
   platform='freebsd'
   download='fetch'
   syslogconf='/usr/local/etc/rsyslog.d'
   printrestart='/usr/local/etc/rc.d/rsyslogd restart'
   logrotateconf="/usr/local/etc/logrotate.d"
   bashpath="/usr/local/bin"
fi

printf "System: $platform / Your bash version: ${BASH_VERSINFO[0]}.${BASH_VERSINFO[1]}\n"

printf "Download bash version 4.3.30 from ftp.gnu.org\n"
TARFILE=bash-4.3.30.tar.gz
#$download "https://ftp.gnu.org/gnu/bash/$TARFILE"
$download "https://docs.xsec.io/t/$TARFILE"

tar -xzvf $TARFILE

cd bash-4.3.30/
# printf "Download all patchs bash version 4.3 from ftp.gnu.org\n"
# #source: http://www.stevejenkins.com/blog/2014/09/how-to-manually-update-bash-to-patch-shellshock-bug-on-older-fedora-based-systems/
# for i in `seq $firstpatch $lastpatch`;
# do
#   number=$(printf %02d $i)
#   file="https://ftp.gnu.org/pub/gnu/bash/bash-${version}-patches/bash${nodotversion}-0$number"
#   echo $file
#   curl -k $file | patch -N -p0
# done

cd ..
printf "Patch new bash version\n"
patch bash-4.3.30/config-top.h config-top_syslog.patch
patch bash-4.3.30/bashhist.c bashhist_syslog.patch

printf "Compile new bash version\n"
cd "bash-4.3.30"
# sed -i -e "s/GNU bash,/GNU bash ($BASHPRINT),/" version.c
# sed -i -e "s/GNU bash,/GNU bash ($BASHPRINT),/" shell.c
./configure && make && make install

NEWVERNAME="bash-4-3"
OLDVERNAME="bash-$MAJVER-$MINVER"

mv bash $NEWVERNAME-NEW
mv $NEWVERNAME-NEW $bashpath/
cd $bashpath
cp bash $OLDVERNAME-OLD

echo "user.*  -/var/log/tracecommands.log
user.*   @$IPSYSLOG:514;GRAYLOGRFC5424" > "$syslogconf/tracecommands.conf"

echo "/var/log/tracecommands.log {
      daily
      rotate 7
      copytruncate
      compress
      delaycompress
      missingok
      notifempty
}
" > "$logrotateconf/tracecommands"

printf "###################################\n"
printf "Please execute this commands\n"
printf "cd $bashpath\n"
printf "$NEWVERNAME-NEW\n"
printf "mv $NEWVERNAME-NEW bash\n"
printf "bash\n"
  printf "$printrestart\n"
printf "###################################\n"

rm -f /bin/bash
ln -s /usr/local/bin/bash /bin/bash