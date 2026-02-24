package dto

type RecordResponse struct {
	Id              int64  `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Level           int    `json:"level"`
	RepeatTimestamp int64  `json:"repeat timestamp"`
}
