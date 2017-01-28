package api

import (
	"os"
	"testing"

	"github.com/thomasylee/GoRaft/global"
	"github.com/thomasylee/GoRaft/state"
)

func TestMain(m *testing.M) {
	state.Node = state.NewNodeStateTestImpl()

	global.SetUpLogger()
	global.SetLogLevel("critical")

	global.TimeoutChannel = make(chan bool, 1)
	go RunServer(8000)

	global.Log.Debug("Tests for api have been set up.")
	os.Exit(m.Run())
}
