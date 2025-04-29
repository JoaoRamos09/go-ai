package user

type User struct {
	ID       int
	Username string
	Email    string
	Password []byte
	Role     *Role   
}
