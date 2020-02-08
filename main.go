package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/dirkarnez/wait2die"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	Sites map[string]struct{
		URL             string `yaml:"url"`
		Evaluate        string `yaml:"evaluate"`
		IntervalSeconds int    `yaml:"interval-seconds"`
	} `yaml:"config"`
}

func MyBeep() {
	beep := syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")
	for i := 0; i < 10; i++ {
		beep.Call(0xffffffff)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	file, err := ioutil.ReadFile("config.yml")
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf))

	for key, s := range config.Sites {
		go Crawl(taskCtx, key, s.URL, s.Evaluate, s.IntervalSeconds)
	}

	wait2die.WaitToDie(nil)
}

func Crawl(taskCtx context.Context, name, url, evaluate string, intervalSeconds int) {
	ctx, _ := chromedp.NewContext(taskCtx)
	var notAvailable = true

	for notAvailable {
		var err error
		var result []byte

		chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Evaluate(evaluate, &result),
		)

		notAvailable, err = strconv.ParseBool(string(result))
		if err != nil {
			notAvailable = true
		}

		if notAvailable {
			log.Printf(`%s is not available`, name)
		} else {
			log.Printf(`%s is available`, name)
		}

		time.Sleep(time.Duration(intervalSeconds) * time.Second)
	}

	MyBeep()

	for { select{ } }
}