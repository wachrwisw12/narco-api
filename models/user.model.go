package models

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"` // ❌ ไม่ส่งกลับ
	Role     string `json:"role" db:"role"`
	Fullname string `json:"fullname,omitempty" db:"fullname"`
	RoleId   int    `json:"role_id,omitempty" db:"role_id"`
}
