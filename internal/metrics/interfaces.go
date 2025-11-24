package metrics

// Metrics интерфейс определяет все методы для сбора метрик приложения
// Соответствует принципам чистой архитектуры - абстракция без зависимостей от конкретной реализации
type Metrics interface {
	HttpMetrics
	GrpcMetrics
	BusinessMetrics
}

// HTTP метрики
type HttpMetrics interface {
	// IncHTTPRequest увеличивает счетчик HTTP запросов
	// method - HTTP метод (GET, POST, PUT, DELETE и т.д.)
	// path - путь запроса (нормализованный, например /api/v1/users/:id)
	// statusCode - HTTP статус код ответа
	IncHTTPRequest(method, path string, statusCode int)

	// ObserveHTTPDuration измеряет продолжительность HTTP запроса
	// method - HTTP метод
	// path - путь запроса
	// statusCode - HTTP статус код
	// duration - продолжительность в секундах
	ObserveHTTPDuration(method, path string, statusCode int, duration float64)

	// Метрики кэша

	// IncCacheHit увеличивает счетчик попаданий в кэш
	// service - имя сервиса
	// cacheType - тип кэша (redis, memory, database и т.д.)
	//IncCacheHit(cacheType string)

	// IncCacheMiss увеличивает счетчик промахов кэша
	// service - имя сервиса
	// cacheType - тип кэша
	//IncCacheMiss(cacheType string)
}

// gRPC метрики
type GrpcMetrics interface {
	// IncGRPCRequest увеличивает счетчик gRPC запросов
	// method - полное имя gRPC метода (package.Service/Method)
	// statusCode - код статуса gRPC (приведенный к int)
	IncGRPCRequest(method string, statusCode int)

	// ObserveGRPCDuration измеряет продолжительность gRPC запроса
	// method - полное имя gRPC метода
	// statusCode - код статуса gRPC
	// duration - продолжительность в секундах
	ObserveGRPCDuration(method string, statusCode int, duration float64)

	// Метрики ошибок

	// IncError увеличивает счетчик ошибок
	// service - имя сервиса, в котором произошла ошибка
	// operation - операция, при выполнении которой произошла ошибка
	// errorType - тип ошибки (validation, database, network, auth, business и т.д.)
	IncError(operation, errorType string)
}

// Бизнес метрики
type BusinessMetrics interface {

	// SetActiveSessions устанавливает количество активных сессий
	// count - количество активных сессий пользователей
	SetActiveSessions(count int)

	// SetDatabasePoolMetrics устанавливает метрики пула соединений БД
	// service - имя сервиса
	// poolType - тип пула (open, idle, active, waiting)
	// count - количество соединений
	SetDatabasePoolMetrics(service string, poolType string, count int)

	// SetErrorMetricsустанавливает метрики пула соединений БД
	// service - имя сервиса
	// operation - тип пула (open, idle, active, waiting)
	// errorType - количество соединений
	SetErrorMetrics(model, method, errorType string)
}

// Дополнительные бизнес-метрики которые могут понадобиться
// SetBusinessMetric устанавливает произвольную бизнес-метрику
type NamedMetrics interface {
	// name - название метрики
	// value - значение
	// labels - дополнительные метки
	SetBusinessMetric(name string, value float64, labels ...string)

	// IncBusinessCounter увеличивает произвольный бизнес-счетчик
	// name - название счетчика
	// labels - дополнительные метки
	IncBusinessCounter(name string, labels ...string)

	// ObserveBusinessDuration измеряет продолжительность бизнес-операции
	// name - название операции
	// duration - продолжительность в секундах
	// labels - дополнительные метки
	ObserveBusinessDuration(name string, duration float64, labels ...string)
}

// MetricsCollector интерфейс для сборщиков метрик, которые работают в фоне
type MetricsCollector interface {
	// CollectHTTPServerMetrics собирает метрики HTTP сервера
	CollectHTTPServerMetrics()

	// CollectGRPCMetrics собирает метрики gRPC сервисов
	CollectGRPCMetrics()

	// CollectSystemMetrics собирает системные метрики (CPU, память, диск)
	CollectSystemMetrics()

	// CollectBusinessMetrics собирает бизнес-метрики
	CollectBusinessMetrics()

	// Start запускает все сборщики метрик
	Start()

	// Stop останавливает все сборщики метрик
	Stop()
}

// MetricsFactory интерфейс для фабрики создания метрик
// Позволяет создавать метрики с определенными labels по умолчанию
type MetricsFactory interface {
	// WithService создает экземпляр метрик с привязкой к конкретному сервису
	WithService(serviceName string) Metrics

	// WithLabels создает экземпляр метрик с дополнительными labels
	WithLabels(labels ...string) Metrics
}

// ErrorClassifier интерфейс для классификации ошибок
type ErrorClassifier interface {
	// ClassifyError определяет тип ошибки для метрик
	ClassifyError(err error) string
}

// Constants для типов ошибок
const (
	ErrorTypeValidation = "validation"
	ErrorTypeDatabase   = "database"
	ErrorTypeNetwork    = "network"
	ErrorTypeAuth       = "auth"
	ErrorTypeBusiness   = "business"
	ErrorTypeExternal   = "external"
	ErrorTypeUnknown    = "unknown"
)

// Constants для типов кэша
const (
	CacheTypeRedis    = "redis"
	CacheTypeMemory   = "memory"
	CacheTypeDatabase = "database"
)

// Constants для типов пула БД
const (
	PoolTypeOpen    = "open"
	PoolTypeIdle    = "idle"
	PoolTypeActive  = "active"
	PoolTypeWaiting = "waiting"
)
