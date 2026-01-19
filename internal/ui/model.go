package ui

import (
	"context"
	"fmt"

	"aicommits/internal/llm"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 定义 UI 的三种状态
type sessionState int

const (
	stateLoading sessionState = iota // 正在调用 LLM
	stateReview                      // 显示结果，等待用户确认
	stateError                       // 出错了
)

// Model 存储 UI 的所有数据状态
type Model struct {
	client llm.Client // LLM 客户端
	diff   string     // Git diff 内容
	ctx    context.Context

	state   sessionState  // 当前状态
	Msg     string        // 生成的 Commit Message
	err     error         // 错误信息
	spinner spinner.Model // 加载动画组件

	Confirmed bool // 用户是否确认提交 (用于返回给主程序)
}

// NewModel 初始化模型
func NewModel(ctx context.Context, client llm.Client, diff string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // 粉色 Loading

	return Model{
		client:  client,
		diff:    diff,
		ctx:     ctx,
		state:   stateLoading, // 初始状态为加载中
		spinner: s,
	}
}

// Init 程序启动时执行的第一个命令
func (m Model) Init() tea.Cmd {
	// 同时启动两个任务：
	// 1. 启动 Spinner 动画
	// 2. 触发 generateCmd 去调用 LLM
	return tea.Batch(m.spinner.Tick, m.generateMsgCmd)
}

// Update 核心逻辑：根据消息更新状态
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// 按键处理
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit // 退出程序
		case "enter":
			if m.state == stateReview {
				m.Confirmed = true
				return m, tea.Quit // 确认并退出
			}
		case "r":
			if m.state == stateReview {
				m.state = stateLoading
				m.Msg = ""
				return m, tea.Batch(m.spinner.Tick, m.generateMsgCmd) // 重新生成
			}
		}

	// 生成完成的消息
	case generatedMsg:
		m.state = stateReview
		m.Msg = string(msg)
		return m, nil

	// 错误消息
	case errMsg:
		m.state = stateError
		m.err = error(msg)
		return m, tea.Quit

	// Spinner 动画帧更新
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View 渲染界面：根据状态返回字符串
func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return fmt.Sprintf("\n %s 正在思考 Git 提交日志...\n\n", m.spinner.View())

	case stateReview:
		// 使用 lipgloss 渲染漂亮的边框
		boxStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)

		tipsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

		return fmt.Sprintf(
			"\n%s\n\n%s\n",
			boxStyle.Render(m.Msg),
			tipsStyle.Render("Confirm: [Enter] | Retry: [r] | Cancel: [Esc]"),
		)

	case stateError:
		return fmt.Sprintf("\n❌ 发生错误: %v\n", m.err)
	}

	return ""
}

// 定义消息类型
type generatedMsg string
type errMsg error

// generateMsgCmd 实际调用 LLM 的异步任务
func (m Model) generateMsgCmd() tea.Msg {
	res, err := m.client.GenerateCommitMessage(m.ctx, m.diff)
	if err != nil {
		return errMsg(err)
	}
	return generatedMsg(res)
}
