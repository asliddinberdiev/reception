package models

type AppointmentCreateInput struct {
	DoctorID        string `json:"doctor_id" validate:"required,uuid"`
	AppointmentDate string `json:"appointment_date" validate:"required"`
	AppointmentTime string `json:"appointment_time" validate:"required"`
}
