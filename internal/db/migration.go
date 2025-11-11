package db

import (
	"context"
	"log"
)

// RunMigrations creates required tables if they don’t exist.
func (db *DB) RunMigrations(ctx context.Context) {
	queries := []string{
		//USERS TABLE
		`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(100) UNIQUE NOT NULL,
			mobile_number BIGINT UNIQUE NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        );
        `,

		`ALTER TABLE users ADD COLUMN IF NOT EXISTS name VARCHAR(100);`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS about TEXT;`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP;`,

		//CONTACTS TABLE
		`
		CREATE TABLE IF NOT EXISTS contacts (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			contact_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, contact_id)
		);
		`,
		`ALTER TABLE contacts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;`,

		//MESSAGES TABLE
		`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			sender_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			receiver_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			status VARCHAR(20) DEFAULT 'sent',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
	}

	for _, q := range queries {
		if _, err := db.Pool.Exec(ctx, q); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	}

	log.Println("✅ Database migrations executed successfully")
}
