package main

import (
	"log"

	"github.com/spf13/cobra"
	buaaclock "github.com/wangjq4214/buaa-clock"
)

var (
	clockUsername string
	clockPassword string
	clockConfig   string
	retry         int
)

var clockCmd = &cobra.Command{
	Use:   "clock",
	Short: "Clock daily report.",
	Run: func(cmd *cobra.Command, args []string) {
		if clockConfig == "" {
			c := buaaclock.NewClock(buaaclock.Config{
				UserName: clockUsername,
				Password: clockPassword,
				Retry:    retry,
			})

			if err := c.Exec(); err != nil {
				log.Println(err)
			}
			return
		}

		cfg, err := readConfigFile(clockConfig)
		if err != nil {
			log.Println(err)
			return
		}

		errs := clockForMutileUser(cfg)
		for _, v := range errs {
			log.Println(v)
		}
	},
}

func init() {
	clockCmd.Flags().StringVarP(&clockUsername, "username", "u", "", "Your buaa username.")

	clockCmd.Flags().StringVarP(&clockPassword, "password", "p", "", "Your buaa password.")

	clockCmd.Flags().StringVarP(&clockConfig, "config", "c", "", "Your config file path")

	clockCmd.Flags().IntVarP(&retry, "retry", "r", 10, "The clock retry times.")

	rootCmd.AddCommand(clockCmd)
}
