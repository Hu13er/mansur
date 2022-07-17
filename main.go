package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Hu13er/telegrus"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	var (
		botToken string
		chatID   int64
	)

	botToken = os.Getenv("MANSUR_TOKEN")
	if botToken == "" {
		log.Fatalln("MANSUR_TOKEN not provided")
	}

	chatIDs := os.Getenv("MANSUR_CHATID")
	if chatIDs == "" {
		log.Fatalln("MANSUR_CHATID not provided")
	}
	_, err := fmt.Sscanf(chatIDs, "%d", &chatID)
	if err != nil {
		log.Fatalln("MANSUR_CHATID must be valid integer, it is:", chatIDs)
	}

	minIdles := os.Getenv("MANSUR_MINIDLE")
	var minIdle int
	if minIdles == "" {
		log.Warnln("MANSUR_MINIDLE not provided")
	} else {
		_, err := fmt.Sscanf(minIdles, "%d", &minIdle)
		if err != nil {
			log.Fatalln("MANSUR_MINIDLE must be valid integer, it is:", minIdles)
		}
	}

	intervals := os.Getenv("MANSUR_INTERVAL")
	var interval int
	if intervals == "" {
		interval = 10
		log.Infoln("MANSUR_INTERVAL not provided, default 10 secs")
	} else {
		_, err := fmt.Sscanf(intervals, "%d", &interval)
		if err != nil {
			log.Fatalln("MANSUR_INTERVAL must be valid integer, it is:", intervals)
		}
	}

	warns := os.Getenv("MANSUR_MENTION")
	warnlst := make([]string, 0)
	if warns == "" {
		log.Warnln("MANSUR_MENTION not provided")
	} else {
		warnlst = strings.SplitN(warns, ",", -1)
	}

	log.AddHook(
		telegrus.NewHooker(botToken, chatID).
			MentionOn(logrus.WarnLevel,
				warnlst...).
			SetLevel(logrus.InfoLevel),
	)

	t := time.NewTicker(time.Second * time.Duration(interval))
	check := func(err error) {
		if err != nil {
			log.Error("Error getting cpu usage:", err)
		}
	}

	before, err := cpu.Get()
	check(err)
	for range t.C {
		now, err := cpu.Get()
		check(err)

		var (
			total = float64(now.Total - before.Total)
			user  = float64(now.User-before.User) / total * 100.0
			sys   = float64(now.System-before.System) / total * 100.0
			idle  = float64(now.Idle-before.Idle) / total * 100.0
		)
		if idle <= float64(minIdle) {
			log.Warnf("user, sys, idle: %f, %f, %f", user, sys, idle)
		} else {
			log.Infof("user, sys, idle: %f, %f, %f", user, sys, idle)
		}
		before = now
	}
}
