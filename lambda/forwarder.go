package lambda

import (
	"context"
	"errors"
	"fmt"

	"github.com/LIVEauctioneers/rabbit-amazon-forwarder/config"
	"github.com/LIVEauctioneers/rabbit-amazon-forwarder/forwarder"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	log "github.com/sirupsen/logrus"
)

const (
	// Type forwarder type
	Type = "Lambda"
)

// Forwarder forwarding client
type Forwarder struct {
	name         string
	lambdaClient IFace
	function     string
}

type IFace interface {
	Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}

// CreateForwarder creates instance of forwarder
func CreateForwarder(entry config.AmazonEntry, lambdaClient ...IFace) forwarder.Client {
	var client IFace
	if len(lambdaClient) > 0 {
		client = lambdaClient[0]
	} else {
		cfg, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithRegion("us-east-1"))
		if err != nil {
			panic(fmt.Sprintf("unable to load SDK config, %v", err))
		}
		client = lambda.NewFromConfig(cfg)
	}
	forwarder := Forwarder{entry.Name, client, entry.Target}
	log.WithField("forwarderName", forwarder.Name()).Info("Created forwarder")
	return forwarder
}

// Name forwarder name
func (f Forwarder) Name() string {
	return f.name
}

// Push pushes message to forwarding infrastructure
func (f Forwarder) Push(message string) error {
	if message == "" {
		return errors.New(forwarder.EmptyMessageError)
	}
	params := &lambda.InvokeInput{
		FunctionName: aws.String(f.function),
		Payload:      []byte(message),
	}
	resp, err := f.lambdaClient.Invoke(context.Background(), params)
	if err != nil {
		log.WithFields(log.Fields{
			"forwarderName": f.Name(),
			"error":         err.Error()}).Error("Could not forward message")
		return err
	}
	if resp.FunctionError != nil {
		log.WithFields(log.Fields{
			"forwarderName": f.Name(),
			"functionError": *resp.FunctionError}).Errorf("Could not forward message")
		return errors.New(*resp.FunctionError)
	}
	log.WithFields(log.Fields{
		"forwarderName": f.Name(),
		"statusCode":    resp.StatusCode}).Info("Forward succeeded")
	return nil
}
