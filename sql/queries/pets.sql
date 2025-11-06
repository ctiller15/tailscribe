-- name: CreatePet :one
INSERT INTO pet(name, imageUrl, created_at, updated_at)
VALUES (
    $1,
    $2,
    NOW(),
    NOW()
)
RETURNING *;

-- name: CreateUserPet :one
INSERT INTO UserPets(userId, petId, permissions_level, active, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
)
RETURNING *;

-- name: DeletePets :exec
DELETE FROM pet;

-- name: DeleteUserPets :exec
DELETE FROM UserPets;