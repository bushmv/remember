package dto

type UpdateRecordLevelRequest struct {
	Id       int64 `json:"id"`
	NewLevel int   `json:"new level"`
}
