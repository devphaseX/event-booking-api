package models

import (
	"fmt"

	"github.com/devphaseX/event-booking-api/db"
)

type Ticket struct {
	ID      int64 `json:"id"`
	UserId  int64 `binding:"required" json:"user_id"`
	EventId int64 `binding:"required" json:"event_id"`
}

type EventRegisterUserInput struct {
	RegisteredUserID int64 `json:"registered_user_id" binding:"required"`
}

func (t *Ticket) Save() error {
	query := `
		INSERT INTO tickets (
			user_id, 
			event_id
		) 	VALUES (?, ?)
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(t.UserId, t.EventId)

	if err != nil {
		return err
	}

	ticketId, err := res.LastInsertId()

	if err != nil {
		return err
	}

	t.ID = ticketId

	return nil
}

func GetRegUserTicket(eventId, regUserId, ownedUserId int64) (*Ticket, error) {
	fmt.Println(eventId, regUserId, ownedUserId)
	query := `
		SELECT tickets.id as id FROM tickets 
		INNER JOIN events ON events.user_id = ?
		WHERE tickets.user_id = ? and event_id = ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(ownedUserId, regUserId, eventId)
	var ticket Ticket

	if row.Err() != nil {
		return nil, err
	}

	err = row.Scan(&ticket.ID)

	if err != nil {
		return nil, err
	}

	ticket.UserId = regUserId
	ticket.EventId = eventId

	return &ticket, nil
}

func DeleteTicketById(id int64) error {
	query := `
		DELETE FROM tickets
		WHERE id = ?
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
