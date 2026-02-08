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

1. Import data to source PostgreSQL:
```bash
docker exec -it etl release/bin/go-etl -http :6080 -c data/import_config.json
```

2. Sync data to destination PostgreSQL:
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

### 3.3 Import Grafana Configuration

The following are detailed steps for adding a Prometheus data source in Grafana:

#### 3.3.1 Login to Grafana

Open your browser and access the Grafana address (default `http://127.0.0.1:3000`), and log in with the administrator account (default username/password: `admin/admin`).

#### 3.3.2 Go to Data Source Configuration Page

1. Click on **⚙️ Configuration (gear icon)** in the left sidebar.
2. Select **Data Sources**.

#### 3.3.3 Add Prometheus Data Source

1. Click **Add data source**.
2. Enter `Prometheus` in the search box and select the **Prometheus** data source type.

#### 3.3.4 Configure Prometheus Connection

- **HTTP Settings**:
  - **URL**: Enter the Prometheus address
    - For local deployment: `http://etl:9090`

#### 3.3.5 Save and Test

1. Click **Save & test** at the bottom.
2. The message **Data source is working** indicates success.

#### 3.3.6 Verify Data Source

1. First click **Dashboards**, then click **New**, and then click **Import**, select `go-etl-grafana.json` to import.
2. Click **datasource** and select `Prometheus`.

This completes the monitoring setup for go-etl with Prometheus and Grafana.