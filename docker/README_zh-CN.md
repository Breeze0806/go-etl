# go-etl监控prometheus部署

## 1 部署prometheus环境

### 1.1 启动prometheus环境

将`prometheus.yml`的`192.168.188.1`替换成你主机的任一网卡地址

```yml
global:
  scrape_interval: 1s

scrape_configs:
  - job_name: 'etl'
    static_configs:
      - targets: ['192.168.188.1:6080']
    metrics_path: '/metrics'
```

使用docker命令启动

```bahs
cd docker
docker compose up -d
```

查看go-etl相关容器


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

1. 连接`postgres_source`，如果没有发现表`source.split`，请执行SQL
```sql
CREATE SCHEMA source;

CREATE TABLE source.split (
	id bigint,
	dt date,
	str varchar(10)
);
```

2. 连接`postgres_dest，如果没有发现表`dest.split`，请执行SQL
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

1. 修改`import_config.json`,将其中的`192.168.188.1`替换成你主机的任一网卡地址

2. 在源postgres导入数据

   ```bash
   docker exec -it etl release/bin/go-etl -http :6080 -c data/import_config.json
   ```

3. 修改`config.json`,将其中的`192.168.188.1`替换成你主机的任一网卡地址

4. 将源postgres的数据同步到目标postgres

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
以下是在 Grafana 中添加 Prometheus 数据源的详细步骤：

#### 3.3.1 登录 Grafana
打开浏览器访问 Grafana 地址（默认 `http://127.0.0.1:3000`），使用管理员账号登录（默认用户名/密码：`admin/admin`）。

#### 3.3.2 进入数据源配置页面
1. 点击左侧菜单栏的 **⚙️ Configuration（齿轮图标）**。
2. 选择 **Data Sources**。

#### 3.3.3 添加 Prometheus 数据源
1. 点击 **Add data source**。
2. 在搜索框中输入 `Prometheus`，选择 **Prometheus** 数据源类型。

#### 3.3.4 配置 Prometheus 连接
- **HTTP Settings**:
  - **URL**: 填写 Prometheus 地址  
    - 本地部署：`http://192.168.188.1:9090`,`192.168.188.1`是你主机的任一网卡的地址

#### 3.3.5 保存并测试
1. 点击底部 **Save & test**。
2. 看到 **Data source is working** 提示表示成功。

3.3.6 验证数据源

1. 先点击**Dashboards**，再点击**New**，然后点击**Import**，选择`go-etl-grafana.json`导入
2. 点击**datasource**，选择`Prometheus`