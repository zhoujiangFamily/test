#!/bin/bash
declare SERVICE=first-test
function cleanup()
{
        local pids=`jobs -p`
        if [[ "$pids" != "" ]]; then
                kill $pids >/dev/null 2>/dev/null
        fi
}

trap cleanup EXIT
cd /$SERVICE
./$SERVICE >> /var/log/go_log/$SERVICE.out 2>&1
