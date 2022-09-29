package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"
	buaaclock "github.com/wangjq4214/buaa-clock"
	"gopkg.in/yaml.v3"
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
	configPath string

	pool *ants.Pool
)

var timingCmd = &cobra.Command{
	Use:   "timing",
	Short: "Clock everyday.",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		pool, err = ants.NewPool(1000)
		if err != nil {
			log.Println(err)
			return
		}

		file, err := os.Open(configPath)
		if err != nil {
			log.Println(err)
			return
		}

		content, err := io.ReadAll(file)
		if err != nil {
			log.Println(err)
			return
		}

		cfg := c{}
		err = yaml.Unmarshal(content, &cfg)
		if err != nil {
			log.Println(err)
			return
		}

		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Day().At("17:10").Do(func() {
			for _, v := range cfg.Users {
				pool.Submit(func() {
					config := buaaclock.Config{
						UserName: v.UserName,
						Password: v.Password,
					}

					if v.Boarder == "0" {
						config.Boarder = "0"
						config.NotBoarderReasen = v.NotBoarderReasen
						config.NotBoarderNote = v.NotBoarderNote
						config.NotBoarderAddress = v.Address
						config.NotBoarderArea = v.Area
						config.NotBoarderCity = v.City
						config.NotBoarderProvince = v.Province
					} else {
						config.Boarder = "1"
						config.BoarderAddress = v.Address
						config.BoarderArea = v.Area
						config.BoarderCity = v.City
						config.BoarderProvince = v.Province
					}

					clock := buaaclock.NewClock(config)
					if err := clock.Exec(); err != nil {
						log.Println(err)
						return
					}
				})
			}
		})

		s.StartBlocking()
	},
}

func init() {
	timingCmd.Flags().StringVarP(&configPath, "config", "c", "config.yml", "Your config path.")
	timingCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(timingCmd)
}
