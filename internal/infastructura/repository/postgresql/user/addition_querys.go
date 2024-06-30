package user_repository

import "fmt"

func (reop *UserRepository) Returning(data string) string {
	return fmt.Sprintf("RETURNING %s", data)
}

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
