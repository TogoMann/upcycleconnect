package users

func IsValidRole(role UserRole) bool {
	switch role {
	case Client, Pro, Intra, Admin:
		return true
	default:
		return false
	}
}
