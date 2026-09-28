ALTER TABLE messages
    DROP CONSTRAINT messages_conversation_id_fkey,
    ADD CONSTRAINT messages_conversation_id_fkey
        FOREIGN KEY (conversation_id) REFERENCES conversations(id)
        DEFERRABLE INITIALLY IMMEDIATE;