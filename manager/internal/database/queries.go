package database

const (
	createRolesQuery = `
		INSERT INTO roles (role_type)
		VALUES ($1), ($2), ($3) ($4)
		ON CONFLICT DO NOTHING`
	getRoleQuery = `
		SELECT * FORM roles 
		WHERE id = ($1)
	`
	createUserQuery = `
		INSERT INTO users (username, password, role_id)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
		RETURNING id`

	getUserQuery = `
		SELECT * FROM users
		WHERE username = $1`

	getUserRoleQuery = `
        SELECT username, role_type FROM users
        JOIN roles ON roles.id = users.role_id         
        WHERE id = $1`

	getRoleIdByTypeQuery = `
		SELECT id FROM roles 
		WHERE role_type = $1`

	getUserRoleByUsernameQuery = `
       SELECT users.id, role_type FROM users
       JOIN roles ON role_id = roles.id
       WHERE users.username = $1`

	updateUserRoleQuery = `
    	UPDATE users 
    	SET role_id = (
    		SELECT id FROM roles 
    		WHERE role_type = $2
    		)
		WHERE id = $1`

	getUsernamesQuery = `
		SELECT username FROM users
		ORDER BY id ASC
		OFFSET $1 LIMIT $2`
	getUsernamesCountQuery = `
		SELECT count(*) FROM users`
)
