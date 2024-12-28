CREATE TABLE IF NOT EXISTS attachments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email_id INT,
    type VARCHAR(255),
    filename VARCHAR(255),
    content TEXT,
    FOREIGN KEY (email_id) REFERENCES emails(id) ON DELETE CASCADE
);