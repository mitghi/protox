#!/bin/bash

DEBUG_LEVEL=DEBUG
NOTIME=yes
TRACELEVEL=4

export RLOG_LOG_LEVEL=$DEBUG_LEVEL
export RLOG_TRACE_LEVEL=$TRACELEVEL
export RLOG_LOG_NOTIME=$NOTIME
export RLOG_LOG_FILE="/tmp/cli.log"

USR=test4
PSWD='$2a$06$9wavlAtmNZ66Whe2wturDO7yIBdE41/Zcn4c5z4ydzJ/ydVJIZwJK'
QOS=0
ROUTE="a/simple/demo"

./cli -username=$USR -password=$PSWD -qos=$QOS -sr=$ROUTE 2>/dev/null
