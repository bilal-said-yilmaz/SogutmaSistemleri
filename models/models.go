package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
}

type Service struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type About struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

type Contact struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	WeekdayHours  string  `json:"weekdayHours"`
	SaturdayHours string  `json:"saturdayHours"`
	SundayHours   string  `json:"sundayHours"`
}

type Hero struct {
	ID              int    `json:"id"`
	Subheading      string `json:"subheading"`
	Heading         string `json:"heading"`
	ButtonText      string `json:"buttonText"`
	BackgroundImage string `json:"backgroundImage"`
}

type Footer struct {
	ID          int               `json:"id"`
	Copyright   string            `json:"copyright"`
	SocialLinks map[string]string `json:"socialLinks"`
	Links       map[string]string `json:"links"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password,omitempty"`
	RoleID   int    `json:"role_id"`
}
