package api

import (
	"os"
	"testing"

	"github.com/op/go-logging"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/state"
)

func TestMain(m *testing.M) {
	replaceGetNodeState()
	global.Log = logging.MustGetLogger("GoRaftTest")
	global.TimeoutChannel = make(chan bool, 1)
	go RunServer(8000)

	global.Log.Debug("Tests for api have been set up.")
	os.Exit(m.Run())
}

func replaceGetNodeState() {
	state.Node = state.NewNodeStateTestImpl()
}
