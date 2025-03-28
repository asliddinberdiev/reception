package models

import "time"

type Appointment struct {
	ID              string       `json:"id"`
	Patient         PatientShort `json:"patient"`
	Doctor          DoctorShort  `json:"doctor"`
	AppointmentDate string       `json:"appointment_date" validate:"required"`
	AppointmentTime string       `json:"appointment_time" validate:"required"`
	Status          string       `json:"status"`
	DoctorApproval  bool         `json:"doctor_approval"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type AppointmentCreateInput struct {
	UserID          string `json:"-"`
	DoctorID        string `json:"doctor_id" validate:"required,uuid"`
	AppointmentDate string `json:"appointment_date" validate:"required"`
	AppointmentTime string `json:"appointment_time" validate:"required"`
}

type AppointmentRangeTime struct {
	DoctorID  string
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
}
