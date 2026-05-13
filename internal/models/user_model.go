package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID	int
	Username	string
	Email	string
	Password	string
	ProfilePicture	sql.NullString
	CreatedAt time.Time
	CreatedAtAgo string
	UpdatedAt *time.Time
	Roles []Role
}

type UserModel struct {
	DB *sql.DB
	Roles *RoleModel
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db, Roles: NewRoleModel(db)}
}

func (model *UserModel) CreateUser(email, username, password string) (int, error) {
	result, err := model.DB.Exec(`
		INSERT INTO users (email, username, password_hash)
		VALUES (?, ?, ?)
	`, email, username, password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	userID := int(id)
	defaultRole, err := model.Roles.GetRoleByName("user")
	if err == nil && defaultRole != nil {
		_ = model.Roles.AddRoleToUser(userID, defaultRole.ID)
	}
	return userID, nil
}

func (user *User) HasRole(roleName string) bool {
	for _, r := range user.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

func (model *UserModel) loadRoles(user *User) error {
	if model.Roles == nil || user == nil {
		return nil
	}
	roles, err := model.Roles.GetRolesForUser(user.ID)
	if err != nil {
		return err
	}
	user.Roles = roles
	return nil
}

func (model *UserModel) FindByEmail(email string) (*User, error) {
	row := model.DB.QueryRow(`
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?
	`, email)
	var user User
	var updatedAt sql.NullTime
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}
	if err := model.loadRoles(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (model *UserModel) FindByUsername(username string) (*User, error) {
    row := model.DB.QueryRow(`
        SELECT id, email, username, password_hash, created_at, updated_at
        FROM users
        WHERE username = ?
    `, username)
    var user User
    var updatedAt sql.NullTime
    err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &updatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    if updatedAt.Valid {
        user.UpdatedAt = &updatedAt.Time
    }
    if err := model.loadRoles(&user); err != nil {
        return nil, err
    }
    return &user, nil
}

func (model *UserModel) FindByID(id int) (*User, error) {
    row := model.DB.QueryRow(`
        SELECT id, email, username, password_hash, profile_picture, created_at, updated_at
        FROM users
        WHERE id = ?
    `, id)
    var user User
    var updatedAt sql.NullTime
    err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.ProfilePicture, &user.CreatedAt, &updatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    if updatedAt.Valid {
        user.UpdatedAt = &updatedAt.Time
    }
    if err := model.loadRoles(&user); err != nil {
        return nil, err
    }
    return &user, nil
}

func (model *UserModel) GetAllUsers() ([]User, error) {
    rows, err := model.DB.Query(`
        SELECT id, username, email, created_at
        FROM users
        ORDER BY id ASC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
            return nil, err
        }
        if err := model.loadRoles(&user); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func (user *User) GetRole() {}

func (model *UserModel) DeleteUser(id int) error {
    _, err := model.DB.Exec("DELETE FROM users WHERE id = ?", id)
    return err
}