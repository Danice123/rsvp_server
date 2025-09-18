package main

import (
	"database/sql"
	"slices"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type db struct {
	conn *sql.DB
}

type Invite struct {
	Id        int
	FirstName string
	LastName  string
	HasEmail  bool
}

func (db *db) getFuzzyInvite(lastName string) ([]Invite, error) {
	rows, err := db.conn.Query("select id, first_name, last_name, email is not null from invite")
	if err != nil {
		return nil, err
	}

	uniqueLastNames := map[string]interface{}{}
	var invites []Invite
	for rows.Next() {
		var i Invite
		err = rows.Scan(&i.Id, &i.FirstName, &i.LastName, &i.HasEmail)
		if err != nil {
			return nil, err
		}
		invites = append(invites, i)
		uniqueLastNames[i.LastName] = nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	type nameDistance struct {
		Name     string
		Distance int
	}

	var lastNames []nameDistance
	for n := range uniqueLastNames {
		lastNames = append(lastNames, nameDistance{
			Name:     n,
			Distance: fuzzy.LevenshteinDistance(lastName, n),
		})
	}

	slices.SortFunc(lastNames, func(i, j nameDistance) int {
		return i.Distance - j.Distance
	})

	var results []Invite
	for i := 0; i < len(lastNames) && lastNames[i].Distance < 3; i++ {
		n := lastNames[i].Name
		for _, invite := range invites {
			if invite.LastName == n {
				results = append(results, invite)
			}
		}
	}
	return results, nil
}

type Person struct {
	Id       int
	InviteId int
	Name     string
	RSVP     *int
	Dietary  *string
	Notes    *string
}

func (db *db) getInvite(id int) ([]Person, error) {
	rows, err := db.conn.Query("select id, invite, name, rsvp, dietary, notes from person where invite = $1", id)
	if err != nil {
		return nil, err
	}

	var persons []Person
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.Id, &p.InviteId, &p.Name, &p.RSVP, &p.Dietary, &p.Notes)
		if err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}

func (db *db) rsvp(p Person) error {
	_, err := db.conn.Exec("update person set rsvp = $1, dietary = $2, notes = $3 where id = $4", p.RSVP, p.Dietary, p.Notes, p.Id)
	return err
}

func (db *db) updateEmail(id int, email string) error {
	_, err := db.conn.Exec("update invite set email = $1 where id = $2", email, id)
	return err
}
