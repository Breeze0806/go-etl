# go-etl用户手册

## 从源码进行编译

### linux

#### 依赖

1. golang 1.16以及以上

#### 构建
```bash
make dependencies
make release
```

### windows

#### 依赖
1.需要mingw-w64 with gcc 7.2.0以上的环境进行编译
2.golang 1.16以及以上
3.最小编译环境为win7 

#### 构建
```bash
release.bat
```

## 如何开始

下载对应操作系统的datax，在linux下如Makefile所示export LD_LIBRARY_PATH=${DB2HOME}/lib，否则无法运行

可以使用[ibm db2](https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/)以及[oracle](https://www.oracle.com/database/technologies/instant-client/downloads.html)下载到对应64位版本odbc依赖，也可以在**QQ群185188648**群共享中中下载到。

### 查看版本

```
datax version
v0.1.0 (git commit: c82eb302218f38cd3851df4b425256e93f85160d) complied by go version go1.16.5 windows/amd64
```

### 使用方式

```bash
Usage of datax:
  -c string
        config (default "config.json")
  -w string
        wizard
```

### 批量生成配置集和执行脚本

```bash
datax -c tools/testData/xlsx.json -w tools/testData/wizard.csv 
```
-c 指定数据源配置文件 -w 指定源目的配置向导文件。

执行结果会在数据源配置文件目录文件生成源目的配置向导文件行数的配置集，分别以指定数据源配置文件1.json,指定数据源配置文件2.json,...,指定数据源配置文件[n].json的配置集。

另外，在当前目录会生成执行脚本run.bat或者run.sh。
#### 数据源配置文件

数据源配置文件是json文件，使用数据源相互组合，如从mysql同步到postgres中
```json
{
    "core" : {
        "container": {
            "job":{
                "id": 1,
                "sleepInterval":100
            },
            "taskGroup":{
                "id": 1,
                "failover":{
                    "retryIntervalInMsec":0
                }
            }
        },
        "transport":{
            "channel":{
                "speed":{
                    "byte": 100,
                    "record":100
                }
            }
        }
    },
    "job":{
        "content":[
            {
                "reader":{
                    "name": "mysqlreader",
                    "parameter": {
                        "username": "test:",
                        "password": "test:",
                        "column": ["*"],
                        "connection":  {
                                "url": "tcp(192.168.15.130:3306)/source?parseTime=false",
                                "table": {
                                    "db":"source",
                                    "name":"type_table"
                                }
                            },
                        "where": ""
                    }
                },
                "writer":{
                    "name": "postgreswriter",
                    "parameter": {
                        "username": "postgres",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "preSql": [],
                        "connection":  {
                                "url": "postgres://192.168.15.130:5432/postgres?sslmode=disable&connect_timeout=2",
                                "table": {
                                    "schema":"destination",
                                    "name":"type_table"
                                }
                         },
                        "batchTimeout": "1s",
                        "batchSize":1000
                    }
                },
               "transformer":[]
            }
        ],
        "setting":{
            "speed":{
                "byte":3000,
                "record":400,
                "channel":4
            }
        }
    }
}
```

#### 源目的配置向导文件
源目的配置向导文件是csv文件，每行配置可以配置如下:
```csv
path[table],path[table]
```
每一列可以是路径或者是表名，注意所有的表要配置库名或者模式名，需要在数据源配置文件配置。

#### 测试结果
可以运行cmd/datax/testData的测试数据
```bash
cd cmd/datax
datax -c testData/xlsx.json -w testData/wizard.csv 
```
结果会在testData下生成wizard.csv行数的配置文件，分别以xlsx1.json,xlsx2.json,...,xlsx[n].json的配置集。

### 数据同步

调用datax十分简单，只要直接调用它即可

```bash
data -c config.json
```
-c 指定数据源配置文件

当返回值是0，并且显示run success,表示执行成功

当返回值是1，并且显示run fail,并告知执行失败的原因

#### 数据库全局配置

```json
{
    "job":{
        "setting":{
            "pool":{
              "maxOpenConns":8,
              "maxIdleConns":8,
              "connMaxIdleTime":"40m",
              "connMaxLifetime":"40m"
            },
            "retry":{
              "type":"ntimes",
              "strategy":{
                "n":3,
                "wait":"1s"
              },
              "ignoreOneByOneError":true
            }
        }
    }
}
```
##### 连接池pool
+ maxOpenConns: 最大连接打开数
+ maxIdleConns: 最大空闲连接打开数
+ connMaxIdleTime： 最大空闲时间
+ connMaxLifetime： 最大生存时间

##### 重试retry
ignoreOneByOneError 是否忽略一个个重试错误
+ 重试类型type和重试策略
1. 类型有`ntimes`,指n次数重复重试策略,`"strategy":{"n":3,"wait":"1s"}`,n代表重试次数，wait代表等待时间
2. 类型有`forever`,指永久重复重试策略,`"strategy":{"wait":"1s"}`,wait代表等待时间
3. 类型有`exponential`,指幂等重复重试策略,`"strategy":{"init":"100ms","max":"4s"}`,init代表开始时间，max代表最大时间

### 使用示例

#### 使用mysql同步

- 使用cmd/datax/examples/mysql/init.sql初始化数据库**用于测试**
- 开启同步mysql命令

```bash
datax -c examples/mysql/config.json
```

#### 使用postgres同步

- 使用cmd/datax/examples/postgres/init.sql初始化数据库**用于测试**
- 开启同步postgres命令

```bash
datax -c examples/postgres/config.json
```

#### 使用db2同步

- 注意使用前请下载相应的db2的odbc库，如linux的make dependencies和release.bat
- 注意在linux下如Makefile所示export LD_LIBRARY_PATH=${DB2HOME}/lib
- 注意在windows下如release.bat所示set path=%path%;%GOPATH%\src\github.com\ibmdb\go_ibm_db\clidriver\bin
- 使用cmd/datax/examples/db2/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/db2/config.json
```

#### 使用oracle同步

- 注意使用前请下载相应的[Oracle Instant Client]( https://www.oracle.com/database/technologies/instant-client/downloads.html)，例如，连接oracle 11g最好下载12.x版本。
- 注意在linux下如export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_1:$LD_LIBRARY_PATH，另需要安装libaio
- 注意在windows下如set path=%path%;%GOPATH%\oracle\instantclient_21_1，
Oracle Instant Client 19不再支持windows7
- 使用cmd/datax/examples/oracle/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/oracle/config.json
```

#### 使用sql server同步

- 使用cmd/datax/examples/sqlserver/init.sql初始化数据库**用于测试**
- 开启同步sql server命令

```bash
datax -c examples/sqlserver/config.json
```

#### 使用csv同步到postgres

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/csvpostgres/config.json
```

#### 使用xlsx同步到postgres

- 使用cmd/examples/datax/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/xlsxpostgres/config.json
```

#### 使用postgres同步csv

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/postgrescsv/config.json
```

#### 使用postgres同步xlsx

- 使用cmd/datax/examples/csvpostgres/init.sql初始化数据库**用于测试**
- 开启同步命令

```bash
datax -c examples/postgresxlsx/config.json
```

#### 使用切分键

这里假设数据按切分键分布是均匀的，合理使用这样的切分键可以使同步更快。

##### 测试方式
- 使用程序生成mysql数据产生split.csv
```bash
cd cmd/datax/examples/split
go run main.go
```
- 使用init.sql建表
- 同步至mysql数据库
```bash
cd ../..
datax -c examples/split/csv.json
```
- 修改examples/split/mysql.json的split的key为id,dt,str
- mysql数据库切分同步整形，日期，字符串类型
```bash
datax -c examples/split/mysql.json
```
#### 使用preSql和postSql

preSql和postSql分别是写入数据前和写入数据后的sql语句组

##### 测试方式
在本例子中，采用了全量导入的方式
1.写入数据前先建立了一个临时表
2.在写入数据后，将原表删除，将临时表重名为新表

```bash
datax -c examples/prePostSql/mysql.json
```

#### 其他同步例子

除了上述例子外，在go-etl特性中所列出的数据源都可以交叉使用，还配置例如mysql到postgresql数据源，mysql到oracle,oracle到db2等等，