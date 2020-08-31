package fileMeta

import (
	"crypto/sha1"
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileMeta struct {
	Id int
	Name string
	Size int64
	UpdatedAt string
	Path string
	Hash string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func Set(fm FileMeta) {
	fileMetas[fm.Name] = fm
}

func (fm FileMeta) Get() FileMeta {
	return fileMetas[fm.Name]
}

func (fm FileMeta) GetSize() int64 {
	file, _ := os.Stat(fm.Path)
	return file.Size()
}

func (fm FileMeta) isExist() bool {
	_, err := os.Stat(fm.Path)
	if (os.IsNotExist(err)) {
		return false
	}
	return true
}

func (fm FileMeta) CreateDirIfNotExist(dir string) error {
	 if !fm.isExist() {
		 err := os.MkdirAll(dir, 0777)
		 if err != nil {
			 return err
		 }
	 }
	 return nil
}

func (fm FileMeta) GetModTime() time.Time {
	file, _ := os.Stat(fm.Path)
	return file.ModTime()
}

func (fm *FileMeta) ToSha1() {
	s := fm.Name + strconv.FormatInt(fm.Size, 10)
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	fm.Hash = fmt.Sprintf("%x", bs)
}