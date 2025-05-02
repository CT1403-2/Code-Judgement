package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"manger/internal"
	"manger/internal/database"
	"strconv"
)

type Manager struct {
	db database.Repository
	proto.UnimplementedManagerServer
}

func NewManager() (*Manager, error) {
	db, err := database.NewRepository()
	return &Manager{db: db}, err
}

func (m *Manager) Start() error {
	panic("implement me")
}

func (m *Manager) Stop() error {
	panic("implement me")
}

func (m *Manager) Register(ctx context.Context, authRequest *proto.AuthenticationRequest) (*proto.AuthenticationResponse, error) {
	username := authRequest.GetUsername()
	if len(username) < usernameMinLength {
		return nil, status.Error(codes.InvalidArgument, "username too short")
	}
	password := authRequest.GetPassword()
	if len(password) < passwordMinLength {
		return nil, status.Error(codes.InvalidArgument, "password too short")
	}
	userId, err := m.db.CreateMember(ctx, username, password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.AlreadyExists, "Username already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	jwtToken, err := internal.GenerateJWT(userId, proto.Role_ROLE_MEMBER.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate jwt token")
	}
	return &proto.AuthenticationResponse{Value: jwtToken, Role: proto.Role_ROLE_MEMBER}, status.New(codes.OK, "Registered Successfully").Err()
}

func (m *Manager) Login(ctx context.Context, authRequest *proto.AuthenticationRequest) (*proto.AuthenticationResponse, error) {
	username := authRequest.GetUsername()
	password := authRequest.GetPassword()
	userId, roleType, err := m.db.Authenticate(ctx, username, password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	jwtToken, err := internal.GenerateJWT(userId, proto.Role_name[roleType])
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate jwt token")
	}
	return &proto.AuthenticationResponse{Value: jwtToken, Role: proto.Role(roleType)}, status.New(codes.OK, "Logged In successfully").Err()
}

func (m *Manager) ChangeRole(ctx context.Context, req *proto.ChangeRoleRequest) (*proto.Empty, error) {
	userId, _, err := authenticate(ctx)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Unauthenticated, err.Error())
	}
	_, requesterRole, err := m.db.GetUserRole(ctx, userId)
	if err != nil || requesterRole == proto.Role_ROLE_UNKNOWN {
		return &proto.Empty{}, status.Error(codes.Internal, "")
	}
	targetUsername := req.GetUsername()
	targetUserId, oldTargetRole, err := m.db.GetUserRoleByUsername(ctx, targetUsername)
	if err != nil || oldTargetRole == proto.Role_ROLE_UNKNOWN {
		return nil, status.Error(codes.Internal, "")
	}
	newTargetRole := req.GetRole()
	if newTargetRole == proto.Role_ROLE_UNKNOWN {
		return &proto.Empty{}, status.Error(codes.InvalidArgument, "target role is unknown")
	}

	if requesterRole == proto.Role_ROLE_MEMBER {
		return &proto.Empty{}, status.Error(codes.PermissionDenied, "change role request is aborted")
	}
	if oldTargetRole == proto.Role_ROLE_ADMIN && requesterRole == proto.Role_ROLE_ADMIN {
		return &proto.Empty{}, status.Error(codes.PermissionDenied, "change role request is aborted")
	}
	if oldTargetRole == proto.Role_ROLE_SUPERUSER {
		return &proto.Empty{}, status.Error(codes.PermissionDenied, "change role request is aborted")
	}

	if newTargetRole == proto.Role_ROLE_SUPERUSER {
		return &proto.Empty{}, status.Error(codes.PermissionDenied, "change role request is aborted")
	}

	err = m.db.UpdateUserRole(ctx, targetUserId, newTargetRole)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, status.Error(codes.OK, "Changed Successfully")
}

func (m *Manager) GetProfile(ctx context.Context, req *proto.ID) (*proto.GetProfileResponse, error) {
	username := req.GetValue()
	if username == "" { //return client username and role
		userId, _, err := authenticate(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		username, role, err := m.db.GetUserRole(ctx, userId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &proto.GetProfileResponse{Username: username, Role: role}, nil

	} else { // return requested username's role
		_, role, err := m.db.GetUserRoleByUsername(ctx, username)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, status.Error(codes.NotFound, "user not found")
			}
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &proto.GetProfileResponse{Username: username, Role: role}, nil
	}
}

func (m *Manager) GetProfiles(ctx context.Context, req *proto.GetProfilesRequest) (*proto.GetProfilesResponse, error) {
	var pageNumber int
	filtersMap := getFiltersMap(req.GetFilters())
	var str string
	str, ok := filtersMap[pageNumberName]
	pageNumber, err := strconv.Atoi(str)

	if !ok || err != nil {
		pageNumber = defaultPageNumber
	}

	usernames, totalPage, err := m.db.GetUsernames(ctx, pageNumber, defaultPageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.GetProfilesResponse{Usernames: usernames, TotalPageSize: int64(totalPage)}, nil
}

func (m *Manager) GetStatsRequest(ctx context.Context, req *proto.ID) (*proto.GetStatsResponse, error) {
	username := req.GetValue()
	userId, _, err := m.db.GetUserRoleByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &proto.GetStatsResponse{}, status.Error(codes.NotFound, "user not found")
		}
		return &proto.GetStatsResponse{}, status.Error(codes.Internal, err.Error())
	}
	triedCount, successCount, err := m.db.GetUserStats(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &proto.GetStatsResponse{}, status.Error(codes.NotFound, "user not found")
		}
		return &proto.GetStatsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.GetStatsResponse{TriedQuestions: triedCount, SolvedQuestions: successCount}, nil
}

func (m *Manager) GetQuestions(ctx context.Context, req *proto.GetQuestionsRequest) (*proto.GetQuestionsResponse, error) {
	var pageNumber, totalPage int
	var str string
	var isOwner, published bool
	var questions []*proto.Question

	userId, _, err := authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	username, role, err := m.db.GetUserRole(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	filtersMap := getFiltersMap(req.GetFilters())
	str, ok := filtersMap[pageNumberName]
	pageNumber, err = strconv.Atoi(str)
	if !ok || err != nil {
		pageNumber = defaultPageNumber
	}
	str, ok = filtersMap[questionOwnerFilter]
	if !ok || str == "false" {
		isOwner = false
	} else {
		isOwner = true
	}
	published = isAdmin(role)

	if !isOwner {
		questions, totalPage, err = m.db.GetQuestions(ctx, published, pageNumber, defaultPageSize)
	} else {
		questions, totalPage, err = m.db.GetUserQuestions(ctx, userId, username, pageNumber, defaultPageSize)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.GetQuestionsResponse{Questions: questions, TotalPageSize: int64(totalPage)}, nil
}

func (m *Manager) GetQuestion(ctx context.Context, req *proto.ID) (*proto.GetQuestionResponse, error) {
	userId, isJudge, err := authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	var username string
	var role proto.Role
	if !isJudge {
		username, role, err = m.db.GetUserRole(ctx, userId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	questionId, err := strconv.Atoi(req.Value)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	question, err := m.db.GetQuestion(ctx, questionId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !isJudge && !isAdmin(role) && username != question.GetOwner() {
		question.Input = nil
		question.Output = nil
	}
	return &proto.GetQuestionResponse{Question: question}, status.Error(codes.OK, "")
}

func (m *Manager) Submit(ctx context.Context, req *proto.SubmitRequest) (*proto.Empty, error) {
	userId, _, err := authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	submission := req.GetSubmission()
	questionIdStr := submission.GetQuestionId()
	questionId, err := strconv.Atoi(questionIdStr)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}
	question, err := m.db.GetQuestion(ctx, questionId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &proto.Empty{}, status.Error(codes.NotFound, "question not found")
		}
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}
	if question.GetState() != proto.QuestionState_QUESTION_STATE_PUBLISHED {
		return &proto.Empty{}, status.Error(codes.FailedPrecondition, "question is not published")
	}

	code := submission.GetCode()

	err = m.db.CreateSubmission(ctx, userId, int32(questionId), code)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, status.Error(codes.OK, "")
}

func (m *Manager) GetSubmissions(ctx context.Context, req *proto.GetSubmissionsRequest) (*proto.GetSubmissionsResponse, error) {
	userId, isJudge, err := authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	filtersMap := getFiltersMap(req.GetFilters())

	str, ok := filtersMap[pageNumberName]
	pageNumber, err := strconv.Atoi(str)

	if !ok || err != nil {
		pageNumber = defaultPageNumber
	}
	if isJudge {
		stateName, ok := filtersMap[stateFilter]
		if ok {
			state := proto.SubmissionState_value[stateName]
			submissions, totalPage, err := m.db.GetSubmissionsWithState(ctx, state, pageNumber, defaultPageSize)
			return &proto.GetSubmissionsResponse{Submissions: submissions, TotalPageSize: int64(totalPage)}, err
		} else {
			return nil, status.Error(codes.Unimplemented, "judge server gets submission without state filter not implemented")
		}
	} else {
		username, usernameOk := filtersMap[usernameFilter]
		questionIdStr, questionIdOk := filtersMap[questionIdFilter]

		if questionId, err := strconv.Atoi(questionIdStr); questionIdOk && err != nil {
			// to see user submissions in a given question
			submissions, totalPage, err := m.db.GetUserSubmissions(ctx, userId, int32(questionId), true,
				pageNumber, defaultPageSize)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return &proto.GetSubmissionsResponse{Submissions: submissions, TotalPageSize: int64(totalPage)}, err
		}
		if usernameOk {
			//to see given username's all submissions in all questions
			userId, _, err := m.db.GetUserRoleByUsername(ctx, username)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, status.Error(codes.NotFound, "user not found")
				}
				return nil, status.Error(codes.Internal, err.Error())
			}
			submissions, totalPage, err := m.db.GetUserSubmissions(ctx, userId, 0, false,
				pageNumber, defaultPageSize)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return &proto.GetSubmissionsResponse{Submissions: submissions, TotalPageSize: int64(totalPage)}, nil
		}
		return nil, status.Error(codes.Unimplemented,
			"case client get questions without question or username filter not implemented ")
	}
}

func (m *Manager) CreateQuestion(ctx context.Context, question *proto.Question) (*proto.ID, error) {
	userId, _, err := authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	questionId, err := m.db.CreateQuestion(ctx, userId, question)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	questionIdStr := fmt.Sprintf("%d", questionId)
	return &proto.ID{Value: questionIdStr}, status.Error(codes.OK, "")
}

func (m *Manager) EditQuestion(ctx context.Context, question *proto.Question) (*proto.Empty, error) {
	userId, _, err := authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	username, _, err := m.db.GetUserRole(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if question.GetOwner() != username {
		return nil, status.Error(codes.PermissionDenied, "you do not have access to this question")
	}

	err = m.db.EditQuestion(ctx, question)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, status.Error(codes.OK, "question edited successfully")
}

func (m *Manager) ChangeQuestionState(ctx context.Context, req *proto.ChangeQuestionStateRequest) (*proto.Empty, error) {
	userId, _, err := authenticate(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	_, role, err := m.db.GetUserRole(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isAdmin(role) {
		return nil, status.Error(codes.PermissionDenied, "you are not an admin")
	}
	questionIdStr := req.QuestionId
	questionId, err := strconv.Atoi(questionIdStr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	newState := req.GetState()
	if newState == proto.QuestionState_QUESTION_STATE_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid state: %s", newState))
	}
	err = m.db.ChangeQuestionState(ctx, questionId, int32(newState))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "question not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, status.Error(codes.OK, "question state changed successfully")
}

func (m *Manager) UpdateSubmission(ctx context.Context, submission *proto.Submission) (*proto.UpdateSubmissionResponse, error) {
	_, isJudge, err := authenticate(ctx)
	if err != nil || !isJudge {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	newState := submission.GetState()
	if newState == proto.SubmissionState_SUBMISSION_STATE_UNKNOWN {
		return nil, status.Errorf(codes.InvalidArgument, "invalid state: %s", newState)
	}
	submissionIdStr := submission.GetId()
	submissionId, err := strconv.Atoi(submissionIdStr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	updated, err := m.db.UpdateSubmissionState(ctx, int32(submissionId), int32(newState))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "submission not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.UpdateSubmissionResponse{Updated: updated}, status.Errorf(codes.Unimplemented, "method UpdateSubmission not implemented")
}

func authenticate(ctx context.Context) (userId int32, isJudge bool, err error) {
	tokenType, token, err := internal.ExtractTokenFromContext(ctx)
	if err != nil {
		return 0, false, err
	}
	if tokenType == "Bearer" {
		userId, _, err := internal.ValidateJWT(token)
		if err != nil {
			return 0, false, err
		}
		return userId, false, nil

	} else if tokenType == "Token" {
		isJudge := internal.IsJudgeServer(token)
		return 0, isJudge, nil
	}
	return 0, false, nil
}

func getFiltersMap(filters []*proto.Filter) map[string]string {
	m := make(map[string]string)
	for _, filter := range filters {
		m[filter.Field] = filter.Value
	}
	return m
}

func isAdmin(role proto.Role) bool {
	return role == proto.Role_ROLE_ADMIN || role == proto.Role_ROLE_SUPERUSER
}
