package cmd

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	var logType string

	command := &cobra.Command{
		Use:   "datadog",
		Short: "datadog interfaces the Datadog API",
		Long:  `datadog provides a feature complete interface to the Datadog API`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if logType == "text" {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{"time", "level"}})
			} else {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
				zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("an action is required")
			}
			return nil
		},
	}

	command.PersistentFlags().StringVarP(&logType, "output", "o", "text", "logging output type ('text', 'json')")

	return command
}

func Execute() {
	root := rootCmd()
	root.AddCommand(logsCmd())

	err := root.Execute()
	if err != nil {
		log.Err(err).Msgf("%v", err)
	}
}
