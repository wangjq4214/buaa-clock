package main

import (
	"log"

	"github.com/spf13/cobra"
	buaaclock "github.com/wangjq4214/buaa-clock"
)

var (
	clockUsername string
	clockPassword string
	retry         int
)

var clockCmd = &cobra.Command{
	Use:   "clock",
	Short: "Clock daily report.",
	Run: func(cmd *cobra.Command, args []string) {
		c := buaaclock.NewClock(buaaclock.Config{
			UserName: clockUsername,
			Password: clockPassword,
			Retry:    retry,
		})

		if err := c.Exec(); err != nil {
			log.Println(err)
		}
	},
}

func init() {
	clockCmd.Flags().StringVarP(&clockUsername, "username", "u", "", "Your buaa username.")
	clockCmd.MarkFlagRequired("username")

	clockCmd.Flags().StringVarP(&clockPassword, "password", "p", "", "Your buaa password.")
	clockCmd.MarkFlagRequired("password")

	clockCmd.Flags().IntVarP(&retry, "retry", "r", 10, "The clock retry times.")

	rootCmd.AddCommand(clockCmd)
}
