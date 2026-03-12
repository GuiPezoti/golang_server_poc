-- name: GetChirp :one
SELECT * FROM chirps where id = $1;