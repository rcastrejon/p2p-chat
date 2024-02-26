package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rcastrejon/p2p-chat/cmd/chat"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <peer-host:port>",
	Short: "Start a chat session with a peer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flagPort, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal("Failed to get port flag:\n", err)
		}

		// Start the peer-to-peer chat client
		c, err := chat.NewChatClient(args[0], flagPort)
		if err != nil {
			log.Fatal("Error initializing chat client: ", err)
		}
		defer c.Close()

		// Initialize and run the chat ui
		p := tea.NewProgram(initialModel(c))
		if _, err := p.Run(); err != nil {
			log.Fatal("Error running chat ui: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("port", "p", "0", "Local port to bind the client session")
}
