package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

func NewBot(token string, jobUri, stateUri string, q *Queue, u *s3Uploader) *bot {
	return &bot{
		token:    token,
		q:        q,
		u:        u,
		stateUri: stateUri,
		jobUri:   jobUri,
	}
}

type bot struct {
	token, jobUri, stateUri string
	q                       *Queue
	u                       *s3Uploader
	tbot                    *tb.Bot
}

func (b *bot) handleState(msg []byte) error {
	s := &statePayload{}
	if err := json.Unmarshal(msg, s); err != nil {
		return err
	}

	log.Info().Msgf("handleState: %v", s)
	return nil
}

func (b *bot) Listen() {
	go b.q.Listen(b.stateUri, b.handleState)

	b.tbot.Handle(tb.OnText, b.HandleText)
	b.tbot.Handle(tb.OnVideo, b.HandleVideo)
	b.tbot.Handle(tb.OnVideoNote, b.HandleVideoNote)

	b.tbot.Start()
}

func (b *bot) Init() error {
	tbot, err := tb.NewBot(tb.Settings{
		Token:  b.token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return err
	}
	b.tbot = tbot
	return nil
}

func (b *bot) HandleVideoNote(m *tb.Message) {
	reader, err := b.tbot.GetFile(&m.VideoNote.File)
	if err != nil {
		log.Error().Err(err).Msg("failed to get file from telegram server")
		return
	}

	filePath := fmt.Sprintf("/%s%s", m.VideoNote.FileID, path.Ext(m.VideoNote.FilePath))
	b.upload(filePath, reader)
}

func (b *bot) HandleVideo(m *tb.Message) {
	reader, err := b.tbot.GetFile(&m.Video.File)
	if err != nil {
		log.Error().Err(err).Msg("failed to get file from telegram server")
		return
	}

	filePath := fmt.Sprintf("/%s%s", m.Video.FileID, path.Ext(m.Video.FilePath))
	b.upload(filePath, reader)
}

func (b *bot) upload(filePath string, reader io.Reader) {
	log.Debug().Msgf("uploaded: %s", filePath)
	if err := b.u.Upload(filePath, reader); err != nil {
		log.Error().Err(err).Msg("failed to upload video file")
		return
	}
}

func (b *bot) HandleText(m *tb.Message) {
	// log.Debug().Msgf("%s", m)
	b.tbot.Reply(m, "okay")
}
