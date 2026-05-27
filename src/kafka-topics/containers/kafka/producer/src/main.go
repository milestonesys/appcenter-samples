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
	bootstrapServer := os.Getenv("KAFKA_CLUSTER_BOOTSTRAP_SERVER")
	if bootstrapServer == "" {
		return errors.New("error getting bootstrap server from environment")
	}

	// Initialize channel we can use for testing the interrupt signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Load CA cert
	caCert, err := os.ReadFile("/usr/local/share/ca-certificates/kafka-ca.crt")
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
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(bootstrapServer),
		kgo.DialTLSConfig(&tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{clientCert},
		}),
	)

	// cl, err := kgo.NewClient(kgo.SeedBrokers(bootstrapServer))
	if err != nil {
		return err
	}
	defer cl.Close()

	topic := "samples.my-topic"

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
