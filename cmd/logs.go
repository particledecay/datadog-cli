package cmd

import (
	"fmt"

	"github.com/particledecay/datadog-cli/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func logsCmd() *cobra.Command {
	opts := &api.LogQueryOpts{}

	command := &cobra.Command{
		Use:   "logs",
		Short: "query logs",
		Long:  `submit a query to retrieve a list of logs with optional parameters`,
		RunE: func(cmd *cobra.Command, args []string) error {
			authed, err := api.Validate()
			if err != nil {
				log.Err(err).Msgf("%v", err)
			}
			if authed == false {
				log.Fatal().Msg("unable to authenticate")
			}
			cmd.SilenceUsage = true

			logs, err := api.Logs(opts)
			if err != nil {
				log.Err(err).Msgf("%v", err)
			}

			for _, logObj := range logs {
				svc, err := logObj.GetTag("service")
				if err != nil {
					log.Err(err).Msgf("%v", err)
				}
				env, err := logObj.GetTag("env")
				if err != nil {
					log.Err(err).Msgf("%v", err)
				}

				msgSvc := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, svc) // yellow
				msgEnv := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 36, env) // cyan
				fmt.Printf("%s(%s): %s\n", msgSvc, msgEnv, logObj.Content.Message)
			}
			return nil
		},
	}

	command.Flags().IntVarP(&opts.Limit, "limit", "l", 100, "maximum number of log messages to retrieve")
	command.Flags().StringVarP(&opts.Query, "query", "q", "", "optional query to filter log messages")
	command.Flags().StringVarP(&opts.Sort, "sort", "s", "desc", "either 'asc' or 'desc' (default 'desc')")
	command.Flags().StringVarP(&opts.Time.From, "from", "f", "", "beginning of time range to filter log messages")
	command.Flags().StringVarP(&opts.Time.To, "to", "t", "", "end of time range to filter log messages")

	return command
}
