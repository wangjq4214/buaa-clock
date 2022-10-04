package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
)

type url struct {
	Login string
	Info  string
	Save  string
}

type item struct {
	UserName string `yml:"username"`
	Password string `yml:"password"`

	Boarder string `yml:"boarder"`

	NotBoarderReasen string `yml:"reason"`

	NotBoarderNote string `yml:"note"`

	Address  string `yml:"address"`
	Area     string `yml:"area"`
	City     string `yml:"city"`
	Province string `yml:"province"`
}

type c struct {
	URl   url    `yml:"url"`
	Users []item `yml:"users"`
}

var (
	timingConfig string
)

var timingCmd = &cobra.Command{
	Use:   "timing",
	Short: "Clock everyday.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := readConfigFile(timingConfig)
		if err != nil {
			log.Println(err)
			return
		}

		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Day().At("17:10").Do(func() {
			errs := clockForMutileUser(cfg)
			for _, v := range errs {
				log.Println(v)
			}
		})

		s.StartBlocking()
	},
}

func init() {
	timingCmd.Flags().StringVarP(&timingConfig, "config", "c", "config.yml", "Your config path.")
	timingCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(timingCmd)
}
