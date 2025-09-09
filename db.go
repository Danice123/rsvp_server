package main

import (
	"database/sql"
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type db struct {
	conn *sql.DB
}

type Invite struct {
	Id        int
	FirstName string
	LastName  string
}

func (db *db) getFuzzyInvite(lastName string) ([]Invite, error) {
	rows, err := db.conn.Query("select id, first_name, last_name from invite")
	if err != nil {
		return nil, err
	}

	uniqueLastNames := map[string]interface{}{}
	var invites []Invite
	for rows.Next() {
		var i Invite
		err = rows.Scan(&i.Id, &i.FirstName, &i.LastName)
		if err != nil {
			return nil, err
		}
		invites = append(invites, i)
		uniqueLastNames[i.LastName] = nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var lastNames []string
	for n := range uniqueLastNames {
		lastNames = append(lastNames, n)
	}

	ranked := fuzzy.RankFindFold(lastName, lastNames)
	sort.Sort(ranked)

	var results []Invite
	for i := 0; i < len(ranked) && ranked[i].Distance < 5; i++ {
		n := ranked[i].Target
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
