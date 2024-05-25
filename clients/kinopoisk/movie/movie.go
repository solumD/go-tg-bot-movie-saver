package movie

type Movie struct {
	Id          int    `json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	Rating      Rating `json:"rating"`
	Age         int    `json:"ageRating"`
	Year        int    `json:"year"`
	Length      int    `json:"movieLength"`
}

type Rating struct {
	KpRating float32 `json:"kp"`
}
