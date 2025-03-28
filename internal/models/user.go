package models

type ProfileShort struct {
	ID          string     `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Specialty   string     `json:"specialty"`
	Description string     `json:"description"`
	WorkTime    []WorkTime `json:"work_time"`
}

type WorkTime struct {
	WeekDay    uint8  `json:"week_day"`
	StartTime  string `json:"start_time"`
	FinishTime string `json:"finish_time"`
}

type GetAllProfileShort struct {
	Data  []ProfileShort `json:"data"`
	Total uint32         `json:"total"`
}

type DoctorShort struct {
	ID        string `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}
