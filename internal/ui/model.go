package ui

import (
	"context"
	"fmt"

	"aicommits/internal/llm"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput" // 1. 引入 textinput
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateLoading sessionState = iota
	stateReview
	stateEditing
	stateError
)

type Model struct {
	client llm.Client
	diff   string
	ctx    context.Context

	state     sessionState
	Msg       string
	err       error
	spinner   spinner.Model
	textInput textinput.Model // 2. 改为 textInput

	Confirmed bool
}

func NewModel(ctx context.Context, client llm.Client, diff string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// 3. 初始化单行输入框
	ti := textinput.New()
	ti.Placeholder = "在此编辑提交信息..."
	ti.Focus()
	ti.CharLimit = 0 // 可以限制长度，或者设为 0 (不限制)
	ti.Width = 100   // 显示宽度

	return Model{
		client:    client,
		diff:      diff,
		ctx:       ctx,
		state:     stateLoading,
		spinner:   s,
		textInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.generateMsgCmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch m.state {

		// --- 预览状态 ---
		case stateReview:
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			case "enter":
				m.Confirmed = true
				return m, tea.Quit
			case "r":
				m.state = stateLoading
				m.Msg = ""
				return m, tea.Batch(m.spinner.Tick, m.generateMsgCmd)
			case "e":
				m.state = stateEditing
				// 进入编辑模式时，把当前消息填进去，并把光标移到最后
				m.textInput.SetValue(m.Msg)
				m.textInput.CursorEnd()
				return m, textinput.Blink
			}

		// --- 编辑状态 ---
		case stateEditing:
			switch msg.String() {
			// 4. 单行模式下，回车(Enter)通常意味着“完成编辑”
			case "enter", "esc":
				m.Msg = m.textInput.Value() // 保存修改
				m.state = stateReview       //以此返回预览界面
				return m, nil
			}
			// 透传按键给输入框
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case generatedMsg:
		m.state = stateReview
		m.Msg = string(msg)
		return m, nil

	case errMsg:
		m.state = stateError
		m.err = error(msg)
		return m, tea.Quit

	case spinner.TickMsg:
		if m.state == stateLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return fmt.Sprintf("\n %s 正在思考...\n\n", m.spinner.View())

	case stateReview:
		boxStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(60)

		tipsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)

		content := m.Msg
		if content == "" {
			content = "(空)"
		}

		return fmt.Sprintf(
			"\n%s\n%s\n",
			boxStyle.Render(content),
			tipsStyle.Render("Confirm: [Enter] | Edit: [e] | Retry: [r] | Cancel: [Ctrl+C or Esc]"),
		)

	case stateEditing:
		// 5. 渲染单行输入框样式
		return fmt.Sprintf(
			"\n 编辑提交信息 (Enter 保存):\n\n %s\n\n",
			m.textInput.View(),
		)

	case stateError:
		return fmt.Sprintf("\n❌ Error: %v\n", m.err)
	}

	return ""
}

// 辅助类型保持不变
type generatedMsg string
type errMsg error

func (m Model) generateMsgCmd() tea.Msg {
	res, err := m.client.GenerateCommitMessage(m.ctx, m.diff)
	if err != nil {
		return errMsg(err)
	}
	return generatedMsg(res)
}
