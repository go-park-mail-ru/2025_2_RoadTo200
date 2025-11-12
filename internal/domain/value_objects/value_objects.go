package domain

type SwipeDirection string

const (
	SwipeLeft  SwipeDirection = "left"
	SwipeRight SwipeDirection = "right"
)

type Swipe struct {
	UserID    string
	TargetID  string
	Direction SwipeDirection
}
