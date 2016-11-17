package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func createProcess(t *testing.T) *exec.Cmd {
	cmd := exec.Command("sleep", "666")
	err := cmd.Start()

	if err != nil {
		t.Error(err)
		t.FailNow()
		return nil
	}

	return cmd
}

func testKill(t *testing.T, safe bool, cmds []*exec.Cmd) {
	pids := make([]int, 0, len(cmds))

	for i := 0; i < len(cmds); i++ {
		pids = append(pids, cmds[i].Process.Pid)
	}

	if errs := kill(pids, safe); len(errs) > 0 {
		for pid, err := range errs {
			t.Errorf("pid %d: error: %s", pid, err)
		}

		return
	}

	terminated := make(chan struct{})

	go func() {
		for i := 0; i < len(cmds); i++ {
			cmds[i].Wait()
		}

		terminated <- struct{}{}
	}()

	select {
	case <-time.After(time.Second):
		t.Errorf("Some processes still running")
		return
	case <-terminated:
	}

	for _, cmd := range cmds {
		if state := cmd.ProcessState; state.Success() {
			t.Errorf("Process finished successfully (wasn't killed): %v", cmd.Process.Pid, state)
			return
		}
	}

}

func TestKillUnsafe(t *testing.T) {
	for i := 0; i < 10; i++ {
		cmd := createProcess(t)

		testName := fmt.Sprintf("kill %s", strconv.Itoa(cmd.Process.Pid))

		t.Run(testName, func(t *testing.T) {
			testKill(t, false, []*exec.Cmd{cmd})
		})
	}

	rand.Seed(666)

	for i := 0; i < 5; i++ {
		npids := rand.Intn(10)
		cmds := make([]*exec.Cmd, 0, npids)

		for j := 0; j < npids; j++ {
			cmds = append(cmds, createProcess(t))
		}

		testName := "kill"

		for j := 0; j < npids; j++ {
			testName += " " + strconv.Itoa(cmds[j].Process.Pid)
		}

		t.Run(testName, func(t *testing.T) {
			testKill(t, false, cmds)
		})
	}
}
