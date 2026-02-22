package main

type Record struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type Database struct {
	users   map[string]*User
	records map[string][]Record
}

func NewDatabase() *Database {
	return &Database{
		users:   make(map[string]*User),
		records: make(map[string][]Record),
	}
}

func (db *Database) AppendRecord(username string, record Record) {
	db.records[username] = append(db.records[username], record)
}
