# confd

## Introduce

confd 是一个配置管理库，支持多种配置形式，多种渠道获取

## Features

+ 支持多种数据来源（file, etcd,...）
+ 支持多种数据格式（json,yaml,...）
+ 支持watch配置变化
+ 支持配置的读管理

## Start

测试结构体定义
```go
type NestStruct struct {
    Nick  string  `json:"nick,omitempty" yaml:"nick,omitempty"`
    Score float64 `json:"score,omitempty" yaml:"score,omitempty"`
}

type TestStruct struct {
    Count int64      `json:"count,omitempty" yaml:"count,omitempty"`
    Str   string     `json:"str,omitempty" yaml:"str,omitempty"`
    Nest  NestStruct `json:"nest,omitempty" yaml:"nest,omitempty"`
}
```

使用方式
```go

cnd := confd.NewConfd(confd.WithSchema("test1.yml"))

err := cnd.LoadConfig()
if err != nil{
// handler err
}

var val TestStruct
// read test1.yml to val
if err = cnd.Read(&val); err != nil{
// handler err
}

// read count to val.Count
if err = cnd.ReadSection("count", &val.Count); err != nil{
    // handler err
}

// read count to val.Count
if err = cnd.ReadSection("nest.nick", &val.Nest.Nick); err != nil{
// handler err
}
```
