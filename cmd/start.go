package cmd

import (
	"log"

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

		c, err := chat.NewChatClient(args[0], flagPort)
		if err != nil {
			log.Fatal("Error initializing chat client: ", err)
		}
		defer c.Close()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("port", "p", "0", "Local port to bind the client session")
}
