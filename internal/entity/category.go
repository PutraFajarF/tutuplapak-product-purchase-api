package entity

type Category struct {
	Type string `json:"type" gorm:"column:type;primaryKey"`
}
