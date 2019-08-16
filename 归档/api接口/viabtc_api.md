# ViaBTC 数据API


## API说明

- endpoint: https://pool.viabtc.com

- API接口均返回如下json对象，错误码`code=0`时成功，其他表示失败，`message`描述错误原因。`data`为json对象或数组。
```
{
  "code": 0,
  "data": {}
  "message": "OK"
}
```

- 对于`GET`和`DELETE`请求，请求参数应放在url中。

- 对于`POST`, `PUT`请求，请求参数应是json字符串，放在请求体中，并带上`Content-Type: application/json`请求头。

- API接口默认返回的日数据统计周期为北京时间0点\~24点，如日算力，日收益等。如需获取以utc0点\~24点为周期统计的日数据，可传入`utc=true`参数。支持`utc`参数的接口会在其参数列表中列出该参数。

| 参数 | 类型 | 必须 | 备注 |
| -- | -- | -- | -- |
| utc | string | no | true/false, 默认false |


## API鉴权

- 每个账户或子账户都有一对`api_key`和`secret_key`，可在[挖矿设置](https://pool.viabtc.com/setting/mining)页面找到。`api_key`用于标识用户身份，`secret_key`用于签名，以验证用户身份，请勿泄露。如：

| key | value |
| -- | -- |
| api_key | 16289e05354c3c3814b8f3045950395f |
| secret_key | d186ababcb0eb1f6af5c1519424f462b84c631f86c06309992ae1f15604668b0 |

- 每个API接口均需要鉴权，鉴权方式为在请求中带上`X-API-KEY: <your_api_key>`的请求头。例如：
```
curl 'https://pool.viabtc.com/res/openapi/v1/hashrate?coin=BTC' -H 'X-API-KEY: 16289e05354c3c3814b8f3045950395f'
```
- 部分接口涉及到用户敏感信息，除鉴权外，还需要签名。签名及请求步骤如下：

1. 在请求参数中添加`tonce`字段，取值为精确到毫秒的时间戳，确保该时间戳小于服务器时间，并且跟服务器时间相差不超过1分钟。例如：

| 参数 | 值 |
| -- | -- |
| coin | BTC |
| amount | 1.0 |
| tonce | 1513746038205 |

2. 如果是`GET`或`DELETE`请求，生成`query string`(参数顺序任意)：
```
coin=BTC&amount=1.0&tonce=1513746038205
```
用密钥`serect_key`计算`query string`的HMAC SHA256签名:
```
echo -n 'coin=BTC&amount=1.0&tonce=1513746038205' | openssl dgst -sha256 -hmac "d186ababcb0eb1f6af5c1519424f462b84c631f86c06309992ae1f15604668b0"
# output: 4a1c9e4c73629b62fd999cbbcd2bc8b87a07f1791ae61ba576427f820dd3bc59
```

如果是`POST`或`PUT`请求，生成json格式的`request body`(参数顺序任意):
```
{"coin": "BTC", "amount": "1.0", "tonce": 1513746038205}
```
用密钥`serect_key`计算`request body`的HMAC SHA256签名:
```
echo -n '{"coin": "BTC", "amount": "1.0", tonce: 1513746038205}' | openssl dgst -sha256 -hmac "d186ababcb0eb1f6af5c1519424f462b84c631f86c06309992ae1f15604668b0"
# output: 2dd2cc09cb5b11c7f24d457b1c01e46f2f2b17f06a62bf29eda35b5fb45ce354
```

3. 在请求中带上`X-API-KEY: <your_api_key>`以及`X-SIGNATURE: <signature>`的请求头，`signature`即是上一步中计算出的签名，并用上一步的`query string`或`request body`发出请求。

```
curl 'https://pool.viabtc.com/res/openapi/v1/hashrate?coin=BTC&amount=1.0&tonce=1513746038205' -H 'X-API-KEY: 16289e05354c3c3814b8f3045950395f' -H 'X-SIGNATURE: 4a1c9e4c73629b62fd999cbbcd2bc8b87a07f1791ae61ba576427f820dd3bc59'
```
或
```
curl 'https://pool.viabtc.com/res/openapi/v1/hashrate' -H 'Content-Type: application/json' -H 'X-API-KEY: 16289e05354c3c3814b8f3045950395f' -H 'X-SIGNATURE: 2dd2cc09cb5b11c7f24d457b1c01e46f2f2b17f06a62bf29eda35b5fb45ce354' -d '{"coin": "BTC", "amount": "1.0", tonce: 1513746038205}' -X POST
```


## 目录

* [账户](#账户)
  * [账户信息](#账户信息)
  * [子账户列表](#子账户列表)
  * [创建子账户](#创建子账户)
  * [观察者列表](#观察者列表)
  * [创建观察者](#创建观察者)
  * [创建子账户-聚合](#创建子账户-聚合)
* [算力](#算力)
  * [账户实时算力](#账户实时算力)
  * [账户历史算力]    (#账户历史算力)
  * [账户算力曲线](#账户算力曲线)
  * [矿工实时算力列表](#矿工实时算力列表)
  * [矿工实时算力](#矿工实时算力)
  * [矿工历史算力](#矿工历史算力)
  * [矿工分组实时算力列表](#分组实时算力列表)
  * [矿工分组实时算力](#分组实时算力)
* [钱包](#钱包)
  * [设置自动提现地址](#设置自动提现地址)
  * [账户支付记录](#账户支付记录)
  * [账户收益汇总](#账户收益汇总)
  * [账户历史收益]    (#账户历史收益)
  * [结算钱包余额](#结算钱包余额)

-----

## 账户

### 账户信息

* 请求：

GET /res/openapi/v1/account

* 响应：

```
{
  "code": 0,
  "data":{
    "account": {
      "id": 45,                     # 用户id
      "parent_user_id": null,       # 主账户id
      "create_time": 1530012667,    # 注册时间
      "account": "test",            # 用户名
      "email": "test@viabtc.com",   # 邮箱
      "country_code": "86",         # 手机国家区号
      "mobile": "17816354561",      # 手机号
      "country": "China"            # 国家
    },
    "observer": [{                  # 观察者列表
      "id": 55,                     # 观察者ID
      "user_id": 45,                # 用户id
      "name": "ob1",                # 观察者名称
      "access_key": "87920fce6b25bcdc6fb86922b935609e",   # 观察者访问key
      "create_time": 1548242387                           # 创建时间
    }],
    "withdraw_address": [{                                # 自动提现地址
      "coin": "BTC",
      "address": "mtRJjPJGVLGs5YDf4VUP5RQXipzHjnjeCe"
    }]，
    "balance": [{                                         # 币种余额
      "coin": "BTC",
      "amount": "0.001"
    }]
  },
  "message": "OK"
}
```

### 子账户列表

* 该接口需要签名。

* 请求：

GET /res/openapi/v1/account/sub

|参数名|类型|必须|备注|
| -----|----|----|----|
| tonce | int | yes | 时间戳 |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |

* 响应：

```
{
  "code": 0,
  "count": 2,
  "curr_page": 1,
  "data": [
    {
      "account": "sub1",                                 # 账户名
      "api_key": "87920fce6b25bcdc6fb86922b935609e",     # 账户api key
      "secret_key": "",                                  # 账户serect key
      "id": 55,                                          # 账户id
      "create_time": 1548242387
    },
    {
      "account": "sub2",
      "api_key": "68ceca48bb090d9d71392fc38b9c96b9",
      "secret_key": "",
      "id": 122,
      "create_time": 1548242387
    }
  ],
  "has_next": false,
  "total": 2,
  "total_page": 1
  "message": "OK"
}
```

### 创建子账户

* 该接口需要签名。

* 请求：

POST /res/openapi/v1/account/sub

|参数名|类型|必须|备注|
| -----|----|----|----|
| tonce | int | yes | 时间戳 |
| account | string | yes | 账户名 |

* 响应：

```
{
  "code": 0,
  "data":{
      "account": "sub1",                                 # 账户名
      "api_key": "87920fce6b25bcdc6fb86922b935609e",     # 账户api key
      "secret_key": "",                                  # 账户serect key
      "id": 55,                                          # 账户id
      "create_time": 1548242387
    },
  "message": "OK"
}
```

### 观察者列表

* 观察者面板页面: https://pool.viabtc.com/observer/bitdeer/dashboard?coin={coin}&access_key={access_key}

* 请求：

GET /res/openapi/v1/account/observer

|参数名|类型|必须|备注|
| -----|----|----|----|
| page | int | no | 页码 |
| limit | int | no | 每页条数 |

* 响应：

```
{
  "code": 0,
  "count": 2,
  "curr_page": 1,
  "data": [
    {
      "id": 55,
      "user_id": 23,
      "name": "ob1",
      "access_key": "87920fce6b25bcdc6fb86922b935609e",
      "create_time": 1548242387
    },
    {
      "id": 56,
      "user_id": 23,
      "name": "ob2",
      "access_key": "c5620fcf4b25bcdc6f568923b9306e87",
      "create_time": 1548242388
    }
  ],
  "has_next": false,
  "total": 2,
  "total_page": 1,
  "message": "OK"
}
```

### 创建观察者

* 该接口需要签名。

* 请求：

POST /res/openapi/v1/account/observer

|参数名|类型|必须|备注|
| -----|----|----|----|
| tonce | int | yes | 时间戳 |
| name | string | yes | 观察者名称 |

* 响应：

```
{
  "code": 0,
  "data":{
      "id": 55,
      "user_id": 23,
      "name": "ob1",
      "access_key": "87920fce6b25bcdc6fb86922b935609e",
      "create_time": 1548242387
    },
  "message": "OK"
}
```

### 创建子账户-聚合

* 该接口需要签名。

* 该接口创建一个子账户和观察者，并设置自动提现地址。

* 请求：

POST /res/openapi/v1/account/sub/aggregate

|参数名|类型|必须|备注|
| -----|----|----|----|
| tonce | int | yes | 时间戳 |
| account | string | yes | 账户名 |
| observer_name | string | yes | 观察者名称 |
| withdraw_address | array | no | 收款地址，格式[{"coin": "BTC", "address": "mtRJjPJGVLGs5YDf4VUP5RQXipzHjnjeCe"}] |

* 响应：

```
{
  "code": 0,
  "data":{
    "account": {
        "account": "sub1",
        "api_key": "87920fce6b25bcdc6fb86922b935609e",
        "secret_key": "",
        "id": 55,
        "create_time": 1548242387
    },
    "observer": {
      "id": 67,
      "user_id": 55,
      "name": "ob1",
      "access_key": "87920fce6b25bcdc6fb86922b935609e",
      "create_time": 1548242387
    },
    "withdraw_address": [{
      "coin": "BTC",
      "address": "mtRJjPJGVLGs5YDf4VUP5RQXipzHjnjeCe"
    }]
  },
  "message": "OK"
}
```

-----

## 算力

### 账户实时算力

* 请求：

GET /res/openapi/v1/hashrate

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |

* 响应：

```
{
  "code": 0,
  "data": {
      "active_workers": 1,                # 有效矿工数
      "coin": "BTC",                      # 币种
      "hashrate_10min": "747667906887",   # 10分钟平均算力，单位hash/s
      "hashrate_1hour": "124611317814",   # 1小时平均算力
      "hashrate_24hour": "0",             # 24小时平均算力
      "unactive_workers": 0               # 无效矿工数
    },
  "message": "OK"
}
```

### 账户历史算力

* 请求：

GET /res/openapi/v1/hashrate/history

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |
| start_date | string | no | 起始日期，格式 2019-01-24 |
| end_date | string | no | 截止日期 |
| utc | string | no | true/false, 默认false |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |

* 响应：

```
{
  "code": 0,
  "data": {
    "count": 2,
    "curr_page": 1,
    "data": [
      {
        "coin": "BCH",
        "date": "2018-10-05",
        "hashrate": "651358832825",
        "reject_rate": "0"
      },
      {
        "coin": "BCH",
        "date": "2018-10-04",
        "hashrate": "663931951901",
        "reject_rate": "0"
      }
    ],
    "has_next": false,
    "total": 2,
    "total_page": 1
  },
  "message": "OK"
}
```

### 账户算力曲线

* 请求：

GET /res/openapi/v1/hashrate/chart

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |
| interval | string | yes | min/hour/day，取样点间隔，每10分钟/1小时/1天 |
| period | int | no | 点数量，不传则返回系统默认的最大点数量 |
| utc | string | no | true/false, 默认false，仅interval=day时有效 |

* 响应：

```
{
  "code": 0,
  "data":[
    {
      "timestamp": 1548241200,      # 时间戳
      "hashrate": "663931951901",   # 算力
      "reject_rate": "0"            # 拒绝率
    },
    {
      "timestamp": 1548237600,
      "hashrate": "663931951901",
      "reject_rate": "0"
    },
  ],
  "message": "OK"
}

```

### 矿工实时算力列表

* 请求：

GET /res/openapi/v1/hashrate/worker

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | int | yes | 币种 |
| group_id | int | no | 分组id，可以从分组实时算力列表API获取 |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |

* 响应：

```
{
  "code": 0,
  "data": {
    "count": 2,
    "curr_page": 1,
    "data": [
      {
        "coin": "BTC",
        "group_id": 14,                     # 分组id，默认分组为null
        "group_name": "group1",             # 分组名，默认分组为null
        "hashrate_10min": "278542945703",
        "hashrate_1hour": "46423824283",
        "hashrate_24hour": "0",
        "last_active": 1545811200,          # 最后提交时间
        "reject_rate": "0",                 # 拒绝率
        "worker_id": 369,                   # 矿工id
        "worker_name": "1x1",               # 矿工名
        "worker_status": "active"           # 矿工状态, active-有效, unactive-无效
      },
      {
        "coin": "BTC",
        "group_id": null,
        "group_name": null,
        "hashrate_10min": "0",
        "hashrate_1hour": "0",
        "hashrate_24hour": "0",
        "last_active": 1545808266,
        "reject_rate": "0",
        "worker_id": 370,
        "worker_name": "1x1",
        "worker_status": "unactive"
      }
    ],
    "has_next": false,
    "total": 2,
    "total_page": 1
  },
  "message": "OK"
}
```

### 矿工实时算力

* 请求：

GET /res/openapi/v1/hashrate/worker/{worker_id}

* 响应：

```
{
  "code": 0,
  "data": {
    "coin": "BTC",
    "group_id": 14,
    "group_name": "group1",
    "hashrate_10min": "747667906887",
    "hashrate_1hour": "124611317814",
    "hashrate_24hour": "0",
    "last_active": 1545811500,
    "reject_rate": "0",
    "worker_id": 369,
    "worker_name": "1x1",
    "worker_status": "active"
  },
  "message": "OK"
}
```

### 矿工历史算力

* 请求：

GET /res/openapi/v1/hashrate/worker/{worker_id}/history

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |
| start_date | string | no | 起始日期 |
| end_date | string | no | 截止日期 |
| utc | string | no | true/false, 默认false |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |

* 响应：

```
{
  "code": 0,
  "data": {
    "count": 2,
    "curr_page": 1,
    "data": [
      {
        "coin": "BTC",
        "date": "2018-07-11",
        "hashrate": "10333373168",
        "reject_rate": "0"
      },
      {
        "coin": "BTC",
        "date": "2018-07-10",
        "hashrate": "7788207363",
        "reject_rate": "0"
      }
    ],
    "has_next": false,
    "total": 2,
    "total_page": 1
  },
  "message": "OK"
}
```

### 矿工分组实时算力列表

* 请求：

GET /res/openapi/v1/hashrate/group

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |

* 响应：

```
{
  "code": 0,
  "data": {
    "count": 2,
    "curr_page": 1,
    "data": [
      {
        "active_workers": 2,              # 有效矿工数
        "coin": "BTC",
        "group_id": 14,                   # 分组id
        "group_name": "group1",           # 分组名
        "hashrate_10min": "498445271257",
        "hashrate_1hour": "296868139498",
        "hashrate_24hour": "0",
        "reject_rate": "0",
        "total_workers": 2                # 总矿工数
      },
      {
        "active_workers": 0,
        "coin": "BTC",
        "group_id": 15,
        "group_name": "group2",
        "hashrate_10min": "0",
        "hashrate_1hour": "0",
        "hashrate_24hour": "0",
        "reject_rate": "0",
        "total_workers": 1
      }
    ],
    "has_next": false,
    "total": 2,
    "total_page": 1
  },
  "message": "OK"
}
```

### 矿工分组实时算力

* 请求：

GET /res/openapi/v1/hashrate/group/{group_id}

* 响应：

```
{
  "code": 0,
  "data": {
    "active_workers": 2,
    "coin": "BTC",
    "group_id": 14,
    "group_name": "group1",
    "hashrate_10min": "505775348776",
    "hashrate_1hour": "313971653709",
    "hashrate_24hour": "0",
    "reject_rate": "0",
    "total_workers": 2
  },
  "message": "OK"
}
```

-----

## 钱包

### 设置自动提现地址

* 该接口需要签名。

* 请求：

POST /res/openapi/v1/wallet/payment/address

|参数名|类型|必须|备注|
| -----|----|----|----|
| tonce | int | yes | 时间戳 |
| coin | string | yes | 币种 |
| address | string | yes | 地址 |
| payment_password | string | no | 支付密码，未设置则忽略 |


* 响应：

```
{
  "code": 0,
  "data":{
    "coin": "BTC",
    "address": "mtRJjPJGVLGs5YDf4VUP5RQXipzHjnjeCe"
  },
  "message": "OK"
}
```

### 账户支付记录

* 请求：

GET /res/openapi/v1/wallet/payment/history

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |
| start_date | string | no | 起始日期 |
| end_date | string | no | 截止日期 |
| utc | string | no | true/false, 默认false |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |

* 响应：

```
{
  "code": 0,
  "count": 2,
  "curr_page": 1,
  "data":[
    {
      "id": 157,
      "coin": "BTC",
      "amount": "0.001",
      "address": "mtRJjPJGVLGs5YDf4VUP5RQXipzHjnjeCe",    # 汇出地址
      "tx": "eaa0597e556ceda83ffe5d3533a4aba93b49e7dbb2fa35895dd08754fb9d62d0",   # txid
      "create_time": 1530704756
    },
    {
      "id": 266,
      "coin": "BTC",
      "amount": "0.001",
      "address": "qr0xs7rk2ku4rkg4kpayvcet3pvw7zcceszsaflkzs",
      "tx": "00802e7d6acf262c9afda8882dc361c705f224985c074f62e57942af9beae8b9",
      "create_time": 1530704757
    },
  ],
  "has_next": false,
  "total": 2,
  "total_page": 1
  "message": "OK"
}
```


### 账户收益汇总

* 请求：

GET /res/openapi/v1/profit

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |


* 响应：

```
{
  "code": 0,
  "data": {
      "coin": "BTC",
      "pplns_profit": "0",            # pplns收益
      "pps_profit": "0.00002148",     # pps收益
      "solo_profit": "0",             # solo收益
      "total_profit": "0.00002148"    # 总收益
    },
  "message": "OK"
}
```

### 账户历史收益

* 请求：

GET /res/openapi/v1/profit/history

|参数名|类型|必须|备注|
| -----|----|----|----|
| coin | string | yes | 币种 |
| start_date | string | no | 起始日期 |
| end_date | string | no | 截止日期 |
| utc | string | no | true/false, 默认false |
| page | int | no | 页码 |
| limit | int | no | 每页条数 |

* 响应：

```
{
  "code": 0,
  "data": {
    "count": 2,
    "curr_page": 1,
    "data": [
      {
        "coin": "BTC",
        "date": "2018-10-05",
        "pplns_profit": "0",            # pplns收益
        "pps_profit": "0.00002148",     # pps收益
        "solo_profit": "0",             # solo收益
        "total_profit": "0.00002148"    # 总收益
      },
      {
        "coin": "BTC",
        "date": "2018-10-06",
        "pplns_profit": "0",
        "pps_profit": "0.00028948",
        "solo_profit": "0",
        "total_profit": "0.00028948"
      }
    ],
    "has_next": false,
    "total": 2,
    "total_page": 1
  },
  "message": "OK"
}
```

### 结算钱包余额

* 该接口需要签名。

* 请求：

POST /res/openapi/v1/wallet/sweep

|参数名|类型|必须|备注|
| -----|----|----|----|
| tonce | int | yes | 时间戳 |
| coin | string | yes | 币种 |

* 响应：

```
{
  "code": 0,
  "data": {
      "coin": "BTC",
      "balance": "0.001"
    },
  "message": "OK"
}
```
