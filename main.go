package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/ymtdzzz/tetra/components/app"
)

var version = "unknown"

func main() {
	rootCmd := &cobra.Command{
		Use:     "tetra",
		Short:   "A TUI SQL IDE",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := app.New()
			if err != nil {
				return err
			}
			defer m.Close()
			_, err = tea.NewProgram(m, tea.WithAltScreen()).Run()
			return err
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Println("Error running program:", err)
		os.Exit(1)
	}
}
