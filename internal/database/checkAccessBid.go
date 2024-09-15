package database

import (
	"net/http"
)

// проверка пользователя на ответственного в организации
func (strg *Storage) IsUserInOrganizationOrCreator(userId string, userNamee string, organizationID string, tenderCreator string) (error, int) {
	db := strg.Mng.Db

	var userName *string
	var isInOrganization, userExists bool

	if userId == "" {
		query := `
		WITH user_exists AS (
			SELECT id, username 
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
			) AS user_exists,
			(SELECT username FROM user_exists) AS user_name;
		`

		err := db.QueryRow(query, userNamee, organizationID).Scan(&isInOrganization, &userExists, &userName)
		if err != nil {
			return err, http.StatusForbidden
		}
	} else {
		query := `
		WITH user_exists AS (
			SELECT id, username 
			FROM employee 
			WHERE id = $1
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
			) AS user_exists,
			(SELECT username FROM user_exists) AS user_name;
		`

		err := db.QueryRow(query, userId, organizationID).Scan(&isInOrganization, &userExists, &userName)
		if err != nil {
			return err, http.StatusForbidden
		}
	}

	// Проверка, существует ли пользователь
	if !userExists {
		return ErrUserNotFound, http.StatusForbidden
	}

	// Проверка, состоит ли пользователь в организации или создатель
	if !isInOrganization && (userName == nil || *userName != tenderCreator) {
		return ErrUserNotInOrganization, http.StatusUnauthorized
	}

	return nil, 200
}
