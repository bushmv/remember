package dto

type RecordShortResponse struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Level int    `json:"level"`
}
