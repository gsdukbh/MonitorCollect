## 数据监控


## build 监控

1. 编译 Linux amd64 版本

```shell
  env GOOS=linux GOARCH=amd64 go build -o monitor_collect 
```
2. 编译 Windows amd64 版本

```shell
  env GOOS=windows GOARCH=amd64 go build -o monitor_collect.exe 
```

3. 编译 Linux arm64 版本
```shell
  env GOOS=linux GOARCH=arm64 go build -o monitor_collect 
```
4. 编译 docker 镜像

```shell
  docker build -t monitor_collect:latest .
```