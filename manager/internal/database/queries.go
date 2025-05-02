package database

const (
	truncateAllTablesQuery = `
		TRUNCATE TABLE submissions, questions, users, roles;`

	createRolesQuery = `
		INSERT INTO roles (role_type)
		VALUES ($1), ($2), ($3), ($4)
		ON CONFLICT DO NOTHING`
	getRoleQuery = `
		SELECT * FROM roles 
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
        WHERE users.id = $1`

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

	getUserStatsQuery = `
		SELECT 
			COUNT(DISTINCT question_id) AS tried_count,
			COUNT(DISTINCT CASE WHEN state = $2 THEN question_id END) AS success_count,
		FROM submission
		WHERE user_id = $1`

	getQuestionsCountQuery = `SELECT count(*) FROM questions`

	getQuestionsQuery = `
		SELECT questions.id, title, state, users.username
		FROM questions
		JOIN users ON users.id = questions.owner
		ORDER BY id ASC
		OFFSET $1 LIMIT $2`

	getQuestionsCountWithStateQuery = `
		SELECT count(*) FROM questions
		WHERE state = $1`

	getQuestionsWithStateQuery = `
		SELECT questions.id, title, state, users.username
		FROM questions
		JOIN users ON users.id = questions.owner
		WHERE state = $1
		ORDER BY id ASC
		OFFSET $2 LIMIT $3`

	getUserQuestionsCountQuery = `
		SELECT count(*) FROM questions
		WHERE owner = $1`
	getUserQuestionsQuery = `
		SELECT questions.id, title, state
		FROM questions
		WHERE owner = $1
		ORDER BY id ASC
		OFFSET $2 LIMIT $3`

	getQuestionQuery = `
		SELECT questions.id, title, statement, "input", "output", memory_limit, time_limit, state, username
		FROM questions 
		JOIN users ON users.id = questions.owner
		WHERE questions.id = $1`

	changeQuestionStateQuery = `
		UPDATE questions 
		SET state = $2
		WHERE id = $1
		`
	createQuestionQuery = `
		INSERT INTO questions (title, statement, owner, input, output, memory_limit, time_limit, state)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	createSubmissionQuery = `
		INSERT INTO submissions (user_id, question_id, code, state)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	selectSubmissionForUpdateQuery = `
		SELECT id, state, retry_count, state_updated_at
		FROM submissions
		WHERE id = $1
		FOR UPDATE`

	updateSubmissionStateQuery = `
		UPDATE submissions
		SET state = $2, retry_count = $3, state_updated_at = now()
		WHERE id = $1
		`

	getSubmissionsWithStateCountQuery = `
		SELECT count(*) FROM submissions
		WHERE state = $1`

	getSubmissionsWithStateQuery = `
		SELECT id, code, question_id, state
		FROM submissions 
		WHERE state = $1
		ORDER BY id ASC
		OFFSET $2 LIMIT $3`

	getUserQuestionSubmissionsCountQuery = `
		SELECT count(*) FROM submissions 
		WHERE user_id = $1 and question_id = $2`
	getUserQuestionSubmissionsQuery = `
		SELECT id, code, question_id, state
		FROM submissions 
		WHERE user_id = $1 and question_id = $2
		ORDER BY id
		OFFSET $3 LIMIT $4`

	getUserAllSubmissionsCountQuery = `
		SELECT count(*) FROM submissions
		WHERE user_id = $1`

	getUserAllSubmissionsQuery = `
		SELECT id, code, question_id, state
		FROM submissions
		WHERE user_id = $1
		ORDER BY id
		OFFSET $2 LIMIT $3`
)
