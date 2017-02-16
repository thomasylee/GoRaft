package main

import (
	"strconv"
	"time"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/rpc"
	"github.com/thomasylee/GoRaft/state"
)

/**
 * The main() method is the initial method that gets run, so it will start
 * the necessary goroutines to become a functional Raft node.
 */
func main() {
	global.SetUpLogger()
	global.Log.Info("GoRaft starting... Logger initialized.")

	global.LoadConfig()

	global.Log.Info("Loaded config:", global.Config)
	global.SetLogLevel(global.Config.LogLevel)

	// Check if state was loaded correctly from previous run.
	global.Log.Debug(state.GetNodeState())

	runNode()
}

/**
 * Runs the infinite loop that keeps the node active.
 */
func runNode() {
	go rpc.RunServer(strconv.Itoa(int(global.Config.Nodes[global.Config.NodeId].ApiPort)))

	// Randomize the election timeout to minimize the risk of two nodes
	// initiating an election at the same time.
	electionTimeout := global.Config.ElectionTimeout
	electionTimeoutJitter := global.Config.ElectionTimeoutJitter
	for {
		timeout := global.GenerateTimeout(electionTimeout, electionTimeoutJitter)
		select {
		case <-global.TimeoutChannel:
			// Do nothing since we didn't time out.
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			// TODO: Start leader election process.
		}
	}
}
