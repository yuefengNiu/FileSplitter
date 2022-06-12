# FileSplitter
> 一段用于分割大文件的程序，先build，再执行，执行之前需要配置`conf.txt`
## 编译不同平台
```bash
# 编译成 Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
# 编译成 Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```
