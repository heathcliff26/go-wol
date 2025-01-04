package wol

import (
	"os"
	"os/exec"
	"testing"
)

func TestCMD(t *testing.T) {
	tMatrix := []struct {
		Name, Broadcast, MAC string
		ExitWithError        bool
		NoMac                bool
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
			Name:          "MissingMAC",
			NoMac:         true,
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
			Broadcast:     "not-a-ip",
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
				if !tCase.NoMac {
					args = append(args, tCase.MAC)
				}
				cmd.SetArgs(args)

				err := cmd.Execute()
				if err != nil {
					t.Logf("Execute failed: %v", err)
					os.Exit(2)
				}

				// Should not reach here, ensure exit with 0 if it does
				os.Exit(0)
			}
			execExitTest(t, "TestCMD/"+tCase.Name, tCase.ExitWithError)
		})
	}
}

func execExitTest(t *testing.T, test string, exitsError bool) {
	cmd := exec.Command(os.Args[0], "-test.run="+test)
	cmd.Env = append(os.Environ(), "RUN_CRASH_TEST=1")
	err := cmd.Run()
	if exitsError && err == nil {
		t.Fatal("Process exited without error")
	} else if !exitsError && err == nil {
		return
	}
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
