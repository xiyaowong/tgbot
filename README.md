# TGBOT

**GO+LUA+TDLIB 跨平台 TG 机器人框架**

# 说明

这是[TGBOT](https://github.com/OPQBOT/TGBOT)项目的二次开发, 功能比较简陋, 相关内容已得到原作者许可

# 配置

在程序所在目录的`config.ini`中填写，支持两项配置

1. `ServerPort`HTTP 监听的端口, 默认 8889,用于下面的 WebApi
2. `HookPostURL`数据上报的地址, 默认为空即不上报(POST), 用于下下面的 WebHook
3. `ShowRawMessage` 是否打印消息的原始数据, 默认为 false

# Web Api

目前只封装了最基本的 api

请求都必须是 json 数据格式,返回数据格式如下:

```jsonc
{
    "code": int, // 成功或出错为0，其他情况为1
    "err": string, // 如果出错，则为错误信息
    "data": // 格式根据api不同而不定
}
```

## 发文字

`POST IP:ServerPort/sendText`

```jsonc
{
  "chat_id": int, // 群组ID, 必填
  "text": string // 文本内容, 不为空
}
```

## 发图(图文)

`POST IP:ServerPort/sendPhoto`

```jsonc
{
  "chat_id": int, // 群组ID, 必填
  "text": string, // 图文文字, 可选
  "base64": string, // 图片的base64编码, 与url二选一
  "url": string // 图片链接, 与base64二选一
  // 如果base64和url都有，优先选用base64
}
```

## 删除消息

`POST IP:ServerPort/deleteMessages`

```jsonc
{
  "chat_id": int, // 群组ID, 必填
  "message_id": int, // 消息ID, 必填
  "revoke": bool // 填 true (暂不知道是什么标志)
}
```

## 获取用户信息

`POST IP:ServerPort/getUser`

```jsonc
{
  "user_id": int // 用户id， 必选
}
```

## 获取群组信息

`POST IP:ServerPort/getChat`

```jsonc
{
  "chat_id": int // 群组 id， 必选
}
```

# Web Hook

收到消息后以 json 格式向 hook 地址 Post 数据

数据内容为, 带\*号表示可能为空

```jsonc
{
  "ChatID": int, // 群组 ID
  "ChatTitle": int, // 群组名称*
  "MessageID": int, // 消息 ID
  "SenderUserID": int, // 发送人 ID
  "SenderUserFirstName": string, // *
  "SenderUserLastName": string, // *
  "SenderUserName": string, // *
  "SenderUserPhoneNumber": string, // *
  "MsgType": string, // 有两种类型 1.MessageText 文本消息 2.MessagePhoto 图片消息
  "Content": string,
  "PhotoBase64": string // 图片的 base64, 仅在消息类型为 MessagePhoto 有
}
```

# 自己编译

将 CrossLib 目录中的链接库拷贝到 TGCLI 目录里，设置软链接
如:

```shell
cp ./CrossLib/linux-amd64/libtdjson.so.1.6.8 ./TGCLI
cd TGCLI
ln -s libtdjson.so.1.6.8 libtdjson.so
export CGO_ENABLED=1
go build
```
