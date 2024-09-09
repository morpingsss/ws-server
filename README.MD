# ws-server

## 快速开始

**下载本项目：**

```shell
git clone https://github.com/morpingsss/ws-server.git
```

**编译：**

```shell
// 编译适用于本机的版本
go build

// 编译Linux版本
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

// 编译Windows 64位版本
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

// 编译MacOS版本
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
```

**执行：**

编译成功之后会得到一个二进制文件`go-websocket`，执行该二进制文件。

```shell
./ws-server -c ./conf/app.ini
```

**连接测试：**

编辑hosts文件  120.79.131.10 websocket.example.com

```shell
# 注册系统
curl -X POST http://websocket.example.com/api/register -H "Content-Type: application/json; charset=UTF-8" -d "{\"systemId\": \"abcd\"}"
# 启动一个客户端
websocat ws://websocket.example.com/ws?systemId=abcd # 这个命令需要新开一个shell
# 将客户端绑定到分组
curl -X POST http://websocket.example.com/api/bind_to_group -H "Content-Type: application/json; charset=UTF-8" -H "systemId: abcd" -d "{\"clientId\": \"$clientId\", \"groupName\": \"abc\"}"
# 获取在线列表
curl -X POST http://websocket.example.com/api/get_online_list -H "Content-Type: application/json; charset=UTF-8" -H "systemId: abcd" -d "{\"groupName\": \"abc\"}" 
```

## 使用docker

### 体验单机

1. 构建镜像

   ```shell
   docker build -t ws-server .
   ```

2. 基于镜像运行容器

   ```shell
   docker run -tid -p 7800:7800 ws-server
   ```

### 体验集群，同时运行ETCD集群
在当前目录下，直接运行` docker-compose up` 即可体验。

## 使用k8s + docker 

   ```shell
   
    kubectl apply -f k8s-backup/
   ```
