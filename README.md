# 影枢 - Cinexus

这是一个基于 Golang 的后台项目使用了以下技术栈：

- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [GORM](https://gorm.io/) - ORM 库，支持 MySQL 和 SQLite
- [Viper](https://github.com/spf13/viper) - 配置管理
- [Zap](https://github.com/uber-go/zap) - 日志库
- [JWT](https://github.com/golang-jwt/jwt) - JSON Web Token 认证

## 项目结构

```
.
├── config              # 配置文件和配置加载
│   ├── config.go
│   └── config.toml
├── internal            # 内部包
│   ├── controller      # 控制器层，处理HTTP请求
│   ├── database        # 数据库连接和管理
│   ├── middleware      # HTTP中间件
│   ├── model           # 数据模型
│   ├── router          # 路由定义
│   └── service         # 业务逻辑层
├── pkg                 # 可重用的包
│   ├── jwt             # JWT工具
│   └── logger          # 日志工具
├── logs                # 日志文件目录
├── main.go             # 应用入口
├── go.mod              # Go模块定义
└── README.md           # 项目说明
```

## 功能特性

- 完整的项目结构和分层设计
- 基于 JWT 的用户认证
- 支持 MySQL 和 SQLite 数据库
- 基于 TOML 的配置管理
- 高性能日志系统，支持按日期分割
- 优雅关闭服务器
- 中间件：日志记录、异常恢复、CORS 支持

## 快速开始

### 前置条件

- Go 1.16+
- MySQL 或 SQLite

### 安装

1. 克隆仓库

```bash
git clone https://github.com/yourusername/cinexus.git
cd cinexus
```

2. 安装依赖

```bash
go mod tidy
```

3. 配置

编辑 `config/config.toml` 文件，根据你的环境配置数据库等信息。

4. 运行

```bash
go run main.go
```

服务器将在配置的端口上启动（默认为 8080）。

## API 文档

### 认证相关

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册

### 用户相关

- `GET /api/v1/user/info` - 获取用户信息
- `PUT /api/v1/user/info` - 更新用户信息
- `PUT /api/v1/user/password` - 更新用户密码

## 许可证

MIT