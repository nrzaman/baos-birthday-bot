-- Birthday database schema

CREATE TABLE IF NOT EXISTS birthdays (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    month INTEGER NOT NULL CHECK(month >= 1 AND month <= 12),
    day INTEGER NOT NULL CHECK(day >= 1 AND day <= 31),
    gender TEXT CHECK(gender IN ('male', 'female', 'nonbinary', 'other', NULL)),  -- For pronoun reference
    discord_id TEXT,  -- Optional: link to Discord user
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name)  -- One birthday per name
);

-- Index for faster birthday lookups
CREATE INDEX IF NOT EXISTS idx_birthdays_date ON birthdays(month, day);
CREATE INDEX IF NOT EXISTS idx_birthdays_discord_id ON birthdays(discord_id);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_birthdays_timestamp
AFTER UPDATE ON birthdays
BEGIN
    UPDATE birthdays SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
