package ui

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

func (ui *UI) md5(filePath string) (string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", errors.Wrap(err, "copy")
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
