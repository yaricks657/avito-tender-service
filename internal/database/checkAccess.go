package database

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserNotInOrganization = errors.New("user not in organization")
)

// проверка пользователя на ответственного в организации
func (strg *Storage) IsUserInOrganization(username string, organizationID string) (error, int) {
	db := strg.Mng.Db

	var isInOrganization, userExists bool
	query := `
        WITH user_exists AS (
            SELECT id 
            FROM employee 
            WHERE username = $1
        )
        SELECT 
            EXISTS (
                SELECT 1 
                FROM organization_responsible AS orp
                WHERE orp.user_id = (SELECT id FROM user_exists) AND orp.organization_id = $2
            ) AS is_in_organization,
            EXISTS (
                SELECT 1 
                FROM user_exists
            ) AS user_exists;
    `

	err := db.QueryRow(query, username, organizationID).Scan(&isInOrganization, &userExists)
	if err != nil {
		return err, http.StatusBadGateway
	}

	// Проверка, существует ли пользователь
	if !userExists {
		return ErrUserNotFound, http.StatusForbidden
	}

	// Проверка, состоит ли пользователь в организации
	if !isInOrganization {
		return ErrUserNotInOrganization, http.StatusUnauthorized
	}
	return nil, 200
}
