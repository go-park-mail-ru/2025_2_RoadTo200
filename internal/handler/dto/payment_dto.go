package dto

//go:generate easyjson -all -no_std_marshalers payment_dto.go

// CreatePaymentRequest запрос на создание платежа
// @Description Запрос для создания ссылки на оплату премиум-подписки
type CreatePaymentRequest struct {
	Amount float64 `json:"amount" binding:"required" example:"1000"` // Сумма платежа
	PlanID string  `json:"planId" binding:"required" example:"week"` // ID плана: week, month, quarter
}

// CreatePaymentResponse ответ со ссылкой на оплату
type CreatePaymentResponse struct {
	PaymentURL string `json:"payment_url" example:"https://yoomoney.ru/quickpay/confirm?receiver=..."` // Ссылка для редиректа на оплату
}
