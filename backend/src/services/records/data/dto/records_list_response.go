package dto

type RecordsListResponse struct {
	CountOfRecords int                   `json:"count of records"`
	Records        []RecordShortResponse `json:"records"`
}
