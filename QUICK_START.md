# go-etl Quick Start (3 Minutes)

## What is go-etl?

A data sync tool for transferring data between MySQL, PostgreSQL, Oracle, SQL Server, CSV, XLSX, and more.

## Install

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

## Start Syncing

### 1. Create config file `config.json`

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

### 2. Run sync

**Linux:**
```bash
./go-etl -c config.json
```

**Windows:**
```powershell
.\go-etl.exe -c config.json
```

See `run success`? You're done!

## Common Configs

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

## Batch Sync

Use a wizard CSV file to batch sync multiple tables.

**1. Create data source config `config.json`** - same as single sync

**2. Create wizard file `wizard.csv`** - each row defines a source-target table pair:
```csv
source_table,target_table
table1,table1_copy
table2,table2_copy
```

**3. Generate batch configs and run script:**

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

## Links

- For more details, please refer to the [User Manual](README_USER.md)
