#!/bin/bash

DEBUG_LEVEL=DEBUG
NOTIME=yes
TRACELEVEL=4

export RLOG_LOG_LEVEL=$DEBUG_LEVEL
export RLOG_TRACE_LEVEL=$TRACELEVEL
export RLOG_LOG_NOTIME=$NOTIME

USR=test
PSWD='$2a$06$lNi8H5kc5Z9T9xJAXwQqyunl2EYhGUi6ct3TgpR1BNb1vpzpp9pzC'
QOS=1
ROUTE="a/simple/demo"
MSG="hello"


./demob -username=$USR -password=$PSWD -qos=$QOS -pr=$ROUTE -pm=$MSG
