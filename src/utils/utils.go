package utils

import "os/exec"

func CommonRun(name string, args []string) error {
	return exec.Command(name, args...).Run()
}
