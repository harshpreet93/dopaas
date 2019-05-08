package key_util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/harshpreet93/dopaas/error_check"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"strings"
)

func GetPubKeySignature() (string, error) {
	sshPubKeyFile, err := homedir.Expand("~/.ssh/id_rsa.pub")
	error_check.ExitOn(err, "error getting pub key in ~/.ssh/id_rsa.pub")
	sshPubKeyContents, err := ioutil.ReadFile(sshPubKeyFile)
	error_check.ExitOn(err, "Error getting pub key file contents")
	log.Println("pubKey is ", string( sshPubKeyContents))
	parts := strings.Fields(string(sshPubKeyContents))
	if len(parts) < 2 {
		log.Fatal("bad key")
	}
	k, err := base64.StdEncoding.DecodeString(parts[1])
	error_check.ExitOn(err, "key decoding error")
	fp := md5.Sum([]byte(k))
	str := ""
	for i, b := range fp {
		str = str + fmt.Sprintf("%02x", b)
		if i < len(fp)-1 {
			str = str + ":"
		}
	}
	return str, nil
}