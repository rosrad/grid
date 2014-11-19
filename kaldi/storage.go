//
package kaldi

import (
	"io"
	"os"
	"path"
)

func ResultFile(tag string) string {
	dir := path.Join(TaskDir(), tag)
	InsureDir(dir)
	return path.Join(dir, Now()+".org")
}

func ResultTo(r io.Reader, tag string) error {
	fw, err := os.OpenFile(ResultFile(tag), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer fw.Close()
	_, cerr := io.Copy(fw, r)
	return cerr
}
