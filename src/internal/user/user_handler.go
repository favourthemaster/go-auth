package user

type UserHandler struct {
	// UserHandler is a struct that contains methods for handling user-related operations.

	UserService UserService
}

// NewUserHandler creates a new UserHandler with the provided UserService.
func NewUserHandler(us UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}
