package entities

import "time"

type ForumResponse struct {
	Id           string         `json:"id"`
	Title        string         `json:"title"`
	Caption      string         `json:"caption"`
	Media        []ForumMedia   `json:"medias"`
	Comment      []ForumComment `json:"comments"`
	Like         []ForumLike    `json:"likes"`
	IsLiked      bool           `json:"is_liked"`
	CommentCount int            `json:"comment_count"`
	LikeCount    int            `json:"like_count"`
	ForumType    ForumType      `json:"type"`
	User         ForumUser      `json:"user"`
	CreatedAt    time.Time      `json:"created_at"`
}

type Forum struct {
	Uid           string    `json:"uid"`
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Caption       string    `json:"caption"`
	Avatar        string    `json:"avatar"`
	Fullname      string    `json:"fullname"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	ForumTypeId   int       `json:"forum_type_id"`
	ForumTypeName string    `json:"forum_type_name"`
	UserId        string    `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type ForumType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ForumMedia struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

type ForumLikeQuery struct {
	Id       string `json:"id"`
	UserId   string `json:"user_id"`
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname"`
}

type CheckIsLike struct {
	IsExist bool `json:"is_exist"`
}

type ForumLike struct {
	Id   string        `json:"id"`
	User ForumLikeUser `json:"user"`
}

type ForumCommentQuery struct {
	Id        string    `json:"id"`
	Comment   string    `json:"comment"`
	Avatar    string    `json:"avatar"`
	UserId    string    `json:"user_id"`
	Fullname  string    `json:"fullname"`
	CreatedAt time.Time `json:"created_at"`
}

type ForumComment struct {
	Id         string              `json:"id"`
	Comment    string              `json:"comment"`
	User       ForumCommentUser    `json:"user"`
	IsLiked    bool                `json:"is_liked"`
	Reply      []ForumCommentReply `json:"replies"`
	ReplyCount int                 `json:"reply_count"`
	CreatedAt  time.Time           `json:"created_At"`
}

type ForumCommentUser struct {
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname"`
}

type ForumCommentReplyQuery struct {
	Id        string    `json:"id"`
	Reply     string    `json:"reply"`
	Avatar    string    `json:"avatar"`
	UserId    string    `json:"user_id"`
	Fullname  string    `json:"fullname"`
	CreatedAt time.Time `json:"created_at"`
}

type ForumCommentReply struct {
	Id        string                `json:"id"`
	Reply     string                `json:"reply"`
	User      ForumCommentReplyUser `json:"user"`
	IsLiked   bool                  `json:"is_liked"`
	CreatedAt time.Time             `json:"created_at"`
}

type ForumCommentReplyUser struct {
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname"`
}

type ForumLikeUser struct {
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname"`
}

type ForumUser struct {
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
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

type ForumStoreLike struct {
	Id      string `json:"id"`
	ForumId string `json:"forum_id"`
	UserId  string `json:"user_id"`
}

type CommentStore struct {
	Id      string `json:"id"`
	ForumId string `json:"forum_id"`
	Comment string `json:"comment"`
	UserId  string `json:"user_id"`
}

type ReplyStore struct {
	Id        string `json:"id"`
	CommentId string `json:"comment_id"`
	Reply     string `json:"reply"`
	UserId    string `json:"user_id"`
}
