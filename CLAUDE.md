# CLAUDE.md

## 1. 项目概览

`cli-mock` 是一个轻量级 mock CLI，用来给 shell 脚本、自动化流程和测试场景提供稳定、可控的命令行行为。它不是通用任务执行器，而是有意只提供少量、边界清晰的模拟动作。

当前核心能力：

- 输出版本信息：`version`
- 延迟执行：`sleep`
- 输出到标准输出或标准错误：`echo`、`stderr`
- 返回指定退出码或失败：`exit`、`fail`
- 处理结构化和列表输出：`json`、`args`、`lines`
- 提交 JSON 驱动的业务表单：`leave`、`expense`、`procurement`
- 读取运行环境：`env`、`stdin`
- 生成带时间间隔的流式输出：`stream`
- 创建和读取 mock XDG 风格环境树：`xdg apply`、`xdg inspect`
- 输出标准化帮助文本：`help` / `<command> --help`

## 2. 技术栈

- 语言：Go
- CLI 框架：`spf13/cobra`
- 构建与测试：`go build`、`go test`
- 版本信息：`internal/buildinfo`
- 发布：仓库内 `scripts/release/build.sh`
- 依赖管理：Go Modules

当前通过 `replace github.com/spf13/cobra => ./third_party/cobra` 使用仓库内的本地依赖副本。

## 3. 架构设计

模块分层很简单：

- `cmd/mock`
  进程入口，只负责接收参数并把退出码交给操作系统。
- `internal/app`
  负责根命令组装、子命令定义、错误与退出码语义。
- `internal/app/xdg.go`
  负责 XDG manifest 校验、目录落盘、tree inspect 和 reveal 逻辑。
- `internal/buildinfo`
  负责版本字符串输出。
- `third_party/cobra`
  本地 vendored CLI 依赖。

核心调用链：

1. `main` 调用 `app.Execute`
2. `Execute` 构建 root command，并注入 `stdin/stdout/stderr`
3. Cobra 完成命令匹配和参数分发
4. 子命令返回普通错误、usage 错误或带退出码的 `exitError`
5. `Execute` 统一把错误映射为进程退出码

退出码语义：

- `0`
  成功
- `1`
  业务失败，例如 `fail` 或 `env` 读取不到变量
- `2`
  用法错误，例如未知命令或非法参数

## 4. 目录结构

- `cmd/mock`
  CLI 程序入口
- `internal/app`
  根命令、子命令、参数解析、退出码语义、测试
- `internal/buildinfo`
  版本信息
- `scripts/release`
  本地发布打包脚本
- `third_party/cobra`
  本地 Cobra 依赖

文档职责：

- `README.md`
  面向使用者的快速开始、示例、简单排查
- `CLAUDE.md`
  项目事实、模块分工、开发约定

## 5. 数据结构

当前项目没有数据库，但现在包含一个轻量级的文件树 manifest 模型，用来描述 mock `.config` / `.local` 环境。

### `exitError`

`internal/app` 定义了：

- `Code int`
- `Err error`

它用来把命令执行错误和进程退出码绑定起来，供 `Execute` 统一处理。

### 版本信息

`internal/buildinfo` 维护三个包级变量：

- `Version`
- `Commit`
- `BuildTime`

默认值分别是 `dev`、`none`、`unknown`，构建时可以通过 `-ldflags -X ...` 覆盖。

### JSON 输出

`json` 和 `args` 命令都使用 JSON 文本作为外部交换格式：

- `json <raw-json>`
  解析任意合法 JSON，然后以紧凑格式重新输出
- `args <arg...>`
  把位置参数编码成 JSON 数组

### XDG Manifest

`xdg apply` 读取一个 JSON manifest：

- 顶层：`entries`
- entry 字段：
  - `path`
  - `type`
  - `format`
  - `content`
  - `mode`

约束：

- `path` 必须是相对路径
- `path` 只能落在 `.config/**` 或 `.local/**`
- `type=dir` 不能带 `content`
- `type=file` 必须带 `content`
- `format=json` 会写成规范 JSON 文本

## 6. API 定义

`cli-mock` 的外部接口只有 CLI，没有 HTTP API。

核心命令面：

- `mock version`
- `mock sleep <duration>`
- `mock echo <text...>`
- `mock stderr <text...>`
- `mock exit <code>`
- `mock fail [message...]`
- `mock json <raw-json>`
- `mock args <arg...>`
- `mock env <key>`
- `mock stdin`
- `mock lines <count>`
- `mock stream <count> --interval <duration>`
- `mock leave --payload <json>`
- `mock expense --payload <json>`
- `mock procurement --payload <json>`
- `mock xdg apply --root <dir> --manifest <path-or-> [--overwrite]`
- `mock xdg inspect --root <dir> [--reveal]`

接口约定：

- 大多数参数错误走 usage 退出码 `2`
- 明确的模拟失败走退出码 `1`
- `stdin` 不接收位置参数
- `env` 要求目标变量存在
- `lines` 和 `stream` 要求 `count > 0`
- `sleep` 和 `stream --interval` 使用 Go duration 语法
- `stream` 在提供自定义内容时，内容条目数必须和 `count` 一致
- `leave`、`expense`、`procurement` 只接受 `--payload` 形式的 JSON object 输入
- JSON 业务表单的语法/缺字段错误走 usage 退出码 `2`
- JSON 业务表单的业务校验失败走退出码 `1`
- `xdg apply` 只接受 JSON manifest，且只允许写入 `.config/**` 和 `.local/**`
- `xdg inspect` 默认只暴露 metadata；只有 `--reveal` 才返回可读文本或 JSON 内容
- `xdg` 子命令始终要求显式 `--root`，不会默认改写当前真实 home
- 帮助输出使用统一 section 格式，按命令能力展示 `Description`、`Flags`、`Args fields`、`Params fields`、`Examples`

## 7. 开发要点

- 新增子命令时，优先保持行为单一、输出稳定、易于脚本断言。
- 对复杂业务 mock，优先用单个 JSON payload 表达层级字段，而不是扩展大量 flags。
- 对外错误信息尽量直接，避免模糊错误文本，因为测试通常会断言它们。
- 新增或修改帮助文本时，优先复用统一帮助元数据，而不是在命令执行逻辑里手写输出。
- 如果新增命令改变用户可见行为，需要同步更新：
  - `README.md`
  - `CLAUDE.md`
  - `internal/app/app_test.go`
- `Execute` 是统一退出码出口，新增错误类型时不要绕开它。
- XDG manifest 保持轻量，优先服务 mock 配置树场景，不要演进成通用配置管理系统。

## 8. 开发流程

本地开发常用命令：

```bash
go build -o ./mock ./cmd/mock
go test ./...
```

修改 CLI 行为后，至少做以下验证：

```bash
./mock help
./mock stream --help
./mock version
./mock echo hello
./mock fail broken
./mock stream 2 --interval 10ms
./mock stream 2 hello world --interval 10ms
./mock xdg apply --root /tmp/mock-home --manifest ./manifest.json
./mock xdg inspect --root /tmp/mock-home --reveal
```

如果修改了退出码、错误文本或帮助信息，优先补或改 `internal/app/app_test.go`，避免行为漂移。

发布流程当前由维护者手工执行：

1. 确认代码和文档已提交
2. 运行 `go test ./...`
3. 创建并推送 tag，例如 `v0.1.0`
4. 执行 `scripts/release/build.sh v0.1.0`
5. 校验 `dist/<version>/mock_<version>_checksums.txt`
6. 分发或上传 release 产物

## 9. 已知约束与注意事项

- 这是 mock 工具，不负责模拟复杂交互式 TTY 行为。
- `stream` 通过 `time.Sleep` 实现，适合轻量测试，不适合高精度计时场景。
- `json` 只做解析和紧凑输出，不保留原始格式或注释。
- 本地依赖通过 `third_party/cobra` 固定；升级 Cobra 时要注意仓库内副本同步。
- `xdg inspect --reveal` 只返回 UTF-8 文本或合法 JSON 文件内容；二进制文件仍只显示 metadata。
