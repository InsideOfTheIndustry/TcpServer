# TcpServer
目的：建立一个tcp服务器
# 需要的功能
* 注册
* 登录
* 聊天
* 身份验证
* 安全性保证
* 性能
# 实现
* 自定义的信息传递结构
# 目前的需求
* 在线情况的修改（需要数据库配合)

# 启动

## 使用docker
请先在配置文件夹中配置`config.yaml`
```shell
# 构建镜像
docker build -t tcpserver:0.1 .
# 启动镜像
docker run -d \
--rm \
--name tcpserver \
-p4000:4000 \
-p4001:4001 \
tcpserver:0.1
```