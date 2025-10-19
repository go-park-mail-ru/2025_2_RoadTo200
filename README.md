# 2025_2_RoadTo200
Репозиторий команды RoadTo200. Проект: Тиндер

## Состав
* Егоров Дмитрий
* Гилязетдинов Кирилл
* Матвеев Илья
* Микулин Михаил

## Менторы
* Fronted: Нигматуллин Алик
* Backend: Кузьмин Ярослав
* UX: Ченцова Дарья
* СУБД: Конопкин Евгений 

## API 
* POST /api/register  – регистрация
* POST /api/login     – вход
* GET  /api/session   – проверка авторизации
* POST /api/logout    – выход


## Структура
#### cmd/server/ - точка входа приложения
    * main.go - инициализация и запуск

#### internal/domain/ - доменный слой
    * entities/ - основные сущности (User, Session, Profile, etc.)
    * value_objects/ - объекты-значения (Email, Password, etc.)
    * errors/ - кастомные ошибки домена
    * constants/ - константы

#### internal/service/ - бизнес-логика
    * interfaces/ - интерфейсы сервисов
    * implementations/ - реализации сервисов

#### internal/handler/ - обработчики HTTP
    * http/ - HTTP хендлеры
    * middleware/ - CORS, аутентификация, логирование

#### internal/repository/ - работа с данными
    * interfaces/ - интерфейсы репозиториев
    * implementations/ - in-memory, postgres реализации

#### internal/config/ - конфигурация приложения
#### internal/logger/ - логгер

#### pkg/ - переиспользуемые пакеты
    * utils/ - утилиты (JSON, валидация)
    * security/ - хэширование, JWT
    * httpserver/ - настройка HTTP сервера

#### api/docs/swagger/ - Swagger документация

#### tests/ - тесты
    * unit/ - юнит-тесты
    * integration/ - интеграционные тесты
    * mocks/ - моки для Gomock

### Остальные папки:

#### migrations/ - миграции БД
#### scripts/ - вспомогательные скрипты
#### deployments/ - docker-compose, k8s манифесты