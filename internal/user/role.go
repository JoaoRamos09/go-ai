package user

type Role int

const (
	RoleUser Role = iota + 1
	RoleAdmin
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "user"
	}
	return ""
}

func (r Role) Value() int {
	return int(r)
}




