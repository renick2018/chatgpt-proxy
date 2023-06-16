- 轻量级提供带有账号调度的chatgpt api
- 调度多个chatgpt账号请求问题，保证每个账号同时只处理一个问题
- 缓存上下文，可自定义conversation key，但没有做持久化缓存
- 需要[chatgpt-api-server](https://github.com/renick2018/chatgpt-api-server.git) 配合使用

--- 
使用时需要使用[chatgpt-api-server](https://github.com/renick2018/chatgpt-api-server.git) (base [chatpgt-api](https://github.com/transitive-bullshit/chatgpt-api.git)) 部署账号，然后使用这个程序调度对外提供API  

.conf.yml配置api密钥与chatgpt-api-server地址

chatgpt-api-server提供http接口，但web方式没有对同时多个请求进行控制

部署chatgpt-api-server时使用server分支

> 忙碌指数 = asking + 0.5 * activeConv + 1 - 1.0/(count + 1)

```http request
/chatgpt/ask

#request
{
    "conversationId": "custom conversation nickname",
    "message": "hello"
}

#response
{
    "code": 0,
    "error": "",
    "response": {
        "response": "Hello! How can I assist you today?",
        "conversationId": "custom conversation nickname"
    }
}

```


---

### 使用docker compose
```shell
cd app # 到你想存放项目的目录
mkdir log && mkdir log/proxy && mkdir log/api
git clone https://github.com/renick2018/chatgpt-api-server.git
git clone https://github.com/renick2018/chatgpt-proxy.git
# 修改配置文件，参考example
vim .env
vim .conf.yml
cp chatgpt-proxy/deploy.sh .
chmod +x deploy.sh
sh deploy.sh
```