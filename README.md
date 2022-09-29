# buaa-clock

北航每日打卡程序，不能代替例行防疫，引发的问题本人概不负责。

## 子命令

### clock

clock 子命令可以完成一次打卡操作，示例命令如下：

```bash
$ buaa-clock clock --username=by2101111 --password=password --retry=20
```

| 选项     | 含义               |
| -------- | ------------------ |
| username | 用户学号           |
| password | 用户密码           |
| retry    | 未登录成功重试次数 |

### timing

timing 子命令可以在每天的 17:10 自动执行打卡命令，但是需要保持程序后台运行。示例命令如下：

```bash
$ buaa-clock timing --config=config.yml
```

| 选项   | 含义         |
| ------ | ------------ |
| config | 配置文件路径 |

配置文件示例如下：

```yml
url:
  # 登录 URL
  login: https://xxx.xxx.xxxx
  # 获取信息 URL
  info: https://xxx.xxx.xxxx
  # 打卡 URL
  save: https://xxx.xxx.xxxx

users:
    # 用户账号
  - username: by2101111
    # 用户密码
    password: password
    # 是否在校，1在校，0不在校，如果为 1，reason 和 note 字段无效
    boarder: "1"
    # 不在校理由
    reasen: "2"
    # 理由为其他时，填充的原因
    note: "回家"

    # 下面4个字段为打卡地址，如果不指定为北航学院路
    address: ""
    area: ""
    city: ""
    province: ""
```
