package seed

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// dirents 读取目录下文件信息
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

// WalkDir 遍历目录下的文件
func WalkDir(dir string, ext string, files chan<- string) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			WalkDir(subdir, ext, files)
		} else {
			files <- entry.Name()
		}
	}
}
