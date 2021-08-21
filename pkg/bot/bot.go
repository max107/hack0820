package bot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

func NewBot(token string, jobUri, stateUri string, q *queue, u *uploader) *bot {
	return &bot{
		token:    token,
		q:        q,
		u:        u,
		stateUri: stateUri,
		jobUri:   jobUri,
	}
}

type bot struct {
	stateUri string
	jobUri   string
	token    string
	q        *queue
	u        *uploader
	tbot     *tb.Bot
}

func (b *bot) handleState(msg []byte) error {
	s := &msgPayload{}
	if err := json.Unmarshal(msg, s); err != nil {
		return err
	}

	log.Info().Msgf("handleState: %v", s)
	replyRecipient := &tb.Message{
		ID:   s.MessageID,
		Chat: &tb.Chat{ID: s.ChatID},
	}

	b.tbot.Reply(replyRecipient, fmt.Sprintf("Видео готово: %s", s.VideoURL))
	return nil
}

func (b *bot) Listen() {
	go b.q.Listen(b.stateUri, b.handleState)

	b.tbot.Handle(tb.OnVideo, func(m *tb.Message) {
		b.handleUpload(m, &m.Video.File)
	})
	b.tbot.Handle(tb.OnVideoNote, func(m *tb.Message) {
		b.handleUpload(m, &m.VideoNote.File)
	})

	b.tbot.Handle("/start", func(m *tb.Message) {
		b.tbot.Reply(m, "Привет друг! Я принимаю только видео файлы и видео сообщения, все остальное будет проигнорировано.")
	})

	b.tbot.Start()
}

func (b *bot) Init() error {
	tbot, err := tb.NewBot(tb.Settings{
		Token:  b.token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize telegram bot")
		return err
	}
	b.tbot = tbot
	return nil
}

func (b *bot) handleUpload(m *tb.Message, f *tb.File) {
	filePath := fmt.Sprintf("/%s", f.UniqueID)
	reader, err := b.tbot.GetFile(f)
	if err != nil {
		log.Error().Err(err).Msg("failed to get file from telegram server")
		return
	}

	log.Debug().Msgf("uploaded: %s", filePath)
	if err := b.u.Upload(filePath, reader); err != nil {
		log.Error().Err(err).Msg("failed to upload video file")
		return
	}

	if _, err := b.tbot.Reply(m, "Видео получено, начата обработка"); err != nil {
		log.Error().Err(err).Msg("Не удалось отправить уведомление")
	}

	payload := &msgPayload{
		ChatID:    m.Chat.ID,
		MessageID: m.ID,
		VideoURL:  filePath,
	}
	if err := b.q.Send(b.jobUri, payload); err != nil {
		log.Error().Err(err).Msg("failed to upload video file")
		return
	}
}
