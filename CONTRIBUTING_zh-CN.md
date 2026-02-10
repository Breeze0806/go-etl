# go-etl 贡献指南

感谢您对 go-etl 项目的兴趣！我们欢迎并感谢各种形式的贡献，包括但不限于提交 Bug 报告、功能建议、代码改进、文档完善等。

## 目录

- [行为准则](#行为准则)
- [开始贡献](#开始贡献)
- [开发环境搭建](#开发环境搭建)
- [项目结构](#项目结构)
- [开发新插件](#开发新插件)
- [代码规范](#代码规范)
- [提交代码](#提交代码)
- [测试](#测试)
- [文档](#文档)
- [获取帮助](#获取帮助)

## 行为准则

请尊重所有参与项目的贡献者，保持友善和专业的沟通态度。我们希望创建一个开放、包容的社区，任何不当行为都是不可接受的。

## 开始贡献

### 报告 Bug

如果您发现了 Bug，请通过 GitHub Issues 进行报告。报告时请包含以下信息：

- 问题描述
- 复现步骤
- 期望行为
- 实际行为
- 环境信息（操作系统、Go 版本等）
- 可能的解决方案

### 提出功能建议

如果您有新功能的建议，欢迎通过 Git Issues 进行讨论。请描述：

- 您希望实现的功能
- 使用场景
- 可能的实现思路
- 任何相关的参考资料

### 提交 Pull Request

1. Fork 本项目
2. 创建您的功能分支：`git checkout -b feature/amazing-feature`
3. 提交您的更改：`git commit -m 'Add some amazing feature'`
4. 推送您的分支：`git push origin feature/amazing-feature`
5. 提交 Pull Request

## 开发环境搭建

### 环境要求

- Go 1.20 或更高版本
- GCC 4.8 或更高版本（Linux）
- MinGW-w64 环境（Windows，GCC 7.2.0 或更高）

### 获取源码

```bash
cd ${GO_PATH}/src
git clone https://github.com/Breeze0806/go-etl.git "github.com/Breeze0806/go-etl"
cd github.com/Breeze0806/go-etl
```

### 安装依赖

**Linux：**
```bash
make dependencies
```

**Windows：**
```bash
release.bat
```

### 编译项目

**Linux：**
```bash
make release
```

**Windows：**
```bash
release.bat
```

### 排除 DB2 依赖

如果不需要 DB2 支持，可以在编译前设置环境变量：

**Linux：**
```bash
export IGNORE_PACKAGES=db2
make dependencies
make release
```

**Windows：**
```bash
set IGNORE_PACKAGES=db2
release.bat
```

## 项目结构

go-etl 项目采用 Framework + 插件架构，主要包含以下模块：

```
go-etl/
├── datax/                    # 数据同步框架
│   ├── plugin/
│   │   ├── reader/          # 读取插件
│   │   └── writer/          # 写入插件
│   └── ...
├── element/                  # 数据类型和类型转换
├── storage/                  # 存储模块
│   ├── database/            # 数据库集成
│   └── stream/             # 数据流处理
├── tools/                   # 工具集合
│   ├── datax/              # 构建和发布工具
│   └── license/            # 许可证工具
└── ...
```

### 核心模块说明

- **datax**: 类似于阿里巴巴 DataX 的离线数据同步框架
- **element**: 数据类型定义和类型转换
- **storage/database**: 数据库基础集成和方言接口
- **storage/stream/file**: 文件解析（CSV、Excel 等）
- **tools/build**: 插件注册和代码生成
- **tools/license**: 自动添加许可证

## 开发新插件

go-etl 支持通过插件机制扩展数据源。以下是开发新插件的指南。

### 创建 Reader 插件

1. 使用模板生成工具创建插件框架：

```bash
cd tools/datax/plugin
go run main.go -t reader -p Mysql
```

2. 修改生成的文件：
   - 更新 `plugin.json` 中的插件信息
   - 实现 `job.go` 中的 Job 接口
   - 实现 `task.go` 中的 Task 接口
   - 实现 `reader.go` 中的 Reader 接口

3. 注册插件：

```bash
go generate ./...
```

### 创建 Writer 插件

1. 使用模板生成工具创建插件框架：

```bash
cd tools/datax/plugin
go run main.go -t writer -p Mysql
```

2. 修改生成的文件：
   - 更新 `plugin.json` 中的插件信息
   - 实现 `job.go` 中的 Job 接口
   - 实现 `task.go` 中的 Task 接口
   - 实现 `writer.go` 中的 Writer 接口

3. 注册插件：

```bash
go generate ./...
```

### 插件接口说明

#### Reader 插件接口

Reader 插件需要实现以下接口：

**Job 接口：**
```go
Init(ctx context.Context) (err error)
Destroy(ctx context.Context) (err error)
Split(ctx context.Context, number int) ([]*config.JSON, error)
Prepare(ctx context.Context) error
Post(ctx context.Context) error
```

**Task 接口：**
```go
Init(ctx context.Context) (err error)
Destroy(ctx context.Context) (err error)
StartRead(ctx context.Context, sender plugin.RecordSender) error
Prepare(ctx context.Context) error
Post(ctx context.Context) error
```

#### Writer 插件接口

Writer 插件需要实现以下接口：

**Job 接口：**
```go
Init(ctx context.Context) (err error)
Destroy(ctx context.Context) (err error)
Split(ctx context.Context, number int) ([]*config.JSON, error)
Prepare(ctx context.Context) error
Post(ctx context.Context) error
```

**Task 接口：**
```go
Init(ctx context.Context) (err error)
Destroy(ctx context.Context) (err error)
StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error
Prepare(ctx context.Context) error
Post(ctx context.Context) error
SupportFailOver() bool
```

### 数据库插件开发

如果实现关系型数据库插件，建议：

1. 参考 [Database Storage Developer Guide](storage/database/README.md)
2. 实现 `Querier` 接口（Reader）或 `Execer` 接口（Writer）
3. 使用 `dbms.StartRead` 或 `dbms.StartWrite` 函数

### 文件流插件开发

如果实现二维表文件插件（如 CSV、Excel），建议：

1. 参考 [Two-dimensional Table File Stream Storage Developer Guide](storage/stream/file/README.md)
2. 实现文件解析和生成逻辑

### 插件配置设计原则

- 使用驼峰命名法
- 配置项正交，无重叠功能
- 合理使用 JSON 类型
- 参考同类插件的约定

## 代码规范

### 代码格式

项目使用 `gofmt` 进行代码格式化：

```bash
gofmt -s -w yourfile.go
```

### 添加许可证

在提交代码前，请运行以下命令自动添加许可证：

```bash
go run tools/license/main.go
```

### 命名规范

- 包名：简洁、有意义
- 函数名：驼峰命名法
- 常量名：全大写，下划线分隔
- 变量名：驼峰命名法

### 注释规范

- 公共函数和类型需要注释
- 复杂逻辑需要详细注释
- 注释使用英文

## 提交代码

### Commit 消息规范

- 使用英文描述
- 简明扼要描述更改内容
- 包含相关的 Issue 编号（如有）

示例：
```
Add MySQL reader plugin support

Implement basic MySQL data reading functionality with batch fetch.

Fixes #123
```

### Pull Request 要求

- 代码必须通过所有测试
- 遵循项目代码规范
- 包含必要的文档更新
- 描述清楚 PR 的目的和内容

## 测试

### 运行测试

```bash
go test ./...
```

### 测试覆盖率

项目鼓励提高测试覆盖率。新功能应包含相应的单元测试。

### 性能测试

对于性能相关的改动，请提供性能测试报告，包括：

- 测试环境（硬件、操作系统等）
- 测试数据特征
- 测试参数配置
- 性能对比数据

## 文档

### 代码注释

- 公共 API 必须有注释
- 复杂逻辑需要详细说明
- 注释使用英文

### README 文档

新插件需要包含以下文档内容：

1. **快速介绍**: 插件功能和使用场景
2. **实现原理**: 底层实现原理
3. **配置说明**: JSON 配置示例和参数说明
4. **类型转换**: 数据类型转换规则
5. **性能报告**: 性能测试数据
6. **约束限制**: 使用限制和注意事项
7. **常见问题**: FAQ

### 更新相关文档

如果您的更改影响了用户使用，请更新：

- `README_USER.md`: 用户手册
- `README_USER_zh-CN.md`: 中文用户手册
- 插件文档

## 获取帮助

如果您在贡献过程中遇到问题，可以通过以下方式获取帮助：

- 查看 [项目文档](README.md)
- 查看 [开发者文档](datax/README.md)
- 查看 [用户手册](README_USER.md)
- 提交 GitHub Issue 讨论
- QQ 群：185188648

## 致谢

感谢所有为 go-etl 项目做出贡献的开发者！
