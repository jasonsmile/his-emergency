package handler

// CreateRegistrationRequest 是挂号接口的请求契约。收费与退费字段先保留，
// 当前阶段只建立接口边界，不代表收费已经完成。
type CreateRegistrationRequest struct {
	PatientID      string `json:"patient_id" binding:"required"`
	DeptCode       string `json:"dept_code" binding:"required"`
	DoctorID       string `json:"doctor_id"`
	ScheduleID     int64  `json:"schedule_id"`
	ChargeType     string `json:"charge_type"`
	IsEmergency    bool   `json:"is_emergency"`
	IsGreenChannel bool   `json:"is_green_channel"`
	PayMethod      string `json:"pay_method"`
	Remark         string `json:"remark"`
}
