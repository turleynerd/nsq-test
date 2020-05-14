#!/bin/bash
shutdown() {
    echo
    echo "Child Proc: Shutdown request found exiting..."
    exit
}
trap shutdown SIGINT SIGTERM

for i in {10..1};
do
    echo "sleeping $i"
    sleep 1
done
echo "done"
