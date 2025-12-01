package repository

import (
	"bilibili-clone/internal/interaction/model"
	"bilibili-clone/pkg/database"
)

type CommentRepository struct{}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return database.DB.Create(comment).Error
}

func (r *CommentRepository) FindByVideoID(videoID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := database.DB.Where("video_id = ?", videoID).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) Delete(commentID, userID uint) error {
	return database.DB.Where("id = ? AND user_id = ?", commentID, userID).Delete(&model.Comment{}).Error
}
