package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
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

// hotReloadTLSDialer returns a dial function that builds a fresh tls.Config,
// reading the CA cert and client cert/key from disk, on every connection rather
// than once at startup. Since franz-go dials afresh whenever it (re)connects,
// this lets the client pick up rotated certs (kept in sync from the Strimzi
// secrets) without needing a restart, while keeping standard TLS verification.
func hotReloadTLSDialer(caCertPath, clientCertPath, clientKeyPath string) func(context.Context, string, string) (net.Conn, error) {
	netDialer := &net.Dialer{Timeout: 10 * time.Second}
	return func(ctx context.Context, network, host string) (net.Conn, error) {
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA cert: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, errors.New("failed to parse CA certificate")
		}
		clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate/key: %w", err)
		}
		serverName, _, err := net.SplitHostPort(host)
		if err != nil {
			return nil, fmt.Errorf("unable to split host:port for dialing: %w", err)
		}
		return (&tls.Dialer{
			NetDialer: netDialer,
			Config: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{clientCert},
				ServerName:   serverName,
			},
		}).DialContext(ctx, network, host)
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

	// Create kafka client. franz-go clones the TLS config on every dial, so the
	// hot-reloading callbacks below are re-run whenever the client (re)connects
	// after Strimzi rotates the cluster/client certs.
	topic := "app.samples.kafka-app.my-topic"
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(bootstrapServer),
		kgo.Dialer(hotReloadTLSDialer(
			"/usr/local/share/ca-certificates/kafka_ca.crt",
			"/opt/kafka/user-certs/user.crt",
			"/opt/kafka/user-certs/user.key",
		)),
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
