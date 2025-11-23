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
	PlanTypePremium  PlanType = "premium"
	PlanTypeGold     PlanType = "gold"
	PlanTypePlatinum PlanType = "platinum"
)

const (
	MaxBioLength     = 500
	MaxNameLength    = 50
	MaxPhotosPerUser = 9
	MaxPhotoSize     = 10 << 20 // 10MB
	AllowedMimeTypes = "image/jpeg,image/png,image/webp"

	MinAge      = 18
	MaxAge      = 99
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
