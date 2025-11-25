CREATE TABLE board_members (
    board_internal_id BIGINT not null REFERENCES boards (internal_id) on delete CASCADE,
    user_internal_id BIGINT NOT NULL REFERENCES users (internal_id) on DELETE CASCADE,
    joined_at TIMESTAMP not null default now(),
    PRIMARY KEY (
        board_internal_id,
        user_internal_id
    )
);