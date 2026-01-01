package repo

import "database/sql"

type Repos struct {
	Users     *UserRepo
	Todos     *TodoRepo
	Statuses  *TodoStatusRepo
	Goodlucks *GoodluckRepo
}

func New(db *sql.DB) *Repos {
	return &Repos{
		Users:     &UserRepo{db: db},
		Todos:     &TodoRepo{db: db},
		Statuses:  &TodoStatusRepo{db: db},
		Goodlucks: &GoodluckRepo{db: db},
	}
}


