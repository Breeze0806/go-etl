# go-etl Monitoring with Prometheus Deployment

## 1. Deploy Prometheus Environment

### 1.1 Start Prometheus Environment

Start using Docker:

```bash
cd docker
docker compose up -d
```

Check go-etl-related containers:

```bash
docker exec -it etl bash
```

### 1.2 Access PostgreSQL

Connect using clients like [DBeaver](https://github.com/dbeaver/dbeaver/releases).

**Connect to postgres_source**:
```
127.0.0.1:5432
db1
user1
password1
```

**Connect to postgres_dest**:
```
127.0.0.1:5433
db2
user2
password2
```

### 1.3 Access Prometheus

Open in browser: `http://127.0.0.1:9090`.

### 1.4 Access Grafana

Open in browser: `http://127.0.0.1:3000`.

## 2. Generate Test Data

### 2.1 Create Tables in PostgreSQL

1. Connect to `postgres_source`. If `source.split` table doesn't exist, execute:
```sql
CREATE SCHEMA source;

CREATE TABLE source.split (
  id bigint,
  dt date,
  str varchar(10)
);
```

2. Connect to `postgres_dest`. If `dest.split` table doesn't exist, execute:
```sql
CREATE SCHEMA dest;

CREATE TABLE dest.split (
  id bigint,
  dt date,
  str varchar(10)
);
```

### 2.2 Prepare Data Files

Place data files and `import_config.json` in ETL's mounted volume:

```bash
cd docker/data
go run main.go
```

## 3. Process Test Data

### 3.1 Run ETL Jobs

1. Modify `import_config.json`.

2. Import data to source PostgreSQL:
```bash
docker exec -it etl release/bin/go-etl -http :6080 -c data/import_config.json
```

3. Modify `config.json`P.

4. Sync data to destination PostgreSQL:
```bash
docker exec -it etl release/bin/go-etl -http :6080 -c data/config.json
```

### 3.2 Monitor in Prometheus

1. Check job status:
```
up{job="etl"}
```

2. Get total sync records:
```
datax_channel_total_record
```

3. Get records per second:
```
rate(datax_channel_total_record{job="etl"}[1m])
```

4. Get total data volume:
```
datax_channel_total_byte
```

5. Get data volume per second:
```
rate(datax_channel_total_byte{job="etl"}[30s])
```

### 3.3 Configure Grafana Dashboard

#### 3.3.1 Login Grafana

Open browser: `http://127.0.0.1:3000` (default credentials: `admin/admin`).

#### 3.3.2 Configure Data Source

1. Click **⚙️ Configuration (gear icon)** → **Data Sources**.
2. Click **Add data source** → Search for **Prometheus**.
3. Configure:
   - **HTTP → URL**: `http://etl:9090` (replace with your IP)
4. Click **Save & Test**.

#### 3.3.3 Import Dashboard

1. Click **Dashboards → New → Import**.
2. Upload `go-etl-grafana.json`.
3. Select **Prometheus** as data source.

This completes the monitoring setup for go-etl with Prometheus and Grafana.