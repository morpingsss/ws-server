# 使用 CentOS 作为基础镜像进行生产环境配置
FROM centos:latest

# 配置工作目录
WORKDIR /data/go-websocket

# 复制已编译好的二进制文件和配置文件到工作目录
COPY go-websocket ./
COPY conf ./conf

# 暴露端口
EXPOSE 7800

# 设置容器启动命令
CMD ["/data/go-websocket/go-websocket", "-c", "./conf/app.ini"]
