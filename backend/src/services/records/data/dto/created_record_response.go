package dto

type CreatedRecordResponse struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
}
