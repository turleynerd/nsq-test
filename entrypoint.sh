#!/bin/bash
# Lets handle the shutdown gracefully provide time to shutdown other services.

# Max ammont of time to wait before shutting down, this should be less than your ecs timeout
maxtimeout=60

# Get message count from nsqd
function messages() {
    curl -k https://127.0.0.1:4152/stats\?format\=json 2>/dev/null | jq -r '.topics[].channels[].depth' | awk '{s+=$1} END {print s}'
}

function shutdown() {
    # Get message count
    q=`messages`
    # Counter for how long we have ran
    i="0"
    # Are there messages to process?
    while [ $q -gt 1 ]; do
        i=$[$i+1]
        if [ $i -gt $maxtimeout ]; then
            # exit loop if we exceed our max timeout
            break
        else
            # Lets wait while the messages drain
            sleep 1
            q=`messages`
            echo "waiting $i/$maxtimeout seconds for $q messages"
        fi
    done
   
    echo "Shutdown request recieved. Sending SIGTERM and waiting for close."
    # Get stats just before shutdown to count messages we may lose
    curl -k https://127.0.0.1:4152/stats?format=json 2>/dev/null
    # SIGTERM
    kill $pid
    # Wait for graceful shutdown of child proc
    wait $pid
    exit
}
# The magic - Intercet SIG events and handle them with our function.
trap shutdown SIGTERM SIGINT

# Run command and wait for exit 
exec "$@" &
pid="$pid $!"
wait $pid
