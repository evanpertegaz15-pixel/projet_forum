package models
// attribuer, retirer un rôle et vérifier si un utilisateur l'a ou non
import (
	"database/sql"
)

type UserRoleModel struct {
	DB *sql.DB
}

func NewUserRoleModel(db *sql.DB) *UserRoleModel {
	return &UserRoleModel{DB: db}
}

// Attribuer un rôle à un utilisateur
func (m *UserRoleModel) AssignRole(userID, roleID int) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
	`

	_, err := m.DB.Exec(query, userID, roleID)
	return err
}

// Retirer un rôle à un utilisateur
func (m *UserRoleModel) RemoveRole(userID, roleID int) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = ? AND role_id = ?
	`

	_, err := m.DB.Exec(query, userID, roleID)
	return err
}

// Vérifier si un utilisateur possède un rôle (par ID)
func (m *UserRoleModel) HasRoleByID(userID, roleID int) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM user_roles
		WHERE user_id = ? AND role_id = ?
	`

	var count int
	err := m.DB.QueryRow(query, userID, roleID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Vérifier si un utilisateur possède un rôle (par nom)
func (m *UserRoleModel) HasRoleByName(userID int, roleName string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = ? AND r.name = ?
	`

	var count int
	err := m.DB.QueryRow(query, userID, roleName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Lister tous les rôles d’un utilisateur
func (m *UserRoleModel) GetRolesByUserID(userID int) ([]Role, error) {
	query := `
		SELECT r.id, r.name, r.permissions
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role

	for rows.Next() {
		var r Role
		err := rows.Scan(&r.ID, &r.Name, &r.Permissions)
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}

// Retirer tous les rôles d’un utilisateur (utile pour reset admin)
func (m *UserRoleModel) ClearUserRoles(userID int) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = ?
	`

	_, err := m.DB.Exec(query, userID)
	return err
}