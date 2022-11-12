package code

//go:generate codegen -type=int

// api: user errors.
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User already exist.
	ErrUserAlreadyExist
)

// api: user errors.
const (
	// ErrMenuNotFound - 404: Menu not found.
	ErrMenuNotFound int = iota + 110101

	// ErrMenuAlreadyExist - 400: Menu already exist.
	ErrMenuAlreadyExist
)

const (
	ErrRoleUserNotFound int = iota + 120101
	ErrRoleUserExist
)

const (
	ErrRoleNotFound int = iota + 130101
	ErrRoleExist
)
