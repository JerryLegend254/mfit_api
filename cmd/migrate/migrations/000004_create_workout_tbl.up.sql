CREATE TYPE workout_difficulty AS ENUM ('beginner', 'intermediate', 'advanced');
CREATE TABLE workout (
    id bigserial PRIMARY KEY,
    name varchar(40) UNIQUE NOT NULL,
    bodypart_id bigint NOT NULL REFERENCES body_part (id) ON DELETE CASCADE,
    equipment_id bigint REFERENCES equipment (id) ON DELETE CASCADE,
    gif_url varchar(255),
    instructions varchar(255) [],
    calories_burned integer,
    duration_minutes integer,
    difficulty workout_difficulty NOT NULL
);
