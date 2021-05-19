
1. 拉取依赖
```
go get -u github.com/sumansoul/aliyun-ahas-go-sdk
```

2. 添加 github.com/aliyun/aliyun-ahas-go-sdk
```
require (
    github.com/aliyun/aliyun-ahas-go-sdk v1.0.4
)
```
3. 替换 github.com/aliyun/aliyun-ahas-go-sdk
```
replace github.com/aliyun/aliyun-ahas-go-sdk => github.com/sumansoul/aliyun-ahas-go-sdk v0.0.0-20210519015343-658efacab2ab
```

4. 配置 github.com/sumansoul/aliyun-ahas-go-sdk 为间接依赖
```
github.com/sumansoul/aliyun-ahas-go-sdk v0.0.0-20210519015343-658efacab2ab // indirect
```