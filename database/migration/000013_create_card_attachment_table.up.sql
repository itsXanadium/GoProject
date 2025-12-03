CREATE TABLE card_attachment (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    file VARCHAR(255) NOT NULL,
    user_internal_id BIGINT NOT NULL REFERENCES users (internal_id) on delete CASCADE,
    card_internal_id BIGINT NOT NULL REFERENCES cards (internal_id) on delete CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    constraint card_attatchment_public_id_unique UNIQUE (public_id)
);