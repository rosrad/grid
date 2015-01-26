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

func ResultSet(tag, set string) string {
	dir := path.Join(TaskDir(), tag)
	InsureDir(dir)
	return path.Join(dir, set+".org")
}

func ResultTo(r io.Reader, tag, set string) error {
	fset, serr := os.OpenFile(ResultSet(tag, set), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if serr != nil {
		return serr
	}
	fw, err := os.OpenFile(ResultFile(tag), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer fw.Close()

	multi := io.MultiWriter(fset, fw)
	_, cerr := io.Copy(multi, r)
	return cerr
}
