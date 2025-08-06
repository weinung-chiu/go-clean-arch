-- Create articles table
CREATE TABLE IF NOT EXISTS articles (
    id VARCHAR(255) PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP WITH TIME ZONE
);

-- Create index on published_at for efficient filtering of published articles
CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at);

-- Create index on author_id for filtering by author
CREATE INDEX IF NOT EXISTS idx_articles_author_id ON articles(author_id);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at);