package entities

type ForumResponse struct {
	Id        string       `json:"id"`
	Title     string       `json:"title"`
	Caption   string       `json:"caption"`
	Media     []ForumMedia `json:"media"`
	ForumType ForumType    `json:"type"`
	User      ForumUser    `json:"user"`
}

type Forum struct {
	Uid           string `json:"uid"`
	Id            string `json:"id"`
	Title         string `json:"title"`
	Caption       string `json:"caption"`
	Fullname      string `json:"fullname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	ForumTypeId   int    `json:"forum_type_id"`
	ForumTypeName string `json:"forum_type_name"`
	UserId        string `json:"user_id"`
}

type ForumType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ForumMedia struct {
	ForumId string `json:"forum_id"`
	Path    string `json:"path"`
	Size    string `json:"size"`
}

type ForumLike struct {
	Id       string    `json:"id"`
	UserId   string    `json:"user_id"`
	Fullname string    `json:"fullname"`
	User     ForumUser `json:"user"`
}

type ForumComment struct {
	Id       string    `json:"id"`
	Comment  string    `json:"comment"`
	UserId   string    `json:"user_id"`
	Fullname string    `json:"fullname"`
	User     ForumUser `json:"user"`
}

type ForumUser struct {
	Id       string `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type ForumCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ForumStore struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Type   int    `json:"type"`
	UserId string `json:"user_id"`
}
