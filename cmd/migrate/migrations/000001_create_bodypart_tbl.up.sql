 CREATE TABLE body_part (
    id bigserial PRIMARY KEY,
    name varchar(40) UNIQUE NOT NULL,
    image_url varchar(255) NOT NULL
);
