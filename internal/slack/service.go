package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ashos/internal/core/event"
	"ashos/internal/storage"
)

const (
	ConfigCollection = "config"
	SlackTokenKey    = "slack_token"
	SlackChannelKey  = "slack_channel"
)

type Config struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

type service struct {
	store storage.Store
	bus   *event.EventBus
}

type Service interface {
	SendMessage(ctx context.Context, text string) error
	GetRecentMessages(ctx context.Context, count int) ([]SlackMessage, error)
	Connect(ctx context.Context, token, channel string) error
	GetConfig(ctx context.Context) (Config, error)
}

type SlackMessage struct {
	Text      string    `json:"text"`
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}

func NewService(store storage.Store, bus *event.EventBus) Service {
	s := &service{
		store: store,
		bus:   bus,
	}

	// Auto-post events to Slack if configured
	bus.Subscribe(event.TaskCompleted{}, func(e any) {
		ev := e.(event.TaskCompleted)
		s.SendMessage(context.Background(), fmt.Sprintf("✅ *Task Completed:* %s", ev.Title))
	})

	bus.Subscribe(event.FocusEnded{}, func(e any) {
		ev := e.(event.FocusEnded)
		s.SendMessage(context.Background(), fmt.Sprintf("🧘 *Focus Session Ended:* %v tracked. Session: %s", ev.Duration.Round(time.Second), ev.Summary))
	})

	return s
}

func (s *service) GetConfig(ctx context.Context) (Config, error) {
	var cfg Config
	err := s.store.Get(ctx, ConfigCollection, "slack_config", &cfg)
	return cfg, err
}

func (s *service) Connect(ctx context.Context, token, channel string) error {
	cfg := Config{Token: token, Channel: channel}
	return s.store.Save(ctx, ConfigCollection, "slack_config", cfg)
}

func (s *service) SendMessage(ctx context.Context, text string) error {
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("slack not configured: use 'ash slack connect'")
	}

	payload := map[string]string{
		"channel": cfg.Channel,
		"text":    text,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API error: status %d", resp.StatusCode)
	}

	var slackResp struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&slackResp); err != nil {
		return err
	}

	if !slackResp.OK {
		return fmt.Errorf("slack API error: %s", slackResp.Error)
	}

	return nil
}

func (s *service) GetRecentMessages(ctx context.Context, count int) ([]SlackMessage, error) {
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("slack not configured")
	}

	url := fmt.Sprintf("https://slack.com/api/conversations.history?channel=%s&limit=%d", cfg.Channel, count)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var slackResp struct {
		OK       bool   `json:"ok"`
		Error    string `json:"error"`
		Messages []struct {
			Text string `json:"text"`
			User string `json:"user"`
			Ts   string `json:"ts"`
		} `json:"messages"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&slackResp); err != nil {
		return nil, err
	}

	if !slackResp.OK {
		return nil, fmt.Errorf("slack API error: %s", slackResp.Error)
	}

	var msgs []SlackMessage
	for _, m := range slackResp.Messages {
		msgs = append(msgs, SlackMessage{
			Text: m.Text,
			User: m.User,
		})
	}

	return msgs, nil
}
