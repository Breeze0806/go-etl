# go-etl 3分钟入门

## 简介

go-etl是一个数据同步工具，可在MySQL、PostgreSQL、Oracle、SQL Server、CSV、XLSX等数据源之间同步数据。

## 安装

**Linux:**
```bash
wget https://github.com/Breeze0806/go-etl/releases/download/v0.2.3/go-etl-linux-amd64.tar.gz
tar -xzf go-etl-linux-amd64.tar.gz
```

**Windows:**
```powershell
# Using PowerShell
Invoke-WebRequest -Uri "https://github.com/Breeze0806/go-etl/releases/download/v0.2.3/go-etl-windows-amd64.tar.gz" -OutFile "go-etl-windows-amd64.tar.gz"
tar -xzf go-etl-windows-amd64.tar.gz
```

**Docker:**
```bash
docker pull go-etl:v0.2.3
docker run -d -p 6080:6080 --name etl -v /data:/usr/local/go-etl/data go-etl:v0.2.3
```

## 开始同步

### 1. 创建配置文件 `config.json`

```json
{
    "job": {
        "content": [
            {
                "reader": {
                    "name": "mysqlreader",
                    "parameter": {
                        "username": "root",
                        "password": "123456",
                        "column": ["*"],
                        "connection": {
                            "url": "tcp(127.0.0.1:3306)/source_db",
                            "table": {"db": "source_db", "name": "my_table"}
                        }
                    }
                },
                "writer": {
                    "name": "mysqlwriter",
                    "parameter": {
                        "username": "root",
                        "password": "123456",
                        "writeMode": "insert",
                        "column": ["*"],
                        "connection": {
                            "url": "tcp(127.0.0.1:3306)/dest_db",
                            "table": {"db": "dest_db", "name": "my_table"}
                        },
                        "batchSize": 1000
                    }
                }
            }
        ]
    }
}
```

### 2. 执行同步

**Linux:**
```bash
./go-etl -c config.json
```

**Windows:**
```powershell
.\go-etl.exe -c config.json
```

看到 `run success` 表示成功！

## 常用配置

**MySQL → PostgreSQL:**
```json
"reader": {
    "name": "mysqlreader",
    "parameter": {
        "url": "tcp(192.168.1.100:3306)/source_db",
        "username": "root",
        "password": "123456"
    }
},
"writer": {
    "name": "postgreswriter",
    "parameter": {
        "url": "postgres://192.168.1.100:5432/dest_db?sslmode=disable",
        "username": "postgres",
        "password": "123456"
    }
}
```

**CSV → MySQL:**
```json
"reader": {
    "name": "csvreader",
    "parameter": {
        "path": ["/data/*.csv"],
        "column": ["*"]
    }
},
"writer": {
    "name": "mysqlwriter",
    "parameter": {
        "url": "tcp(127.0.0.1:3306)/dest_db",
        "username": "root",
        "password": "123456"
    }
}
```

## 批量同步

使用向导 CSV 文件批量同步多个表。

**1. 创建数据源配置文件 `config.json`** - 与单表同步相同

**2. 创建向导文件 `wizard.csv`** - 每行定义一个源-目标表对：
```csv
source_table,target_table
table1,table1_copy
table2,table2_copy
```

**3. 生成批量配置并执行：**

**Linux:**
```bash
./go-etl -c config.json -w wizard.csv
./run.sh
```

**Windows:**
```powershell
.\go-etl.exe -c config.json -w wizard.csv
run.bat
```

## 文档链接

- 详细请见[用户文档](README_USER_zh-CN.md)
