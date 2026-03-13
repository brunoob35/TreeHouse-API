package authentication

// Permission represents a permission flag used by the application.
// IMPORTANT:
// Each permission value must match exactly the ID stored in the database table
// "permissoes". Since the system uses bit flags, all permission IDs must be
// powers of two.
type Permission uint64

const (
	PermGestao Permission = 1 << iota
	PermProfessor
	PermGestaoMaster
	PermFinanceiro
	PermFinanceiroView
)

// BuildPermissionMask receives a list of permission IDs and aggregates them
// into a single numeric bitmask using the bitwise OR operator.
//
// Example:
//
//	[1, 4] => 1 | 4 = 5
//
// This function assumes the permission IDs coming from the database are already
// valid bit flags.
func BuildPermissionMask(permissionIDs []uint64) uint64 {
	var mask uint64

	for _, permissionID := range permissionIDs {
		mask |= permissionID
	}

	return mask
}

// HasPermission checks whether the user's permission mask contains
// the required permission flag.
func HasPermission(userPermissions uint64, required Permission) bool {
	return userPermissions&uint64(required) != 0
}

// HasAnyPermission checks whether the user's permission mask contains
// at least one of the required permission flags.
func HasAnyPermission(userPermissions uint64, required ...Permission) bool {
	for _, permission := range required {
		if HasPermission(userPermissions, permission) {
			return true
		}
	}

	return false
}

// HasAllPermissions checks whether the user's permission mask contains
// all required permission flags.
func HasAllPermissions(userPermissions uint64, required ...Permission) bool {
	for _, permission := range required {
		if !HasPermission(userPermissions, permission) {
			return false
		}
	}

	return true
}
