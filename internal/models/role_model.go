package models

import (
    "database/sql"
)

type Role struct {
    ID    int
    Name  string // valeur interne : "admin", "moderator", "ranger", "user", "blocked"
    Label string // affichage : "Administrateur", "Modérateur", etc.
}

type UserRole struct {
    UserID int
    RoleID int
}

type RoleModel struct {
    DB *sql.DB
}

func NewRoleModel(db *sql.DB) *RoleModel {
    return &RoleModel{DB: db}
}

func (model *RoleModel) GetRolesForUser(userID int) ([]Role, error) {
    rows, err := model.DB.Query(`
        SELECT r.id, r.name, r.label
        FROM roles r
        INNER JOIN user_roles ur ON ur.role_id = r.id
        WHERE ur.user_id = ?
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    roles := []Role{}
    for rows.Next() {
        var r Role
        if err := rows.Scan(&r.ID, &r.Name, &r.Label); err != nil {
            return nil, err
        }
        roles = append(roles, r)
    }
    return roles, nil
}

func (model *RoleModel) AddRoleToUser(userID, roleID int) error {
    _, err := model.DB.Exec(`
        INSERT INTO user_roles (user_id, role_id)
        VALUES (?, ?)
    `, userID, roleID)
    return err
}

func (model *RoleModel) RemoveRoleFromUser(userID, roleID int) error {
    _, err := model.DB.Exec(`
        DELETE FROM user_roles
        WHERE user_id = ? AND role_id = ?
    `, userID, roleID)
    return err
}

func (model *RoleModel) GetRoleByName(roleName string) (*Role, error) {
    row := model.DB.QueryRow(`
        SELECT id, name, label
        FROM roles
        WHERE name = ?
    `, roleName)
    var role Role
    err := row.Scan(&role.ID, &role.Name, &role.Label)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &role, nil
}

func (model *RoleModel) UserHasRole(userID int, roleName string) (bool, error) {
    row := model.DB.QueryRow(`
        SELECT COUNT(*)
        FROM roles r
        INNER JOIN user_roles ur ON ur.role_id = r.id
        WHERE ur.user_id = ? AND r.name = ?
    `, userID, roleName)
    var count int
    if err := row.Scan(&count); err != nil {
        return false, err
    }
    return count > 0, nil
}