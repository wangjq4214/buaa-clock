package main

import (
	"log"
	"os"

	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"
)

var (
	pool *ants.Pool

	rootCmd = &cobra.Command{
		Use:   "buaa-clock",
		Short: "Use the command line to clock daily report.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	var err error
	pool, err = ants.NewPool(1000)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
