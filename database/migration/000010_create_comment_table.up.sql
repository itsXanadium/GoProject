CREATE TABLE comments (
    internal_id bigserial PRIMARY KEY,
    public_id UUID not null default gen_random_uuid (),
    card_internal_id bigint not null REFERENCES cards (internal_id) ON DELETE CASCADE,
    user_internal_id bigint not null REFERENCES users (internal_id) ON DELETE CASCADE,
    message text not null,
    created_at TIMESTAMP not null default now(),
    constraint comments_public_id_unique UNIQUE (public_id)
);