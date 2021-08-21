package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/max107/hack0820/pkg/bot"
	"github.com/rs/zerolog/log"
)

var c = new(cfg)

func init() {
	LoadConfig(c)
}

func createSess(c *cfg) (*session.Session, error) {
	awsCfg := aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(c.AwsRegion),
		Credentials:      credentials.NewStaticCredentials(c.AwsAccessKey, c.AwsSecretKey, ""),
		// Credentials: credentials.NewSharedCredentials("", "hack0820"),
	}
	return session.NewSession(&awsCfg)
}

func main() {
	sess, err := createSess(c)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create aws session")
		return
	}

	u := bot.NewS3Uploader(c.AwsBucket, s3.New(sess))
	q := bot.NewSQSQueue(sqs.New(sess))
	b := bot.NewBot(c.BotToken, c.JobURL, c.StateURL, q, u)
	if err := b.Init(); err != nil {
		log.Fatal().Err(err).Msg("failed to init telegram bot")
	}
	b.Listen()
}
