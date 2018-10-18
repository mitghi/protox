#!/bin/bash

export PRTXBIN=demoe

DEBUG_LEVEL=DEBUG
NOTIME=yes
TRACE_LEVEL=4
CALLER_INFO=no

export RLOG_LOG_LEVEL=$DEBUG_LEVEL
export RLOG_TRACE_LEVEL=$TRACE_LEVEL
export RLOG_LOG_NOTIME=$NOTIME
export RLOG_CALLER_INFO=$CALLER_INFO
export RLOG_LOG_FILE="/tmp/demoe.log"

# uncomment to redirect logs
# ./$PRTXBIN 2>/dev/null

./$PRTXBIN 
