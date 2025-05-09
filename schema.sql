CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);


CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL,
    rec_id BIGINT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_sender
        FOREIGN KEY(sender_id) 
        REFERENCES users(id)
        ON DELETE CASCADE,
        
    CONSTRAINT fk_receiver
        FOREIGN KEY(rec_id) 
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Indexes for faster message retrieval
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_receiver ON messages(rec_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
