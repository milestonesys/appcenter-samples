package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	bootstrapServer := os.Getenv("KAFKA_CLUSTER_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		return errors.New("error getting bootstrap server from environment")
	}

	// Initialize channel we can use for testing the interrupt signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Load CA cert
	caCert, err := os.ReadFile("/usr/local/share/ca-certificates/kafka_ca.crt")
	if err != nil {
		log.Fatalf("Failed to read CA cert: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to parse CA certificate")
	}

	// Load client certificate and key for mTLS
	clientCertPath := "/opt/kafka/user-certs/user.crt"
	clientKeyPath := "/opt/kafka/user-certs/user.key"
	clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalf("Failed to load client certificate/key: %v", err)
	}

	// Create kafka client
	topic := "samples.my-topic"
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(bootstrapServer),
		kgo.DialTLSConfig(&tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{clientCert},
		}),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),
		kgo.ConsumeStartOffset(kgo.NewOffset().AtEnd()),
		kgo.ConsumeTopics(topic))

	if err != nil {
		return err
	}
	defer cl.Close()

	// Consume messages until told otherwise
	ctx, cancel := context.WithCancel(context.Background())
	wgrp := &sync.WaitGroup{}
	wgrp.Add(1)
	go func() {
		defer wgrp.Done()
		for {
			// Wait for fetches to be available
			fetches := cl.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				fmt.Println(err)
				fmt.Println("no more messages will be consumed")
				return
			}

			// Iterate through all records in fetches and print them to standard output
			fetches.EachRecord(func(record *kgo.Record) {
				fmt.Println("consuming '" + string(record.Value) + "'")
			})
		}
	}()

	// Wait for interrupt signal
	<-signals
	cancel()
	wgrp.Wait()

	// Success
	return nil
}
