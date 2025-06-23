package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ravikantchauhan246/ospy/internal/storage"
)

// TelegramNotifier handles Telegram notifications
type TelegramNotifier struct {
	botToken string
	chatID   string
	client   *http.Client
	enabled  bool
}

// TelegramMessage represents a Telegram message
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier(botToken, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
		client:   &http.Client{Timeout: 10 * time.Second},
		enabled:  botToken != "" && chatID != "",
	}
}

// IsEnabled returns whether Telegram notifications are enabled
func (t *TelegramNotifier) IsEnabled() bool {
	return t.enabled
}

// SendDownAlert sends an alert when a website goes down
func (t *TelegramNotifier) SendDownAlert(websiteName, url, message string) error {
	if !t.enabled {
		return nil
	}

	text := fmt.Sprintf(`🚨 *Website Down Alert*

*Website:* %s
*URL:* %s
*Status:* DOWN
*Message:* %s
*Time:* %s`, 
		escapeMarkdown(websiteName), 
		escapeMarkdown(url), 
		escapeMarkdown(message), 
		time.Now().Format("2006-01-02 15:04:05"))

	return t.sendMessage(text)
}

// SendUpAlert sends an alert when a website comes back up
func (t *TelegramNotifier) SendUpAlert(websiteName, url string, downtime time.Duration) error {
	if !t.enabled {
		return nil
	}

	text := fmt.Sprintf(`✅ *Website Restored*

*Website:* %s
*URL:* %s
*Status:* UP
*Downtime:* %v
*Time:* %s`, 
		escapeMarkdown(websiteName), 
		escapeMarkdown(url), 
		downtime, 
		time.Now().Format("2006-01-02 15:04:05"))

	return t.sendMessage(text)
}

// SendSummaryReport sends a periodic summary report
func (t *TelegramNotifier) SendSummaryReport(stats []storage.WebsiteStats) error {
	if !t.enabled {
		return nil
	}

	var text strings.Builder
	text.WriteString("📊 *Weekly Summary Report*\n\n")
	
	for _, stat := range stats {
		status := "🟢"
		if stat.LastStatus == "DOWN" {
			status = "🔴"
		}
		
		text.WriteString(fmt.Sprintf("%s *%s*\n", status, escapeMarkdown(stat.WebsiteName)))
		text.WriteString(fmt.Sprintf("   Uptime: %.2f%%\n", stat.UptimePercent))
		text.WriteString(fmt.Sprintf("   Avg Response: %.0fms\n", stat.AvgResponseTime))
		text.WriteString(fmt.Sprintf("   Total Checks: %d\n\n", stat.TotalChecks))
	}
	
	text.WriteString(fmt.Sprintf("*Report time:* %s", time.Now().Format("2006-01-02 15:04:05")))

	return t.sendMessage(text.String())
}

// sendMessage sends a message to Telegram
func (t *TelegramNotifier) sendMessage(text string) error {
	message := TelegramMessage{
		ChatID:    t.chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)
	
	resp, err := t.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}

// escapeMarkdown escapes special characters for Telegram markdown
func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}