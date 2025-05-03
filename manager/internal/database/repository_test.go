package database

import (
	"fmt"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"log"
	"manger/internal"
	"os"
	"strconv"
	"testing"
)

func newRepository(truncate bool) (*postgresqlRepository, error) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := getTestDbConfig()
	if config.password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable not set")
	}
	connStr := getConnStrFromConfig(config)
	p, err := newRepositoryWithDsn(connStr)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if truncate {
		err = p.truncateTables()
	}
	return p, err
}

func (p *postgresqlRepository) truncateTables() error {
	_, err := p.pool.Exec(p.ctx, truncateAllTablesQuery)
	return err
}

func TestUserAndRole(t *testing.T) {
	repo, err := newRepository(true)
	require.NoError(t, err)
	err = repo.SetUpRoles()
	require.NoError(t, err)

	username := "username"
	password := "password"
	require.NoError(t, err)

	role := proto.Role_ROLE_MEMBER

	adminUsername := "admin"
	adminPassword := "admin"

	duplicateUsername := "username"
	wrongUsername := "wrong_username"
	wrongPassword := "wrong_password"

	t.Run("create user success", func(t *testing.T) {

		userId, err := repo.createUser(repo.ctx, username, password, role)
		require.NoError(t, err)
		require.NotZero(t, userId)

	})

	t.Run("get user and role success", func(t *testing.T) {
		userId, role, err := repo.GetUserRoleByUsername(repo.ctx, username)
		require.NoError(t, err)
		require.Equal(t, proto.Role_ROLE_MEMBER, role)

		userName, role, err := repo.GetUserRole(repo.ctx, userId)
		require.NoError(t, err)
		require.Equal(t, username, userName)
		require.Equal(t, proto.Role_ROLE_MEMBER, role)

		user, err := repo.getUser(repo.ctx, username)
		require.NoError(t, err)
		require.NotZero(t, user.id)
		require.Equal(t, username, user.username)
		passwordIsValid := internal.CheckPasswordHash(password, user.password)
		require.True(t, passwordIsValid)
		require.NotZero(t, user.roleId)

		r, err := repo.getRole(repo.ctx, user.roleId)
		require.NoError(t, err)
		require.NotZero(t, r.id)
		require.Equal(t, int32(proto.Role_ROLE_MEMBER), r.roleType)

	})

	t.Run("get wrong username fail", func(t *testing.T) {
		userId, role, err := repo.GetUserRoleByUsername(repo.ctx, wrongUsername)
		require.Error(t, err)
		require.Equal(t, pgx.ErrNoRows, err)
		require.Zero(t, userId)
		require.Equal(t, proto.Role_ROLE_UNKNOWN, role)
	})

	t.Run("authenticate success", func(t *testing.T) {
		userId, roleType, err := repo.Authenticate(repo.ctx, username, password)
		require.NoError(t, err)
		require.Equal(t, int32(proto.Role_ROLE_MEMBER), roleType)
		require.NotZero(t, userId)
	})

	t.Run("authenticate fail: wrong username", func(t *testing.T) {
		userId, roleType, err := repo.Authenticate(repo.ctx, wrongUsername, password)
		require.Error(t, err)
		require.Equal(t, pgx.ErrNoRows, err)
		require.Zero(t, userId)
		require.Equal(t, int32(proto.Role_ROLE_UNKNOWN), roleType)
	})

	t.Run("authenticate fail: wrong password", func(t *testing.T) {
		userId, roleType, err := repo.Authenticate(repo.ctx, username, wrongPassword)
		require.Error(t, err)
		require.Equal(t, pgx.ErrNoRows, err)
		require.Zero(t, userId)
		require.Equal(t, int32(proto.Role_ROLE_UNKNOWN), roleType)
	})

	t.Run("create admin success", func(t *testing.T) {
		err := os.Setenv("ADMIN_USERNAME", adminUsername)
		require.NoError(t, err)
		err = os.Setenv("ADMIN_PASSWORD", adminPassword)
		require.NoError(t, err)
		userId, err := repo.CreateSuperUserIfNotExists()
		require.NoError(t, err)
		require.NotZero(t, userId)
		userName, role, err := repo.GetUserRole(repo.ctx, userId)
		require.NoError(t, err)
		require.Equal(t, adminUsername, userName)
		require.Equal(t, proto.Role_ROLE_SUPERUSER, role)
	})

	t.Run("create duplicate username fail", func(t *testing.T) {
		userId, err := repo.CreateMember(repo.ctx, duplicateUsername, password)
		require.Error(t, err)
		require.Equal(t, err, pgx.ErrNoRows)
		require.Zero(t, userId)
	})

	t.Run("change user role success", func(t *testing.T) {
		userId, role, err := repo.GetUserRoleByUsername(repo.ctx, username)
		require.NoError(t, err)
		require.NotZero(t, userId)
		require.Equal(t, proto.Role_ROLE_MEMBER, role)
		err = repo.UpdateUserRole(repo.ctx, userId, proto.Role_ROLE_ADMIN)
		require.NoError(t, err)

		_, role, err = repo.GetUserRoleByUsername(repo.ctx, username)
		require.NoError(t, err)
		require.Equal(t, proto.Role_ROLE_ADMIN, role)
	})

	t.Run("get all usernames", func(t *testing.T) {
		usernames, totalPage, err := repo.GetUsernames(repo.ctx, 1, 10)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Contains(t, usernames, username)
		require.Contains(t, usernames, adminUsername)
		require.Len(t, usernames, 2)
	})
}

func TestQuestion(t *testing.T) {
	repo, err := newRepository(true)
	require.NoError(t, err)
	err = repo.SetUp()
	username := "username"
	password := "password"
	userId, err := repo.CreateMember(repo.ctx, username, password)
	require.NoError(t, err)
	require.NotZero(t, userId)

	username2 := "username2"
	userId2, err := repo.CreateMember(repo.ctx, username2, password)
	require.NoError(t, err)
	require.NotZero(t, userId2)

	title := "Test Question"
	statement := "this is test question statement"
	input := "this is test question input"
	output := "this is test question output"
	memoryLimit := int64(1024) //MB
	timeLimit := int64(1)      //seconds
	question := &proto.Question{
		Title:     title,
		Statement: statement,
		Input:     &input,
		Output:    &output,
		Limitations: &proto.Limitations{
			Memory:   memoryLimit,
			Duration: timeLimit,
		},
	}
	title2 := "Test Question2"
	question2 := &proto.Question{
		Title: title2,
	}
	pageNumber := 1
	pageSize := 10

	t.Run("create and get question success", func(t *testing.T) {

		id, err := repo.CreateQuestion(repo.ctx, userId, question)
		require.NoError(t, err)
		require.NotZero(t, id)

		q, err := repo.GetQuestion(repo.ctx, int(id))
		require.NoError(t, err)
		require.Equal(t, question.Title, q.Title)
		require.Equal(t, question.Statement, q.Statement)
		require.Equal(t, question.Input, q.Input)
		require.Equal(t, question.Output, q.Output)
		require.Equal(t, question.Limitations.Duration, q.Limitations.Duration)
		require.Equal(t, question.Limitations.Memory, q.Limitations.Memory)
		require.Equal(t, proto.QuestionState_QUESTION_STATE_DRAFT, q.State)
		require.Equal(t, username, q.Owner)
	})

	questionId2, err := repo.CreateQuestion(repo.ctx, userId2, question2)
	require.NoError(t, err)
	require.NotZero(t, questionId2)

	t.Run("change question state success", func(t *testing.T) {
		err := repo.ChangeQuestionState(repo.ctx, int(questionId2), int32(proto.QuestionState_QUESTION_STATE_PUBLISHED))
		require.NoError(t, err)
		q, err := repo.GetQuestion(repo.ctx, int(questionId2))
		require.NoError(t, err)
		require.Equal(t, proto.QuestionState_QUESTION_STATE_PUBLISHED, q.State)
	})

	t.Run("get all questions success", func(t *testing.T) {
		questions, totalPage, err := repo.GetQuestions(repo.ctx, false, pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Len(t, questions, 2)
		q0Id, err := strconv.Atoi(*questions[0].Id)
		require.NoError(t, err)
		q1Id, err := strconv.Atoi(*questions[1].Id)
		require.NoError(t, err)
		require.Greater(t, q1Id, q0Id)

		q0 := questions[0]
		require.Equal(t, question.Title, q0.Title)
		require.Equal(t, proto.QuestionState_QUESTION_STATE_DRAFT, q0.State)
		require.Equal(t, username, q0.Owner)

		q1 := questions[1]
		require.Equal(t, question2.Title, q1.Title)
		require.Equal(t, proto.QuestionState_QUESTION_STATE_PUBLISHED, q1.State)
		require.Equal(t, username2, q1.Owner)

	})

	t.Run("get published questions success", func(t *testing.T) {
		questions, totalPage, err := repo.GetQuestions(repo.ctx, true, pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Len(t, questions, 1)
		q := questions[0]
		require.Equal(t, question2.Title, q.Title)
		require.Equal(t, username2, q.Owner)
	})

	t.Run("get user questions success", func(t *testing.T) {
		questions, totalPage, err := repo.GetUserQuestions(repo.ctx, userId, username, pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Len(t, questions, 1)
		q := questions[0]
		require.Equal(t, question.Title, q.Title)
		require.Equal(t, username, q.Owner)
		require.Equal(t, proto.QuestionState_QUESTION_STATE_DRAFT, q.State)
	})

	t.Run("test edit question success", func(t *testing.T) {
		newStatement := "this is new statement"
		oldTitle := question2.Title
		newTitle := ""

		qIdStr := fmt.Sprintf("%v", questionId2)
		q := &proto.Question{Id: &qIdStr, Statement: newStatement, Title: newTitle}
		err := repo.EditQuestion(repo.ctx, q)
		require.NoError(t, err)
		q, err = repo.GetQuestion(repo.ctx, int(questionId2))
		require.NoError(t, err)
		require.Equal(t, newStatement, q.Statement)
		require.Equal(t, oldTitle, q.Title)
	})

	t.Run("test questions pagination", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			q := &proto.Question{
				Title: fmt.Sprintf("Question %d", i),
			}
			qId, err := repo.CreateQuestion(repo.ctx, userId, q)
			require.NoError(t, err)
			require.NotZero(t, qId)
		}

		questions1, totalPageSize1, err1 := repo.GetQuestions(repo.ctx, false, 1, pageSize)
		questions2, totalPageSize2, err2 := repo.GetQuestions(repo.ctx, false, 2, pageSize)
		questions3, totalPageSize3, err3 := repo.GetQuestions(repo.ctx, false, 3, pageSize)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		require.Equal(t, 3, totalPageSize1)
		require.Equal(t, 3, totalPageSize2)
		require.Equal(t, 3, totalPageSize3)

		require.Len(t, questions1, 10)
		require.Len(t, questions2, 10)
		require.Len(t, questions3, 2)

		m := make(map[string]struct{})
		for _, q := range append(questions1, append(questions2, questions3...)...) {
			if _, exists := m[*q.Id]; exists {
				require.Failf(t, "duplicate question id %v", *q.Id)
			}
			m[*q.Id] = struct{}{}
		}
	})
}

func TestSubmission(t *testing.T) {
	repo, err := newRepository(true)
	require.NoError(t, err)
	err = repo.SetUp()
	username := "username"
	password := "password"
	userId, err := repo.CreateMember(repo.ctx, username, password)
	require.NoError(t, err)
	require.NotZero(t, userId)

	title := "Test Question"
	statement := "this is test question statement"
	input := "this is test question input"
	output := "this is test question output"
	memoryLimit := int64(1024) //MB
	timeLimit := int64(1)      //seconds
	question := &proto.Question{
		Title:     title,
		Statement: statement,
		Input:     &input,
		Output:    &output,
		Limitations: &proto.Limitations{
			Memory:   memoryLimit,
			Duration: timeLimit,
		},
	}
	questionId, err := repo.CreateQuestion(repo.ctx, userId, question)
	require.NoError(t, err)
	require.NotZero(t, questionId)
	questionIdStr := fmt.Sprintf("%v", questionId)

	pageNumber := 1
	pageSize := 10
	//code := make([]byte, 0)
	var code []byte
	t.Run("test submit fail, question not found", func(t *testing.T) {
		wrongQuestionId := int32(-1)
		err := repo.CreateSubmission(repo.ctx, userId, wrongQuestionId, code)
		require.Error(t, err)
	})

	t.Run("test submit fail, user not found", func(t *testing.T) {
		wrongUserId := int32(-1)
		err := repo.CreateSubmission(repo.ctx, wrongUserId, questionId, code)
		require.Error(t, err)
	})

	t.Run("test submit success", func(t *testing.T) {
		err := repo.CreateSubmission(repo.ctx, userId, questionId, code)
		require.NoError(t, err)
	})

	t.Run("test get user submissions fail, question not found", func(t *testing.T) {
		wrongQuestionId := int32(-1)
		submissions, totalPage, err := repo.GetUserSubmissions(repo.ctx, userId, wrongQuestionId, true,
			pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 0, totalPage)
		require.Len(t, submissions, 0)
	})

	t.Run("test get user submissions success", func(t *testing.T) {
		submissions, totalPage, err := repo.GetUserSubmissions(repo.ctx, userId, questionId, true,
			pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Len(t, submissions, 1)
		s := submissions[0]
		require.Equal(t, questionIdStr, s.QuestionId)
		require.Equal(t, proto.SubmissionState_SUBMISSION_STATE_PENDING, *s.State)
	})

	question2 := &proto.Question{}
	qId2, err := repo.CreateQuestion(repo.ctx, userId, question2)
	require.NoError(t, err)
	require.NotZero(t, qId2)

	err = repo.CreateSubmission(repo.ctx, userId, qId2, code)
	require.NoError(t, err)

	t.Run("test get user all submissions success", func(t *testing.T) {
		submissions, totalPage, err := repo.GetUserSubmissions(repo.ctx, userId, 0, false,
			pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Len(t, submissions, 2)
	})

	t.Run("test change submission state success", func(t *testing.T) {
		submissions, totalPage, err := repo.GetUserSubmissions(repo.ctx, userId, 0, false,
			pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Len(t, submissions, 2)

		s := submissions[0]
		require.Equal(t, proto.SubmissionState_SUBMISSION_STATE_PENDING, *s.State)
		sId, err := strconv.Atoi(*s.Id)
		require.NoError(t, err)
		updated, err := repo.UpdateSubmissionState(repo.ctx, int32(sId), int32(proto.SubmissionState_SUBMISSION_STATE_OK))
		require.NoError(t, err)
		require.True(t, updated)
	})

	t.Run("test update submission state fail, submission already in target state", func(t *testing.T) {
		submissions, totalPage, err := repo.GetUserSubmissions(repo.ctx, userId, 0, false,
			pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Len(t, submissions, 2)

		s := submissions[0]
		require.Equal(t, proto.SubmissionState_SUBMISSION_STATE_OK, *s.State)
		sId, err := strconv.Atoi(*s.Id)
		require.NoError(t, err)
		updated, err := repo.UpdateSubmissionState(repo.ctx, int32(sId), int32(proto.SubmissionState_SUBMISSION_STATE_OK))
		require.NoError(t, err)
		require.False(t, updated)
	})

	t.Run("test get all pending submissions success", func(t *testing.T) {
		pending := int32(proto.SubmissionState_SUBMISSION_STATE_PENDING)
		submissions, totalPage, err := repo.GetSubmissionsWithState(repo.ctx, pending, pageNumber, pageSize)
		require.NoError(t, err)
		require.Equal(t, 1, totalPage)
		require.Len(t, submissions, 1)
		for _, s := range submissions {
			require.Equal(t, proto.SubmissionState_SUBMISSION_STATE_PENDING, *s.State)
		}
	})

}
