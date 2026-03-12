package models

import (
	"database/sql"
)

type Event struct {
	ID         string
	Central    int
	Link       int
	DeviceId   int
	EventType  string
	Local      string
	Device     string
	DeviceType string
	TsUnixMs   int64
	InstId     string
}

type EventModel struct {
	DB *sql.DB
}

func (m *EventModel) All(page, pageSize int, eventType string, centralID int) ([]*Event, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, central, link, device_id, event_type, local, device, device_type, ts_unix_ms, inst_id 
		FROM event 
		WHERE (? = '' OR event_type LIKE ?)
		AND (? = 0 OR central = ?)
		ORDER BY ts_unix_ms DESC 
		LIMIT ? OFFSET ?`

	searchPattern := "%" + eventType + "%"
	rows, err := m.DB.Query(query, eventType, searchPattern, centralID, centralID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		e := &Event{}
		err := rows.Scan(&e.ID, &e.Central, &e.Link, &e.DeviceId, &e.EventType, &e.Local, &e.Device, &e.DeviceType, &e.TsUnixMs, &e.InstId)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (m *EventModel) Count(eventType string, centralID int) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM event 
		WHERE (? = '' OR event_type LIKE ?)
		AND (? = 0 OR central = ?)`

	searchPattern := "%" + eventType + "%"
	err := m.DB.QueryRow(query, eventType, searchPattern, centralID, centralID).Scan(&count)
	return count, err
}
