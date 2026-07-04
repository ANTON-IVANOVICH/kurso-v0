package domain

// AdminRole is a platform administrator's permission level.
type AdminRole string

const (
	AdminSuperadmin AdminRole = "superadmin"
	AdminModerator  AdminRole = "moderator"
)

// Admin is a platform administrator (admin.kurso.io).
type Admin struct {
	ID           string
	Email        string
	PasswordHash string
	Role         AdminRole
	TOTPSecret   *string
	TOTPEnabled  bool
	Status       string // active | disabled
}
