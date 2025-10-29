-- name: CreateUser :one
INSERT INTO users(created_at, updated_at, email, password)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByLoginCredentials :one
SELECT *
FROM users
WHERE email = $1
AND password = $2;