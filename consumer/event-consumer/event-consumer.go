package event_consumer

import (
	"log"
	"sync"
	"time"

	"main.go/events"
)


type Consumer struct {
    fetcher events.Fetcher
    processor events.Processor
    bathSize int
}

func New(fetcher events.Fetcher, processor events.Processor, bathSize int)Consumer {
    return Consumer{
    fetcher: fetcher,
    processor: processor,
    bathSize: bathSize,
    } 
}

func (c Consumer) Start() error {
    for {
        gotEvents, err := c.fetcher.Fetch(c.bathSize)
        if err != nil {
            log.Printf("[ERR] consumer: %s", err.Error())
            continue
        }

        if len(gotEvents) == 0 {
            time.Sleep(1 * time.Second)
            continue
        }
        
        if err := c.handleEvents(gotEvents); err != nil {
            log.Print(err)
            continue
        }
    }
}

func (c *Consumer) handleEvents(events []events.Event) error {
    var wg sync.WaitGroup
    for _, event := range events {
        log.Println("got new event")
        tmp := event
        wg.Add(1)
        go func() {
            defer wg.Done()
            if err:= c.processor.Process(tmp);err != nil {
                log.Printf("can't handle event: %s", err.Error())
            }
        }()
    }
    wg.Wait()
    return nil
}
