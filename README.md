# cli-mock

## 1. 项目简介

`cli-mock` 是一个给脚本、自动化流程和测试用例使用的 mock CLI。它提供一组可预测的命令，用来模拟常见命令行程序行为，例如：

- 输出标准输出或标准错误
- 返回指定退出码
- 读取环境变量和标准输入
- 生成 JSON、逐行输出和延迟流式输出
- 通过无状态 CRUD 表单命令模拟较重的业务场景
- 创建和检查可直接使用的 `.config` / `.local` mock 环境树

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
./mock create-leave --payload '{"applicant_id":"E1001","department_id":"engineering","leave_type":"annual","start_date":"2026-04-20","end_date":"2026-04-22","days":2.5,"reason":"family_trip"}'
./mock expense add --payload-file ./expense.json --result approved --output json
printf '{"requester_id":"E1001","department":"engineering","budget_code":"RD-2026-001","reason":"team expansion","delivery_city":"Shanghai","items":[{"name":"MacBook Pro","quantity":2,"unit_price":18999,"vendor":"Apple"}],"approvers":["MGR100","FIN200"],"requested_at":"2026-04-14T11:00:00+08:00"}' | ./mock procurement update --request-id PR-BA08D42C31 --payload-stdin --result rejected
./mock expense get --request-id EX-14C0A7B992 --result not_found
./mock delete-leave --request-id LV-7B0A3D4F10
./mock xdg apply --root /tmp/mock-home --manifest ./manifest.json
./mock xdg inspect --root /tmp/mock-home
```

### 测试

```bash
go test ./...
```

## 3. 配置说明

`cli-mock` 仍然是以命令参数、环境变量和标准输入为主的 mock CLI，但现在额外支持通过 `mock xdg` 创建和读取一个显式 root 下的 `.config` / `.local` 环境树。

- 命令参数：
  `sleep`、`exit`、`lines`、`stream` 等命令通过位置参数控制行为。
- JSON 业务表单：
  业务 mock 采用 3 种不同命令风格：请假继续使用平铺命令，例如 `create-leave`；报销使用资源分组，例如 `expense add`；采购使用资源分组加标准 CRUD，例如 `procurement create`。
- 表单输入来源：
  业务 `create` / `add` / `update` 命令支持 `--payload`、`--payload-file`、`--payload-stdin` 三种输入方式，但一次只能使用一种。
- 结果分支：
  业务 CRUD 命令支持 `--result` 显式控制返回状态，例如 `submitted`、`approved`、`rejected`、`found`、`not_found`、`deleted`。
- 输出格式：
  业务命令默认输出稳定的结构化文本；显式传 `--output json` 时返回原有 JSON 响应。
- 环境变量：
  `env <key>` 会读取指定环境变量；变量不存在时返回退出码 `1`。
- 标准输入：
  `stdin` 会把输入内容原样复制到标准输出。
- XDG 环境树：
  `xdg apply` 根据 JSON manifest 写入 `.config/**` 和 `.local/**`，`xdg inspect` 读取该树并输出 JSON 摘要。

`mock xdg` 的 v1 约定：

- 只接受 JSON manifest
- 必须显式传 `--root`
- 只允许创建或读取 `.config/**` 和 `.local/**`
- `inspect` 默认只输出 metadata；需要内容时再加 `--reveal`

最小 manifest 示例：

```json
{
  "entries": [
    {
      "path": ".config/demo/config.toml",
      "type": "file",
      "format": "text",
      "content": "token = \"demo\"\n"
    },
    {
      "path": ".local/share/demo/secret.json",
      "type": "file",
      "format": "json",
      "content": { "api_key": "demo-key" }
    },
    {
      "path": ".local/state/demo",
      "type": "dir",
      "mode": "0700"
    }
  ]
}
```

常见调用：

```bash
./mock xdg apply --root /tmp/mock-home --manifest ./manifest.json
./mock xdg inspect --root /tmp/mock-home
./mock xdg inspect --root /tmp/mock-home --reveal
```

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

如果你是普通使用者，优先使用维护者提供的 release 压缩包：

- macOS Apple Silicon：`mock_vX.Y.Z_darwin_arm64.tar.gz`
- macOS Intel：`mock_vX.Y.Z_darwin_amd64.tar.gz`
- Linux ARM64：`mock_vX.Y.Z_linux_arm64.tar.gz`
- Linux AMD64：`mock_vX.Y.Z_linux_amd64.tar.gz`

维护者可以在仓库根目录执行本地打包脚本：

```bash
scripts/release/build.sh v0.1.0
```

产物会生成到：

```bash
dist/v0.1.0/
```

其中包含 4 个平台压缩包和一个校验文件：

- `mock_v0.1.0_darwin_amd64.tar.gz`
- `mock_v0.1.0_darwin_arm64.tar.gz`
- `mock_v0.1.0_linux_amd64.tar.gz`
- `mock_v0.1.0_linux_arm64.tar.gz`
- `mock_v0.1.0_checksums.txt`

解压后可直接验证：

```bash
tar -xzf mock_v0.1.0_darwin_arm64.tar.gz
./mock version
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
./mock xdg inspect --root /tmp/mock-home
```

### 常见排查

- 命令报 `unknown command`：先执行 `./mock help` 确认子命令名。
- `sleep` 或 `stream --interval` 报 duration 错误：使用 Go duration 格式，例如 `20ms`、`1s`。
- `stream` 传了自定义内容却失败：确认内容条目数和 `count` 完全一致。
- 业务 `create` / `add` / `update` 命令报 payload 输入错误：确认只传了 `--payload`、`--payload-file`、`--payload-stdin` 其中一种。
- `update-leave` 失败：确认 payload 里包含 `request_id`，且前缀是 `LV-`。
- 业务 `get` / `delete` 失败：确认传了 `--request-id`，且前缀和业务类型匹配，例如 `EX-` 对应 `expense`。
- JSON 表单命令报字段错误：确认顶层是 JSON object，必填字段齐全，日期格式分别使用 `YYYY-MM-DD` 或 RFC3339。
- `expense` 报金额不匹配：确认 `total_amount` 等于所有 `items[].amount` 之和。
- `procurement` 报 `budget exceeded`：当前 mock 会在采购总金额超过 `50000` 时返回业务失败。
- `--result` 报非法值：使用命令帮助中列出的动作级枚举，例如 create/update 用 `submitted|approved|rejected`。
- 需要机器可解析输出时：为业务命令显式加 `--output json`。
- `xdg apply` 报路径错误：确认 manifest 中的 `path` 是相对路径，并且以 `.config/` 或 `.local/` 开头。
- `xdg apply` 报覆盖错误：已有文件默认不会覆盖，重写时加 `--overwrite`。
- `xdg inspect` 输出里看不到文件内容：默认只返回 metadata，需要时加 `--reveal`。
- 发布脚本报 `invalid version`：使用 `v0.1.0` 这种 tag 风格版本号。
- `exit` 报参数非法：退出码必须在 `0` 到 `255` 之间。
- `env` 失败：确认目标环境变量已经导出到当前 shell。
- `json` 失败：传入的字符串必须是合法 JSON。

## 7. 进一步阅读

- [CLAUDE.md](./CLAUDE.md)
  项目事实、架构分层和开发约定
