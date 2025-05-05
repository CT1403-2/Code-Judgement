SELECT 'CREATE DATABASE judge_db'
WHERE NOT EXISTS (
    SELECT FROM pg_database WHERE datname = 'judge_db'
)\gexec
