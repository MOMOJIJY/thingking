# Thinking
Thinking 是一个微信公众号后台的学习项目。
使用golang语言开发，利用docker部署。

## 部署
1. 安装docker
2. 运行 `sudo service docker start`
3. 修改 docker-compose.yaml 文件中的`APPID`,`APPSECRET`和`BGKEY`
3. 运行 `sudo docker compose -f docker-compose.yaml up -d`
4. 访问 `127.0.0.1:8080`
