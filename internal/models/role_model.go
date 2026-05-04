package models
// lister les rôles, les charger pour l'admin, vérifier les perms
import (
	"database/sql"
)

type Role struct {
	ID          int
	Name        string
	Permissions string // simple version (JSON ou CSV selon ton choix)
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

// Lister tous les rôles (admin panel)
func (m *RoleModel) GetAllRoles() ([]Role, error) {
	query := `
		SELECT id, name, permissions
		FROM roles
	`

	rows, err := m.DB.Query(query)
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

// Récupérer un rôle par ID
func (m *RoleModel) GetRoleByID(roleID int) (*Role, error) {
	query := `
		SELECT id, name, permissions
		FROM roles
		WHERE id = ?
	`

	var r Role

	err := m.DB.QueryRow(query, roleID).Scan(&r.ID, &r.Name, &r.Permissions)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Récupérer les rôles d’un utilisateur
func (m *RoleModel) GetUserRoles(userID int) ([]Role, error) {
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

// Vérifier si un utilisateur a un rôle précis
func (m *RoleModel) UserHasRole(userID int, roleName string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.name = ?
	`

	var count int
	err := m.DB.QueryRow(query, userID, roleName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Vérifier une permission (version simple string contains)
func (m *RoleModel) UserHasPermission(userID int, permission string) (bool, error) {
	query := `
		SELECT r.permissions
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?
	`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var perms string
		if err := rows.Scan(&perms); err != nil {
			return false, err
		}

		// Version simple (tu peux upgrader en JSON plus tard)
		if containsPermission(perms, permission) {
			return true, nil
		}
	}

	return false, nil
}

// helper simple (CSV style: "ban_user,delete_post,edit_role")
func containsPermission(perms string, permission string) bool {
	// simple check sans dépendance JSON
	// améliorable en strings.Split
	for i := 0; i < len(perms)-len(permission); i++ {
		if perms[i:i+len(permission)] == permission {
			return true
		}
	}
	return false
}
