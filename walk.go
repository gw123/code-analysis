package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 忽略的目录列表
var ignoredDirs = []string{"vendor", "testdata", ".git", "test/", "mocks/"}

func isIgnored(dir string) bool {
	// 获取目录名称
	dirName := filepath.Base(dir)
	for _, ignored := range ignoredDirs {
		if dirName == ignored {
			return true
		}
	}
	for _, ignored := range ignoredDirs {
		if strings.Contains(dir, ignored) {
			//fmt.Println("Ignored:", dir, ignored)
			return true
		}
	}
	return false
}

// WalkDir 遍历目录并输出所有的 .go 文件
func WalkDir(dir string, callback func(path string)) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//fmt.Println("Error:", err)
			return err
		}

		if strings.HasSuffix(path, "_test.go") {
			//fmt.Println("Skip:", path)
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".go" && !isIgnored(path) {
			fmt.Println(path)
			callback(path)
		}
		return nil
	})
}
