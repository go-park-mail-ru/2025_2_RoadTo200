package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/entities"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	bot     *tgbotapi.BotAPI
	logger  logger.Log
	chatID  int64
	enabled bool
}

func NewTelegramService(token string, l logger.Log, chat string, enabled bool) (*TelegramService, error) {
	token = os.Getenv(token)
	chatID, err := strconv.ParseInt(os.Getenv(chat), 10, 64)
	if err != nil {
		return nil, err
	}
	if !enabled || token == "" {
		l.Error("Telegram service disabled or token not provided")
		return &TelegramService{enabled: false}, nil
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	l.Debugf("Telegram bot authorized on account %s", bot.Self.UserName)

	return &TelegramService{
		bot:     bot,
		logger:  l,
		chatID:  chatID,
		enabled: true,
	}, nil
}

func (s *TelegramService) SendNewTicketNotification(ticket *domain.Report) (int, error) {
	if !s.enabled {
		return 0, nil
	}

	message := s.formatTicketMessage(ticket)
	keyboard := s.createTicketKeyboard(ticket.ID.String())

	msg := tgbotapi.NewMessage(s.chatID, message)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = keyboard

	sentMsg, err := s.bot.Send(msg)
	if err != nil {
		return 0, fmt.Errorf("failed to send telegram message: %w", err)
	}

	s.logger.Debugf("Telegram notification sent for ticket %s, message ID: %d", ticket.ID, sentMsg.MessageID)
	return sentMsg.MessageID, nil
}

func (s *TelegramService) UpdateTicketMessage(chatID int64, messageID int, ticket *domain.Report) error {
	if !s.enabled {
		return nil
	}

	s.logger.Debugf("Updating ticket Message %v", ticket)
	message := s.formatTicketMessage(ticket)
	keyboard := s.createTicketKeyboard(ticket.ID.String())

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, message)
	editMsg.ParseMode = "HTML"
	editMsg.ReplyMarkup = &keyboard

	_, err := s.bot.Send(editMsg)
	if err != nil {
		return fmt.Errorf("failed to update telegram message: %w", err)
	}

	s.logger.Debugf("Telegram message updated for ticket %s", ticket.ID)
	return nil
}

func (s *TelegramService) ProcessCallback(callbackData string) (ticketID string, action string, err error) {
	// Формат callback_data: "ticket_{action}_{ticketID}"
	parts := strings.Split(callbackData, "_")
	if len(parts) != 3 || parts[0] != "ticket" {
		return "", "", fmt.Errorf("invalid callback data format")
	}

	action = parts[1]
	ticketID = parts[2]

	s.logger.Debugf("Processed telegram callback: ticket=%s, action=%s", ticketID, action)
	return ticketID, action, nil
}

func (s *TelegramService) GetBot() *tgbotapi.BotAPI {
	return s.bot
}

func (s *TelegramService) createTicketKeyboard(ticketID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ В работу", fmt.Sprintf("ticket_work_%s", ticketID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Закрыть", fmt.Sprintf("ticket_close_%s", ticketID)),
		),
	)
}

func (s *TelegramService) formatTicketMessage(ticket *domain.Report) string {
	s.logger.Debugf("formatTicketMessage: ticket ID: %d", ticket.ID)
	emoji := s.getCategoryEmoji(string(ticket.Theme))
	themeName := s.getThemeDisplayName(string(ticket.Theme))
	statusEmoji := s.getStatusEmoji(string(ticket.Status))

	s.logger.Debugf("emoji %s, %s %s", emoji, statusEmoji, themeName)

	return fmt.Sprintf(
		`%s <b>ОБРАЩЕНИЕ</b> %s

📋 <b>Тема:</b> %s
👤 <b>Пользователь:</b> <code>%s</code>
📅 <b>Создано:</b> %s
🆔 <b>ID:</b> <code>%s</code>
📊 <b>Статус:</b> %s

📝 <b>Текст обращения:</b>
%s

—————————————
<code>#%s</code>`,
		emoji,
		statusEmoji,
		themeName,
		ticket.Contact,
		ticket.CreatedAt.Format("02.01.2006 15:04"),
		ticket.ID.String(),
		s.getStatusDisplayName(string(ticket.Status)),
		s.truncateText(ticket.Problem, 300),
		strings.ToUpper(string(ticket.Theme)),
	)
}

func (s *TelegramService) getStatusEmoji(status string) string {
	emojis := map[string]string{
		"open":   "🆕",
		"work":   "🔄",
		"closed": "✅",
	}
	if emoji, exists := emojis[status]; exists {
		return emoji
	}
	return "📄"
}

func (s *TelegramService) getStatusDisplayName(status string) string {
	statuses := map[string]string{
		"open":   "Новое",
		"work":   "В работе",
		"closed": "Закрыто",
	}
	if displayName, exists := statuses[status]; exists {
		return displayName
	}
	return status
}

func (s *TelegramService) getCategoryEmoji(theme string) string {
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

func (s *TelegramService) getThemeDisplayName(theme string) string {
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

func (s *TelegramService) truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}
