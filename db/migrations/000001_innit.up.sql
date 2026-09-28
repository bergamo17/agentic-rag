CREATE TABLE "conversations" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "title" text NOT NULL DEFAULT 'Percakapan Baru',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "messages" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "conversation_id" uuid NOT NULL,
  "role" text NOT NULL,
  "content" text NOT NULL,
  "widgets" jsonb,
  "documents" jsonb,
  "pages" jsonb,
  "is_partial" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "messages" ("conversation_id", "created_at");

ALTER TABLE "messages" ADD FOREIGN KEY ("conversation_id") REFERENCES "conversations" ("id") DEFERRABLE INITIALLY IMMEDIATE;
