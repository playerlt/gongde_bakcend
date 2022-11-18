# gongde_bakcend
## 可有可无的简介
使用redis递减功德，websocket向客户端推送剩余功德，mysql每小时备份一次redis数据

## 使用
1. 执行model/sql中的sql脚本创建数据库和用户
2. 在config/config.go文件中配置mysql和redis连接