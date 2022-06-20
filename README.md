# http

盘古框架`Http`集成，方便用户在`盘古`框架中接入`Http`功能

> 本框架对`Http`做了一定的封装，使用参看使用文档

## 快速开始

`Http`使用非常简单，只需要定义`配置`和依赖项

> `配置`有很多，但是大部分都有默认值，可以参考[配置文档](https://http.pangum.tech/config)

`配置`代码如下

```yaml
http:
  broker:
    - tcp://192.168.95.102:31883
    - ws://192.168.95.102:38083
  options:
    username: test_username
    password: test_password
    clientid: ${HOSTNAME}
```

`依赖项`的代码如下

```go
package main

import (
  `github.com/pangum/http`
)

type agent struct {
  client *http.Client
}

func newAgent(client *http.Client) *agent {
  return &agent{
    client: client,
  }
}

func (a *agent) subscribe() error {
  return a.client.Subscribe(`topic`, opts...)
}
```

> `Http`有非常多的配置项，请参看[**使用文档**](https://http.pangum.tech/guide)

## 文档

[点击这里查看最新文档](https://http.pangum.tech)

## 使用示例

[点击这里查看最新的代码示例](example)

## 交流

![微信群](doc/.vuepress/public/communication/wxwork.jpg)

## 捐助

![支持宝](https://github.com/storezhang/donate/raw/master/alipay-small.jpg)
![微信](https://github.com/storezhang/donate/raw/master/weipay-small.jpg)
