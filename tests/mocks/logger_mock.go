package mocks

import (
	"io"
	"sync"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/logger"
)

// MockLogger - мок-реализация интерфейса Log для тестирования
type MockLogger struct {
	mu          sync.RWMutex
	calls       []Call
	level       logger.Level
	output      io.Writer
	fatalPanics bool // Определяет, должен ли Fatal вызывать panic в тестах
}

// Call представляет один вызов метода логгера
type Call struct {
	Method string
	Args   []interface{}
}

// NewMockLogger создает новый экземпляр MockLogger
func NewMockLogger() *MockLogger {
	return &MockLogger{
		calls:       make([]Call, 0),
		level:       logger.INFO,
		fatalPanics: true, // По умолчанию Fatal вызывает panic в тестах
	}
}

// Trace записывает трассировочное сообщение
func (m *MockLogger) Trace(args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, Call{
		Method: "Trace",
		Args:   args,
	})
}

// Tracef записывает форматированное трассировочное сообщение
func (m *MockLogger) Tracef(format string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	allArgs := make([]interface{}, 0, len(args)+1)
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	m.calls = append(m.calls, Call{
		Method: "Tracef",
		Args:   allArgs,
	})
}

// Debug записывает отладочное сообщение
func (m *MockLogger) Debug(args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, Call{
		Method: "Debug",
		Args:   args,
	})
}

// Debugf записывает форматированное отладочное сообщение
func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	allArgs := make([]interface{}, 0, len(args)+1)
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	m.calls = append(m.calls, Call{
		Method: "Debugf",
		Args:   allArgs,
	})
}

// Info записывает информационное сообщение
func (m *MockLogger) Info(args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, Call{
		Method: "Info",
		Args:   args,
	})
}

// Infof записывает форматированное информационное сообщение
func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	allArgs := make([]interface{}, 0, len(args)+1)
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	m.calls = append(m.calls, Call{
		Method: "Infof",
		Args:   allArgs,
	})
}

// Warn записывает предупреждение
func (m *MockLogger) Warn(args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, Call{
		Method: "Warn",
		Args:   args,
	})
}

// Warnf записывает форматированное предупреждение
func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	allArgs := make([]interface{}, 0, len(args)+1)
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	m.calls = append(m.calls, Call{
		Method: "Warnf",
		Args:   allArgs,
	})
}

// Error записывает сообщение об ошибке
func (m *MockLogger) Error(args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, Call{
		Method: "Error",
		Args:   args,
	})
}

// Errorf записывает форматированное сообщение об ошибке
func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	allArgs := make([]interface{}, 0, len(args)+1)
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	m.calls = append(m.calls, Call{
		Method: "Errorf",
		Args:   allArgs,
	})
}

// Fatal записывает критическую ошибку и вызывает panic (если включено)
func (m *MockLogger) Fatal(args ...interface{}) {
	m.mu.Lock()
	m.calls = append(m.calls, Call{
		Method: "Fatal",
		Args:   args,
	})
	m.mu.Unlock()

	if m.fatalPanics {
		panic("fatal log called")
	}
}

// Fatalf записывает форматированную критическую ошибку и вызывает panic (если включено)
func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.mu.Lock()
	allArgs := make([]interface{}, 0, len(args)+1)
	allArgs = append(allArgs, format)
	allArgs = append(allArgs, args...)
	m.calls = append(m.calls, Call{
		Method: "Fatalf",
		Args:   allArgs,
	})
	m.mu.Unlock()

	if m.fatalPanics {
		panic("fatal log called")
	}
}

// SetFatalPanic устанавливает, должен ли Fatal вызывать panic в тестах
func (m *MockLogger) SetFatalPanic(panic bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fatalPanics = panic
}

// GetCalls возвращает все записанные вызовы
func (m *MockLogger) GetCalls() []Call {
	m.mu.RLock()
	defer m.mu.RUnlock()
	calls := make([]Call, len(m.calls))
	copy(calls, m.calls)
	return calls
}

// ClearCalls очищает историю вызовов
func (m *MockLogger) ClearCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = make([]Call, 0)
}

// WasCalled проверяет, был ли вызван метод с указанным именем
func (m *MockLogger) WasCalled(method string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, call := range m.calls {
		if call.Method == method {
			return true
		}
	}
	return false
}

// GetCallsByMethod возвращает все вызовы указанного метода
func (m *MockLogger) GetCallsByMethod(method string) []Call {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []Call
	for _, call := range m.calls {
		if call.Method == method {
			result = append(result, call)
		}
	}
	return result
}

// CallCount возвращает количество вызовов указанного метода
func (m *MockLogger) CallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, call := range m.calls {
		if call.Method == method {
			count++
		}
	}
	return count
}
