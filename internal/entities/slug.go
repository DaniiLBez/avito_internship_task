package entities

import "time"

type Slug struct {
	ID        int       `db:"id"`
	Name      string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Users     []*User   `db:"users"`
}
