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
}

type TopicModel struct {
    DB *sql.DB
}

func NewTopicModel(db *sql.DB) *TopicModel {
    return &TopicModel{DB: db}
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