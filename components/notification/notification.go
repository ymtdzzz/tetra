package notification

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	notificationDuration = time.Second * 4
)

type notification struct {
	message   string
	createdAt time.Time
	// TODO: error, warning, info, success
}

type Model struct {
	styles        styles
	notifications []notification
}

func New() Model {
	return Model{
		styles:        defaultStyles(),
		notifications: []notification{},
	}
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		_    tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case NotificationMsg:
		notification := notification{
			message:   msg.Message,
			createdAt: time.Now(),
		}
		m.notifications = append(m.notifications, notification)
	case notificationCleanTickMsg:
		m.cleanNotifications(msg.t)
		cmds = append(cmds, m.tick())
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	notifications := []string{}

	for _, n := range m.notifications {
		notifications = append(notifications, m.styles.notificationBase.Render(n.message))
	}

	return lipgloss.JoinVertical(lipgloss.Top, notifications...)
}

func (m *Model) cleanNotifications(t time.Time) {
	var kept []notification
	for _, n := range m.notifications {
		if t.Sub(n.createdAt) <= notificationDuration {
			kept = append(kept, n)
		}
	}
	m.notifications = kept
}

func (m *Model) tick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return notificationCleanTickMsg{t}
	})
}

func (m *Model) UpdateLayout(width, height int) {
	m.styles.notificationBase = m.styles.notificationBase.Width(width)
}
