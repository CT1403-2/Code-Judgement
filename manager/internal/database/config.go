package database

import (
	"fmt"
	"github.com/CT1403-2/Code-Judgement/manager/internal"
)

const (
	JudgeTimeout     = 120 //seconds
	MaxJudgeTryCount = 5
)

type dbConfig struct {
	scheme, username, password, host, port, name string
}

func getConnStrFromConfig(c dbConfig) string {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s", c.scheme, c.username, c.password, c.host, c.port, c.name)
	return connStr
}

func getMainDbConfig() dbConfig {
	return dbConfig{
		scheme:   internal.GetEnv("DB_SCHEME", "postgres"),
		username: internal.GetEnv("DB_USERNAME", "username"),
		password: internal.GetEnv("DB_PASSWORD", ""),
		host:     internal.GetEnv("DB_HOST", "localhost"),
		port:     internal.GetEnv("DB_PORT", "5432"),
		name:     internal.GetEnv("DB_NAME", "judge_db"),
	}
}

func getTestDbConfig() dbConfig {
	return dbConfig{
		scheme:   internal.GetEnv("TEST_DB_SCHEME", "postgres"),
		username: internal.GetEnv("TEST_DB_USERNAME", "username"),
		password: internal.GetEnv("TEST_DB_PASSWORD", "password"),
		host:     internal.GetEnv("TEST_DB_HOST", "localhost"),
		port:     internal.GetEnv("TEST_DB_PORT", "5432"),
		name:     internal.GetEnv("TEST_DB_NAME", "judge_db_test"),
	}
}
