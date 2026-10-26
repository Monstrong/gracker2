package domain

type Role int8

const (
	RoleUnspecified Role = iota
	BannedRole
	UserRole
	ModeratorRole
	AdminRole
	OwnerRole // admin who can ban admins

	BillingRole
)

func (r Role) Weight() int {
    switch r {
	case BannedRole:    return 0 
    case UserRole:      return 1
    case ModeratorRole: return 2
    case AdminRole:     return 3
    case OwnerRole:     return 4
    
	case BillingRole:   return 100
    default:            return -1 // unspecified / неизвестное
    }
}

type PremiumType int8

const (
	PremiumUnspecified PremiumType = iota
	Free
	Premium
)
