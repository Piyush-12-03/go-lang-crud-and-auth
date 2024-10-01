package request

type CreateNecheRequest struct {
	Name  string `validate:"required,min=1,max=200" json:"name"`
	TagID    int    `gorm:"not null" json:"tagId"`
}
