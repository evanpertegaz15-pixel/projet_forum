package models

import (
    "database/sql"
    "time"
)

type Post struct {
    ID         int
    TopicID    int
    UserID     int
    Username   string
    TopicTitle string
    Content    string
    ParentID   *int
    CreatedAt  time.Time
    CreatedAtAgo string
    LikesCount    int
    DislikesCount int
    Images     []Image
}

type PostWithReplies struct {
    Post    Post
    Replies []Post
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
        WHERE p.topic_id = ? AND p.parent_id IS NULL
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

func (model *PostModel) GetFilteredPosts(categoryID, authorID, likedByUserID int) ([]Post, error) {
    query := `
        SELECT p.id, p.topic_id, p.user_id, u.username, p.content, p.parent_id, p.created_at, t.title
        FROM posts p
        JOIN topics t ON p.topic_id = t.id
        JOIN users u ON p.user_id = u.id
        WHERE 1=1`
    args := []interface{}{}
    if categoryID > 0 {
        query += ` AND t.category_id = ?`
        args = append(args, categoryID)
    }
    if authorID > 0 {
        query += ` AND p.user_id = ?`
        args = append(args, authorID)
    }
    if likedByUserID > 0 {
        query += ` AND p.id IN (SELECT post_id FROM likes WHERE user_id = ? AND value = 1 AND post_id IS NOT NULL)`
        args = append(args, likedByUserID)
    }
    query += ` ORDER BY p.created_at DESC`
    rows, err := model.DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.TopicID, &post.UserID, &post.Username, &post.Content, &post.ParentID, &post.CreatedAt, &post.TopicTitle)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func (model *PostModel) Delete(postID int) error {
	_, err := model.DB.Exec(`DELETE FROM posts WHERE id = ?`, postID)
	return err
}

func (model *PostModel) CreatePost(topicID, userID int, content string) (int, error) {
    result, err := model.DB.Exec(`INSERT INTO posts (topic_id, user_id, content, parent_id) VALUES (?, ?, ?, NULL)`, topicID, userID, content)
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    return int(id), nil
}

func (model *PostModel) CreateReply(topicID, userID, parentID int, content string) (int, error) {
    result, err := model.DB.Exec(`INSERT INTO posts (topic_id, user_id, content, parent_id) VALUES (?, ?, ?, ?)`, topicID, userID, content, parentID)
    if err != nil {
        return 0, err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    return int(id), nil
}

func (model *PostModel) GetReplies(postID int) ([]Post, error) {
    rows, err := model.DB.Query(`
        SELECT p.id, p.topic_id, p.user_id, u.username, p.content, p.parent_id, p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.parent_id = ?
        ORDER BY p.created_at ASC
    `, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var replies []Post
    for rows.Next() {
        var p Post
        err := rows.Scan(&p.ID, &p.TopicID, &p.UserID, &p.Username, &p.Content, &p.ParentID, &p.CreatedAt)
        if err != nil {
            return nil, err
        }
        replies = append(replies, p)
    }
    return replies, nil
}

func (model *PostModel) GetPostsWithRepliesByTopic(topicID int) ([]PostWithReplies, error) {
    mainPosts, err := model.GetPostsByTopic(topicID)
    if err != nil {
        return nil, err
    }
    var postsWithReplies []PostWithReplies
    for _, mainPost := range mainPosts {
        replies, err := model.GetReplies(mainPost.ID)
        if err != nil {
            return nil, err
        }
        postsWithReplies = append(postsWithReplies, PostWithReplies{
            Post:    mainPost,
            Replies: replies,
        })
    }
    return postsWithReplies, nil
}

func (model *PostModel) GetPostByID(postID int) (Post, error) {
    row := model.DB.QueryRow(`
        SELECT p.id, p.topic_id, p.user_id, u.username, p.content, p.parent_id, p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = ?
    `, postID)
    var post Post
    var parentID sql.NullInt64
    err := row.Scan(&post.ID, &post.TopicID, &post.UserID, &post.Username, &post.Content, &parentID, &post.CreatedAt)
    if err != nil {
        return Post{}, err
    }
    if parentID.Valid {
        parentIDInt := int(parentID.Int64)
        post.ParentID = &parentIDInt
    }
    return post, nil
}