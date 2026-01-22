# 🚀 AICommits

**AICommits** 是一个基于大语言模型（LLM）的智能 CLI 工具，旨在通过 AI 自动化生成符合语义规范（Conventional Commits）的 Git 提交日志。

## ✨ 核心特性

* **⚡️ 极速启动**：基于 Go 语言编写，编译为二进制文件，启动无需等待运行时。
* **🧠 多模型支持**：完美支持 **DeepSeek (V3/R1)**、OpenAI (GPT-4o/mini) 以及本地模型 (**Ollama**)。
* **🎨 交互式 UI**：基于 Bubble Tea 构建的现代化终端界面，支持加载动画、表单配置和确认流程。
* **📝 高度可配**：支持多语言（中/英）、详细描述模式及自定义 Prompt 规则。

## 📦 安装

### 方式 1: 直接下载 (推荐)

`curl -sS https://raw.githubusercontent.com/zshnb/aicommits/main/install.sh | bash`

### 方式 2: 使用 Go 安装

```bash
go install github.com/zshnb/aicommits@latest

```

### 方式 3: 源码编译

```bash
git clone https://github.com/zshnb/aicommits.git
cd aicommits
go build -o aicommits

```

## ⚙️ 配置向导

首次使用，请运行配置命令。工具会通过交互式表单引导你完成设置：

```bash
aicommits config

```

你将看到如下界面：

```text
? 选择 AI 提供商
> DeepSeek
  OpenAI
  Ollama (本地)
  自定义

? API Key
> sk-xxxxxxxxxxxxxxxx

? 提交日志语言
> 中文 (Chinese)
  English

```

配置文件将保存在 `~/.aicommits.yaml`。

## 🚀 使用指南

### 1. 基础生成

当你已经 `git add` 了文件后，直接运行：

```bash
aicommits

```

工具将读取暂存区的 Diff，生成提交日志，并等待你确认。

* 按 `Enter`：确认并提交。
* 按 `r`：重新生成。
* 按 `Esc`：取消。

### 2. 自动暂存并生成 (`--add`)

如果你想一次性提交所有变动（相当于 `git commit -a`）：

```bash
aicommits --add
# 或者简写
aicommits -a

```

## 💻 本地开发

如果你想参与贡献：

1. 克隆仓库。
2. 安装依赖：`go mod tidy`。
3. 运行测试（需自行补充测试文件）：`go test ./...`。
4. 提交 PR！