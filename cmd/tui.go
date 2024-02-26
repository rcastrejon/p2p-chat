package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rcastrejon/p2p-chat/cmd/chat"
	"github.com/rcastrejon/p2p-chat/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type replyMsg struct {
	alias string
	body  string
}

// Waits for a pb.Message from the chat client and returns a replyMsg. This is
// a command that can be used with tea.Batch to wait for a reply, once the
// reply is received, it should be called again to wait for the next message.
func waitForReply(conn *chat.ChatClient) tea.Cmd {
	return func() tea.Msg {
		m := conn.Receive()
		return replyMsg{m.GetAlias(), m.GetBody()}
	}
}

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	peerStyle   lipgloss.Style
	conn        *chat.ChatClient
	err         error
}

func initialModel(conn *chat.ChatClient) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 140

	ta.SetWidth(30)
	ta.SetHeight(2)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		peerStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		conn:        conn,
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForReply(m.conn),
		textarea.Blink,
	)
}

// TODO: Take care of the order of messages.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			// this will send a message to the peer with the alias of "Peer", so that the
			// receiver displays the message as coming from "Peer".
			message := &pb.Message{
				Alias:     "Peer",
				Body:      m.textarea.Value(),
				Timestamp: timestamppb.Now(),
			}
			m.conn.Send(message)

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	case replyMsg:
		m.messages = append(m.messages, m.peerStyle.Render(msg.alias+": ")+msg.body)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, tea.Batch(waitForReply(m.conn), tiCmd, vpCmd)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
		"(esc to quit)",
	) + "\n"
}
