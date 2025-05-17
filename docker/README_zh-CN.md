# go-etl prometheus docker

## 1 部署prometheus环境

```bash
cd docker
docker compose up -d
```

### 1.1 查看go-etl容器


```bash
docker exec -it etl bash
```

### 1.2 查看postgres

通过例如[DBeaver](https://github.com/dbeaver/dbeaver/releases)的客户端连接

连接postgres_source
```
127.0.0.1:5432
db1
user1
password1
```

连接postgres_dest
```
127.0.0.1:5433
db2
user2
password2
```

### 1.3 查看prometheus

使用http://127.0.0.1:9090在浏览器中查看

### 1.4 查看granafa!

使用http://127.0.0.1:3000在浏览器中查看

## 2 生成测试数据

### 2.1 在postgres建表

1. 连接postgres_source执行SQL
```sql
CREATE SCHEMA source;

CREATE TABLE source.split (
	id bigint,
	dt date,
	str varchar(10)
);
```

2. 连接postgres_dest执行SQL
```sql
CREATE SCHEMA dest;

CREATE TABLE dest.split (
	id bigint,
	dt date,
	str varchar(10)
);
```

### 2.2 获取表数据

将数据以及`import_config.json`放入etl的挂载盘

```bash
cd docker/data
go run main.go
```

## 3 获取测试数据

### 3.1 运行etl

1. 在源postgres导入数据
```bash
docker exec -it etl release/bin/go-etl -http :6080 -c data/import_config.json
```

2. 将源postgres的数据同步到目标postgres
```bash
docker exec -it etl release/bin/go-etl -http :6080 -c data/config.json
```

### 3.2 在prometheus查看

1.查看job是否在线，看看是不是1


```
up{job="etl"}
```

2.获取同步记录总数


```
datax_channel_total_record
```

3.获取每秒同步的记录数

```
rate(datax_channel_total_record{job="etl"}[1m])
```

4.获取同步数据量

```
datax_channel_total_byte
```

5.获取每秒同步的数据量

```
rate(datax_channel_total_byte{job="etl"}[30s])
```

### 3.3 导入grafana的配置

在http://127.0.0.1:3000导入`go-etl-grafana.json`