package models

import (
	"database/sql"
)

// HomePageData holds all data needed to render the home page.
type HomePageData struct {
	Users       []User
	Topics      []Topic
	CategoryMap map[int]string
}

// HomeModel centralises the queries required by the home page.
type HomeModel struct {
	DB *sql.DB
}

func NewHomeModel(db *sql.DB) *HomeModel {
	return &HomeModel{DB: db}
}

// GetHomePageData fetches users, the N latest topics (with like counts),
// and a category-ID→name map in one place.
func (m *HomeModel) GetHomePageData(limit int) (HomePageData, error) {
	var data HomePageData

	// ── Users ────────────────────────────────────────────────────────────────
	userRows, err := m.DB.Query(`
		SELECT id, username, email, created_at
		FROM users
		ORDER BY id ASC
	`)
	if err != nil {
		return data, err
	}
	defer userRows.Close()
	for userRows.Next() {
		var u User
		if err := userRows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return data, err
		}
		data.Users = append(data.Users, u)
	}

	// ── Latest topics with like counts ───────────────────────────────────────
	topicRows, err := m.DB.Query(`
		SELECT t.id, t.category_id, t.user_id, t.title, t.created_at,
		       COUNT(l.id) AS likes_count
		FROM topics t
		LEFT JOIN likes l ON l.topic_id = t.id AND l.value = 1
		GROUP BY t.id
		ORDER BY t.created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return data, err
	}
	defer topicRows.Close()
	for topicRows.Next() {
		var t Topic
		if err := topicRows.Scan(&t.ID, &t.CategoryID, &t.UserID, &t.Title, &t.CreatedAt, &t.LikesCount); err != nil {
			return data, err
		}
		data.Topics = append(data.Topics, t)
	}

	// ── Category map ─────────────────────────────────────────────────────────
	catRows, err := m.DB.Query(`SELECT id, name FROM categories ORDER BY id ASC`)
	if err != nil {
		return data, err
	}
	defer catRows.Close()
	data.CategoryMap = make(map[int]string)
	for catRows.Next() {
		var id int
		var name string
		if err := catRows.Scan(&id, &name); err != nil {
			return data, err
		}
		data.CategoryMap[id] = name
	}

	return data, nil
}
