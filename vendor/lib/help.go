package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//
// WriteFile 写入文件，写之前要先建立文件
//
func WriteFile(filename string, data []byte) {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		fmt.Println("create dir error:", err)
	}
	ioutil.WriteFile(filename, data, 0777)
}

//
// GetRootPath 根据url里的/得到rootpath
//
func GetRootPath(url string) []byte {
	count := strings.Count(url, "/")
	str := ""
	for i := 0; i < count; i++ {
		str += "../"
	}
	return []byte(str)
}

//
// GetPageTitle 生成页面的标题
//
func GetPageTitle(title []byte) []byte {
	arr := make([][]byte, 0)
	arr = append(arr, []byte(title))
	arr = append(arr, []byte("-"))
	arr = append(arr, []byte(Conf.Title))
	return bytes.Join(arr, []byte(" "))
}
