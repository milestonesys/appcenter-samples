package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	os.Exit(mainReturnWithCode())
}

func mainReturnWithCode() int {
	if err := mainReturnWithError(); err == nil {
		return 0
	} else {
		log.Println(err)
		return 1
	}
}

func mainReturnWithError() error {
	// Get bootstrap server from environment
	bootstrapServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		return errors.New("error getting bootstrap server from environment")
	}

	// Initialize channel we can use for testing the interrupt signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Create kafka client
	topic := "samples.my-topic"
	cl, err := kgo.NewClient(kgo.SeedBrokers(bootstrapServer))
	if err != nil {
		return err
	}
	defer cl.Close()

	// Produce messages until told otherwise
	ctx, cancel := context.WithCancel(context.Background())
	wgrp := &sync.WaitGroup{}
	wgrp.Add(1)
	go func() {
		defer wgrp.Done()
		for i := range 99999999 {
			record := &kgo.Record{Topic: topic, Value: []byte("This is a test message number " + strconv.Itoa(int(i)))}
			fmt.Println("producing '" + string(record.Value) + "'")
			if err := cl.ProduceSync(ctx, record).FirstErr(); err != nil {
				fmt.Println(err)
				fmt.Println("no more messages will be produced")
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Wait for interrupt signal
	<-signals
	cancel()
	wgrp.Wait()

	// Success
	return nil
}
