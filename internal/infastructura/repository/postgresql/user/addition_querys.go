package user_repository

import "fmt"

func (reop *UserRepository) Returning(data string) string {
	return fmt.Sprintf("RETURNING %s", data)
}

// CREATE TABLE IF NOT EXISTS users (
// id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// first_name VARCHAR(255) NOT NULL,
// last_name VARCHAR(255) NOT NULL,
// email VARCHAR(255) UNIQUE NOT NULL,
// password VARCHAR(255) NOT NULL,
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// deleted_at TIMESTAMP DEFAULT NULL,
// is_registered BOOLEAN DEFAULT FALSE,
// secret_code VARCHAR(255) NOT NULL
// );
func (reop *UserRepository) SelectQuery() string {
	return `
	id,
	first_name,
	last_name,
	email,
	password,
	created_at,
	updated_at,
	deleted_at,
	is_registered
`
}
