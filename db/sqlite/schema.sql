-- schema.sql
CREATE TABLE IF NOT EXISTS courses (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assignments (
    id INTEGER PRIMARY KEY,
    course_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    due_date DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    difficulty INTEGER CHECK(difficulty BETWEEN 1 AND 10),  -- New column for difficulty with valid range
    length INTEGER,    -- New column for length
    FOREIGN KEY(course_id) REFERENCES courses(id) 
    ON DELETE CASCADE
);


-- Index for faster lookups (optional in SQLite, but can improve performance)
CREATE INDEX IF NOT EXISTS idx_assignments_course_id ON assignments(course_id);
CREATE INDEX IF NOT EXISTS idx_assignments_due_date ON assignments(due_date);

-- CREATE TABLE courses (
--     id INTEGER PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );
--
-- CREATE TABLE assignments (
--     id INTEGER PRIMARY KEY,
--     course_id INTEGER NOT NULL,
--     name VARCHAR(255) NOT NULL,
--     due_date TIMESTAMP WITH TIME ZONE,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--
--     CONSTRAINT fk_course 
--         FOREIGN KEY(course_id) 
--         REFERENCES courses(id) 
--         ON DELETE CASCADE
-- );
--
-- -- index on course_id for faster lookups
-- CREATE INDEX idx_assignments_course_id ON assignments(course_id);
--
-- -- index on due_date for filtering assignments
-- CREATE INDEX idx_assignments_due_date ON assignments(due_date);
