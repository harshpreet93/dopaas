package userdata

import (
	"bytes"
	"github.com/harshpreet93/dopaas/errorcheck"
	"text/template"
)

func Generate() string {
	userData := `#!/bin/bash
				
				`
	tmpl, err := template.New("userdata").Parse(userData)
	errorcheck.ExitOn(err, "error creating userdata template")
	tmplVars := template.FuncMap{}
	compiledUserData := &bytes.Buffer{}
	err = tmpl.Execute(compiledUserData, tmplVars)
	errorcheck.ExitOn(err, "error compiling userdata template")
	return compiledUserData.String()
}
