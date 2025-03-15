 CREATE TABLE target (
    id bigserial PRIMARY KEY,
    name varchar(40) UNIQUE NOT NULL,
    bodypart_id bigint NOT NULL REFERENCES body_part (id) ON DELETE CASCADE
);
