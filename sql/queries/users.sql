-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, password = $3
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;

-- name: DowngradeUser :exec
UPDATE users
SET is_chirpy_red = false
WHERE id = $1;
