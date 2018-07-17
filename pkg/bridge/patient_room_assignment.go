package bridge

type patientRoomAssignment struct {
	ID        int  `gorm:"primary_key,column:WZJID"`
	PatientID uint `gorm:"column:PID"`
}

func (ass *patientRoomAssignment) TableName() string {
	return "PDS6_WZ"
}
