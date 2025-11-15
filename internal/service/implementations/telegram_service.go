package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
)

type TelegramService interface {
	SendNewTicketNotification(ticket *domain.Report) error
}

type telegramService struct {
	bot     *tgbotapi.BotAPI
	chatID  int64
	enabled bool
}

func NewTelegramService(token string, chatID int64, enabled bool) (TelegramService, error) {
	if !enabled || token == "" {
		log.Println("Telegram service disabled or token not provided")
		return &telegramService{enabled: false}, nil
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	log.Printf("Telegram bot authorized on account %s", bot.Self.UserName)

	return &telegramService{
		bot:     bot,
		chatID:  chatID,
		enabled: true,
	}, nil
}

func (s *telegramService) SendNewTicketNotification(ticket *domain.Report) error {
	if !s.enabled {
		return nil
	}

	message := s.formatTicketMessage(ticket)

	msg := tgbotapi.NewMessage(s.chatID, message)
	msg.ParseMode = "HTML"

	_, err := s.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	log.Printf("Telegram notification sent for ticket %s", ticket.ID)
	return nil
}

func (s *telegramService) formatTicketMessage(ticket *domain.Report) string {
	emoji := s.getCategoryEmoji(string(ticket.Theme))
	themeName := s.getThemeDisplayName(string(ticket.Theme))

	return fmt.Sprintf(
		`%s <b>НОВОЕ ОБРАЩЕНИЕ</b>

📋 <b>Тема:</b> %s
👤 <b>Пользователь:</b> <code>%s</code>
📅 <b>Создано:</b> %s
🆔 <b>ID:</b> <code>%s</code>

📝 <b>Текст обращения:</b>
%s

—————————————
<code>#%s</code>`,
		emoji,
		themeName,
		ticket.Contact,
		ticket.CreatedAt.Format("02.01.2006 15:04"),
		ticket.ID.String(),
		s.truncateText(ticket.Problem, 300),
		strings.ToUpper(string(ticket.Theme)),
	)
}

func (s *telegramService) getCategoryEmoji(theme string) string {
	emojis := map[string]string{
		"technical": "🔧",
		"feature":   "💡",
		"question":  "❓",
		"security":  "🛡️",
		"billing":   "💳",
		"device":    "📱",
	}
	if emoji, exists := emojis[theme]; exists {
		return emoji
	}
	return "📄"
}

func (s *telegramService) getThemeDisplayName(theme string) string {
	themes := map[string]string{
		"technical": "Технические проблемы",
		"feature":   "Предложения по улучшению",
		"question":  "Вопросы по использованию",
		"security":  "Проблемы с безопасностью",
		"billing":   "Вопросы по оплате",
		"device":    "Проблемы с устройством",
	}

	if displayName, exists := themes[theme]; exists {
		return displayName
	}
	return theme
}

func (s *telegramService) truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}
