#!/bin/bash

if [ -z $GOPATH ]; then
    echo "FAIL: GOPATH environment variable is not set"
    exit 1
fi

if [ -n "$(go version | grep 'darwin/amd64')" ]; then    
    GOOS="darwin_amd64"
elif [ -n "$(go version | grep 'linux/amd64')" ]; then
    GOOS="linux_amd64"
else	
    echo "FAIL: only 64-bit Mac OS X and Linux operating systems are supported"
    exit 1
fi

# Build the test binary to use to test the student's libstore implementation.
# Exit immediately if there was a compile-time error.
go install github.com/cmu440/tests/paxostest
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi
go install github.com/cmu440/runners
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

# Pick random port between [10000, 20000).
PAXOS_PORT=$(((RANDOM % 10000) + 10000))
PAXOS_SERVER=$GOPATH/bin/runners
PAXOS_TEST=$GOPATH/bin/paxostest

function startPaxosServers {
    N=${#PAXOS_ID[@]}
    # Start master paxos server.
    ${PAXOS_SERVER} -N=${N} -id=${PAXOS_ID[0]} -port=${PAXOS_PORT} &> /dev/null &
    PAXOS_SERVER_PID[0]=$!
    # Start slave paxos servers.
    if [ "$N" -gt 1 ]
    then
        for i in `seq 1 $((N - 1))`
        do
	    PAXOS_SLAVE_PORT=$(((RANDOM % 10000) + 10000))
            ${PAXOS_SERVER} -port=${PAXOS_SLAVE_PORT} -id=${PAXOS_ID[$i]} -N=${N} -master="localhost:${PAXOS_PORT}" &> /dev/null &
            PAXOS_SERVER_PID[$i]=$!
        done
    fi
    sleep 2
    ${PAXOS_TEST} -master="localhost:${PAXOS_PORT}" &> /dev/null &
}

function stopPaxosServers {
    N=${#PAXOS_ID[@]}
    for i in `seq 0 $((N - 1))`
    do
        kill -9 ${PAXOS_SERVER_PID[$i]}
        wait ${PAXOS_SERVER_PID[$i]} 2> /dev/null
    done
}


# Testing single client, single tribserver, multiple storageserver.
function startTestThreeNodes {
    echo "Running testPaxosBasic:"
    PAXOS_ID=('0' '1' '2')
    TIMEOUT=15
    startPaxosServers
    sleep 10
    stopPaxosServers
}

startTestThreeNodes

