CREATE TYPE target_type AS ENUM ('primary', 'secondary');
CREATE TABLE workout_target (
    workout_id bigint NOT NULL REFERENCES workout (id) ON DELETE CASCADE,
    target_id bigint REFERENCES target (id) ON DELETE CASCADE,
    type target_type NOT NULL,
    UNIQUE (workout_id, target_id)
);
