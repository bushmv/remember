package repository

import (
	"errors"
	"time"

	"github.com/bushmv/remember/src/services/records/data/dto"
)

type RecordRow struct {
	id              int64
	title           string
	description     string
	level           int
	repeatTimestamp int64
	ownerId         int64
}

type InMemoryRepository struct {
	table []RecordRow
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		table: make([]RecordRow, 0),
	}
}

func (r *InMemoryRepository) InsertRecord(userId int64, record *dto.NewRecordRequest) *dto.RecordResponse {
	row := RecordRow{
		id:              int64(len(r.table)),
		title:           record.Title,
		description:     record.Description,
		ownerId:         userId,
		repeatTimestamp: time.Now().Unix(),
		level:           0,
	}
	r.table = append(r.table, row)
	return &dto.RecordResponse{
		Id:              row.id,
		Title:           row.title,
		Description:     row.description,
		Level:           0,
		RepeatTimestamp: row.repeatTimestamp,
	}
}

func (r *InMemoryRepository) GetAllRecordsForUser(userId int64) []dto.RecordShortResponse {
	result := make([]dto.RecordShortResponse, 0)
	for _, row := range r.table {
		if row.ownerId == userId {
			record := dto.RecordShortResponse{
				Id:    row.id,
				Title: row.title,
				Level: row.level,
			}
			result = append(result, record)
		}
	}
	return result
}
func (r *InMemoryRepository) GetRecordsForLearninigForUser(userId int64) []dto.RecordShortResponse {
	result := make([]dto.RecordShortResponse, 0)
	now := time.Now().Unix()
	for _, row := range r.table {
		if row.ownerId == userId && now >= row.repeatTimestamp {
			record := dto.RecordShortResponse{
				Id:    row.id,
				Title: row.title,
				Level: row.level,
			}
			result = append(result, record)
		}
	}
	return result

}
func (r *InMemoryRepository) UpdateRecordLevel(recordId int64, newLevel int) (*dto.UpdateRecordLevelResponse, error) {
	if recordId < 0 || recordId >= int64(len(r.table)) {
		return nil, errors.New("Record not exists")
	}
	row := &r.table[recordId]
	oldLevel := row.level
	row.level = newLevel
	record := &dto.UpdateRecordLevelResponse{
		Id:       row.id,
		OldLevel: oldLevel,
		NewLevel: newLevel,
		Title:    row.title,
	}
	return record, nil
}

func (r *InMemoryRepository) GetOwnerForRecord(recordId int64) (int64, error) {
	if recordId < 0 || recordId >= int64(len(r.table)) {
		return 0, errors.New("Record no exists")
	}
	row := r.table[recordId]
	return row.ownerId, nil
}

func (r *InMemoryRepository) GetRecord(recordId int64) (*dto.RecordResponse, error) {
	if recordId < 0 || recordId >= int64(len(r.table)) {
		return nil, errors.New("Record no exists")
	}
	row := r.table[recordId]
	record := &dto.RecordResponse{
		Id:              row.id,
		Title:           row.title,
		Description:     row.description,
		Level:           row.level,
		RepeatTimestamp: row.repeatTimestamp,
	}
	return record, nil
}
