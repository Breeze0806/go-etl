Here's the translated version of your technical documentation:

# go-etl prometheus docker

## 1. Deploy Prometheus Environment

```bash
cd docker
docker compose up -d
```

### 1.1 Check the go-etl Container

```bash
docker exec -it etl bash
```

### 1.2 Access PostgreSQL

Connect using a client like [DBeaver](https://github.com/dbeaver/dbeaver/releases):

1. Connecting to `postgres_source`
```
127.0.0.1:5432
db1
user1
password1
```

2. Connecting to `postgres_dest`
```
127.0.0.1:5433
db2
user2
password2
```

### 1.3 Check Prometheus

Access via browser: http://127.0.0.1:9090

### 1.4 Check Grafana

Access via browser: http://127.0.0.1:3000

## 2. Generate Test Data

### 2.1 Create Tables in PostgreSQL

1. Connecting to `postgres_source` and executing SQL:
```sql
CREATE SCHEMA source;

CREATE TABLE source.split (
	id bigint,
	dt date,
	str varchar(10)
);
```

2. Connecting to `postgres_dest` and executing SQL:
```sql
CREATE SCHEMA dest;

CREATE TABLE dest.split (
	id bigint,
	dt date,
	str varchar(10)
);
```

### 2.2 Prepare Data Files

Place data files and `import_config.json` in the ETL mounted directory:

```bash
cd docker/data
go run test/main.go
```

## 3. Run ETL Operations

### 3.1 Execute ETL Job

1. Import data into source PostgreSQL
```bash
docker exec -it etl release/bin/go-etl -http :6080 -c data/import_config.json
```

2. Migrate data from source to target PostgreSQL
```bash
docker exec -it etl release/bin/go-etl -http :6080 -c data/config.json
```

### 3.2 Monitor in Prometheus

1. Check job status:
```
up{job="etl"}
```

2. Total synchronized records:
```
datax_channel_total_record
```

3. Records per second:
```
rate(datax_channel_total_record{job="etl"}[1m])
```

4. Total data volume:
```
datax_channel_total_byte
```

5. Data throughput:
```
rate(datax_channel_total_byte{job="etl"}[30s])
```

### 3.3 Import Grafana Configuration

Import `go-etl-grafana.json` in Grafana at http://127.0.0.1:3000
