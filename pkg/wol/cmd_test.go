package wol

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestCMD(t *testing.T) {
	tMatrix := []struct {
		Name, Broadcast, MAC string
		ExitWithError        bool
	}{
		{
			Name: "MacOnly",
			MAC:  "ff:ff:ff:ff:ff:ff",
		},
		{
			Name:          "InvalidMAC",
			MAC:           "not-a-mac",
			ExitWithError: true,
		},
		{
			Name:          "EmptyMAC",
			MAC:           "",
			ExitWithError: true,
		},
		{
			Name:      "BroadcastAddress",
			MAC:       "ff:ff:ff:ff:ff:ff",
			Broadcast: "127.0.0.1",
		},
		{
			Name:          "InvalidBroadcastAddress",
			MAC:           "ff:ff:ff:ff:ff:ff",
			Broadcast:     "not-an-ip",
			ExitWithError: true,
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			if os.Getenv("RUN_CRASH_TEST") == "1" {
				cmd := NewCommand()

				args := make([]string, 0, 3)
				if tCase.Broadcast != "" {
					args = append(args, "--"+flagNameBroadcastAddress, tCase.Broadcast)
				}
				args = append(args, tCase.MAC)
				cmd.SetArgs(args)

				err := cmd.Execute()
				if err != nil {
					fmt.Printf("Execute failed: %v\n", err)
					os.Exit(0) // Exit without error here so the test fails
				}

				// Should not reach here, ensure exit with 0 if it does
				os.Exit(0)
			}
			execExitTest(t, "TestCMD/"+tCase.Name, tCase.ExitWithError)
		})
	}
}

func execExitTest(t *testing.T, test string, exitsError bool) {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run="+test)
	cmd.Env = append(os.Environ(), "RUN_CRASH_TEST=1")
	buf, err := cmd.Output()
	if exitsError && err == nil {
		t.Log(string(buf))
		t.Fatal("Process exited without error")
	} else if !exitsError && err == nil {
		return
	}
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Log(string(buf))
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
