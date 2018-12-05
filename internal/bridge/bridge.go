package bridge

import (
	"sort"

	bridgeModel "github.com/pagient/pagient-server/internal/bridge/model"
	"github.com/pagient/pagient-server/internal/config"
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
)

// DB interface
type DB interface {
	Begin() (Tx, error)
}

// Tx interface
type Tx interface {
	Commit() error
	Rollback() error

	GetRoomAssignments(string, ...uint) ([]*bridgeModel.RoomAssignment, error)
}

// DefaultBridge struct encapsulates the surgery software bridge
type DefaultBridge struct {
	db              DB
	lastAssignments []*bridgeModel.RoomAssignment
}

// NewBridge returns a surgery software bridge struct
func NewBridge(db DB) *DefaultBridge {
	return &DefaultBridge{db, nil}
}

func (b *DefaultBridge) GetToBeExaminedPatients() ([]*model.Patient, error) {
	assignments, err := b.getRoomAssignments()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	patients := mapAssignmentsToPatients(assignments)

	return patients, nil
}

func (b *DefaultBridge) GetExaminedPatients() ([]*model.Patient, error) {
	assignments, err := b.getRoomAssignments()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	removedAssignments := disjointSet(b.lastAssignments, assignments)

	patients := mapAssignmentsToPatients(removedAssignments)

	// temporary store patients to retrieve finished/examined patients
	b.lastAssignments = make([]*bridgeModel.RoomAssignment, len(assignments))
	copy(b.lastAssignments, assignments)

	return patients, nil
}

func (b *DefaultBridge) getRoomAssignments() ([]*bridgeModel.RoomAssignment, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "create transaction failed")
	}

	assignments, err := tx.GetRoomAssignments(config.Bridge.CallActionWZ, config.Bridge.CallActionQueuePosition)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "get patients by room assignment failed")
	}

	err = tx.Commit()
	return assignments, errors.Wrap(err, "get room assignments failed")
}

func disjointSet(assignmentsA, assignmentsB []*bridgeModel.RoomAssignment) []*bridgeModel.RoomAssignment {
	sortAssignmentsByPID(assignmentsB)

	disjointSet := make([]*bridgeModel.RoomAssignment, max(len(assignmentsA), len(assignmentsB)))

	for _, assignmentA := range assignmentsA {
		found := false
		for _, assignmentB := range assignmentsB {
			if assignmentA.PID < assignmentB.PID {
				break
			}

			if assignmentA.PID == assignmentB.PID {
				found = true
				break
			}
		}

		if !found {
			disjointSet = append(disjointSet, assignmentA)
		}
	}

	return disjointSet
}

func mapAssignmentsToPatients(assignments []*bridgeModel.RoomAssignment) []*model.Patient {
	patients := make([]*model.Patient, len(assignments))
	for _, assignment := range assignments {
		patients = append(patients, &model.Patient{ID: assignment.PID})
	}

	return patients
}

func sortAssignmentsByPID(assignments []*bridgeModel.RoomAssignment) {
	sort.Slice(assignments, func(i, j int) bool {
		return assignments[i].PID < assignments[j].PID
	})
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
