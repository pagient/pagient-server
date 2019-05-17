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

// GetToBeExaminedPatients returns all patients that are queued to be examined next
func (b *DefaultBridge) GetToBeExaminedPatients() ([]*model.Patient, error) {
	assignments, err := b.db.GetRoomAssignments(config.Bridge.CallActionWZ, config.Bridge.CallActionQueuePosition)
	if err != nil {
		return nil, errors.Wrap(err, "get patients by room assignment failed")
	}

	patients := mapAssignmentsToPatients(assignments)

	return patients, nil
}

// GetExaminedPatients returns all patients that have been examined and are finished now since last call
func (b *DefaultBridge) GetExaminedPatients() ([]*model.Patient, error) {
	assignments, err := b.db.GetRoomAssignments(config.Bridge.CallActionWZ, config.Bridge.CallActionQueuePosition)
	if err != nil {
		return nil, errors.Wrap(err, "get patients by room assignment failed")
	}

	removedAssignments := subtractSet(b.lastAssignments, assignments)

	patients := mapAssignmentsToPatients(removedAssignments)

	// temporary store patients to retrieve finished/examined patients
	b.lastAssignments = make([]*bridgeModel.RoomAssignment, len(assignments))
	copy(b.lastAssignments, assignments)

	return patients, nil
}

func subtractSet(assignmentsA, assignmentsB []*bridgeModel.RoomAssignment) []*bridgeModel.RoomAssignment {
	sortAssignmentsByPID(assignmentsB)

	subtractSet := make([]*bridgeModel.RoomAssignment, 0, len(assignmentsA))
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
			subtractSet = append(subtractSet, assignmentA)
		}
	}

	return subtractSet
}

func mapAssignmentsToPatients(assignments []*bridgeModel.RoomAssignment) []*model.Patient {
	patients := make([]*model.Patient, 0, len(assignments))
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
