-- name: CreateConversation :one
INSERT INTO conversations (
    title
) VALUES (
    $1
) RETURNING *;

-- name: GetConversation :one
SELECT * FROM conversations
WHERE id = $1 LIMIT 1;

-- name: UpdateConversationTime :one
UPDATE conversations
SET updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ListConversations :many
SELECT c.*
FROM conversations c
ORDER BY updated_at DESC;