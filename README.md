# restful api demo

+ 程序配置管理
+ impl 基于MySQL存储
+ http协议暴露


### 流程
1. 抽象接口层
2. 配置层
3. 业务逻辑层
4. HTTP API层

### 各级目录注释
```
|-cmd                         # 程序Cli工具包
|-conf                        # 程序配置代码，从程序配置文件读取，生成全局配置
|-etc                         # 程序配置文件
|-protocol                    # 程序监听的协议
|-version                     # 程序自身的版本信息
|-app                         # 业务包
| |-ioc.go                    # IOC控制反转层  
| |-host                      # 业务对象
|   |-model.go                # 业务需要的数据模型
|   |-interface.go            # 业务接口
|   |-impl                    # 业务具体实现
|   |-http                    # http请求处理
|-main.go                     # 程序入口文件
```

