package models

type User struct {
	UID       string   `json:"uid"`
	Email     string   `json:"email"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Age       int      `json:"age"`
	Gender    string   `json:"gender"`
	CreatedAt int64    `json:"created_at"`
	HighScore float64  `json:"high_score"`
}
