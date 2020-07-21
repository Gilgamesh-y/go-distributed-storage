package fileMeta

import (
	"os"
	"time"
)

type FileMeta struct {
	Name string
	Size int64
	UpdateAt string
	Path string
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

func (fm FileMeta) GetModTime() time.Time {
	file, _ := os.Stat(fm.Path)
	return file.ModTime()
}