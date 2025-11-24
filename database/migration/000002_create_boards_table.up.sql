CREATE TABLE boards (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID not NULL default gen_random_uuid (),
    title VARCHAR(255) not NULL,
    description TEXT,
    owner_internal_id BIGINT not NULL REFERENCES users (internal_id),
    owner_public_id UUID not NULL,
    created_at TIMESTAMP not null default NOW(),
    CONSTRAINT boards_public_id_unique UNIQUE (public_id),
    CONSTRAINT fk_boards_owner_public_id Foreign Key (owner_public_id) REFERENCES users (public_id) ON DELETE CASCADE
);