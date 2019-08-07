package conch

import (
	"fmt"
	"os"
	"time"
)

type Log struct {
	Path	string
}

func (log Log)Out(format string, arg ...interface{}) {
	result := time.Now().Format("[2006-01-02 15:04:05] ")
	result += fmt.Sprintf(format + "\n", arg...)
	fmt.Print(result)
	if !Exists(log.Path) {
		file, err := os.Create(log.Path)
		defer file.Close()
		if err != nil {
			fmt.Println("Create log file error.")
		}
	}
	file, err := os.OpenFile(log.Path, os.O_APPEND|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("Log error.")
	}
	_, _ = file.WriteString(result)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}