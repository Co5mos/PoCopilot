# 使用最新版本的官方 Golang 镜像作为基础镜像
FROM golang:1.21.4

# 设置环境变量
ENV GO111MODULE=on


# 安装 wails 最新版
RUN go install github.com/wailsapp/wails/v2/cmd/wails@v2.8.1

# 安装必要软件
RUN apt-get update && apt-get install -y curl gnupg xz-utils wget
RUN apt-get update && apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev libgtk-3-0 libwebkit2gtk-4.0-37 mingw-w64

# 安装 Node.js 和 npm
RUN curl -sL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs

# 安装 antd 组件
RUN npm i react

# 安装 UPX
RUN wget -P /tmp/ https://github.com/upx/upx/releases/download/v4.1.0/upx-4.1.0-amd64_linux.tar.xz
RUN tar -xvf /tmp/upx-4.1.0-amd64_linux.tar.xz -C /tmp/ \
    && mv /tmp/upx-4.1.0-amd64_linux/upx /usr/local/bin/ \
    && rm -rf /tmp/upx-4.1.0-amd64_linux.tar.xz /tmp/upx-4.1.0-amd64_linux

# 定义容器启动时的默认命令
CMD ["bash"]
