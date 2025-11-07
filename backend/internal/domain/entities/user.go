package domain

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/domain/constants"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	Email     string           `json:"email" db:"email"`
	Phone     *string          `json:"phone,omitempty" db:"phone"`
	Name      string           `json:"name" db:"name"`
	Password  string           `json:"-" db:"password"`
	BirthDate time.Time        `json:"birth_date" db:"birth_date"`
	Gender    constants.Gender `json:"gender" db:"gender"`
	Bio       *string          `json:"bio,omitempty" db:"bio"`
	// Новые поля интересов
	Workout    bool `db:"workout" json:"workout"`
	Fun        bool `db:"fun" json:"fun"`
	Party      bool `db:"party" json:"party"`
	Chill      bool `db:"chill" json:"chill"`
	Love       bool `db:"love" json:"love"`
	Relax      bool `db:"relax" json:"relax"`
	Yoga       bool `db:"yoga" json:"yoga"`
	Friendship bool `db:"friendship" json:"friendship"`
	Culture    bool `db:"culture" json:"culture"`
	Cinema     bool `db:"cinema" json:"cinema"`
	// Остальные поля
	Latitude   *float64  `json:"latitude,omitempty" db:"latitude"`
	Longitude  *float64  `json:"longitude,omitempty" db:"longitude"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	LastActive time.Time `json:"last_active" db:"last_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(email, passwordHash string) *User {
	return &User{
		Email:    email,
		Password: passwordHash,
	}
}
