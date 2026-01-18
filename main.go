package main

import (
	"github.com/zjyl1994/cherries/infra/startup"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := startup.Startup(); err != nil {
		logrus.Fatalln("Startup failed:", err.Error())
	}
}
