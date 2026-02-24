package records

import (
	"github.com/bushmv/remember/src/services/records/data/dto"
	"github.com/bushmv/remember/src/services/records/data/repository"
)

type RecordsInteractor struct {
	repo repository.RecordRepository
}

func NewRecordsInteractor(repo repository.RecordRepository) *RecordsInteractor {
	return &RecordsInteractor{
		repo: repo,
	}
}

func (i *RecordsInteractor) AddRecord(userId int64, newRecord *dto.NewRecordRequest) *dto.RecordResponse {
	return i.repo.InsertRecord(userId, newRecord)
}

func (i *RecordsInteractor) AllRecords(userId int64) []dto.RecordShortResponse {
	return i.repo.GetAllRecordsForUser(userId)
}

func (i *RecordsInteractor) RecordsForLearning(userId int64) []dto.RecordShortResponse {
	return i.repo.GetAllRecordsForUser(userId)
}

func (i *RecordsInteractor) RecordOwner(recordId int64) (int64, error) {
	return i.repo.GetOwnerForRecord(recordId)
}

func (i *RecordsInteractor) UpdateRecordLevel(recordId int64, newLevel int) (*dto.UpdateRecordLevelResponse, error) {
	return i.repo.UpdateRecordLevel(recordId, newLevel)
}

func (i *RecordsInteractor) Record(recordId int64) (*dto.RecordResponse, error) {
	return i.repo.GetRecord(recordId)
}
