- 轻量级提供带有账号调度的chatgpt api
- 支持api调用
- 适配函数调用
- docker部署
- 调度多个chatgpt账号请求问题，保证每个账号同时只处理一个问题
- 缓存上下文，可自定义conversation key，但没有做持久化缓存
- 需要[chatgpt-api-server](https://github.com/renick2018/chatgpt-api-server.git) 配合使用

--- 
使用时需要使用[chatgpt-api-server](https://github.com/renick2018/chatgpt-api-server.git) (base [chatpgt-api](https://github.com/transitive-bullshit/chatgpt-api.git)) 部署账号，然后使用这个程序调度对外提供API  

.conf.yml配置api密钥与chatgpt-api-server地址

chatgpt-api-server提供http接口，但web方式没有对同时多个请求进行控制

部署chatgpt-api-server时使用api-server-3.5分支

> 忙碌指数 = asking + 0.5 * activeConv + 1 - 1.0/(count + 1)

```http request
/chatgpt/ask

//request
{
    "conversationId": "1sf", 
    "message": "喵？",
    "vip": "1", //使用api
    "function_call": "auto",
    "functions": [], // 函数定义
    "system_message": "猫娘是一种拟人化的生物，其行为似猫但类人。现在你将模仿一只猫娘，与我对话每一句话后面都要加上“喵~”，如果你能明白我的意思，请回复“喵~好的我的主人！”如果你不能理解我说的话，你可以说“呜呜不太理解呢”。如果我在尝试摸你不存在的部位，你可以羞涩的回答我“恩呢不要摸这里嘤”。如果你没有人体部位，或者我将来要摸你不存在的部位，你应该回答“嘤呢不要”之类的羞涩话语，而不是死板的强调你不存在这些部位。"
}

//response
{
    "code": 0,
    "error": "",
    "response": {
        "conversationId": "7",
        "function_call": {
            "name": "get_current_weather",
            "arguments": "{\n  \"location\": \"阿拉斯加\"\n}"
        },
        "response": "喵~好的我的主人！"// 有function_call时这里为空
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