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
	TypeId     string
}

type Installation struct {
	InstId     string
	EventCount int
}

type EventModel struct {
	DB *sql.DB
}

func (m *EventModel) Installations() ([]*Installation, error) {
	query := `
		SELECT inst_id, COUNT(*) as event_count 
		FROM event 
		GROUP BY inst_id 
		ORDER BY inst_id ASC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var installations []*Installation
	for rows.Next() {
		i := &Installation{}
		err := rows.Scan(&i.InstId, &i.EventCount)
		if err != nil {
			return nil, err
		}
		installations = append(installations, i)
	}
	return installations, nil
}

func (m *EventModel) All(page, pageSize int, eventType string, centralID int, instID string, device string) ([]*Event, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, central, link, device_id, event_type, local, device, device_type, ts_unix_ms, inst_id, type_id 
		FROM event 
		WHERE (? = '' OR event_type LIKE ?)
		AND (? = 0 OR central = ?)
		AND (? = '' OR inst_id LIKE ?)
		AND (? = '' OR device LIKE ?)
		ORDER BY inst_id ASC, ts_unix_ms DESC 
		LIMIT ? OFFSET ?`

	searchPattern := "%" + eventType + "%"
	instPattern := "%" + instID + "%"
	devicePattern := "%" + device + "%"
	rows, err := m.DB.Query(query, eventType, searchPattern, centralID, centralID, instID, instPattern, device, devicePattern, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		e := &Event{}
		err := rows.Scan(&e.ID, &e.Central, &e.Link, &e.DeviceId, &e.EventType, &e.Local, &e.Device, &e.DeviceType, &e.TsUnixMs, &e.InstId, &e.TypeId)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (m *EventModel) Count(eventType string, centralID int, instID string, device string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM event 
		WHERE (? = '' OR event_type LIKE ?)
		AND (? = 0 OR central = ?)
		AND (? = '' OR inst_id LIKE ?)
		AND (? = '' OR device LIKE ?)`

	searchPattern := "%" + eventType + "%"
	instPattern := "%" + instID + "%"
	devicePattern := "%" + device + "%"
	err := m.DB.QueryRow(query, eventType, searchPattern, centralID, centralID, instID, instPattern, device, devicePattern).Scan(&count)
	return count, err
}

func (m *EventModel) GetForExport(instID string, startMs, endMs int64) ([]*Event, error) {
	var query string
	var args []interface{}

	if instID == "" {
		query = `
			SELECT id, central, link, device_id, event_type, local, device, device_type, ts_unix_ms, inst_id 
			FROM event 
			WHERE ts_unix_ms >= ? AND ts_unix_ms <= ?
			ORDER BY ts_unix_ms ASC`
		args = []interface{}{startMs, endMs}
	} else {
		query = `
			SELECT id, central, link, device_id, event_type, local, device, device_type, ts_unix_ms, inst_id 
			FROM event 
			WHERE inst_id = ? 
			AND ts_unix_ms >= ? 
			AND ts_unix_ms <= ?
			ORDER BY ts_unix_ms ASC`
		args = []interface{}{instID, startMs, endMs}
	}

	rows, err := m.DB.Query(query, args...)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
