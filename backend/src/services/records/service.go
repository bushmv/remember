package records

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bushmv/remember/src/services/records/data/dto"
	"github.com/gin-gonic/gin"
)

type RecordsService struct {
	interactor RecordsInteractor
}

func NewRecordsService(interactor *RecordsInteractor) *RecordsService {
	return &RecordsService{
		interactor: *interactor,
	}
}

func (s *RecordsService) AddRecord(c *gin.Context) {
	userId := c.GetInt64("userId")
	username := c.GetString("username")
	var r dto.NewRecordRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	record := s.interactor.AddRecord(userId, &r)
	resp := dto.CreatedRecordResponse{
		Id:          record.Id,
		Title:       record.Title,
		Description: record.Description,
		Owner:       username,
	}
	c.JSON(http.StatusCreated, resp)
}

func (s *RecordsService) AllRecords(c *gin.Context) {
	userId := c.GetInt64("userId")
	records := s.interactor.AllRecords(userId)
	resp := dto.RecordsListResponse{
		CountOfRecords: len(records),
		Records:        records,
	}
	c.JSON(http.StatusOK, resp)
}

func (s *RecordsService) Learn(c *gin.Context) {
	userId := c.GetInt64("userId")
	records := s.interactor.RecordsForLearning(userId)
	resp := dto.RecordsListResponse{
		CountOfRecords: len(records),
		Records:        records,
	}
	c.JSON(http.StatusOK, resp)
}

func (s *RecordsService) UpdateRecordLevel(c *gin.Context) {
	var r dto.UpdateRecordLevelRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userId := c.GetInt64("userId")
	ownerId, err := s.interactor.RecordOwner(r.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if userId != ownerId {
		err := errors.New("No permission to change this record")
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	resp, err := s.interactor.UpdateRecordLevel(r.Id, r.NewLevel)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *RecordsService) GetRecord(c *gin.Context) {
	recordIdParam := c.Param("recordId")
	userId := c.GetInt64("userId")
	recordId, err := strconv.ParseInt(recordIdParam, 10, 64)
	if err != nil || recordId < 0 {
		err = errors.New("Not valid record id (should be positive integer)")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ownerId, err := s.interactor.RecordOwner(recordId)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if ownerId != userId {
		err := errors.New("No permission to change this record")
		c.JSON(http.StatusForbidden, errorResponse(err))
		return

	}
	resp, err := s.interactor.Record(recordId)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *RecordsService) BindRoutes(r *gin.Engine, m gin.HandlerFunc) {
	sr := r.Group("/api/v1/records")
	sr.GET("/all", m, s.AllRecords)
	sr.GET("/learn", m, s.Learn)
	sr.GET("/:recordId", m, s.GetRecord)

	sr.POST("/add", m, s.AddRecord)

	sr.PATCH("/update", m, s.UpdateRecordLevel)
	// sr.DELETE("/delete", m, s.DeleteRecord)
}

func errorResponse(err error) map[string]any {
	return map[string]any{
		"error": err.Error(),
	}
}
