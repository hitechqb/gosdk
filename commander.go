package sdk

import (
	"github.com/spf13/cobra"
	"os"
)

type commander struct {
	rootCommand *cobra.Command
	service     Service
}

func NewCommander(service Service) *commander {
	var rootCmd = &cobra.Command{
		Short: "Choose below options",
	}

	return &commander{rootCommand: rootCmd, service: service}
}

func (s *commander) AddCommand(use, desc string, run func(Service)) {
	s.rootCommand.AddCommand(&cobra.Command{
		Use:   use,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			run(s.service)
		},
	})
}

func (s *commander) Execute() {
	if err := s.rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
