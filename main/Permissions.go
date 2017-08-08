package main

// Checks whether you are allowed to remove messages or have a higher level permission.
func isAllowedToPrune(permissionLevel int) bool {
	return permissionLevel&0x00002000 == 0x00002000 || isAdmin(permissionLevel)
}

// Checks whether an user is admin.
// Sadly also has to check for manage roles role, since the owner does not have the admin permission.
func isAdmin(permissionLevel int) bool {
	return permissionLevel&0x08 == 0x08 || permissionLevel&0x10000000 == 0x10000000
}
