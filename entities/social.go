package entities

type SocialMediaStore struct {
	Icon string `json:"icon"`
	Link string `json:"link"`
	Name string `json:"name"`
}

type SocialMediaStoreResponse struct {
	Id   int64  `json:"id"`
	Icon string `json:"icon"`
	Link string `json:"link"`
	Name string `json:"name"`
}

type SocialMediaList struct {
	Id   int64  `json:"id"`
	Icon string `json:"icon"`
	Link string `json:"link"`
	Name string `json:"name"`
}
