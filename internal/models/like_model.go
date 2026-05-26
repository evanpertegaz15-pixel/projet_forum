package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Like struct {
	ID        int
	UserID    int
	PostID    sql.NullInt64
	CommentID sql.NullInt64
	TopicID   sql.NullInt64
	Value     int
	CreatedAt time.Time
}

type LikeModel struct {
	DB *sql.DB
}

func NewLikeModel(db *sql.DB) *LikeModel {
	return &LikeModel{DB: db}
}

func (model *LikeModel) CreateLike(userID, postID, commentID, topicID int) error {
	query := `
		INSERT OR IGNORE INTO likes (user_id, post_id, comment_id, topic_id, value, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	var postVal interface{} = nil
	var commentVal interface{} = nil
	var topicVal interface{} = nil
	if postID > 0 {
		postVal = postID
	}
	if commentID > 0 {
		commentVal = commentID
	}
	if topicID > 0 {
		topicVal = topicID
	}
	_, err := model.DB.Exec(query, userID, postVal, commentVal, topicVal, 1, time.Now())
	return err
}

func (model *LikeModel) DeleteLike(userID, postID, commentID, topicID int) error {
	where, params, err := model.likeWhere(postID, commentID, topicID)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`
		DELETE FROM likes
		WHERE user_id = ? AND %s
	`, where)
	params = append([]interface{}{userID}, params...)
	_, err = model.DB.Exec(query, params...)
	return err
}

func (model *LikeModel) HasLiked(userID, postID, commentID, topicID int) (bool, error) {
	where, params, err := model.likeWhere(postID, commentID, topicID)
	if err != nil {
		return false, err
	}
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM likes
		WHERE user_id = ? AND %s AND value = 1
	`, where)
	params = append([]interface{}{userID}, params...)
	var count int
	err = model.DB.QueryRow(query, params...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (model *LikeModel) CountLikes(postID, commentID, topicID int) (int, error) {
	where, params, err := model.likeWhere(postID, commentID, topicID)
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM likes
		WHERE %s AND value = 1
	`, where)
	var count int
	err = model.DB.QueryRow(query, params...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (model *LikeModel) AddLike(userID, postID int) error {
	return model.CreateLike(userID, postID, 0, 0)
}

func (model *LikeModel) RemoveLike(userID, postID int) error {
	return model.DeleteLike(userID, postID, 0, 0)
}

// UpsertVote inserts or replaces a vote (value: 1=like, -1=dislike).
// If the user already has the same value, it removes the vote (toggle off).
func (model *LikeModel) UpsertVote(userID, postID, commentID, topicID, value int) error {
	where, params, err := model.likeWhere(postID, commentID, topicID)
	if err != nil {
		return err
	}

	// Check existing vote value
	var existing sql.NullInt64
	selectQuery := fmt.Sprintf(`SELECT value FROM likes WHERE user_id = ? AND %s LIMIT 1`, where)
	selectParams := append([]interface{}{userID}, params...)
	_ = model.DB.QueryRow(selectQuery, selectParams...).Scan(&existing)

	// Delete any existing vote first
	delQuery := fmt.Sprintf(`DELETE FROM likes WHERE user_id = ? AND %s`, where)
	delParams := append([]interface{}{userID}, params...)
	if _, err := model.DB.Exec(delQuery, delParams...); err != nil {
		return err
	}

	// If same value as before, just remove (toggle off)
	if existing.Valid && int(existing.Int64) == value {
		return nil
	}

	// Insert new vote
	var postVal, commentVal, topicVal interface{}
	if postID > 0 { postVal = postID }
	if commentID > 0 { commentVal = commentID }
	if topicID > 0 { topicVal = topicID }
	_, err = model.DB.Exec(
		`INSERT INTO likes (user_id, post_id, comment_id, topic_id, value, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		userID, postVal, commentVal, topicVal, value, time.Now(),
	)
	return err
}

func (model *LikeModel) HasDisliked(userID, postID, commentID, topicID int) (bool, error) {
	where, params, err := model.likeWhere(postID, commentID, topicID)
	if err != nil {
		return false, err
	}
	query := fmt.Sprintf(`SELECT COUNT(*) FROM likes WHERE user_id = ? AND %s AND value = -1`, where)
	params = append([]interface{}{userID}, params...)
	var count int
	err = model.DB.QueryRow(query, params...).Scan(&count)
	return count > 0, err
}

func (model *LikeModel) CountDislikes(postID, commentID, topicID int) (int, error) {
	where, params, err := model.likeWhere(postID, commentID, topicID)
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf(`SELECT COUNT(*) FROM likes WHERE %s AND value = -1`, where)
	var count int
	err = model.DB.QueryRow(query, params...).Scan(&count)
	return count, err
}

func (model *LikeModel) likeWhere(postID, commentID, topicID int) (string, []interface{}, error) {
	count := 0
	var clause string
	var params []interface{}
	if postID > 0 {
		count++
		clause = "post_id = ? AND comment_id IS NULL AND topic_id IS NULL"
		params = append(params, postID)
	}
	if commentID > 0 {
		count++
		clause = "comment_id = ? AND post_id IS NULL AND topic_id IS NULL"
		params = append(params, commentID)
	}
	if topicID > 0 {
		count++
		clause = "topic_id = ? AND post_id IS NULL AND comment_id IS NULL"
		params = append(params, topicID)
	}
	if count != 1 {
		return "", nil, errors.New("invalid like target")
	}
	return clause, params, nil
}