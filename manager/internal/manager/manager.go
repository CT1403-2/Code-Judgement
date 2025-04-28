package manager

import (
	"context"
	"errors"
	"github.com/CT1403-2/Code-Judgement/proto"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"manger/internal"
	"manger/internal/database"
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
	password := authRequest.GetPassword()
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
	userId, err := GetUserIdFromToken(ctx)
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
		userId, err := GetUserIdFromToken(ctx)
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

func (m *Manager) GetProfiles(ctx context.Context, req *proto.Empty) (*proto.GetProfilesResponse, error) {
	var pageNumber, pageSize int32
	usernames, totalPage, err := m.db.GetUsernames(ctx, pageNumber, pageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	return &proto.GetProfilesResponse{Usernames: usernames, TotalPageSize: int64(totalPage)}, nil
}
func (m *Manager) GetStatsRequest(ctx context.Context, req *proto.ID) (*proto.GetStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatsRequest not implemented")
}
func (m *Manager) GetQuestions(ctx context.Context, req *proto.GetQuestionsRequest) (*proto.GetQuestionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuestions not implemented")
}
func (m *Manager) GetQuestion(ctx context.Context, req *proto.ID) (*proto.GetQuestionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuestion not implemented")
}
func (m *Manager) Submit(ctx context.Context, req *proto.SubmitRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Submit not implemented")
}
func (m *Manager) GetSubmissions(ctx context.Context, req *proto.GetSubmissionsRequest) (*proto.GetSubmissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubmissions not implemented")
}
func (m *Manager) CreateQuestion(ctx context.Context, req *proto.Question) (*proto.ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateQuestion not implemented")
}
func (m *Manager) EditQuestion(ctx context.Context, req *proto.Question) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditQuestion not implemented")
}
func (m *Manager) ChangeQuestionState(ctx context.Context, req *proto.ChangeQuestionStateRequest) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeQuestionState not implemented")
}
func (m *Manager) UpdateSubmission(ctx context.Context, req *proto.Submission) (*proto.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSubmission not implemented")
}

func GetUserIdFromToken(ctx context.Context) (int32, error) {
	token, err := internal.ExtractJWTFromContext(ctx)
	if err != nil {
		return 0, err
	}
	userId, _, err := internal.ValidateJWT(token)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
