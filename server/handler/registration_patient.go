package handler

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	apperrors "emergency-his/server/errors"
	"emergency-his/server/logger"
	"emergency-his/server/middleware"
	"emergency-his/server/response"
	"emergency-his/server/utils"
	"github.com/gin-gonic/gin"
)

type SearchPatientRequest struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=20"`
}
type PatientModel struct {
	ID                                                                                                                    int64
	PatientID, Name, Gender, IDCardNo, Phone, ChargeType, IdentityType, Address, EmergencyContact, EmergencyPhone, CardNo string
	BirthDate                                                                                                             time.Time
}
type PatientInfo struct {
	ID               int64  `json:"id"`
	PatientID        string `json:"patient_id"`
	Name             string `json:"name"`
	Gender           string `json:"gender"`
	BirthDate        string `json:"birth_date"`
	Age              string `json:"age"`
	IDCardNo         string `json:"id_card_no"`
	Phone            string `json:"phone"`
	ChargeType       string `json:"charge_type"`
	IdentityType     string `json:"identity_type"`
	Address          string `json:"address"`
	EmergencyContact string `json:"emergency_contact"`
	EmergencyPhone   string `json:"emergency_phone"`
	CardNo           string `json:"card_no"`
}
type DeptModel struct {
	ID                                          int64
	DeptCode, DeptName, DeptShortName, DeptType string
	IsOutpatient, IsInpatient, IsEmergency      int
}
type DeptInfo struct {
	ID            int64  `json:"id"`
	DeptCode      string `json:"dept_code"`
	DeptName      string `json:"dept_name"`
	DeptShortName string `json:"dept_short_name"`
	DeptType      string `json:"dept_type"`
	IsOutpatient  int    `json:"is_outpatient"`
	IsEmergency   int    `json:"is_emergency"`
}
type DoctorModel struct {
	ID                                                                      int64
	DoctorCode, Name, Title, TitleCode, DeptCode, Specialty, Room, RoomCode string
}
type DoctorInfo struct {
	ID         int64  `json:"id"`
	DoctorCode string `json:"doctor_code"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	DeptCode   string `json:"dept_code"`
	Specialty  string `json:"specialty"`
	Room       string `json:"room"`
}
type GetScheduleRequest struct {
	ClinicDate string `form:"clinic_date"`
	DeptCode   string `form:"dept_code"`
	DoctorID   string `form:"doctor_id"`
	TimeDesc   string `form:"time_desc"`
}
type ScheduleModel struct {
	ID                                                                int64
	ClinicDate                                                        time.Time
	ClinicDept, DeptName, ClinicLabel, TimeDesc, DoctorID, DoctorName string
	RegistrationLimits, RegistrationNum                               int
	RegistPrice                                                       float64
	ClinicType, States                                                string
}
type ScheduleInfo struct {
	ID                 int64   `json:"id"`
	ClinicDate         string  `json:"clinic_date"`
	ClinicDept         string  `json:"clinic_dept"`
	DeptName           string  `json:"dept_name"`
	ClinicLabel        string  `json:"clinic_label"`
	TimeDesc           string  `json:"time_desc"`
	DoctorID           string  `json:"doctor_id"`
	DoctorName         string  `json:"doctor_name"`
	RegistrationLimits int     `json:"registration_limits"`
	RegistrationNum    int     `json:"registration_num"`
	AvailableNum       int     `json:"available_num"`
	RegistPrice        float64 `json:"regist_price"`
	ClinicType         string  `json:"clinic_type"`
	States             string  `json:"states"`
}
type EncounterModel struct {
	VisitNo, PatientID, PatientName, PatientAge, Gender, IDCardNo, Phone, ChargeType, EncounterType, DeptID, DeptName, DoctorID, DoctorName, DoctorTitle, Remark string
	IsEmergency, IsGreenChannel                                                                                                                                  int
}
type GetEncounterListRequest struct {
	Status    string `form:"status"`
	DeptCode  string `form:"dept_code"`
	DoctorID  string `form:"doctor_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Keyword   string `form:"keyword"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"pageSize,default=20"`
}
type EncounterInfo struct {
	ID              int64      `json:"id"`
	VisitNo         string     `json:"visit_no"`
	PatientID       string     `json:"patient_id"`
	PatientName     string     `json:"patient_name"`
	PatientAge      string     `json:"patient_age"`
	Gender          string     `json:"gender"`
	Phone           string     `json:"phone"`
	ChargeType      string     `json:"charge_type"`
	EncounterType   string     `json:"encounter_type"`
	DeptID          string     `json:"dept_id"`
	DeptName        string     `json:"dept_name"`
	DoctorID        string     `json:"doctor_id"`
	DoctorName      string     `json:"doctor_name"`
	DoctorTitle     string     `json:"doctor_title"`
	RegTime         time.Time  `json:"reg_time"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	EncounterStatus string     `json:"encounter_status"`
	TotalAmount     float64    `json:"total_amount"`
	PaidAmount      float64    `json:"paid_amount"`
	PaymentStatus   string     `json:"payment_status"`
	IsEmergency     int        `json:"is_emergency"`
	IsGreenChannel  int        `json:"is_green_channel"`
	Remark          string     `json:"remark"`
	CreatedAt       time.Time  `json:"created_at"`
}
type RegistrationResult struct {
	EncounterID   int64   `json:"encounter_id"`
	VisitNo       string  `json:"visit_no"`
	PatientID     string  `json:"patient_id"`
	PatientName   string  `json:"patient_name"`
	DeptName      string  `json:"dept_name"`
	DoctorName    string  `json:"doctor_name"`
	RegTime       string  `json:"reg_time"`
	RegistPrice   float64 `json:"regist_price"`
	ChargeNo      string  `json:"charge_no"`
	ChargeStatus  string  `json:"charge_status"`
	PaymentMethod string  `json:"payment_method"`
}
type RegistrationRepository struct{ db *sql.DB }
type RegistrationService struct{ repo *RegistrationRepository }

func NewRegistrationService(db *sql.DB) *RegistrationService {
	return &RegistrationService{repo: &RegistrationRepository{db: db}}
}

func (r *RegistrationRepository) SearchPatients(keyword string, page, pageSize int) ([]*PatientModel, int64, error) {
	where, args := " WHERE 1=1 ", []interface{}{}
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		like := "%" + keyword + "%"
		where += " AND (name LIKE ? OR id_card_no LIKE ? OR phone LIKE ? OR card_no LIKE ? OR patient_id LIKE ?) "
		args = append(args, like, like, like, like, like)
	}
	var total int64
	countQuery := "SELECT COUNT(*) FROM sync_patient" + where
	logger.SQL(countQuery, args...)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, pageSize, (page-1)*pageSize)
	query := `SELECT id,patient_id,name,gender,birth_date,id_card_no,phone,charge_type,identity_type,address,emergency_contact,emergency_phone,card_no FROM sync_patient` + where + ` ORDER BY updated_at DESC LIMIT ? OFFSET ?`
	logger.SQL(query, args...)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	result := make([]*PatientModel, 0)
	for rows.Next() {
		p := &PatientModel{}
		var birth sql.NullTime
		if err := rows.Scan(&p.ID, &p.PatientID, &p.Name, &p.Gender, &birth, &p.IDCardNo, &p.Phone, &p.ChargeType, &p.IdentityType, &p.Address, &p.EmergencyContact, &p.EmergencyPhone, &p.CardNo); err != nil {
			return nil, 0, err
		}
		if birth.Valid {
			p.BirthDate = birth.Time
		}
		result = append(result, p)
	}
	return result, total, rows.Err()
}

func (r *RegistrationRepository) GetPatientByID(patientID string) (*PatientModel, error) {
	p := &PatientModel{}
	var birthDate sql.NullTime
	query := `SELECT id, patient_id, name, gender, birth_date, id_card_no, phone,
		charge_type, identity_type, address, emergency_contact, emergency_phone, card_no
		FROM sync_patient WHERE patient_id = ?`
	logger.SQL(query, patientID)
	err := r.db.QueryRow(query, patientID).Scan(
		&p.ID, &p.PatientID, &p.Name, &p.Gender, &birthDate,
		&p.IDCardNo, &p.Phone, &p.ChargeType, &p.IdentityType, &p.Address,
		&p.EmergencyContact, &p.EmergencyPhone, &p.CardNo,
	)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewBusinessError(40401, "患者不存在")
	}
	if err != nil {
		return nil, err
	}
	if birthDate.Valid {
		p.BirthDate = birthDate.Time
	}
	return p, nil
}

func (r *RegistrationRepository) GetDepts(isOutpatientOnly bool) ([]*DeptModel, error) {
	query := `SELECT id, dept_code, dept_name, dept_short_name, dept_type, is_outpatient, is_inpatient, is_emergency FROM sync_dept WHERE status = 1`
	if isOutpatientOnly {
		query += " AND is_outpatient = 1"
	}
	query += " ORDER BY dept_code"
	logger.SQL(query)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	depts := make([]*DeptModel, 0)
	for rows.Next() {
		d := &DeptModel{}
		if err := rows.Scan(&d.ID, &d.DeptCode, &d.DeptName, &d.DeptShortName, &d.DeptType, &d.IsOutpatient, &d.IsInpatient, &d.IsEmergency); err != nil {
			return nil, err
		}
		depts = append(depts, d)
	}
	return depts, rows.Err()
}

func (r *RegistrationRepository) GetDoctors(deptCode string) ([]*DoctorModel, error) {
	query := `SELECT id, doctor_code, name, title, title_code, dept_code, specialty, room, room_code FROM sync_doctor WHERE status = 1`
	args := []interface{}{}
	if deptCode != "" {
		query += " AND dept_code = ?"
		args = append(args, deptCode)
	}
	query += " ORDER BY doctor_code"
	logger.SQL(query, args...)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	doctors := make([]*DoctorModel, 0)
	for rows.Next() {
		d := &DoctorModel{}
		if err := rows.Scan(&d.ID, &d.DoctorCode, &d.Name, &d.Title, &d.TitleCode, &d.DeptCode, &d.Specialty, &d.Room, &d.RoomCode); err != nil {
			return nil, err
		}
		doctors = append(doctors, d)
	}
	return doctors, rows.Err()
}

func (r *RegistrationRepository) GetDoctorByCode(doctorCode string) (*DoctorModel, error) {
	d := &DoctorModel{}
	query := `SELECT id, doctor_code, name, title, title_code, dept_code, specialty, room, room_code FROM sync_doctor WHERE doctor_code = ? AND status = 1`
	logger.SQL(query, doctorCode)
	err := r.db.QueryRow(query, doctorCode).Scan(&d.ID, &d.DoctorCode, &d.Name, &d.Title, &d.TitleCode, &d.DeptCode, &d.Specialty, &d.Room, &d.RoomCode)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewBusinessError(40402, "医生不存在或已停用")
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *RegistrationRepository) GetSchedules(clinicDate, deptCode, doctorID, timeDesc string) ([]*ScheduleModel, error) {
	query := `SELECT id, clinic_date, clinic_dept, dept_name, clinic_label, time_desc, doctor_id, doctor_name, registration_limits, registration_num, regist_price, clinic_type, states FROM sync_schedule WHERE clinic_date = ? AND states = '正常'`
	args := []interface{}{clinicDate}
	for _, filter := range []struct{ value, column string }{{deptCode, "clinic_dept"}, {doctorID, "doctor_id"}, {timeDesc, "time_desc"}} {
		if filter.value != "" {
			query += " AND " + filter.column + " = ?"
			args = append(args, filter.value)
		}
	}
	query += " ORDER BY clinic_dept, time_desc, doctor_id"
	logger.SQL(query, args...)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	schedules := make([]*ScheduleModel, 0)
	for rows.Next() {
		s := &ScheduleModel{}
		if err := rows.Scan(&s.ID, &s.ClinicDate, &s.ClinicDept, &s.DeptName, &s.ClinicLabel, &s.TimeDesc, &s.DoctorID, &s.DoctorName, &s.RegistrationLimits, &s.RegistrationNum, &s.RegistPrice, &s.ClinicType, &s.States); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, rows.Err()
}

func (r *RegistrationRepository) GetScheduleByID(id int64) (*ScheduleModel, error) {
	s := &ScheduleModel{}
	query := `SELECT id, clinic_date, clinic_dept, dept_name, clinic_label, time_desc, doctor_id, doctor_name, registration_limits, registration_num, regist_price, clinic_type, states FROM sync_schedule WHERE id = ?`
	logger.SQL(query, id)
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.ClinicDate, &s.ClinicDept, &s.DeptName, &s.ClinicLabel, &s.TimeDesc, &s.DoctorID, &s.DoctorName, &s.RegistrationLimits, &s.RegistrationNum, &s.RegistPrice, &s.ClinicType, &s.States)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewBusinessError(40403, "排班不存在")
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *RegistrationRepository) IncrementScheduleRegNum(scheduleID int64) error {
	query := `UPDATE sync_schedule SET registration_num = registration_num + 1 WHERE id = ? AND registration_num < registration_limits`
	logger.SQL(query, scheduleID)
	result, err := r.db.Exec(query, scheduleID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return apperrors.NewBusinessError(40901, "号源已满")
	}
	return nil
}

func (r *RegistrationRepository) GetTodayEncounterCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM emergency_encounter WHERE DATE(reg_time) = CURDATE()`
	logger.SQL(query)
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *RegistrationRepository) CreateEncounter(e *EncounterModel) (int64, error) {
	query := `INSERT INTO emergency_encounter (visit_no, patient_id, patient_name, patient_age, gender, id_card_no, phone, charge_type, encounter_type, dept_id, dept_name, doctor_id, doctor_name, doctor_title, reg_time, encounter_status, is_emergency, is_green_channel, payment_status, remark) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), 'REGISTERED', ?, ?, 'UNPAID', ?)`
	args := []interface{}{e.VisitNo, e.PatientID, e.PatientName, e.PatientAge, e.Gender, e.IDCardNo, e.Phone, e.ChargeType, e.EncounterType, e.DeptID, e.DeptName, e.DoctorID, e.DoctorName, e.DoctorTitle, e.IsEmergency, e.IsGreenChannel, e.Remark}
	logger.SQL(query, args...)
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *RegistrationRepository) HasActiveRegistration(patientID, deptID, doctorID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM emergency_encounter WHERE patient_id = ? AND dept_id = ? AND doctor_id = ? AND DATE(reg_time) = CURDATE() AND encounter_status IN ('REGISTERED', 'IN_TREATMENT')`
	logger.SQL(query, patientID, deptID, doctorID)
	err := r.db.QueryRow(query, patientID, deptID, doctorID).Scan(&count)
	return count > 0, err
}

func (r *RegistrationRepository) GetEncounterList(req *GetEncounterListRequest) ([]*EncounterInfo, int64, error) {
	where, args := " WHERE 1=1 ", []interface{}{}
	for _, f := range []struct{ v, col string }{{req.Status, "encounter_status"}, {req.DeptCode, "dept_id"}, {req.DoctorID, "doctor_id"}} {
		if f.v != "" {
			where += " AND " + f.col + " = ?"
			args = append(args, f.v)
		}
	}
	if req.StartDate != "" {
		where += " AND DATE(reg_time) >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		where += " AND DATE(reg_time) <= ?"
		args = append(args, req.EndDate)
	}
	if req.Keyword != "" {
		where += " AND (patient_name LIKE ? OR visit_no LIKE ?)"
		like := "%" + req.Keyword + "%"
		args = append(args, like, like)
	}
	var total int64
	countQuery := "SELECT COUNT(*) FROM emergency_encounter" + where
	logger.SQL(countQuery, args...)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)
	query := `SELECT id, visit_no, patient_id,
		COALESCE(patient_name, ''), COALESCE(patient_age, ''), COALESCE(gender, ''),
		COALESCE(phone, ''), COALESCE(charge_type, ''), COALESCE(encounter_type, ''),
		COALESCE(dept_id, ''), COALESCE(dept_name, ''), COALESCE(doctor_id, ''),
		COALESCE(doctor_name, ''), COALESCE(doctor_title, ''), reg_time,
		COALESCE(encounter_status, ''), COALESCE(total_amount, 0), COALESCE(paid_amount, 0),
		COALESCE(payment_status, ''), COALESCE(is_emergency, 0), COALESCE(is_green_channel, 0),
		COALESCE(remark, ''), created_at
		FROM emergency_encounter` + where + ` ORDER BY reg_time DESC LIMIT ? OFFSET ?`
	logger.SQL(query, args...)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]*EncounterInfo, 0)
	for rows.Next() {
		e := &EncounterInfo{}
		if err := rows.Scan(&e.ID, &e.VisitNo, &e.PatientID, &e.PatientName, &e.PatientAge, &e.Gender, &e.Phone, &e.ChargeType, &e.EncounterType, &e.DeptID, &e.DeptName, &e.DoctorID, &e.DoctorName, &e.DoctorTitle, &e.RegTime, &e.EncounterStatus, &e.TotalAmount, &e.PaidAmount, &e.PaymentStatus, &e.IsEmergency, &e.IsGreenChannel, &e.Remark, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, e)
	}
	return out, total, rows.Err()
}

func (r *RegistrationRepository) GetEncounterByID(id int64) (*EncounterInfo, error) {
	e := &EncounterInfo{}
	var startTime, endTime sql.NullTime
	query := `SELECT id, visit_no, patient_id,
		COALESCE(patient_name, ''), COALESCE(patient_age, ''), COALESCE(gender, ''),
		COALESCE(phone, ''), COALESCE(charge_type, ''), COALESCE(encounter_type, ''),
		COALESCE(dept_id, ''), COALESCE(dept_name, ''), COALESCE(doctor_id, ''),
		COALESCE(doctor_name, ''), COALESCE(doctor_title, ''), reg_time, start_time, end_time,
		COALESCE(encounter_status, ''), COALESCE(total_amount, 0), COALESCE(paid_amount, 0),
		COALESCE(payment_status, ''), COALESCE(is_emergency, 0), COALESCE(is_green_channel, 0),
		COALESCE(remark, ''), created_at
		FROM emergency_encounter WHERE id=?`
	logger.SQL(query, id)
	err := r.db.QueryRow(query, id).Scan(&e.ID, &e.VisitNo, &e.PatientID, &e.PatientName, &e.PatientAge, &e.Gender, &e.Phone, &e.ChargeType, &e.EncounterType, &e.DeptID, &e.DeptName, &e.DoctorID, &e.DoctorName, &e.DoctorTitle, &e.RegTime, &startTime, &endTime, &e.EncounterStatus, &e.TotalAmount, &e.PaidAmount, &e.PaymentStatus, &e.IsEmergency, &e.IsGreenChannel, &e.Remark, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewBusinessError(40404, "就诊记录不存在")
	}
	if err != nil {
		return nil, err
	}
	if startTime.Valid {
		e.StartTime = &startTime.Time
	}
	if endTime.Valid {
		e.EndTime = &endTime.Time
	}
	return e, nil
}

func (r *RegistrationRepository) UpdateEncounterStatus(id int64, status string) error {
	query := `UPDATE emergency_encounter SET encounter_status = ? WHERE id = ?`
	logger.SQL(query, status, id)
	_, err := r.db.Exec(query, status, id)
	return err
}

func (s *RegistrationService) SearchPatients(req *SearchPatientRequest) ([]*PatientInfo, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	patients, total, err := s.repo.SearchPatients(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*PatientInfo, 0, len(patients))
	for _, p := range patients {
		info := &PatientInfo{ID: p.ID, PatientID: p.PatientID, Name: p.Name, Gender: p.Gender, IDCardNo: p.IDCardNo, Phone: p.Phone, ChargeType: p.ChargeType, IdentityType: p.IdentityType, Address: p.Address, EmergencyContact: p.EmergencyContact, EmergencyPhone: p.EmergencyPhone, CardNo: p.CardNo}
		if !p.BirthDate.IsZero() {
			info.BirthDate = p.BirthDate.Format("2006-01-02")
			info.Age = utils.CalculateAge(p.BirthDate)
		}
		result = append(result, info)
	}
	return result, total, nil
}

func (s *RegistrationService) GetPatientDetail(patientID string) (*PatientInfo, error) {
	p, err := s.repo.GetPatientByID(patientID)
	if err != nil {
		return nil, err
	}
	info := &PatientInfo{ID: p.ID, PatientID: p.PatientID, Name: p.Name, Gender: p.Gender,
		IDCardNo: p.IDCardNo, Phone: p.Phone, ChargeType: p.ChargeType, IdentityType: p.IdentityType,
		Address: p.Address, EmergencyContact: p.EmergencyContact, EmergencyPhone: p.EmergencyPhone, CardNo: p.CardNo}
	if !p.BirthDate.IsZero() {
		info.BirthDate = p.BirthDate.Format("2006-01-02")
		info.Age = utils.CalculateAge(p.BirthDate)
	}
	return info, nil
}

func (s *RegistrationService) GetDepts(isOutpatientOnly bool) ([]*DeptInfo, error) {
	depts, err := s.repo.GetDepts(isOutpatientOnly)
	if err != nil {
		return nil, err
	}
	result := make([]*DeptInfo, 0, len(depts))
	for _, d := range depts {
		result = append(result, &DeptInfo{ID: d.ID, DeptCode: d.DeptCode, DeptName: d.DeptName, DeptShortName: d.DeptShortName, DeptType: d.DeptType, IsOutpatient: d.IsOutpatient, IsEmergency: d.IsEmergency})
	}
	return result, nil
}

func (s *RegistrationService) GetDoctors(deptCode string) ([]*DoctorInfo, error) {
	doctors, err := s.repo.GetDoctors(deptCode)
	if err != nil {
		return nil, err
	}
	result := make([]*DoctorInfo, 0, len(doctors))
	for _, d := range doctors {
		result = append(result, &DoctorInfo{ID: d.ID, DoctorCode: d.DoctorCode, Name: d.Name, Title: d.Title, DeptCode: d.DeptCode, Specialty: d.Specialty, Room: d.Room})
	}
	return result, nil
}

func (s *RegistrationService) GetSchedules(req *GetScheduleRequest) ([]*ScheduleInfo, error) {
	clinicDate := req.ClinicDate
	if clinicDate == "" {
		clinicDate = time.Now().Format("2006-01-02")
	} else if _, err := time.Parse("2006-01-02", clinicDate); err != nil {
		return nil, apperrors.NewBusinessError(40002, "日期格式应为 YYYY-MM-DD")
	}
	schedules, err := s.repo.GetSchedules(clinicDate, req.DeptCode, req.DoctorID, req.TimeDesc)
	if err != nil {
		return nil, err
	}
	result := make([]*ScheduleInfo, 0, len(schedules))
	for _, sc := range schedules {
		available := sc.RegistrationLimits - sc.RegistrationNum
		if available < 0 {
			available = 0
		}
		result = append(result, &ScheduleInfo{ID: sc.ID, ClinicDate: sc.ClinicDate.Format("2006-01-02"), ClinicDept: sc.ClinicDept, DeptName: sc.DeptName, ClinicLabel: sc.ClinicLabel, TimeDesc: sc.TimeDesc, DoctorID: sc.DoctorID, DoctorName: sc.DoctorName, RegistrationLimits: sc.RegistrationLimits, RegistrationNum: sc.RegistrationNum, AvailableNum: available, RegistPrice: sc.RegistPrice, ClinicType: sc.ClinicType, States: sc.States})
	}
	return result, nil
}

func (s *RegistrationService) CreateRegistration(req *CreateRegistrationRequest, operatorID, operatorName string) (*RegistrationResult, error) {
	patient, err := s.repo.GetPatientByID(req.PatientID)
	if err != nil {
		return nil, err
	}
	depts, err := s.repo.GetDepts(false)
	if err != nil {
		return nil, err
	}
	deptName := ""
	for _, dept := range depts {
		if dept.DeptCode == req.DeptCode {
			deptName = dept.DeptName
			break
		}
	}
	if deptName == "" {
		return nil, apperrors.NewBusinessError(40010, "科室不存在或已停用")
	}
	var doctorName, doctorTitle string
	if req.DoctorID != "" {
		doctor, err := s.repo.GetDoctorByCode(req.DoctorID)
		if err != nil {
			return nil, err
		}
		if doctor.DeptCode != req.DeptCode {
			return nil, apperrors.NewBusinessError(40011, "医生不属于所选科室")
		}
		doctorName, doctorTitle = doctor.Name, doctor.Title
	}
	var price float64
	if req.ScheduleID > 0 {
		schedule, err := s.repo.GetScheduleByID(req.ScheduleID)
		if err != nil {
			return nil, err
		}
		if schedule.States != "正常" {
			return nil, apperrors.NewBusinessError(40012, "排班已停诊")
		}
		if schedule.ClinicDept != req.DeptCode {
			return nil, apperrors.NewBusinessError(40013, "排班与科室不匹配")
		}
		if req.DoctorID != "" && schedule.DoctorID != "" && req.DoctorID != schedule.DoctorID {
			return nil, apperrors.NewBusinessError(40014, "排班与医生不匹配")
		}
		if req.DoctorID == "" {
			req.DoctorID, doctorName = schedule.DoctorID, schedule.DoctorName
		}
		price = schedule.RegistPrice
		if err := s.repo.IncrementScheduleRegNum(req.ScheduleID); err != nil {
			return nil, err
		}
	}
	active, err := s.repo.HasActiveRegistration(patient.PatientID, req.DeptCode, req.DoctorID)
	if err != nil {
		return nil, err
	}
	if active {
		return nil, apperrors.NewBusinessError(40902, "患者今日已存在相同科室和医生的有效挂号")
	}
	count, err := s.repo.GetTodayEncounterCount()
	if err != nil {
		return nil, err
	}
	encounterType, emergency, green := "OPD", 0, 0
	if req.IsEmergency {
		encounterType, emergency = "EMG", 1
	}
	if req.IsGreenChannel {
		green = 1
	}
	encounter := &EncounterModel{VisitNo: utils.GenerateNoWithSeq("E", time.Now().Format("20060102"), count+1), PatientID: patient.PatientID, PatientName: patient.Name, PatientAge: utils.CalculateAge(patient.BirthDate), Gender: patient.Gender, IDCardNo: patient.IDCardNo, Phone: patient.Phone, ChargeType: req.ChargeType, EncounterType: encounterType, DeptID: req.DeptCode, DeptName: deptName, DoctorID: req.DoctorID, DoctorName: doctorName, DoctorTitle: doctorTitle, IsEmergency: emergency, IsGreenChannel: green, Remark: req.Remark}
	id, err := s.repo.CreateEncounter(encounter)
	if err != nil {
		return nil, err
	}
	_ = operatorID
	_ = operatorName
	return &RegistrationResult{EncounterID: id, VisitNo: encounter.VisitNo, PatientID: patient.PatientID, PatientName: patient.Name, DeptName: deptName, DoctorName: doctorName, RegTime: time.Now().Format("2006-01-02 15:04:05"), RegistPrice: price, ChargeNo: utils.GenerateNo("EC"), ChargeStatus: "UNPAID", PaymentMethod: req.PayMethod}, nil
}

func (s *RegistrationService) GetEncounterList(req *GetEncounterListRequest) ([]*EncounterInfo, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	return s.repo.GetEncounterList(req)
}

func (s *RegistrationService) GetEncounterDetail(id int64) (*EncounterInfo, error) {
	return s.repo.GetEncounterByID(id)
}

func (s *RegistrationService) CancelEncounter(id int64, operatorID, operatorName, reason string) error {
	encounter, err := s.repo.GetEncounterByID(id)
	if err != nil {
		return err
	}
	if encounter.EncounterStatus == "CANCELLED" {
		return apperrors.NewBusinessError(40015, "该挂号已取消")
	}
	if encounter.EncounterStatus != "REGISTERED" {
		return apperrors.NewBusinessError(40016, "只能取消已挂号状态的记录")
	}
	if encounter.PaymentStatus == "PAID" {
		return apperrors.NewBusinessError(40017, "已收费，请先退费后再退号")
	}
	if err := s.repo.UpdateEncounterStatus(id, "CANCELLED"); err != nil {
		return err
	}
	// TODO: emergency_encounter 尚未保存 schedule_id，后续补充后在此释放排班号源。
	// TODO: 接入 sys_audit_log 后记录操作人、取消原因及操作时间。
	_ = operatorID
	_ = operatorName
	_ = reason
	return nil
}

func SearchPatients(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		var req SearchPatientRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		patients, total, err := svc.SearchPatients(&req)
		if err != nil {
			response.InternalError(c, "查询失败: "+err.Error())
			return
		}
		response.SuccessWithPage(c, patients, total, req.Page, req.PageSize)
	}
}

func GetPatientDetail(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		patientID := c.Query("id")
		if patientID == "" {
			response.BadRequest(c, "患者ID不能为空")
			return
		}
		patient, err := svc.GetPatientDetail(patientID)
		if err != nil {
			if be, ok := err.(*apperrors.BusinessError); ok {
				response.Error(c, 404, be.Code, be.Message)
				return
			}
			response.InternalError(c, "查询失败: "+err.Error())
			return
		}
		response.Success(c, patient)
	}
}

func GetDepts(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		depts, err := svc.GetDepts(c.Query("outpatient_only") == "true")
		if err != nil {
			response.InternalError(c, "查询失败: "+err.Error())
			return
		}
		response.Success(c, depts)
	}
}

func GetDoctors(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		doctors, err := svc.GetDoctors(c.Query("dept_code"))
		if err != nil {
			response.InternalError(c, "查询失败: "+err.Error())
			return
		}
		response.Success(c, doctors)
	}
}

func GetSchedules(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		var req GetScheduleRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		schedules, err := svc.GetSchedules(&req)
		if err != nil {
			if be, ok := err.(*apperrors.BusinessError); ok {
				response.Error(c, 400, be.Code, be.Message)
			} else {
				response.InternalError(c, "查询失败: "+err.Error())
			}
			return
		}
		response.Success(c, schedules)
	}
}

func CreateRegistration(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		var req CreateRegistrationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		operatorID, operatorName, _ := middleware.GetCurrentUser(c)
		result, err := svc.CreateRegistration(&req, operatorID, operatorName)
		if err != nil {
			if be, ok := err.(*apperrors.BusinessError); ok {
				response.Error(c, 400, be.Code, be.Message)
			} else {
				response.InternalError(c, "挂号失败: "+err.Error())
			}
			return
		}
		response.SuccessWithMessage(c, "挂号成功", result)
	}
}

func GetEncounterList(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		var req GetEncounterListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		encounters, total, err := svc.GetEncounterList(&req)
		if err != nil {
			response.InternalError(c, "查询失败: "+err.Error())
			return
		}
		response.SuccessWithPage(c, encounters, total, req.Page, req.PageSize)
	}
}

func GetEncounterDetail(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil || id <= 0 {
			response.BadRequest(c, "ID格式错误")
			return
		}
		encounter, err := svc.GetEncounterDetail(id)
		if err != nil {
			if be, ok := err.(*apperrors.BusinessError); ok {
				response.Error(c, 404, be.Code, be.Message)
			} else {
				response.InternalError(c, "查询失败: "+err.Error())
			}
			return
		}
		response.Success(c, encounter)
	}
}

func CancelEncounter(db *sql.DB) gin.HandlerFunc {
	svc := NewRegistrationService(db)
	return func(c *gin.Context) {
		var req struct {
			ID     int64  `json:"id"`
			Reason string `json:"reason"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		if req.ID <= 0 {
			response.BadRequest(c, "ID格式错误")
			return
		}
		operatorID, operatorName, _ := middleware.GetCurrentUser(c)
		if err := svc.CancelEncounter(req.ID, operatorID, operatorName, req.Reason); err != nil {
			if be, ok := err.(*apperrors.BusinessError); ok {
				response.Error(c, 400, be.Code, be.Message)
			} else {
				response.InternalError(c, "取消失败: "+err.Error())
			}
			return
		}
		response.SuccessWithMessage(c, "取消成功", nil)
	}
}
