package server

import (
	"os"
	"os/exec"
	"testing"
)

func TestServerCMD(t *testing.T) {
	tMatrix := []struct {
		Name, Config, LogLevel string
	}{
		{
			Name:   "MissingConfig",
			Config: "not-a-file",
		},
		{
			Name:     "InvalidLogLevel",
			LogLevel: "invalid",
		},
		{
			Name:   "InvalidStorageBackend",
			Config: "testdata/invalid-storage.yaml",
		},
	}

	t.Cleanup(func() {
		_ = os.Remove("hosts.yaml")
	})

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			if os.Getenv("RUN_CRASH_TEST") == "1" {
				cmd := NewCommand()

				args := make([]string, 0, 3)
				if tCase.Config != "" {
					args = append(args, "--"+flagNameConfig, tCase.Config)
				}
				if tCase.LogLevel != "" {
					args = append(args, "--"+flagNameLogLevel, tCase.LogLevel)
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
			execExitTest(t, "TestServerCMD/"+tCase.Name, true)
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
