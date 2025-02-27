package entities

type ForumResponse struct {
	Id           string         `json:"id"`
	Title        string         `json:"title"`
	Caption      string         `json:"caption"`
	Media        []ForumMedia   `json:"medias"`
	Comment      []ForumComment `json:"comments"`
	Like         []ForumLike    `json:"likes"`
	CommentCount int            `json:"comment_count"`
	LikeCount    int            `json:"like_count"`
	ForumType    ForumType      `json:"type"`
	User         ForumUser      `json:"user"`
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
	Id   int    `json:"id"`
	Path string `json:"path"`
}

type ForumLikeQuery struct {
	Id       string `json:"id"`
	UserId   string `json:"user_id"`
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname"`
}

type ForumLike struct {
	Id   string        `json:"id"`
	User ForumLikeUser `json:"user"`
}

type ForumCommentQuery struct {
	Id       string `json:"id"`
	Comment  string `json:"comment"`
	Avatar   string `json:"avatar"`
	UserId   string `json:"user_id"`
	Fullname string `json:"fullname"`
}

type ForumComment struct {
	Id         string              `json:"id"`
	Comment    string              `json:"comment"`
	User       ForumCommentUser    `json:"user"`
	Reply      []ForumCommentReply `json:"replies"`
	ReplyCount int                 `json:"reply_count"`
}

type ForumCommentUser struct {
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname"`
}

type ForumCommentReplyQuery struct {
	Id       string `json:"id"`
	Reply    string `json:"reply"`
	Avatar   string `json:"avatar"`
	UserId   string `json:"user_id"`
	Fullname string `json:"fullname"`
}

type ForumCommentReply struct {
	Id    string                `json:"id"`
	Reply string                `json:"reply"`
	User  ForumCommentReplyUser `json:"user"`
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
