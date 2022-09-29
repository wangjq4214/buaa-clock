package main

import (
	"fmt"

	buaaclock "github.com/wangjq4214/buaa-clock"
)

func main() {
	clock := buaaclock.NewClock(buaaclock.Config{
		UserName: "by2106105",
		Password: "wjq35113616",
	})

	fmt.Println(clock.Exec())
}
