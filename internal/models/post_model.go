package models

import (
    "database/sql"
    "time"
)

type Post struct {
    ID        int
    TopicID   int
    UserID    int
    Username  string
    Content   string
    CreatedAt time.Time
}

type PostModel struct {
    DB *sql.DB
}

func NewPostModel(db *sql.DB) *PostModel {
    return &PostModel{DB: db}
}

func (model *PostModel) GetPostsByTopic(topicID int) ([]Post, error) {
    rows, err := model.DB.Query(`
        SELECT p.id, p.topic_id, p.user_id, u.username, p.content, p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.topic_id = ?
        ORDER BY p.created_at ASC
    `, topicID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var posts []Post
    for rows.Next() {
        var post Post
        if err := rows.Scan(&post.ID, &post.TopicID, &post.UserID, &post.Username, &post.Content, &post.CreatedAt); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func (model *PostModel) CreateTopic(topicID, userID int, content string) (int, error) {
    result, err := model.DB.Exec(`INSERT INTO posts (topic_id, user_id, content) VALUES (?, ?, ?)`, topicID, userID, content)
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    return int(id), nil
}