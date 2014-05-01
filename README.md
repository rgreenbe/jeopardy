Hi there!

Playing our game is very simple. We have made a few different scrips so that you can 
both test our paxos implementation and our Jeopardy! implementation. The file in the path
paxostest.sh will, upon execution, run paxostest.go with 5 different scenarios:

1) A regular run of paxos where there are 3 nodes
2) A run of paxos where there is a dead node
3) A run of paxos where odd numbered nodes (node 1) drop all odd numbered command slots
4) A run of paxos where the first 50 command slot commits are dropped by odd numbered nodes (node 1)
5) A run of paxos where a node is killed, the test is started, a node is added in and caught up, and
	then paxos is resumed. After this a handful more proposals are made to ensure it restarts properly

In the runner directory, there are two different shell files. The first just starts up all of the paxos node
and a client, and the second does the same while killing one of the paxos processes. This is so we can demonstrate
that our Paxos implementation handles a node failure with Jeopardy (as was requested during our review). 
Note that the non-master node is killed because our client doesn't have logic to fail over to a non-master
node when the master dies.

Upon running one of the runner scripts for the game, running the java main file will allow you to start running
our game. If you start 3 instances of it and join in each of them, you can begin to play. We did briefly 
verify that our paxos/Jeopardy! implementation works over a local network. Due to issues with port forwarding
and firewalls we were not successfull in getting it to work over the Internet, but the fact that it works
over the local network implies our implementation works over the internet when proper IP addresses are supplied. 

Please enjoy!

Jeopardy!
========
