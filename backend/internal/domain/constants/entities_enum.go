package constants

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type GenderPreference string

const (
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
