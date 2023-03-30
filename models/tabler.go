package models

// for overriding gorm table name
type Tabler interface {
	TableName() string
}
