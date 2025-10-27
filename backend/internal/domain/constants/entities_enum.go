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
