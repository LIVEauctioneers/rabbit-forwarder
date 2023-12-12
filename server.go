package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/LIVEauctioneers/rabbit-amazon-forwarder/config"
	"github.com/LIVEauctioneers/rabbit-amazon-forwarder/mapping"
	"github.com/LIVEauctioneers/rabbit-amazon-forwarder/supervisor"
	"github.com/LIVEauctioneers/service-utils/awsutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

const (
	LogLevel = "LOG_LEVEL"
)

func main() {
	createLogger()

	var err error
	config.MappingJSON, err = downloadMapping()
	if err != nil {
		log.WithField("error", err.Error()).Fatalf("Could not download mapping file from S3")
	}

	consumerForwarderMapping, err := mapping.New().Load()
	if err != nil {
		log.WithField("error", err.Error()).Fatalf("Could not load consumer - forwarder pairs")
	}
	supervisor := supervisor.New(consumerForwarderMapping)
	if err := supervisor.Start(); err != nil {
		log.WithField("error", err.Error()).Fatal("Could not start supervisor")
	}
	http.HandleFunc("/restart", supervisor.Restart)
	http.HandleFunc("/health", supervisor.Check)
	log.Info("Starting http server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func downloadMapping() ([]byte, error) {
	s3Client, err := awsutils.NewS3Client(context.Background())
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s/rabbit-amazon-forwarder/mapping.json", os.Getenv("DEPLOYMENT"))

	obj, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("la-config"),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	b, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func createLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	if logLevel := os.Getenv(LogLevel); logLevel != "" {
		if level, err := log.ParseLevel(logLevel); err != nil {
			log.Fatal(err)
		} else {
			log.SetLevel(level)
		}
	}
}
