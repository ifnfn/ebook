package util

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

// CheckEmail 检查是否是邮件
func CheckEmail(s string) bool {
	re := regexp.MustCompile(`^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(\\.([a-zA-Z0-9_-])+)+$`)
	fmt.Println(re)
	return len(s) > 3 && re.MatchString(s)
}

// GetKindleEmail 从字符中提取 kindle 邮箱地址
func GetKindleEmail(s string) string {
	emailRegexp := regexp.MustCompile("[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@kindle.*")
	v := emailRegexp.FindStringSubmatch(s)

	if len(v) > 0 {
		return v[0]
	}

	return ""
}

// CopyFile 拷贝文件
func CopyFile(srcFile, destFile string) error {
	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer file.Close()

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	io.Copy(dest, file)

	return nil
}
