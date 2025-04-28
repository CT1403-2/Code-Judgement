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
	GetUsernames(ctx context.Context, pageNumber, pageSize int32) ([]string, int32, error)
}

type postgresqlRepository struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewRepository() (Repository, error) {
	scheme := internal.GetEnv("DB_SCHEME", "postgres")
	username := internal.GetEnv("DB_USERNAME", "username")
	password := internal.GetEnv("DB_PASSWORD", "")
	host := internal.GetEnv("DB_HOST", "localhost")
	port := internal.GetEnv("DB_PORT", "5432")
	name := internal.GetEnv("DB_NAME", "judge_db")
	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable not set")
	}
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s", scheme, username, password, host, port, name)
	pool, err := pgxpool.New(context.Background(), connStr)
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
	hashedPassword, err := internal.HashPassword(password)
	if err != nil {
		return 0, 0, err
	}
	u, err := p.getUser(ctx, username)
	if err != nil {
		return 0, 0, err
	}
	if u.password != hashedPassword {
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

func (p *postgresqlRepository) GetUsernames(ctx context.Context, pageNumber, pageSize int32) ([]string, int32, error) {
	offset := (pageNumber - 1) * pageSize

	var count int32
	err := p.pool.QueryRow(ctx, getUsernamesCountQuery).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan row: %v", err)
	}

	if pageSize <= 0 {
		return nil, 0, errors.New("negative page size")
	}

	totalPage := int32(math.Ceil(float64(count) / float64(pageSize)))

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
