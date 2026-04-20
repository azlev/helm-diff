package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"k8s.io/client-go/util/homedir"
)

const (
	helm2TestSuccessHook = "test-success"
	helm3TestHook        = "test"
)

var (
	// DefaultHelmHome to hold default home path of .helm dir
	DefaultHelmHome = filepath.Join(homedir.HomeDir(), ".helm")
)

func isDebug() bool {
	return os.Getenv("HELM_DEBUG") == "true"
}
func debugPrint(format string, a ...interface{}) {
	if isDebug() {
		fmt.Printf(format+"\n", a...)
	}
}

func outputWithRichError(cmd *exec.Cmd) ([]byte, error) {
	debugPrint("Executing %s", strings.Join(cmd.Args, " "))
	fmt.Fprintf(os.Stderr, "[v3.8.2-debug] CMD: %s\n", strings.Join(cmd.Args, " "))
	for _, env := range cmd.Env {
		if strings.HasPrefix(env, "HELM_") {
			fmt.Fprintf(os.Stderr, "[v3.8.2-debug] ENV: %s\n", env)
		}
	}
	if cmd.Env == nil {
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "HELM_") {
				fmt.Fprintf(os.Stderr, "[v3.8.2-debug] INHERITED-ENV: %s\n", env)
			}
		}
	}
	output, err := cmd.Output()
	if exitError, ok := err.(*exec.ExitError); ok {
		fmt.Fprintf(os.Stderr, "[v3.8.2-debug] CMD-FAILED: %s stderr=%s\n", exitError.Error(), string(exitError.Stderr))
		return output, fmt.Errorf("%s: %s", exitError.Error(), string(exitError.Stderr))
	}
	return output, err
}
