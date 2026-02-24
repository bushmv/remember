package dto

type UpdateRecordLevelResponse struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	OldLevel int    `json:"old level"`
	NewLevel int    `json:"new level"`
}
