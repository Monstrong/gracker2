package domain

type Role int8

const (
	RoleUnspecified Role = iota
	UserRole
	BannedRole
	ModeratorRole
	AdminRole
)

type PremiumType int8

const (
	PremiumUnspecified PremiumType = iota
	Free
	Premium
)
