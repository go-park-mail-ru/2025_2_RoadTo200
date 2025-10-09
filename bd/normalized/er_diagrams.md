# ER-диаграмма схемы данных

Данная диаграмма иллюстрирует структуру базы данных, включая таблицы (отношения), их атрибуты и связи между ними.

*   **Хранение сессий:** Для хранения сессий (`token`, `user_email`, `expires_at`) рекомендуется использовать `in-memory` СУБД, такую как **Redis**, для обеспечения максимальной производительности при аутентификации пользователей. Эти данные не являются частью реляционной схемы.

```mermaid
erDiagram
    user {
        UUID id PK
        TEXT email UK
        TEXT phone UK
        TEXT name
        DATE birth_date
        gender_enum gender
        TEXT bio
        DECIMAL latitude
        DECIMAL longitude
        BOOLEAN is_verified
        BOOLEAN is_premium
        TIMESTAMPTZ last_active
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    user_photo {
        UUID id PK
        UUID user_id FK
        TEXT photo_url
        INTEGER display_order
        BOOLEAN is_approved
        TIMESTAMPTZ created_at
    }

    user_preference {
        UUID user_id PK, FK
        gender_preference_enum show_gender
        INTEGER age_min
        INTEGER age_max
        INTEGER max_distance
        BOOLEAN global_search
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    swipe {
        UUID swiper_user_id PK, FK
        UUID target_user_id PK, FK
        swipe_type_enum swipe_type
        TIMESTAMPTZ created_at
    }

    match {
        UUID id PK
        UUID user1_id FK
        UUID user2_id FK
        BOOLEAN is_active
        TIMESTAMPTZ matched_at
    }

    message {
        UUID id PK
        UUID match_id FK
        UUID sender_id FK
        TEXT message_text
        BOOLEAN is_read
        BOOLEAN is_delivered
        TIMESTAMPTZ created_at
    }

    subscription {
        UUID id PK
        UUID user_id FK
        plan_type_enum plan_type
        TIMESTAMPTZ start_date
        TIMESTAMPTZ end_date
        BOOLEAN is_active
        TIMESTAMPTZ created_at
    }

    user ||--|{ user_photo : "has"
    user ||--|| user_preference : "defines"
    user }o--|| swipe : "swipes (swiper)"
    user }o--|| swipe : "is swiped on (target)"
    user }o--o{ match : "participates in"
    user ||--|{ subscription : "has"
    match ||--|{ message : "contains"
    user }o--|| message : "sends"
```