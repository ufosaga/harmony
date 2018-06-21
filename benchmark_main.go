package main

import (
	"bufio"
	"flag"
	"fmt"
	"harmony-benchmark/consensus"
	"harmony-benchmark/log"
	"harmony-benchmark/node"
	"harmony-benchmark/p2p"
	"os"
	"strings"
	"time"
)

func getShardId(myIp, myPort string, config *[][]string) string {
	for _, node := range *config {
		ip, port, shardId := node[0], node[1], node[3]
		if ip == myIp && port == myPort {
			return shardId
		}
	}
	return "N/A"
}

func getLeader(myShardId string, config *[][]string) p2p.Peer {
	var leaderPeer p2p.Peer
	for _, node := range *config {
		ip, port, status, shardId := node[0], node[1], node[2], node[3]
		if status == "leader" && myShardId == shardId {
			leaderPeer.Ip = ip
			leaderPeer.Port = port
		}
	}
	return leaderPeer
}

func getPeers(myIp, myPort, myShardId string, config *[][]string) []p2p.Peer {
	var peerList []p2p.Peer
	for _, node := range *config {
		ip, port, status, shardId := node[0], node[1], node[2], node[3]
		if status == "leader" || ip == myIp && port == myPort || myShardId != shardId {
			continue
		}
		peer := p2p.Peer{Port: port, Ip: ip}
		peerList = append(peerList, peer)
	}
	return peerList
}

func readConfigFile(configFile string) [][]string {
	file, _ := os.Open(configFile)
	fscanner := bufio.NewScanner(file)

	result := [][]string{}
	for fscanner.Scan() {
		p := strings.Split(fscanner.Text(), " ")
		result = append(result, p)
	}
	return result
}

func main() {
	// Setup a logger to stdout and log file.
	logFileName := fmt.Sprintf("./%v.log", time.Now().Format("20060102150405"))
	h := log.MultiHandler(
		log.Must.FileHandler(logFileName, log.LogfmtFormat()),
		log.StdoutHandler)
	// In cases where you just want a stdout logger, use the following one instead.
	// h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)

	ip := flag.String("ip", "127.0.0.1", "IP of the node")
	port := flag.String("port", "9000", "port of the node.")
	configFile := flag.String("config_file", "config.txt", "file containing all ip addresses")
	flag.Parse()

	config := readConfigFile(*configFile)
	shardId := getShardId(*ip, *port, &config)
	peers := getPeers(*ip, *port, shardId, &config)
	leader := getLeader(shardId, &config)

	consensus := consensus.NewConsensus(*ip, *port, shardId, peers, leader)

	node := node.NewNode(&consensus)
	// Assign closure functions to the consensus object
	consensus.BlockVerifier = node.VerifyNewBlock
	consensus.OnConsensusDone = node.AddNewBlockToBlockchain

	// Temporary testing code, to be removed.
	node.AddMoreFakeTransactions()

	if consensus.IsLeader {
		// Let consensus run
		go func() {
			consensus.WaitForNewBlock(node.BlockChannel)
		}()
		// Node waiting for consensus readiness to create new block
		go func() {
			node.WaitForConsensusReady(consensus.ReadySignal)
		}()
	} else {
		// Node waiting to add new block to the blockchain
		go func() {
			node.WaitForConsensusReady(consensus.ReadySignal)
		}()
	}

	node.StartServer(*port)
}