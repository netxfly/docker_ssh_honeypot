#!/bin/bash

#Add this script in all servers
#Add this line in snmpd.conf
#exec bashversion /usr/local/libexec/check_bash_version.sh # .1.3.6.1.4.1.2021.8.1.101.1
#Nagios command can be : check_bash_version with just one argument: OID

VER=`bash --version`
STATE_OK=0
STATE_CRITICAL=2


#VER="test"

if [[ $VER =~ "GNU bash (MIKA)" ]]
then
        echo "OK: bash version $BASH_VERSION (MIKA)"
        exit $STATE_OK
else
        echo "CRITICAL: bash version $BASH_VERSION"
        exit $STATE_CRITICAL
fi
