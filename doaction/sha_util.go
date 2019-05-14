package doaction

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/harshpreet93/dopaas/conf"
	"github.com/harshpreet93/dopaas/errorcheck"
	"github.com/sfreiberg/simplessh"
	"io"
	"log"
	"os"
	"time"
)

func GetFileSha(filepath string) string {
	f, err := os.Open(filepath)
	errorcheck.ExitOn(err, "error opening artifact file")
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		errorcheck.ExitOn(err, "error getting sha of artifact")
	}
	return hex.EncodeToString(h.Sum(nil))
}

func GetDropletArtifactSha(dropletID int) {
	sha, err := getCurrSHA(dropletID)
	errorcheck.ExitOn(err, "Unable to get IP for droplet id "+string( dropletID))

}

func getCurrSHA(dropletID int) (string, error) {
	done := make(chan struct{string, err})
	go executeWithTimeout( dropletID, done)
	select {
	case err := <-done:
		return "", err
	case <-time.After(30 * time.Second):
	}
	close(done)
	return "", nil
}

func executeWithTimeout( dropletID int, done chan struct{ string, error }) {
	ip, err := tryToGetIPForId(dropletID)
	errorcheck.ExitOn(err, "Error getting IP for droplet id")
	client, err := simplessh.ConnectWithKeyFile(ip+":22", "root", "")
	errorcheck.ExitOn(err, "error establishing connection to "+ip)
	defer client.Close()
	output, err := client.Exec("cat /root/artifact_sha")
	log.Println("start script output", output)
	done <- err
	close(done)
}