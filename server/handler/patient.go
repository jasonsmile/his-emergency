package handler

import (
	"database/sql"
	"strings"

	"emergency-his/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type Patient struct {
	PatientID string `json:"patientId"`
	Name      string `json:"name"`
	Gender    string `json:"gender,omitempty"`
	BirthDate string `json:"birthDate,omitempty"`
	Phone     string `json:"phone,omitempty"`
	CardNo    string `json:"cardNo,omitempty"`
}

func ListPatients(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := strings.TrimSpace(c.Query("keyword"))
		const query = `SELECT patient_id, name, gender, DATE_FORMAT(birth_date, '%Y-%m-%d'), phone, card_no
			FROM sync_patient
			WHERE (? = '' OR patient_id = ? OR name LIKE CONCAT('%', ?, '%') OR phone = ?)
			ORDER BY id DESC LIMIT 100`
		rows, err := db.QueryContext(c.Request.Context(), query, keyword, keyword, keyword, keyword)
		if err != nil {
			response.InternalError(c, "查询患者失败")
			return
		}
		defer rows.Close()
		patients := make([]Patient, 0)
		for rows.Next() {
			var p Patient
			if err := rows.Scan(&p.PatientID, &p.Name, &p.Gender, &p.BirthDate, &p.Phone, &p.CardNo); err != nil {
				response.InternalError(c, "读取患者数据失败")
				return
			}
			patients = append(patients, p)
		}
		if err := rows.Err(); err != nil {
			response.InternalError(c, "读取患者数据失败")
			return
		}
		response.Success(c, patients)
	}
}
