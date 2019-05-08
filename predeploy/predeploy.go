package predeploy

import (
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/error_check"
	"os"
	"os/exec"
)

func Execute() {
	cmd := exec.Command(conf.GetConfig().GetString("pre_deploy"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	error_check.ExitOn(err, "error executing prebuild step")
}