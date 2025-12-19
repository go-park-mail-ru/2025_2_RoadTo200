package constants

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type GenderPreference string

const (
	GenderPrefBoth   GenderPreference = "both"
	GenderPrefMale   GenderPreference = "male"
	GenderPrefFemale GenderPreference = "female"
)

type SwipeType string

const (
	SwipeTypeLike      SwipeType = "like"
	SwipeTypeDislike   SwipeType = "dislike"
	SwipeTypeSuperLike SwipeType = "super_like"
)

type PlanType string

const (
	PlanTypeWeek    PlanType = "week"
	PlanTypeMonth   PlanType = "month"
	PlanTypeQuarter PlanType = "quarter"
)

const (
	MaxBioLength     = 500
	MaxNameLength    = 50
	MaxPhotosPerUser = 9
	MaxPhotoSize     = 10 << 20 // 10MB
	AllowedMimeTypes = "image/jpeg,image/png,image/webp"

	MinAge      = 18
	MaxAge      = 65
	MinDistance = 1
	MaxDistance = 500
)

type InterestType string

const (
	InterestTypeWorkout    InterestType = "workout"
	InterestTypeFun        InterestType = "fun"
	InterestTypeParty      InterestType = "party"
	InterestTypeChill      InterestType = "chill"
	InterestTypeLove       InterestType = "love"
	InterestTypeRelax      InterestType = "relax"
	InterestTypeYoga       InterestType = "yoga"
	InterestTypeFriendship InterestType = "friendship"
	InterestTypeCulture    InterestType = "culture"
	InterestTypeCinema     InterestType = "cinema"
)

var MaxStrikeReasonLength = 250

type StrikeType string

const (
	StrikeTypeSpam          StrikeType = "spam"
	StrikeTypeFakeProfile   StrikeType = "fake_profile"
	StrikeTypeOffensive     StrikeType = "offensive_content"
	StrikeTypeHarassment    StrikeType = "harassment"
	StrikeTypeInappropriate StrikeType = "inappropriate_content"
	StrikeTypeUnderage      StrikeType = "underage"
	StrikeTypeCopyright     StrikeType = "copyright_violation"
	StrikeTypeOther         StrikeType = "other"
)

type StrikeStatus string

const (
	StrikeStatusPending  StrikeStatus = "pending"
	StrikeStatusApproved StrikeStatus = "approved"
	StrikeStatusRejected StrikeStatus = "rejected"
	StrikeStatusResolved StrikeStatus = "resolved"
)

type NotificationType string

const (
	NotificationTypeMatch     NotificationType = "match"
	NotificationTypeSuperLike NotificationType = "super_like"
	NotificationTypeLike      NotificationType = "like"
	NotificationTypeMessage   NotificationType = "message"
)

// ReportTheme представляет тему обращения в поддержку
type ReportTheme string

const (
	ReportThemeTechnical ReportTheme = "technical"
	ReportThemeFeature   ReportTheme = "feature"
	ReportThemeQuestion  ReportTheme = "question"
	ReportThemeSecurity  ReportTheme = "security"
	ReportThemeBilling   ReportTheme = "billing"
	ReportThemeDevice    ReportTheme = "device"
)

// ReportStatus представляет статус обращения
type ReportStatus string

const (
	ReportStatusOpen   ReportStatus = "open"
	ReportStatusWork   ReportStatus = "work"
	ReportStatusClosed ReportStatus = "close"
)

// CategoryMapping маппинг человекочитаемых названий в технические
var CategoryMapping = map[string]ReportTheme{
	"Технические проблемы":     ReportThemeTechnical,
	"Предложения по улучшению": ReportThemeFeature,
	"Вопросы по использованию": ReportThemeQuestion,
	"Проблемы с безопасностью": ReportThemeSecurity,
	"Вопросы по оплате":        ReportThemeBilling,
	"Проблемы с устройством":   ReportThemeDevice,
}

// StatusDisplayNames маппинг технических статусов в человекочитаемые
var StatusDisplayNames = map[ReportStatus]string{
	ReportStatusOpen:   "Открыто",
	ReportStatusWork:   "В работе",
	ReportStatusClosed: "Закрыто",
}

// IsValidHumanCategory проверяет валидность человекочитаемой категории
func IsValidHumanCategory(category string) bool {
	_, exists := CategoryMapping[category]
	return exists
}

// HumanCategoryToTheme конвертирует человекочитаемую категорию в тему
func HumanCategoryToTheme(category string) (ReportTheme, error) {
	theme := CategoryMapping[category]
	return theme, nil
}

// ThemeToHumanCategory конвертирует тему в человекочитаемую категорию
func ThemeToHumanCategory(theme ReportTheme) string {
	// Обратный маппинг
	for human, technical := range CategoryMapping {
		if technical == theme {
			return human
		}
	}
	return "Неизвестная категория"
}

// DefaultReportStatus возвращает статус по умолчанию
func DefaultReportStatus() ReportStatus {
	return ReportStatusOpen
}
