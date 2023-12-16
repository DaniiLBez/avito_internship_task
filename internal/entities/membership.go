package entities

import "time"

type Membership struct {
	UserId    int       `db:"user_id"`
	SlugId    int       `db:"slug_id"`
	CreatedAt time.Time `db:"created_at"`
}
