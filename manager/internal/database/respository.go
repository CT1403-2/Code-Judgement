package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"manger/internal"
	"math"
	"strings"
	"time"
)

type Repository interface {
	SetUp() error
	SetUpRoles() error
	CreateSuperUserIfNotExists() (int32, error)
	CreateMember(ctx context.Context, username string, password string) (int32, error)
	createUser(ctx context.Context, username string, password string, role proto.Role) (int32, error)
	getUser(ctx context.Context, username string) (user, error)
	getRole(ctx context.Context, roleId int32) (role, error)
	getRoleIdByType(ctx context.Context, role proto.Role) (int32, error)
	Authenticate(ctx context.Context, username string, password string) (int32, int32, error)
	GetUserRole(ctx context.Context, userId int32) (string, proto.Role, error)
	GetUserRoleByUsername(ctx context.Context, username string) (int32, proto.Role, error)
	UpdateUserRole(ctx context.Context, userId int32, role proto.Role) error
	GetUsernames(ctx context.Context, pageNumber, pageSize int) ([]string, int, error)
	GetUserStats(ctx context.Context, userId int32) (int64, int64, error)
	GetQuestions(ctx context.Context, published bool, pageNumber, pageSize int) ([]*proto.Question, int, error)
	GetUserQuestions(ctx context.Context, userId int32,
		username string, pageNumber, pageSize int) ([]*proto.Question, int, error)
	GetQuestion(ctx context.Context, questionId int) (*proto.Question, error)
	ChangeQuestionState(ctx context.Context, questionId int, state int32) error
	CreateQuestion(ctx context.Context, owner int32, question *proto.Question) (int32, error)
	EditQuestion(ctx context.Context, question *proto.Question) error
	CreateSubmission(ctx context.Context, userId int32, questionId int32, code []byte) error
	UpdateSubmissionState(ctx context.Context, submissionId int32, state int32) (bool, error)
	GetSubmissionsWithState(ctx context.Context, state int32, pageNumber, pageSize int) ([]*proto.Submission, int, error)
	GetUserSubmissions(ctx context.Context, userId int32,
		questionId int32, filterQuestion bool, pageNumber, pageSize int) ([]*proto.Submission, int, error)
}

type postgresqlRepository struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewRepository() (Repository, error) {
	config := getMainDbConfig()
	if config.password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable not set")
	}
	connStr := getConnStrFromConfig(config)
	p, err := newRepositoryWithDsn(connStr)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func newRepositoryWithDsn(dsn string) (*postgresqlRepository, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("unable to connect to database")
		return nil, err
	}
	ctx := context.Background()
	p := &postgresqlRepository{ctx: ctx, pool: pool}
	err = p.SetUp()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *postgresqlRepository) SetUp() error {
	err := p.SetUpRoles()
	if err != nil {
		return err
	}
	_, err = p.CreateSuperUserIfNotExists()
	return err
}

func (p *postgresqlRepository) SetUpRoles() error {
	_, err := p.pool.Exec(p.ctx, createRolesQuery,
		proto.Role_ROLE_UNKNOWN, proto.Role_ROLE_MEMBER, proto.Role_ROLE_ADMIN, proto.Role_ROLE_SUPERUSER)
	return err
}

func (p *postgresqlRepository) CreateSuperUserIfNotExists() (int32, error) {
	username := internal.GetEnv("ADMIN_USERNAME", "admin")
	password := internal.GetEnv("ADMIN_PASSWORD", "admin")
	user, err := p.getUser(p.ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p.createUser(p.ctx, username, password, proto.Role_ROLE_SUPERUSER)
		}
		return 0, err
	}
	return user.id, nil
}

func (p *postgresqlRepository) CreateMember(ctx context.Context, username string, password string) (int32, error) {
	return p.createUser(ctx, username, password, proto.Role_ROLE_MEMBER)
}

func (p *postgresqlRepository) createUser(ctx context.Context, username string, password string, role proto.Role) (int32, error) {
	var userId int32
	hashedPassword, err := internal.HashPassword(password)
	if err != nil {
		return userId, err
	}
	roleId, err := p.getRoleIdByType(ctx, role)
	if err != nil {
		return userId, err
	}

	err = p.pool.QueryRow(ctx, createUserQuery, username, hashedPassword, roleId).Scan(&userId)

	return userId, err
}

func (p *postgresqlRepository) getUser(ctx context.Context, username string) (user, error) {
	var u user
	err := p.pool.QueryRow(ctx, getUserQuery, username).Scan(
		&u.id, &u.username, &u.password, &u.roleId, &u.createdAt)
	return u, err
}

func (p *postgresqlRepository) getRoleIdByType(ctx context.Context, r proto.Role) (int32, error) {
	var roleId int32
	err := p.pool.QueryRow(ctx, getRoleIdByTypeQuery, int32(r)).Scan(&roleId)
	return roleId, err
}

func (p *postgresqlRepository) getRole(ctx context.Context, roleId int32) (role, error) {
	var r role
	err := p.pool.QueryRow(ctx, getRoleQuery, roleId).Scan(&r.id, &r.roleType, &r.createdAt)
	return r, err
}

func (p *postgresqlRepository) Authenticate(ctx context.Context, username string, password string) (int32, int32, error) {
	u, err := p.getUser(ctx, username)
	if err != nil {
		return 0, 0, err
	}
	if !internal.CheckPasswordHash(password, u.password) {
		return 0, 0, pgx.ErrNoRows
	}
	r, err := p.getRole(ctx, u.roleId)
	if err != nil {
		return 0, 0, err
	}
	return u.id, r.roleType, nil
}

func (p *postgresqlRepository) GetUserRole(ctx context.Context, userId int32) (string, proto.Role, error) {
	var i int32
	var username string
	err := p.pool.QueryRow(ctx, getUserRoleQuery, userId).Scan(&username, &i)
	if err != nil {
		return "", proto.Role_ROLE_UNKNOWN, err
	}
	return username, proto.Role(i), nil
}

func (p *postgresqlRepository) GetUserRoleByUsername(ctx context.Context, username string) (int32, proto.Role, error) {
	var userId, r int32
	err := p.pool.QueryRow(ctx, getUserRoleByUsernameQuery, username).Scan(&userId, &r)
	if err != nil {
		return userId, proto.Role_ROLE_UNKNOWN, err
	}
	return userId, proto.Role(r), nil
}

func (p *postgresqlRepository) UpdateUserRole(ctx context.Context, userId int32, role proto.Role) error {
	_, err := p.pool.Exec(ctx, updateUserRoleQuery, userId, int32(role))
	return err
}

func (p *postgresqlRepository) GetUsernames(ctx context.Context, pageNumber, pageSize int) ([]string, int, error) {
	offset := (pageNumber - 1) * pageSize

	var count int
	err := p.pool.QueryRow(ctx, getUsernamesCountQuery).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan row: %v", err)
	}

	if pageSize <= 0 {
		return nil, 0, errors.New("negative page size")
	}

	totalPage := int(math.Ceil(float64(count) / float64(pageSize)))

	if totalPage == 0 {
		return nil, 0, pgx.ErrNoRows
	}

	if pageNumber > totalPage {
		return nil, 0, errors.New("out of bounds page number")
	}

	rows, err := p.pool.Query(ctx, getUsernamesQuery, offset, pageSize)
	if err != nil {
		return nil, totalPage, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, totalPage, fmt.Errorf("failed to scan row: %v", err)
		}
		usernames = append(usernames, username)
	}

	if err := rows.Err(); err != nil {
		return nil, totalPage, fmt.Errorf("rows iteration error: %v", err)
	}

	return usernames, totalPage, nil
}

func (p *postgresqlRepository) GetUserStats(ctx context.Context, userId int32) (int64, int64, error) {
	var triedCount, successCount int64
	err := p.pool.QueryRow(ctx, getUserStatsQuery, userId, int32(proto.SubmissionState_SUBMISSION_STATE_OK)).Scan(&triedCount, &successCount)
	return triedCount, successCount, err
}

func (p *postgresqlRepository) GetQuestions(ctx context.Context, publishedOnly bool, pageNumber, pageSize int) (
	[]*proto.Question, int, error) {
	offset := (pageNumber - 1) * pageSize
	var count int

	publishedState := int(proto.QuestionState_QUESTION_STATE_PUBLISHED)
	var countQuery, mainQuery string
	var countQueryArgs, mainQueryArgs []interface{}
	if publishedOnly {
		countQuery = getQuestionsCountWithStateQuery
		countQueryArgs = append(countQueryArgs, publishedState)
		mainQuery = getQuestionsWithStateQuery
		mainQueryArgs = append(mainQueryArgs, publishedState, offset, pageSize)
	} else {
		countQuery = getQuestionsCountQuery
		mainQuery = getQuestionsQuery
		mainQueryArgs = append(mainQueryArgs, offset, pageSize)
	}

	err := p.pool.QueryRow(ctx, countQuery, countQueryArgs...).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan row: %v", err)
	}

	if pageSize <= 0 {
		return nil, 0, errors.New("negative page size")
	}

	totalPage := int(math.Ceil(float64(count) / float64(pageSize)))

	if totalPage == 0 {
		return nil, 0, pgx.ErrNoRows
	}

	if pageNumber > totalPage {
		return nil, 0, errors.New("out of bounds page number")
	}
	var questions []*proto.Question

	rows, err := p.pool.Query(ctx, mainQuery, mainQueryArgs...)
	if err != nil {
		return nil, totalPage, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var question proto.Question
		err := rows.Scan(&question.Id, &question.Title, &question.State, &question.Owner)
		if err != nil {
			return nil, totalPage, fmt.Errorf("failed to scan row: %v", err)
		}
		questions = append(questions, &question)
	}
	if err := rows.Err(); err != nil {
		return nil, totalPage, fmt.Errorf("rows iteration error: %v", err)
	}
	return questions, totalPage, nil
}

func (p *postgresqlRepository) GetUserQuestions(ctx context.Context, userId int32,
	username string, pageNumber, pageSize int) ([]*proto.Question, int, error) {
	offset := (pageNumber - 1) * pageSize
	var count int
	err := p.pool.QueryRow(ctx, getUserQuestionsCountQuery, userId).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan row: %v", err)
	}
	if pageSize <= 0 {
		return nil, 0, errors.New("negative page size")
	}

	totalPage := int(math.Ceil(float64(count) / float64(pageSize)))

	if totalPage == 0 {
		return nil, 0, pgx.ErrNoRows
	}

	if pageNumber > totalPage {
		return nil, 0, errors.New("out of bounds page number")
	}
	var questions []*proto.Question
	rows, err := p.pool.Query(ctx, getUserQuestionsQuery, userId, offset, pageNumber)
	if err != nil {
		return nil, totalPage, fmt.Errorf("failed to execute query: %v", err)
	}
	for rows.Next() {
		question := proto.Question{Owner: username}
		err := rows.Scan(&question.Id, &question.Title, &question.State)
		if err != nil {
			return nil, totalPage, fmt.Errorf("failed to scan row: %v", err)
		}
		questions = append(questions, &question)
	}
	if err := rows.Err(); err != nil {
		return nil, totalPage, fmt.Errorf("rows iteration error: %v", err)
	}
	return questions, totalPage, nil
}

func (p *postgresqlRepository) GetQuestion(ctx context.Context, questionId int) (*proto.Question, error) {
	question := &proto.Question{}
	limitations := &proto.Limitations{}
	question.Limitations = limitations
	err := p.pool.QueryRow(ctx, getQuestionQuery, questionId).Scan(&question.Id, &question.Title,
		&question.Statement, &question.Input, &question.Output, &limitations.Memory, &limitations.Duration,
		&question.State, &question.Owner)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (p *postgresqlRepository) ChangeQuestionState(ctx context.Context, questionId int, state int32) error {
	cmdTag, err := p.pool.Exec(ctx, changeQuestionStateQuery, questionId, state)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func (p *postgresqlRepository) CreateQuestion(ctx context.Context, owner int32, question *proto.Question) (int32, error) {
	title := question.GetTitle()
	statement := question.GetStatement()
	input := question.GetInput()
	output := question.GetOutput()
	limitations := question.GetLimitations()
	if limitations == nil {
		limitations = &proto.Limitations{}
	}
	timeLimit := limitations.GetDuration()
	memoryLimit := limitations.GetMemory()
	state := proto.QuestionState_QUESTION_STATE_DRAFT

	var questionId int32
	args := []interface{}{title, statement, owner, input, output, memoryLimit, timeLimit, state}

	err := p.pool.QueryRow(ctx, createQuestionQuery, args...).Scan(&questionId)
	return questionId, err
}

func (p *postgresqlRepository) EditQuestion(ctx context.Context, question *proto.Question) error {
	var (
		setClauses []string
		args       []interface{}
		argIdx     = 1
	)

	if title := question.GetTitle(); title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, title)
		argIdx++
	}
	if statement := question.GetStatement(); statement != "" {
		setClauses = append(setClauses, fmt.Sprintf("statement = $%d", argIdx))
		args = append(args, statement)
		argIdx++
	}
	if input := question.GetInput(); input != "" {
		setClauses = append(setClauses, fmt.Sprintf("input = $%d", argIdx))
		args = append(args, input)
		argIdx++
	}
	if output := question.GetOutput(); output != "" {
		setClauses = append(setClauses, fmt.Sprintf("output = $%d", argIdx))
		args = append(args, output)
		argIdx++
	}

	if timeLimit := question.GetLimitations().GetDuration(); timeLimit != 0 {
		setClauses = append(setClauses, fmt.Sprintf("time_limit = $%d", argIdx))
		args = append(args, timeLimit)
		argIdx++
	}
	if memoryLimit := question.GetLimitations().GetMemory(); memoryLimit != 0 {
		setClauses = append(setClauses, fmt.Sprintf("memory_limit = $%d", argIdx))
		args = append(args, memoryLimit)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClause := strings.Join(setClauses, ", ")
	args = append(args, question.GetId()) // assuming ID is always provided
	query := fmt.Sprintf("UPDATE questions SET %s WHERE id = $%d", setClause, argIdx)

	cmdTag, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (p *postgresqlRepository) CreateSubmission(ctx context.Context, userId int32,
	questionId int32, code []byte) error {
	var submissionId int32
	state := proto.SubmissionState_SUBMISSION_STATE_PENDING
	err := p.pool.QueryRow(ctx, createSubmissionQuery, userId, questionId, code, state).Scan(&submissionId)
	return err
}

func (p *postgresqlRepository) UpdateSubmissionState(ctx context.Context, submissionId int32,
	state int32) (bool, error) {
	tx, err := p.pool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return false, err
	}
	sub := &submission{}
	err = tx.QueryRow(ctx, selectSubmissionForUpdateQuery, submissionId).Scan(&sub.id, &sub.state, &sub.retryCount,
		&sub.stateUpdatedAt)
	if err != nil {
		return false, err
	}
	if state == sub.state {
		return false, nil
	}

	cmdTag, err := tx.Exec(ctx, updateSubmissionStateQuery, submissionId, state, sub.retryCount)
	if err != nil {
		return false, err
	}
	if cmdTag.RowsAffected() == 0 {
		return false, pgx.ErrNoRows
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	go func() {
		if err := p.handleSubmissionJudgingTimeout(ctx, submissionId); err != nil {
			log.Printf("timeout handling failed for submission %d: %v", submissionId, err)
		}
	}()

	return true, nil
}

func (p *postgresqlRepository) handleSubmissionJudgingTimeout(ctx context.Context, submissionId int32) error {
	select {
	case <-ctx.Done():
		return nil
	case <-time.After(JudgeTimeout * time.Second):
		tx, err := p.pool.Begin(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)
		sub := &submission{}

		err = tx.QueryRow(ctx, selectSubmissionForUpdateQuery, submissionId).Scan(
			&sub.id, &sub.state, &sub.retryCount, &sub.stateUpdatedAt)
		if err != nil {
			return err
		}
		if sub.stateUpdatedAt.After(time.Now().Add(-JudgeTimeout * time.Second)) {
			return nil
		}
		if sub.state == int32(proto.SubmissionState_SUBMISSION_STATE_JUDGING) {
			sub.retryCount++
			newState := proto.SubmissionState_SUBMISSION_STATE_PENDING
			if sub.retryCount == MaxJudgeTryCount {
				newState = proto.SubmissionState_SUBMISSION_STATE_FAILED
			}
			_, err := tx.Exec(ctx, updateSubmissionStateQuery, submissionId, newState, sub.retryCount)
			if err != nil {
				return err
			}
		}
		err = tx.Commit(ctx)
		return err
	}
}

func (p *postgresqlRepository) GetSubmissionsWithState(ctx context.Context, state int32, pageNumber, pageSize int) (
	[]*proto.Submission, int, error) {
	offset := (pageNumber - 1) * pageSize
	var count int
	err := p.pool.QueryRow(ctx, getSubmissionsWithStateCountQuery, state).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan row: %v", err)
	}
	if pageSize <= 0 {
		return nil, 0, errors.New("negative page size")
	}

	totalPage := int(math.Ceil(float64(count) / float64(pageSize)))

	if totalPage == 0 {
		return nil, 0, pgx.ErrNoRows
	}

	if pageNumber > totalPage {
		return nil, 0, errors.New("out of bounds page number")
	}

	var submissions []*proto.Submission
	rows, err := p.pool.Query(ctx, getSubmissionsWithStateQuery, state, offset, pageNumber)
	if err != nil {
		return nil, totalPage, fmt.Errorf("failed to execute query: %v", err)
	}
	for rows.Next() {
		submission := proto.Submission{}
		err := rows.Scan(&submission.Id, &submission.Code, &submission.QuestionId, &submission.State)
		if err != nil {
			return nil, totalPage, fmt.Errorf("failed to scan row: %v", err)
		}
		submissions = append(submissions, &submission)
	}
	if err := rows.Err(); err != nil {
		return nil, totalPage, fmt.Errorf("rows iteration error: %v", err)
	}
	return submissions, totalPage, nil
}

func (p *postgresqlRepository) GetUserSubmissions(ctx context.Context, userId int32,
	questionId int32, filterQuestion bool, pageNumber, pageSize int) ([]*proto.Submission, int, error) {

	offset := (pageNumber - 1) * pageSize
	var count int
	var countQuery, mainQuery string
	var countQueryArgs, mainQueryArgs []interface{}
	if filterQuestion {
		countQuery = getUserQuestionSubmissionsCountQuery
		mainQuery = getUserQuestionSubmissionsQuery
		countQueryArgs = []interface{}{userId, questionId}
		mainQueryArgs = []interface{}{userId, questionId, offset, pageSize}
	} else {
		countQuery = getUserAllSubmissionsCountQuery
		mainQuery = getUserAllSubmissionsQuery
		countQueryArgs = []interface{}{userId}
		mainQueryArgs = []interface{}{userId, offset, pageSize}
	}
	err := p.pool.QueryRow(ctx, countQuery, countQueryArgs...).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan row: %v", err)
	}
	if pageSize <= 0 {
		return nil, 0, errors.New("negative page size")
	}

	totalPage := int(math.Ceil(float64(count) / float64(pageSize)))

	if totalPage == 0 {
		return nil, 0, pgx.ErrNoRows
	}

	if pageNumber > totalPage {
		return nil, 0, errors.New("out of bounds page number")
	}

	var submissions []*proto.Submission
	rows, err := p.pool.Query(ctx, mainQuery, mainQueryArgs...)
	if err != nil {
		return nil, totalPage, fmt.Errorf("failed to execute query: %v", err)
	}
	for rows.Next() {
		submission := proto.Submission{}
		err := rows.Scan(&submission.Id, &submission.Code, &submission.QuestionId, &submission.State)
		if err != nil {
			return nil, totalPage, fmt.Errorf("failed to scan row: %v", err)
		}
		submissions = append(submissions, &submission)
	}
	if err := rows.Err(); err != nil {
		return nil, totalPage, fmt.Errorf("rows iteration error: %v", err)
	}
	return submissions, totalPage, nil
}
