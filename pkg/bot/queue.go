package bot

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NewSQSQueue(s *sqs.SQS) *Queue {
	return &Queue{
		s:        s,
		messages: make(chan *sqs.Message),
	}
}

type Queue struct {
	s        *sqs.SQS
	messages chan *sqs.Message
}

func (w *Queue) Send(uri string, msg interface{}) (err error) {
	body, err := json.Marshal(&msg)
	if err != nil {
		log.Error().Err(err).Msg("json marshal failed")
		return err
	}

	if err := w.Publish(uri, body); err != nil {
		log.Error().Err(err).Msg("sqs publish failed")
		return err
	}
	return nil
}

func (w *Queue) Publish(uri string, msg []byte) (err error) {
	_, err = w.s.SendMessage(&sqs.SendMessageInput{
		MessageBody:            aws.String(string(msg)),
		QueueUrl:               aws.String(uri),
		MessageDeduplicationId: aws.String(uuid.New().String()),
		MessageGroupId:         aws.String(uuid.New().String()),
	})
	return
}

func (w *Queue) Receive(uri string) error {
	input := &sqs.ReceiveMessageInput{QueueUrl: &uri}

	for {
		output, err := w.s.ReceiveMessage(input)
		if err != nil {
			return err
		}

		for _, message := range output.Messages {
			log.Info().Msgf("Received message: %s", message)
			w.messages <- message
		}
	}
}

func (w *Queue) Listen(uri string, handle func([]byte) error) error {
	go func() {
		if err := w.Receive(uri); err != nil {
			log.Fatal().Err(err).Msg("error while receive message from sqs")
		}
	}()

	for message := range w.messages {
		log.Info().Msgf("Message body: %s", *message.Body)
		if err := handle([]byte(*message.Body)); err != nil {
			log.Err(err).Msg("error while handle message body")
		} else if err := w.remove(uri, message); err != nil {
			log.Err(err).Msg("error while remove message")
		}
	}

	return nil
}

func (w *Queue) remove(uri string, message *sqs.Message) error {
	_, err := w.s.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(uri),
		ReceiptHandle: message.ReceiptHandle,
	})
	return err
}
