package predeploy

import (
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/errorcheck"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Execute() {
	setPredeployVars()
	cmd := exec.Command(conf.GetConfig().GetString("pre_deploy"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	errorcheck.ExitOn(err, "error executing pre_deploy step")
	verifyArtifactExistence()
}

func setPredeployVars()  {
	for _, key := range conf.GetConfig().AllKeys() {
		if strings.HasPrefix(key, "pre_deploy_var_") {
			os.Setenv(strings.TrimPrefix(key, "pre_deploy_var_"), conf.GetConfig().GetString(key))
		}
	}
}

func verifyArtifactExistence()  {
	finfo, err:= os.Stat(conf.GetConfig().GetString("artifact_file"))
	errorcheck.ExitOn(err, "error checking for artifact file existence")
	log.Println("artifact info ", finfo)
}