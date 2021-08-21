package bot

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NewQueue(s *sqs.SQS) *queue {
	return &queue{s, make(chan *sqs.Message)}
}

type queue struct {
	s        *sqs.SQS
	messages chan *sqs.Message
}

func (q *queue) Send(uri string, msg interface{}) (err error) {
	body, err := json.Marshal(&msg)
	if err != nil {
		log.Error().Err(err).Msg("json marshal failed")
		return err
	}

	if _, err := q.s.SendMessage(&sqs.SendMessageInput{
		MessageBody:            aws.String(string(body)),
		QueueUrl:               aws.String(uri),
		MessageDeduplicationId: aws.String(uuid.New().String()),
		MessageGroupId:         aws.String(uuid.New().String()),
	}); err != nil {
		log.Error().Err(err).Msg("sqs publish failed")
		return err
	}
	return nil
}

func (q *queue) receive(uri string) error {
	input := &sqs.ReceiveMessageInput{QueueUrl: &uri}

	for {
		output, err := q.s.ReceiveMessage(input)
		if err != nil {
			return err
		}

		for _, message := range output.Messages {
			log.Info().Msgf("Received message: %s", message)
			q.messages <- message
		}
	}
}

func (q *queue) Listen(uri string, handle func([]byte) error) error {
	go func() {
		if err := q.receive(uri); err != nil {
			log.Fatal().Err(err).Msg("error while receive message from sqs")
		}
	}()

	for message := range q.messages {
		log.Info().Msgf("Message body: %s", *message.Body)
		if err := handle([]byte(*message.Body)); err != nil {
			log.Err(err).Msg("error while handle message body")
		} else if err := q.remove(uri, message); err != nil {
			log.Err(err).Msg("error while remove message")
		}
	}

	return nil
}

func (q *queue) remove(uri string, message *sqs.Message) error {
	_, err := q.s.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(uri),
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}
