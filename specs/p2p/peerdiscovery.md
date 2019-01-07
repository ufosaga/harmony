# Peer Discovery

## Introduction
For any decentralized blockchain network, peer discovery is the essential step to create a real p2p mesh network.
New node uses peer discovery to join the blockchain network.
In Harmony network, the new node will also need to find the right shard to join.
The proposed peer discovery process works as illustrated below.

## Two overlays of p2p Network
Each node in the Harmony network will have two overlays of p2p networks.
Layer one is at the global level used for inter-shard communication.
Layer two is for intra-shard communications.
In each layer, the node can have up to 16 outgoing links.
Each node will have up to 128 incoming links.
The 16/128 numbers are chosen towards the goal to maintain a well-connected p2p network and keep reasonable resource consumption.

## Bootnodes
New nodes contact bootnodes to get a list of 32 random nodes at the global level.
The IP addresses of the bootnodes are hard-coded in the node software.
Alternatively, the new node may use DNSseed (ex. bootnode.harmony.one) to find a list of long running bootnodes if the hard-coded bootnodes are not accepting connections due to network congestion.

## Register
To prevent Sybil and DoS attacks, new nodes need to solve a PoW puzzle in order to join the p2p network.
The difficulty of the puzzle can be adjusted to within 1 minute on a typical server machine (4 cores, 2.7 GHz).

Bootnodes use the Kademlia based DHT to return a list of peer nodes to the new node.
The Kademlia based DHT calculates the XOR distances of the public keys of the nodes to keep nodes into different buckets.

For a full node that will participate in the consensus, it needs to gossip the network via the peers returned by bootnodes for Proof-of-Stake verification.
The first message is to submit a transaction of deposit of a certain amount of Harmony tokens as the stake to join the network as validators.
Beacon chain will be responsible for validating the stakes and the transaction.

For a light node that will not join in the consensus, it has no requirement of the PoS verification.

Beacon chain will assign the shard id number to the new node via gossip as well.
New node’s public key and deposit will be recorded in the beacon chain blocks.

## Shard Discovery
New node received the shard info from beacon chain.
It then broadcast/gossip the shard info together with its public key to the network using the shard id as a topic.
The peers within the same shard as the new node will respond to the topic based on the XOR distance in the DHT.
New node can have up to 16 peers within the same shard connected.
We are investigating the rendezvous feature of libp2p to implement the layer 2 p2p network.

## LibP2P
The new node are connected to the two layers of the p2p networks.
Harmony utilizes libp2p as the underlying networking layer for peer discovery and gossiping.
It is still a crucial task to understand the protocol messages and how the libp2p handles the messages.
