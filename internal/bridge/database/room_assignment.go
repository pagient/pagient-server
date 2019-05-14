package database

import (
	"github.com/pagient/pagient-server/internal/bridge/model"

	"github.com/pkg/errors"
)

// GetRoomAssignments returns current assignments of patients to surgery rooms
func (t *tx) GetRoomAssignments(roomSymbol string, limit ...uint) ([]*model.RoomAssignment, error) {
	top := 0
	if len(limit) > 0 {
		top = int(limit[0])
	}

	rows, err := t.Query("SELECT pds6_wz.PID FROM pds6_wz JOIN pds6_stwz ON pds6_wz.wzid = pds6_stwz.wzid WHERE pds6_stwz.code = @p1 ORDER BY pds6_wz.flgnr ASC", roomSymbol)
	if top > 0 {
		rows, err = t.Query("SELECT TOP(@p1) pds6_wz.PID FROM pds6_wz JOIN pds6_stwz ON pds6_wz.wzid = pds6_stwz.wzid WHERE pds6_stwz.code = @p2 ORDER BY pds6_wz.flgnr ASC", top, roomSymbol)
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not query database")
	}
	defer rows.Close()

	var assignments []*model.RoomAssignment
	for rows.Next() {
		entry := &model.RoomAssignment{}
		err := rows.Scan(&entry.PID)
		if err != nil {
			return nil, errors.Wrap(err, "could not scan database row")
		}
		assignments = append(assignments, entry)
	}

	err = rows.Err()
	return assignments, errors.Wrap(err, "row got an error")
}
