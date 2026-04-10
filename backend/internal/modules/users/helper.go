package users

func IsValidRole(role UserRole) bool {
	switch role {
	case Client, Pro, Interne, Admin:
		return true
	default:
		return false
	}
}
