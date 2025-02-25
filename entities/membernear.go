package entities

type Membernear struct {
	Uid		string `json:"id"`
	Lat		float64 `json:"lat"`
	Lng		float64 `json:"lng"`
}

type MembernearAssign struct {
	Lat		 float64 `json:"lat"`
	Lng		 float64 `json:"lng"`
	Distance string  `json:"distance"`
}

type MembernearResponse struct {
	Lat		 float64 `json:"lat"`
	Lng		 float64 `json:"lng"`
	Distance string  `json:"distance"`
}