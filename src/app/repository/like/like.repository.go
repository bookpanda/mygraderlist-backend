package like

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/emoji"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/like"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindLikeByUserId(userId string, result *[]*like.Like) error {
	return r.db.Model(&emoji.Emoji{}).Find(result, "user_id = ?", userId).Error
}

func (r *Repository) Create(in *emoji.Emoji) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&emoji.Emoji{}).Error
}
