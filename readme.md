# table2go

把数据库表结构转换成golang的struct，目前支持postgresql 和 mysql两种数据库

> 我们程序员的宗旨就是宁愿花时间开发工具也不做重复劳动，但是开发过程中设计完表结构后，又要写一遍struct，完全是重复劳动。所以这个工具就诞生了。

## 安装

```go
go get -u -v github.com/akkagao/table2go
```

## 配置

在项目根目录创建`t2g.yaml`配置文件

**可以复制 t2g.yaml.example 中的内容修改**

```yaml
databasez:
  #数据库配置
  #database:
  # type 取值范围 postgresql 和 mysql
  type: mysql
  # type 取值为mysql 需要配置 mysqlConn
  mysqlConn: {user}:{password}@tcp({ip}:{port})/{dbname}
  # type postgresql，postgresqlConn按如下格式配置
  postgresqlConn: user={user} password={password} dbname={dbname} host={ip} port={port} sslmode=disable

#数据库表名前缀
tableNameHandler: ["t_","tb_"]
```

## 使用

执行下面命令，t_user 为表名

```shell
table2go -t t_user
```

执行结果

```go
type User struct {
	ID           int64     `json:"id"`            // ID主键
	Mobile       string    `json:"mobile"`        // 用户手机号
	Mail         string    `json:"mail"`          // 用户邮箱账号
	Password     string    `json:"password"`      // 用户密码
	Age          int       `json:"age"`           // 年龄
	Sex          int       `json:"sex"`           // 用户性别 0：未知，1：男，2：女
	ProvinceID   int       `json:"province_id"`   // 省份ID
	CityID       int       `json:"city_id"`       // 城市ID
	DistinctID   int       `json:"distinct_id"`   // 用户区县ID
	Nick         string    `json:"nick"`          // 昵称
	Realname     string    `json:"realname"`      // 用户真实姓名
	Header       string    `json:"header"`        // 用户头像
	Birthday     string    `json:"birthday"`      // 出生日期
	Level        int       `json:"level"`         // 用户等级
	LastLogin    time.Time `json:"last_login"`    // 最后一次登录时间
	CreatedAt    time.Time `json:"created_at"`    // 用户注册时间
}
```

- [ ] changeType 数据库类型转换为go类型，现在类型还不是很完善，有些类型没做转换，如果报错了请issues留言
- [ ] 目前直接输出到控制台了，后续有必要的话再直接写入文件吧