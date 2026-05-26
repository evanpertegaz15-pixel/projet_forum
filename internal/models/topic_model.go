package models

import (
    "database/sql"
    "time"
)

type Topic struct {
    ID         int
    CategoryID int
    UserID     int
    Title      string
    CreatedAt  time.Time
    CreatedAtAgo string
    LikesCount int
}

type TopicModel struct {
    DB *sql.DB
}

func NewTopicModel(db *sql.DB) *TopicModel {
    return &TopicModel{DB: db}
}

func (model *TopicModel) GetTopicByID(id int) (Topic, error) {
    row := model.DB.QueryRow(`
        SELECT id, category_id, user_id, title, created_at
        FROM topics
        WHERE id = ?
    `, id)
    var topic Topic
    err := row.Scan(&topic.ID, &topic.CategoryID, &topic.UserID, &topic.Title, &topic.CreatedAt)
    if err != nil {
        return Topic{}, err
    }
    return topic, nil
}

func (model *TopicModel) GetTopicsByCategory(categoryID int) ([]Topic, error) {
    rows, err := model.DB.Query(`
        SELECT id, category_id, user_id, title, created_at
        FROM topics
        WHERE category_id = ?
        ORDER BY created_at DESC
    `, categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var topics []Topic
    for rows.Next() {
        var topic Topic
        if err := rows.Scan(&topic.ID, &topic.CategoryID, &topic.UserID, &topic.Title, &topic.CreatedAt); err != nil {
            return nil, err
        }
        topics = append(topics, topic)
    }
    return topics, nil
}

func (model *TopicModel) CreateTopic(categoryID, userID int, title string) (int, error) {
    result, err := model.DB.Exec(`INSERT INTO topics (category_id, user_id, title) VALUES (?, ?, ?)`, categoryID, userID, title)
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    return int(id), nil
}

func (model *TopicModel) DeleteTopic(topicID int) error {
    _, err := model.DB.Exec(`DELETE FROM posts WHERE topic_id = ?`, topicID)
    if err != nil {
        return err
    }
    _, err = model.DB.Exec(`DELETE FROM topics WHERE id = ?`, topicID)
    return err
}

func (model *TopicModel) GetLikedTopicsByUser(userID, categoryID int) ([]Topic, error) {
    query := `SELECT t.id, t.category_id, t.user_id, t.title, t.created_at
        FROM topics t
        WHERE t.id IN (SELECT topic_id FROM likes WHERE user_id = ? AND value = 1 AND topic_id IS NOT NULL)`
    args := []interface{}{userID}
    if categoryID > 0 {
        query += ` AND t.category_id = ?`
        args = append(args, categoryID)
    }
    query += ` ORDER BY t.created_at DESC`
    rows, err := model.DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var topics []Topic
    for rows.Next() {
        var t Topic
        if err := rows.Scan(&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.CreatedAt); err != nil {
            return nil, err
        }
        topics = append(topics, t)
    }
    return topics, nil
}