-- name: InsertMessage :one
INSERT INTO messages(
    conversation_id, role, content, widgets, documents, pages, is_partial
) VALUES(
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListMessageByConversationID :many
SELECT * 
FROM messages
WHERE conversation_id = $1
ORDER BY created_at ASC;

-- name: ListRecentMessages :many
SELECT *
FROM messages
WHERE conversation_id = $1
ORDER BY created_at DESC
LIMIT $2;