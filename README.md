# cli-mock

## 1. 项目简介

`cli-mock` 是一个给脚本、自动化流程和测试用例使用的 mock CLI。它提供一组可预测的命令，用来模拟常见命令行程序行为，例如：

- 输出标准输出或标准错误
- 返回指定退出码
- 读取环境变量和标准输入
- 生成 JSON、逐行输出和延迟流式输出

如果你要看命令边界、目录分工和开发约定，请看 [CLAUDE.md](./CLAUDE.md)。

## 2. 快速开始

### 前置要求

- Go 1.26+

### 本地编译

```bash
go build -o ./mock ./cmd/mock
./mock version
```

### 常用示例

```bash
./mock help
./mock stream --help
./mock echo hello world
./mock stderr warning message
./mock exit 7
./mock json '{"ok":true,"count":2}'
printf 'first\nsecond\n' | ./mock stdin
./mock env HOME
./mock lines 3
./mock stream 3 --interval 100ms
./mock stream 3 hello world done --interval 100ms
```

### 测试

```bash
go test ./...
```

## 3. 配置说明

`cli-mock` 默认不依赖配置文件，也没有根目录 `.env.example` 契约。运行时配置主要来自命令参数、环境变量和标准输入。

- 命令参数：
  `sleep`、`exit`、`lines`、`stream` 等命令通过位置参数控制行为。
- 环境变量：
  `env <key>` 会读取指定环境变量；变量不存在时返回退出码 `1`。
- 标准输入：
  `stdin` 会把输入内容原样复制到标准输出。

当前没有额外的外部配置目录、YAML 或 TOML 文件。

## 4. 帮助输出格式

`cli-mock` 的帮助输出使用统一的分段格式，方便脚本测试和文本断言。根据命令能力不同，帮助页会按需展示这些 section：

- `Usage`
- `Description`
- `Available Commands`
- `Flags`
- `Args fields`
- `Params fields`
- `Examples`

示例：

```bash
./mock stream --help
./mock help env
```

## 5. 发布与分发

当前仓库没有独立的发布脚本，默认分发方式是本地构建二进制：

```bash
go build -o ./mock ./cmd/mock
```

如果后续需要发布版本，建议通过 `-ldflags` 注入版本号，覆盖默认的 `dev`：

```bash
go build -ldflags="-X github.com/linlay/cli-mock/internal/buildinfo.version=v0.1.0" -o ./mock ./cmd/mock
```

## 6. 简单验证与排查

### 简单验证

```bash
./mock help
./mock version
./mock echo hello
./mock json '{"name":"cli-mock"}'
./mock lines 2
printf 'demo\n' | ./mock stdin
```

### 常见排查

- 命令报 `unknown command`：先执行 `./mock help` 确认子命令名。
- `sleep` 或 `stream --interval` 报 duration 错误：使用 Go duration 格式，例如 `20ms`、`1s`。
- `stream` 传了自定义内容却失败：确认内容条目数和 `count` 完全一致。
- `exit` 报参数非法：退出码必须在 `0` 到 `255` 之间。
- `env` 失败：确认目标环境变量已经导出到当前 shell。
- `json` 失败：传入的字符串必须是合法 JSON。

## 7. 进一步阅读

- [CLAUDE.md](./CLAUDE.md)
  项目事实、架构分层和开发约定
