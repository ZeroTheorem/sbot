-- name: CreateRow :exec
INSERT INTO sex (month, year) VALUES (?, ?);


-- name: GetAllByMonth :one
SELECT COUNT(*) FROM sex
WHERE month = ?;

-- name: GetAllByYear :one
SELECT COUNT(*) FROM sex
WHERE year = ?;

-- name: DeleteLast :exec
DELETE FROM sex
WHERE id = (SELECT MAX(id) FROM sex)
