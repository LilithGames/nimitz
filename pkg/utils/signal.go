package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/oklog/pkg/group"
)

// HandleSignal just sits and waits for ctrl-C.
func HandleSignal(g *group.Group) {
	wait, terminate := Wait2Terminate()
	g.Add(func() error {
		sig := wait()
		return fmt.Errorf("received signal %s", sig)
	}, func(error) {
		terminate()
	})
}

// Wait2Terminate waits for a manual termination or a user initiated termination IE: Ctrl+Break.
// waitForFunc() will wait indefinitely for a signal.
// terminateFunc() will trigger waitForFunc() to complete immediately.
func Wait2Terminate() (waitForFunc func() os.Signal, terminateFunc func()) {
	// Exit when we see a signal
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	waitForFunc = func() os.Signal {
		return <-terminate
	}
	terminateFunc = func() {
		terminate <- os.Interrupt
		signal.Stop(terminate)
		close(terminate)
	}
	return waitForFunc, terminateFunc
}

//WaitForSignal waits for ctrl-C
func WaitForSignal() os.Signal {
	listener := make(chan os.Signal, 1)
	signal.Notify(listener, syscall.SIGINT, syscall.SIGTERM)
	return <-listener
}

//CreateGroup create a group listening system signal
func CreateGroup() *group.Group {
	var g = group.Group{}
	HandleSignal(&g)
	return &g
}
