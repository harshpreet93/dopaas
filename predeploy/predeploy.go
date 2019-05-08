package predeploy

import (
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/errorcheck"
	"os"
	"os/exec"
)

func Execute() {
	cmd := exec.Command(conf.GetConfig().GetString("pre_deploy"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	errorcheck.ExitOn(err, "error executing prebuild step")
}
