package util

import (
	"bytes"
	"os/exec"

	"github.com/lspaccatrosi16/lbt/lib/log"
)

func RunCmd(eCmd *exec.Cmd, stdout, stderr bytes.Buffer, ml *log.Logger, dir string) bool {
	eCmd.Dir = dir
	eCmd.Stdout = &stdout
	eCmd.Stderr = &stderr

	err := eCmd.Run()
	if err != nil {
		ml.Logln(log.Error, err.Error())
		ml.Logln(log.Error, stdout.String())
		ml.Logln(log.Error, stderr.String())
		return false
	}

	stdout.Reset()
	stderr.Reset()

	return true
}
