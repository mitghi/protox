#!/bin/bash

DEBUG_LEVEL=INFO
NOTIME=yes
TRACELEVEL=

export RLOG_LOG_LEVEL=$DEBUG_LEVEL
export RLOG_TRACE_LEVEL=$TRACELEVEL
export RLOG_LOG_NOTIME=$NOTIME
export RLOG_LOG_FILE="/tmp/demod.log"

USR=test5
PSWD='$2a$06$9wavlAtmNZ6623re2wturDO7yIBdEewfZcn4c5z4ydzJ/ydVJIZwJK'
QOS=1
ROUTE="a/simple/demo"

./cli -username=$USR -password=$PSWD -qos=$QOS -sr=$ROUTE 2>/dev/null
