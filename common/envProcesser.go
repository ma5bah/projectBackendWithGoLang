package common


import (
"bufio"
"os"
"strings"
)

type envName struct {
	names map[string]string
}

func newEnv() *envName {
	name := make(map[string]string)
	file, _ := os.Open(".env")
	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		interData := strings.Split(s.Text(), "=")
		data:=strings.Join(interData[1:],"=")
		name[interData[0]] = data
	}
	return &envName{name}
}
var name = newEnv()

func LocalGetEnv(i string) string {
	return (*name).names[i]
}
