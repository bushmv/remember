package repository

import (
	"github.com/bushmv/remember/src/services/records/data/dto"
)

type RecordRepository interface {
	InsertRecord(userId int64, record *dto.NewRecordRequest) *dto.RecordResponse
	GetAllRecordsForUser(userId int64) []dto.RecordShortResponse
	GetRecordsForLearninigForUser(userId int64) []dto.RecordShortResponse
	UpdateRecordLevel(recordId int64, newLevel int) (*dto.UpdateRecordLevelResponse, error)
	GetOwnerForRecord(recordId int64) (int64, error)
	GetRecord(recordId int64) (*dto.RecordResponse, error)
}
