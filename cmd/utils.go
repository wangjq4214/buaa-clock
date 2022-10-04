package main

import (
	"io"
	"os"
	"sync"

	buaaclock "github.com/wangjq4214/buaa-clock"
	"gopkg.in/yaml.v3"
)

func readConfigFile(path string) (*c, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg := c{}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func clockForMutileUser(cfg *c) []error {
	w := sync.WaitGroup{}
	ch := make(chan error)
	res := make([]error, 0)

	for _, user := range cfg.Users {
		w.Add(1)

		func(v item) {
			pool.Submit(func() {
				defer w.Done()

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
					ch <- err
				}
			})
		}(user)
	}

	go func() {
		w.Wait()
		close(ch)
	}()

	for k := range ch {
		res = append(res, k)
	}

	return res
}
