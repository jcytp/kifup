package aws

import (
	"context"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/jcytp/kifup-api/common/env"
)

var sesClient *ses.Client

func SesInitialize() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(env.AwsRegion()))
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}
	sesClient = ses.NewFromConfig(cfg)
}

func SesSendEmailOne(fromAddress, toAddress, subject, message string) {
	input := &ses.SendEmailInput{
		Source: aws.String(fromAddress),
		Destination: &types.Destination{
			ToAddresses: []string{toAddress},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String(subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data:    aws.String(message),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	_, err := sesClient.SendEmail(context.TODO(), input)
	if err != nil {
		slog.Error("failed to send email", "error", err)
	}
}
