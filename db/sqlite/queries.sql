-- queries.sql
-- name: UpsertCourse :one
INSERT INTO courses (id, name)
VALUES (?1, ?2)
ON CONFLICT(id) DO UPDATE SET 
    name = excluded.name
RETURNING *;

-- name: UpsertAssignment :one
INSERT INTO assignments (id, course_id, name, due_date)
VALUES (?1, ?2, ?3, ?4)
ON CONFLICT(id) DO UPDATE SET 
    course_id = excluded.course_id,
    name = excluded.name,
    due_date = excluded.due_date
RETURNING *;

-- name: ListAssignmentsByCourse :many
SELECT id, name, due_date
FROM assignments
WHERE course_id = ?1
ORDER BY due_date;

-- name: DeleteAssignmentsByCourse :exec
DELETE FROM assignments 
WHERE course_id = ?1;

-- name: DeleteAssignment :exec
DELETE FROM assignments
WHERE course_id = ?1 AND id = ?2;

-- name: DeleteCourse :exec
DELETE FROM courses
WHERE id = ?1;

-- name: UpdateAssignment :exec
UPDATE assignments
SET 
    name = ?2,
    due_date = ?3
WHERE id = ?1;

-- name: GetAssignmentCountsByCourse :many
WITH assignment_counts AS (
    SELECT course_id, COUNT(*) as total_assignments
    FROM assignments
    GROUP BY course_id
)
SELECT 
    c.id AS course_id,
    c.name AS course_name,
    COALESCE(ac.total_assignments, 0) AS assignment_count
FROM courses c
LEFT JOIN assignment_counts ac ON c.id = ac.course_id
ORDER BY assignment_count DESC;

-- name: GetAssignment :one
SELECT * FROM assignments
WHERE id = ?1 and course_id = ?2;

-- name: ListAllAssignments :many
SELECT * FROM assignments;

-- name: ListAllCourses :many
SELECT * FROM courses;

-- -- name: UpsertCourse :one
-- INSERT INTO courses (id, name)
-- VALUES ($1, $2)
-- ON CONFLICT (id) DO UPDATE 
-- SET name = EXCLUDED.name
-- RETURNING *;
--
-- -- name: UpsertAssignment :one
-- INSERT INTO assignments (id, course_id, name, due_date)
-- VALUES ($1, $2, $3, $4)
-- ON CONFLICT (id) DO UPDATE
-- SET 
--     course_id = EXCLUDED.course_id,
--     name = EXCLUDED.name,
--     due_date = EXCLUDED.due_date
-- RETURNING *;
--
-- -- name: ListAssignmentsByCourse :many
-- SELECT id, name, due_date
-- FROM assignments
-- WHERE course_id = $1
-- ORDER BY due_date;
--
-- -- name: DeleteAssignmentsByCourse :exec
-- DELETE FROM assignments 
-- WHERE course_id = $1;
--
-- -- name: UpdateAssignment :exec
-- UPDATE assignments
-- SET 
--     name = $2,
--     due_date = $3
-- WHERE id = $1;
--
-- -- name: GetAssignmentCountsByCourse :many
-- WITH assignment_counts AS (
--     SELECT course_id, COUNT(*) as total_assignments
--     FROM assignments
--     GROUP BY course_id
-- )
-- SELECT 
--     c.id AS course_id,
--     c.name AS course_name,
--     COALESCE(ac.total_assignments, 0) AS assignment_count
-- FROM courses c
-- LEFT JOIN assignment_counts ac ON c.id = ac.course_id
-- ORDER BY assignment_count DESC;
