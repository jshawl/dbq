-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create posts table
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for faster queries
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_published ON posts(published);

-- Insert seed data for users
INSERT INTO users (username, email, first_name, last_name) VALUES
    ('john_doe', 'john.doe@example.com', 'John', 'Doe'),
    ('jane_smith', 'jane.smith@example.com', 'Jane', 'Smith'),
    ('bob_wilson', 'bob.wilson@example.com', 'Bob', 'Wilson'),
    ('alice_johnson', 'alice.johnson@example.com', 'Alice', 'Johnson'),
    ('charlie_brown', 'charlie.brown@example.com', 'Charlie', 'Brown');

-- Insert seed data for posts
INSERT INTO posts (user_id, title, content, published) VALUES
    (1, 'Getting Started with PostgreSQL', 'PostgreSQL is a powerful, open-source object-relational database system...', true),
    (1, 'Docker Best Practices', 'When working with Docker containers, there are several best practices to follow...', true),
    (2, 'Introduction to Node.js', 'Node.js is a JavaScript runtime built on Chrome''s V8 JavaScript engine...', true),
    (2, 'Understanding REST APIs', 'REST (Representational State Transfer) is an architectural style for designing web services...', false),
    (3, 'Python for Beginners', 'Python is a high-level, interpreted programming language known for its simplicity...', true),
    (3, 'Data Structures and Algorithms', 'Understanding data structures and algorithms is crucial for any programmer...', true),
    (4, 'Web Development Trends 2024', 'The web development landscape continues to evolve rapidly...', false),
    (4, 'Frontend vs Backend Development', 'Web development is typically divided into two main areas: frontend and backend...', true),
    (5, 'Machine Learning Basics', 'Machine learning is a subset of artificial intelligence that focuses on algorithms...', true),
    (5, 'Database Design Principles', 'Good database design is fundamental to building efficient and scalable applications...', true);

-- Create a simple view for posts with user information
CREATE VIEW posts_with_users AS
SELECT 
    p.id,
    p.title,
    p.content,
    p.published,
    p.created_at,
    u.username,
    u.email,
    CONCAT(u.first_name, ' ', u.last_name) AS author_name
FROM posts p
JOIN users u ON p.user_id = u.id;