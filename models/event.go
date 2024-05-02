package models

import (
	"time"

	"github.com/devphaseX/event-booking-api/db"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `binding:"required" json:"name"`
	Description string    `binding:"required" json:"description"`
	Location    string    `binding:"required" json:"location"`
	DateTime    time.Time `binding:"required" json:"dateTime"`
	UserID      int64     `json:"user_id"`
}

func (e *Event) Save() error {
	query := `
		INSERT INTO events (
			name, description, location, dateTime,  user_id
		) 	VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	e.ID = id
	return nil
}

func GetAllEvents(userId int64) ([]Event, error) {
	query := `
		SELECT * FROM events
		where user_id = ?
	`

	rows, err := db.DB.Query(query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []Event{}

	for rows.Next() {
		var event Event
		err = rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventById(id, userId int64) (*Event, error) {
	query := `
	SELECT * FROM  events
	where events.id = ? and user_id = ?
`

	row := db.DB.QueryRow(query, id, userId)

	var event Event

	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func DeleteEventById(id int64) error {
	query := `
		DELETE FROM events
		WHERE events.id = ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

func (ev Event) Update(userId int64) error {
	query := `
		UPDATE events
		SET	name=?, description=?, location=?, dateTime=?
		WHERE events.id = ? and user_id = ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ev.Name, ev.Description, ev.Location, ev.DateTime, ev.ID, userId)
	return err

}
