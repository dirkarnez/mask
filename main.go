package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/dirkarnez/wait2die"
	"log"
	"time"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	go crawl(allocCtx)
	go crawl(allocCtx)

	wait2die.WaitToDie(nil)
}

func crawl(context context.Context) {
	taskCtx, cancel := chromedp.NewContext(
		context,
		chromedp.WithLogf(log.Printf))
	defer cancel()

	var example string
	for example != "false" {
		var result []byte
		chromedp.Run(taskCtx,
			chromedp.Navigate(`https://www.bonjourhk.com/tc/search/%E5%8F%A3%E7%BD%A9/1030503`),
			chromedp.Evaluate(`document.evaluate("(//div[@id='content']/div[@class='row']//div[@class='white']/div[@class='row']/div)[2]", document, null, XPathResult.ANY_TYPE, null).iterateNext().innerHTML.indexOf("暫時缺貨") > -1`, &result),
		)

		example = string(result)
		log.Println(example)
		time.Sleep(5 * time.Second)
	}
}