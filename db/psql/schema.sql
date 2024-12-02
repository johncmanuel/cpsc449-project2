CREATE TABLE courses (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE assignments (
    id INTEGER PRIMARY KEY,
    course_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_course 
        FOREIGN KEY(course_id) 
        REFERENCES courses(id) 
        ON DELETE CASCADE
);

-- index on course_id for faster lookups
CREATE INDEX idx_assignments_course_id ON assignments(course_id);

-- index on due_date for filtering assignments
CREATE INDEX idx_assignments_due_date ON assignments(due_date);
