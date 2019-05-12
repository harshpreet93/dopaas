package doaction

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/harshpreet93/dopaas/errorcheck"
	"io"
	"os"
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
