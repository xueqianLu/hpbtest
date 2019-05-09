#!/bin/bash

TESTPATH=/home/hpb/luxqtest/data
TXPENDING=""
TXQUEUED=""
function txpoolstatus()
{
	cd $TESTPATH
	status=`./ghpb attach http://127.0.0.1:18545 <<EOM
txpool.status
EOM`
	TXPENDING=`echo $status | grep -Eo 'pending: ([0-9]+)' | awk '{print $2}'`
	TXQUEUED=`echo $status | grep -Eo 'queued: ([0-9]+)' | awk '{print $2}'`
	if [ "$TXPENDING" == "$1" ] && [ "$TXQUEUED" == "0" ] ;
	then
		return 1
	else
		return 0
	fi
}

## -p num   wait until txpending==num and txqueued==0 and display timestamp
## -l       display txstats per second
pending=0
display=false
while getopts "hp:l" arg 
do
        case $arg in
             h)
		echo "$0 [hlp:]"
		echo "	-l		display txstats per second"
		echo "	-p  num		wait until txpending==num and txqueued==0 and display timestamp"
		exit 0
                ;;
             p)
		pending=$OPTARG
                ;;
             l)
		display=true
                ;;
             ?)  #当有不认识的选项的时候arg为?
	     echo "unkonw argument $arg"
		exit 1
		;;
        esac
done

while true; do
	txpoolstatus $pending
	ret=$?
	timestamp=`date +%m-%d/%H:%M:%S`
	if $display ;then
		echo "$timestamp :"
		echo "		pending: $TXPENDING"
		echo "		queued : $TXQUEUED"
	fi
	if [ "$pending" -gt "0" ] && [ "$ret" == "1" ];
	then
		echo "waitted pending $TXPENDING, time $timestamp"
		exit 0
	fi
	sleep 1
done
