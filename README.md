# xmstep
zeep life (小米运动) 刷步数 api，支持邮箱/手机账号，支持指定时间


# 配置

以下参数为默认值，需要修改可查看 `main.go` 文件内容
请求路径：`/run`
启动端口：`:3000`

# 部署服务
自己编译启动

或者使用已有容器（不保证后续更新）
```bash
    docker run -d --name xmstep -e "PORT=3000" -p 3000:3000 epurs/xmstep
```

# 请求方法

用法：
```
    POST/GET '127.0.0.1:3000/run?username=<phone number|email address>&password=<your password>&step=< >&date=2022-12-08%2021:30:00'
```
date 可选，格式例如`2022-12-08%2021:30:00`（中间空格用`%20`替换）

使用 `curl`：
```bash
    curl '127.0.0.1:3000/run?username=16611110000&password=woshimima&step=8866&date=2022-12-08%2021:30:00'
```

也可以 POST JSON/form 格式请求 API，自己测试