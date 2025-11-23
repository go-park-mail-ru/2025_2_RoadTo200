package errors

import "errors"

var (
	// User errors
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidEmail      = errors.New("invalid email")

	// Password errors
	ErrPasswordTooShort   = errors.New("password is too short")
	ErrPasswordsDontMatch = errors.New("passwords don't match")
	ErrInvalidPassword    = errors.New("invalid password")

	// Session errors
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")

	// Auth errors
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")

	// General errors
	ErrInternalError = errors.New("internal error")

	// Profile errors
	ErrProfileNotFound    = errors.New("profile not found")
	ErrPhotoNotFound      = errors.New("photo not found")
	ErrInvalidPhoto       = errors.New("invalid photo")
	ErrPhotoLimitExceeded = errors.New("photo limit exceeded")
	ErrInvalidProfileData = errors.New("invalid profile data")
	ErrPhotoNotOwned      = errors.New("photo does not belong to user")
	ErrInvalidPreferences = errors.New("invalid preferences data")

	// Swipe errors
	ErrInvalidSwipeAction = errors.New("invalid swipe action")
	ErrCannotSwipeSelf    = errors.New("cannot swipe yourself")

	// Match errors
	ErrMatchNotFound       = errors.New("match not found")
	ErrMatchNotActive      = errors.New("match is not active")
	ErrNotMatchParticipant = errors.New("user is not a participant of this match")
	ErrMatchNotFound = errors.New("match not found")

	// Strike errors
	ErrStrikeNotFound       = errors.New("strike not found")
	ErrStrikeAlreadyDeleted = errors.New("strike already deleted")
	ErrDuplicateStrike      = errors.New("duplicate active strike")
	ErrSelfStrikeNotAllowed = errors.New("self strike not allowed")
	ErrInvalidReporterID    = errors.New("invalid reporter ID")
	ErrInvalidTargetUserID  = errors.New("invalid target user ID")
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrDateRangeTooLarge    = errors.New("date range too large")
)
