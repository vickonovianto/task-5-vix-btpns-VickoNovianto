package models

type (
	Photo struct {
		ID       int    `gorm:"column:id"`
		UserId   int    `gorm:"column:user_id"`
		User     *User  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		Title    string `gorm:"column:title;size:50"`
		Caption  string `gorm:"column:caption;size:50"`
		PhotoUrl string `gorm:"column:photo_url;size:255;not null"`
	}

	PhotoResponse struct {
		ID       int    `json:"id_photo"`
		PhotoUrl string `json:"photo_url"`
	}
)

// override gorm table name
func (Photo) TableName() string {
	return "photos"
}
