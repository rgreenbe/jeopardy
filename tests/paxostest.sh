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

go install -race github.com/cmu440/tests/paxostest
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi
go install -race github.com/cmu440/runners/prunner
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

# Pick random port between [10000, 20000).
PAXOS_PORT=$(((RANDOM % 10000) + 10000))
PAXOS_SERVER=$GOPATH/bin/prunner
PAXOS_TEST=$GOPATH/bin/paxostest

function startPaxosServers {
    N=${#PAXOS_ID[@]}
    # Start master paxos server.
    ${PAXOS_SERVER} -N=${N} -id=${PAXOS_ID[0]} -port=${PAXOS_PORT} &> output.txt &
    PAXOS_SERVER_PID[0]=$!
    # Start slave paxos servers.
    if [ "$N" -gt 1 ]
    then
        for i in `seq 1 $((N - 1))`
        do
	    PAXOS_SLAVE_PORT=$(((RANDOM % 10000) + 10000))
            ${PAXOS_SERVER} -port=${PAXOS_SLAVE_PORT} -id=${PAXOS_ID[$i]} -N=${N} -master="localhost:${PAXOS_PORT}" &> output.txt &
            PAXOS_SERVER_PID[$i]=$!
        done
    fi
    sleep 2
}

function stopPaxosServers {
    N=${#PAXOS_ID[@]}
    for i in `seq 0 $((N - 1))`
    do
        kill -9 ${PAXOS_SERVER_PID[$i]}
        wait ${PAXOS_SERVER_PID[$i]} 2> /dev/null
    done
}

function startPaxosTest {
    ${PAXOS_TEST} -master="localhost:${PAXOS_PORT}" -nodes=${N} &> output.txt &
}

function startPaxosTestDeadNode {
    ${PAXOS_TEST} -master="localhost:${PAXOS_PORT}" -nodes=${N} -type="dead" &> output.txt &
}

function startPaxosTestReplaceNode {
    ${PAXOS_TEST} -master="localhost:${PAXOS_PORT}" -nodes=${N} -type="replace" &> output.txt &
}

# Test with three nodes
function startTestThreeNodes {
    echo "Running paxostest with all nodes:"
    PAXOS_ID=('0' '1' '2')
    TIMEOUT=15
    startPaxosServers
    startPaxosTest
    sleep 5
    stopPaxosServers
}

# Test with three nodes but one dies before paxos starts..this should run to completion
function startTestOneDeadNode {
    echo "Running paxostest with one dead node:"
    PAXOS_ID=('0' '1' '2')
    TIMEOUT=15
    startPaxosServers
    kill -9 ${PAXOS_SERVER_PID[2]}
    startPaxosTestDeadNode
    sleep 5
    stopPaxosServers
}

function startTestReplaceNode {
    echo "Running paxostest with dead node to be replaced:"
    PAXOS_ID=('0' '1' '2')
    TIMEOUT=15
    startPaxosServers
    kill -9 ${PAXOS_SERVER_PID[2]}
    startPaxosTestReplaceNode
    sleep 5
    stopPaxosServers
}

startTestThreeNodes
startTestOneDeadNode
startTestReplaceNode


