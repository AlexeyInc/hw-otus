-- name: CreateEvent :one
INSERT INTO events (
  title, start_event, end_event, description, id_user, notification
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetEvent :one
SELECT * FROM events
WHERE id = $1;

-- name: GetDayEvents :many
SELECT * FROM events
WHERE start_event ::date = cast($1 as date)
ORDER BY start_event;

-- name: GetWeekEvents :many
SELECT * FROM events
WHERE start_event::date >= cast($1 as date) AND start_event::date < cast($1 as date)::date + interval '7 day'
ORDER BY start_event;

-- name: GetMonthEvents :many
SELECT * FROM events
WHERE start_event::date >= cast($1 as date) AND start_event::date < cast($1 as date)::date + interval '1 month'
ORDER BY start_event;

-- name: UpdateEvent :one
UPDATE events 
SET title =  $2, start_event = $3, end_event = $4, description = $5, id_user = $6, notification = $7
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;

-- name: DeleteTestEvents :exec
DELETE FROM events WHERE title like '%_test';

-- name: GetNotifyEvents :many
SELECT * FROM events
WHERE notification <= cast($1 as timestamp) 
  AND start_event > cast($1 as timestamp)
  AND notificationStatus = 0
ORDER BY id;

-- name: DeleteExpiredEvents :exec
DELETE FROM events 
WHERE now() > end_event + INTERVAL '1 year';

-- name: UpdateEventNotificationStatus :one
UPDATE events 
SET notificationStatus = $1
WHERE id = $2
RETURNING *;