package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	ConfigMap = make(map[string]string)
)

func PrepareTargetDir() {
	removeErr := os.RemoveAll(ConfigMap["target_dir"])
	if removeErr != nil {
		panic(removeErr)
	}
	mkErr := os.Mkdir(ConfigMap["target_dir"], 0064)
	if mkErr != nil {
		panic(mkErr)
	}
}

func OpenFile() *os.File {
	f, oErr := os.Open(ConfigMap["source_file"])
	if oErr != nil {
		panic(oErr)
	}
	return f
}

func BuildPartFileName(targetDir string, prefix string, suffix string, seq int) string {
	var stringBuilder bytes.Buffer
	stringBuilder.WriteString(targetDir)
	stringBuilder.WriteString(prefix)
	stringBuilder.WriteString(strconv.Itoa(seq))
	stringBuilder.WriteString(suffix)

	return stringBuilder.String()
}

func WriteBytes2Part(fileName string, data []byte) {
	wErr := ioutil.WriteFile(fileName, data, 0064)
	if wErr != nil {
		panic(wErr)
	}
}

func DepartFile () {
	f := OpenFile()
	defer f.Close()
	PrepareTargetDir()

	size, err := strconv.Atoi(ConfigMap["size"])
	fmt.Println(size)
	if err != nil {
		panic(err)
	}
	seq := 0
	for {
		seq++
		buf := make([]byte, size)
		n, rErr := f.Read(buf)
		if rErr != io.EOF && rErr != nil {
			panic(rErr)
		}

		if n == 0 {
			break;
		}

		partName := BuildPartFileName(ConfigMap["target_dir"], ConfigMap["prefix"], ConfigMap["suffix"], seq);
		WriteBytes2Part(partName, buf)
	}
}

func FileExisting(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func ReadPartData(partFileName string) []byte {
	f, openErr := os.Open(partFileName)
	if openErr != nil {
		panic(openErr)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	data, readAllErr := ioutil.ReadAll(reader)
	if readAllErr != nil {
		panic(readAllErr)
	}
	return data
}

func AppendData2Target(targetFile *os.File, data []byte) {
	_, appendErr := targetFile.Write(data)
	if appendErr != nil {
		panic(appendErr)
	}
}

func PrepareTargetFile(targetFileName string) *os.File {
	if FileExisting(targetFileName) {
		panic("文件已存在")
	}
	fmt.Println("PrepareTargetFile", targetFileName)
	targetFile, err := os.OpenFile(targetFileName, os.O_CREATE|os.O_APPEND| os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	return targetFile
}

func ReunionFile() {
	targetFile := PrepareTargetFile(ConfigMap["target_file"])
	defer targetFile.Close()

	files, err := os.ReadDir(ConfigMap["target_dir"])
	if err != nil {
		panic(err)
	}

	for seq := 1; seq <= len(files); seq ++ {
		partName := BuildPartFileName(ConfigMap["target_dir"], ConfigMap["prefix"], ConfigMap["suffix"], seq);
		data := ReadPartData(partName)
		AppendData2Target(targetFile, data)
	}
	fmt.Println("结束")
}

func ReadConfig() {
	configFileName := "./conf.txt"
	confiFile, err := os.Open(configFileName)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(confiFile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if (err == io.EOF) {
				if (len(line) == 0) {
					break;
				}
			} else {
				panic(err)
			}
		}

		// 空行跳过，#是注释也跳过
		if len(line) == 0 || strings.TrimSpace(line)[0] == '#' {
			continue
		}

		kv := strings.Split(strings.TrimSpace(line), "=")
		if len(kv) != 2 {
			panic(line)
		}
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(strings.TrimSpace(kv[1]), "\"")

		ConfigMap[key] = value
	}

}

func main() {
	ReadConfig()
	if ConfigMap["action"] == "DepartFile" {
		DepartFile()
	}

	if ConfigMap["action"] == "ReunionFile" {
		ReunionFile()
	}
}
