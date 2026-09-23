package handler

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"emergency-his/server/logger"
	"emergency-his/server/middleware"
	"emergency-his/server/response"
	"emergency-his/server/utils"

	"github.com/gin-gonic/gin"
)

const chargeColumns = `id, charge_no, encounter_id, patient_id, patient_name, item_type, item_code, item_name, item_spec, qty, unit, unit_price, amount, discount_amount, mi_amount, self_amount, charge_status, pay_method, pay_time, cashier_id, cashier_name, invoice_no, is_refund, original_charge_id, remark, created_at`

// 以下辅助函数统一记录收费模块 SQL 模板与绑定参数；不拼接参数，避免改变预编译执行方式。
func chargeQueryRow(c *gin.Context, db *sql.DB, statement string, args ...any) *sql.Row {
	logger.SQL(statement, args...)
	return db.QueryRowContext(c, statement, args...)
}

func chargeQuery(c *gin.Context, db *sql.DB, statement string, args ...any) (*sql.Rows, error) {
	logger.SQL(statement, args...)
	return db.QueryContext(c, statement, args...)
}

func chargeTxQueryRow(c *gin.Context, tx *sql.Tx, statement string, args ...any) *sql.Row {
	logger.SQL(statement, args...)
	return tx.QueryRowContext(c, statement, args...)
}

func chargeTxExec(c *gin.Context, tx *sql.Tx, statement string, args ...any) (sql.Result, error) {
	logger.SQL(statement, args...)
	result, err := tx.ExecContext(c, statement, args...)
	if err != nil {
		logger.Error.Printf("sql execution failed: %v", err)
	}
	return result, err
}

func GetPendingCharges(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GetPendingChargesRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		id := req.EncounterID
		var result PendingChargesResponse
		err := chargeQueryRow(c, db, `SELECT id, visit_no, patient_id, patient_name FROM emergency_encounter WHERE id=?`, id).Scan(&result.EncounterID, &result.VisitNo, &result.PatientID, &result.PatientName)
		if err == sql.ErrNoRows {
			response.NotFound(c, "就诊记录不存在")
			return
		}
		if err != nil {
			response.InternalError(c, "查询待收费项目失败")
			return
		}
		result.Items = make([]*PendingChargeItem, 0)
		var regPaid int
		if err = chargeQueryRow(c, db, `SELECT COUNT(*) FROM emergency_charge WHERE encounter_id=? 
                                        AND item_type='REG' AND charge_status='PAID'`, id).Scan(&regPaid); err != nil {
			response.InternalError(c, "查询待收费项目失败")
			return
		}
		if regPaid == 0 {
			addPending(&result, &PendingChargeItem{ItemType: "REG", ItemCode: "REG_FEE", ItemName: "挂号费", Qty: 1, Unit: "次", UnitPrice: 10, Amount: 10, SelfAmount: 10, SourceType: "挂号", SourceNo: result.VisitNo, SourceStatus: "REGISTERED"})
		}
		rows, err := chargeQuery(c, db, `SELECT i.id,p.prescription_no,p.status,i.item_class,i.drug_code,i.drug_name,i.drug_spec,i.total_qty,i.total_unit,i.unit_price,i.amount
										FROM emergency_prescription p JOIN emergency_prescription_item i ON i.prescription_id=p.id WHERE p.encounter_id=? AND p.status='SIGNED' AND 
										NOT EXISTS (SELECT 1 FROM emergency_charge ec WHERE ec.encounter_id=p.encounter_id AND ec.item_id=i.id AND ec.charge_status='PAID') ORDER BY p.id,i.item_no`, id)
		if err != nil {
			response.InternalError(c, "查询待收费项目失败")
			return
		}
		defer rows.Close()
		for rows.Next() {
			item := &PendingChargeItem{}
			var class string
			if err := rows.Scan(&item.ItemID, &item.SourceNo, &item.SourceStatus, &class, &item.ItemCode, &item.ItemName, &item.ItemSpec, &item.Qty, &item.Unit, &item.UnitPrice, &item.Amount); err != nil {
				response.InternalError(c, "读取待收费项目失败")
				return
			}
			item.ItemType = "DRUG"
			if class != "A" {
				item.ItemType = "ITEM"
			}
			item.SelfAmount = item.Amount
			item.SourceType = "处方"
			addPending(&result, item)
		}
		if err = rows.Err(); err != nil {
			response.InternalError(c, "读取待收费项目失败")
			return
		}
		response.Success(c, &result)
	}
}

func addPending(r *PendingChargesResponse, item *PendingChargeItem) {
	r.Items = append(r.Items, item)
	r.TotalAmount += item.Amount
	r.MiAmount += item.MiAmount
	r.SelfAmount += item.SelfAmount
}

func CreateCharge(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateChargeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		if !validPayMethod(req.PayMethod) || req.DiscountAmount < 0 {
			response.BadRequest(c, "支付方式或优惠金额不合法")
			return
		}
		if total := sumItems(req.Items); total <= 0 || req.DiscountAmount > total {
			response.BadRequest(c, "优惠金额不能超过收费总额")
			return
		}
		userID, userName, _ := middleware.GetCurrentUser(c)
		now := time.Now()
		tx, err := db.BeginTx(c, nil)
		if err != nil {
			response.InternalError(c, "收费失败")
			return
		}
		defer tx.Rollback()
		var patientID, patientName string
		if err = chargeTxQueryRow(c, tx, `SELECT patient_id,patient_name FROM emergency_encounter WHERE id=? FOR UPDATE`, req.EncounterID).Scan(&patientID, &patientName); err == sql.ErrNoRows {
			response.NotFound(c, "就诊记录不存在")
			return
		}
		if err != nil {
			response.InternalError(c, "收费失败")
			return
		}
		invoiceNo := utils.GenerateNo("EI")
		chargeNo := ""
		var total, mi, self float64
		for _, item := range req.Items {
			if err = validateChargeItem(item); err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			if item.ItemType == "REG" {
				var paidCount int
				err = chargeTxQueryRow(c, tx, `SELECT COUNT(*) FROM emergency_charge WHERE encounter_id=? AND item_type='REG' AND charge_status='PAID' FOR UPDATE`, req.EncounterID).Scan(&paidCount)
				if err != nil {
					response.InternalError(c, "校验挂号费失败")
					return
				}
				if paidCount > 0 {
					response.BadRequest(c, "挂号费已收费")
					return
				}
			} else {
				var sourceCode, sourceName string
				var sourceQty, sourcePrice, sourceAmount float64
				err = chargeTxQueryRow(c, tx, `SELECT i.drug_code,i.drug_name,i.total_qty,i.unit_price,i.amount FROM emergency_prescription_item i JOIN emergency_prescription p ON p.id=i.prescription_id WHERE i.id=? AND p.encounter_id=? AND p.status='SIGNED' AND NOT EXISTS (SELECT 1 FROM emergency_charge ec WHERE ec.item_id=i.id AND ec.charge_status='PAID')`, item.ItemID, req.EncounterID).Scan(&sourceCode, &sourceName, &sourceQty, &sourcePrice, &sourceAmount)
				if err == nil && (item.ItemCode != sourceCode || item.ItemName != sourceName || math.Abs(item.Qty-sourceQty) > 0.001 || math.Abs(item.UnitPrice-sourcePrice) > 0.0001 || math.Abs(item.Amount-sourceAmount) > 0.01) {
					response.BadRequest(c, "收费项目金额与处方不一致")
					return
				}
				if err == sql.ErrNoRows {
					response.BadRequest(c, "收费项目已收费或不存在")
					return
				}
				if err != nil {
					response.InternalError(c, "校验收费项目失败")
					return
				}
			}
			discount := 0.0
			if len(req.Items) > 0 {
				discount = req.DiscountAmount * item.Amount / sumItems(req.Items)
			}
			lineChargeNo := utils.GenerateNo("EC")
			if chargeNo == "" {
				chargeNo = lineChargeNo
			}
			_, err = chargeTxExec(c, tx, `INSERT INTO emergency_charge (charge_no,encounter_id,patient_id,patient_name,item_type,item_id,item_code,item_name,item_spec,qty,unit,unit_price,amount,discount_amount,mi_amount,self_amount,charge_status,pay_method,pay_time,cashier_id,cashier_name,invoice_no,invoice_prefix,remark) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,'PAID',?,?,?,?,?,?,?)`, lineChargeNo, req.EncounterID, patientID, patientName, item.ItemType, nullID(item.ItemID), item.ItemCode, item.ItemName, item.ItemSpec, item.Qty, item.Unit, item.UnitPrice, item.Amount, discount, item.MiAmount, item.SelfAmount, req.PayMethod, now, userID, userName, invoiceNo, "EI", req.Remark)
			if err != nil {
				response.InternalError(c, "保存收费记录失败")
				return
			}
			total += item.Amount
			mi += item.MiAmount
			self += item.SelfAmount - discount
		}
		if _, err = chargeTxExec(c, tx, `UPDATE emergency_encounter SET total_amount=total_amount+?, paid_amount=paid_amount+?, payment_status='PAID' WHERE id=?`, total, total-req.DiscountAmount, req.EncounterID); err != nil {
			response.InternalError(c, "更新就诊支付状态失败")
			return
		}
		if err = tx.Commit(); err != nil {
			response.InternalError(c, "收费失败")
			return
		}
		response.SuccessWithMessage(c, "收费成功", &ChargeResult{ChargeNo: chargeNo, EncounterID: req.EncounterID, TotalAmount: total, MiAmount: mi, SelfAmount: self, PayMethod: req.PayMethod, PayTime: now.Format("2006-01-02 15:04:05"), InvoiceNo: invoiceNo, CashierName: userName})
	}
}

func RefundCharge(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefundChargeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		userID, userName, _ := middleware.GetCurrentUser(c)
		tx, err := db.BeginTx(c, nil)
		if err != nil {
			response.InternalError(c, "退费失败")
			return
		}
		defer tx.Rollback()
		var chargeNo, patientID, patientName, itemType, itemCode, itemName, itemSpec, payMethod string
		var encounterID, itemID int64
		var qty, unitPrice, amount, mi, self float64
		var unit string
		var status string
		err = chargeTxQueryRow(c, tx, `SELECT charge_no,encounter_id,patient_id,patient_name,item_type,COALESCE(item_id,0),item_code,item_name,item_spec,qty,unit,unit_price,amount,mi_amount,self_amount,pay_method,charge_status FROM emergency_charge WHERE id=? FOR UPDATE`, req.ChargeID).Scan(&chargeNo, &encounterID, &patientID, &patientName, &itemType, &itemID, &itemCode, &itemName, &itemSpec, &qty, &unit, &unitPrice, &amount, &mi, &self, &payMethod, &status)
		if err == sql.ErrNoRows {
			response.NotFound(c, "收费记录不存在")
			return
		}
		if err != nil {
			response.InternalError(c, "退费失败")
			return
		}
		if status != "PAID" {
			response.BadRequest(c, "该收费记录不能退费")
			return
		}
		refundQty := qty
		if len(req.RefundItems) > 0 {
			found := false
			for _, ri := range req.RefundItems {
				if (ri.ItemID == 0 || ri.ItemID == itemID) && (ri.ItemCode == "" || ri.ItemCode == itemCode) {
					refundQty = ri.RefundQty
					found = true
					break
				}
			}
			if !found || refundQty <= 0 || refundQty > qty {
				response.BadRequest(c, "退费项目或数量不合法")
				return
			}
		}
		ratio := refundQty / qty
		refundAmount := round2(amount * ratio)
		refundMi := round2(mi * ratio)
		refundSelf := round2(self * ratio)
		now := time.Now()
		refundNo := utils.GenerateNo("ER")
		_, err = chargeTxExec(c, tx, `INSERT INTO emergency_charge (charge_no,encounter_id,patient_id,patient_name,item_type,item_id,item_code,item_name,item_spec,qty,unit,unit_price,amount,mi_amount,self_amount,charge_status,pay_method,pay_time,cashier_id,cashier_name,is_refund,original_charge_id,remark) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,'REFUNDED',?,?,?, ?,1,?,?)`, refundNo, encounterID, patientID, patientName, itemType, nullID(itemID), itemCode, itemName, itemSpec, refundQty, unit, unitPrice, refundAmount, refundMi, refundSelf, payMethod, now, userID, userName, req.ChargeID, req.RefundReason)
		if err != nil {
			response.InternalError(c, "保存退费记录失败")
			return
		}
		if refundQty == qty {
			_, err = chargeTxExec(c, tx, `UPDATE emergency_charge SET charge_status='REFUNDED', is_refund=1 WHERE id=?`, req.ChargeID)
		} else {
			_, err = chargeTxExec(c, tx, `UPDATE emergency_charge SET qty=qty-?, amount=amount-?, mi_amount=mi_amount-?, self_amount=self_amount-? WHERE id=?`, refundQty, refundAmount, refundMi, refundSelf, req.ChargeID)
		}
		if err != nil {
			response.InternalError(c, "更新原收费记录失败")
			return
		}
		_, err = chargeTxExec(c, tx, `UPDATE emergency_encounter SET paid_amount=GREATEST(0,paid_amount-?), payment_status=CASE WHEN paid_amount-?<=0 THEN 'REFUNDED' ELSE 'PARTIAL' END WHERE id=?`, refundAmount, refundAmount, encounterID)
		if err != nil {
			response.InternalError(c, "更新就诊支付状态失败")
			return
		}
		if err = tx.Commit(); err != nil {
			response.InternalError(c, "退费失败")
			return
		}
		response.SuccessWithMessage(c, "退费成功", &RefundResult{RefundChargeNo: refundNo, OriginalChargeNo: chargeNo, RefundAmount: refundAmount, RefundTime: now.Format("2006-01-02 15:04:05"), CashierName: userName})
	}
}

func GetChargeRecords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GetChargeRecordsRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		if req.Page < 1 {
			req.Page = 1
		}
		if req.PageSize < 1 || req.PageSize > 100 {
			req.PageSize = 20
		}
		where, args := chargeWhere(req)
		var total int64
		if err := chargeQueryRow(c, db, "SELECT COUNT(*) FROM emergency_charge"+where, args...).Scan(&total); err != nil {
			response.InternalError(c, "查询失败")
			return
		}
		args = append(args, req.PageSize, (req.Page-1)*req.PageSize)
		rows, err := chargeQuery(c, db, "SELECT "+chargeColumns+" FROM emergency_charge"+where+" ORDER BY created_at DESC LIMIT ? OFFSET ?", args...)
		if err != nil {
			response.InternalError(c, "查询失败")
			return
		}
		defer rows.Close()
		records := make([]*ChargeRecord, 0)
		for rows.Next() {
			r := &ChargeRecord{}
			var pay sql.NullTime
			var original sql.NullInt64
			if err := rows.Scan(&r.ID, &r.ChargeNo, &r.EncounterID, &r.PatientID, &r.PatientName, &r.ItemType, &r.ItemCode, &r.ItemName, &r.ItemSpec, &r.Qty, &r.Unit, &r.UnitPrice, &r.Amount, &r.DiscountAmount, &r.MiAmount, &r.SelfAmount, &r.ChargeStatus, &r.PayMethod, &pay, &r.CashierID, &r.CashierName, &r.InvoiceNo, &r.IsRefund, &original, &r.Remark, &r.CreatedAt); err != nil {
				response.InternalError(c, "读取收费记录失败")
				return
			}
			if pay.Valid {
				r.PayTime = &pay.Time
			}
			if original.Valid {
				r.OriginalChargeID = &original.Int64
			}
			records = append(records, r)
		}
		response.SuccessWithPage(c, records, total, req.Page, req.PageSize)
	}
}

func GetDailyReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GetDailyReportRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}
		if req.Date == "" {
			req.Date = time.Now().Format("2006-01-02")
		}
		report := &DailyReport{Date: req.Date, CashierID: req.CashierID}
		where := " WHERE DATE(pay_time)=?"
		args := []any{req.Date}
		if req.CashierID != "" {
			where += " AND cashier_id=?"
			args = append(args, req.CashierID)
		}
		q := `SELECT COALESCE(SUM(is_refund=0 AND charge_status='PAID'),0),COALESCE(SUM(CASE WHEN is_refund=0 THEN amount ELSE 0 END),0),COALESCE(SUM(CASE WHEN is_refund=0 AND pay_method='CASH' THEN amount ELSE 0 END),0),COALESCE(SUM(CASE WHEN is_refund=0 AND pay_method='WECHAT' THEN amount ELSE 0 END),0),COALESCE(SUM(CASE WHEN is_refund=0 AND pay_method='ALIPAY' THEN amount ELSE 0 END),0),COALESCE(SUM(CASE WHEN is_refund=0 AND pay_method='CARD' THEN amount ELSE 0 END),0),COALESCE(SUM(CASE WHEN is_refund=0 THEN mi_amount ELSE 0 END),0),COALESCE(SUM(CASE WHEN is_refund=0 THEN self_amount ELSE 0 END),0),COALESCE(SUM(is_refund=1),0),COALESCE(SUM(CASE WHEN is_refund=1 THEN amount ELSE 0 END),0),COALESCE(MIN(NULLIF(invoice_no,'')),''),COALESCE(MAX(NULLIF(invoice_no,'')),'') FROM emergency_charge` + where
		if err := chargeQueryRow(c, db, q, args...).Scan(&report.TotalCount, &report.TotalAmount, &report.CashAmount, &report.WechatAmount, &report.AlipayAmount, &report.CardAmount, &report.MiAmount, &report.SelfAmount, &report.RefundCount, &report.RefundAmount, &report.InvoiceStartNo, &report.InvoiceEndNo); err != nil {
			response.InternalError(c, "查询日结报表失败")
			return
		}
		report.NetAmount = report.TotalAmount - report.RefundAmount
		response.Success(c, &DailyReportDetail{Summary: report, Records: []*ChargeRecord{}})
	}
}

func validPayMethod(v string) bool {
	switch v {
	case "CASH", "WECHAT", "ALIPAY", "CARD":
		return true
	}
	return false
}
func validateChargeItem(i ChargeItemRequest) error {
	if i.ItemType != "REG" && i.ItemType != "DRUG" && i.ItemType != "ITEM" {
		return fmt.Errorf("项目类型不合法")
	}
	if i.ItemType != "REG" && i.ItemID <= 0 {
		return fmt.Errorf("收费项目ID不合法")
	}
	if math.Abs(i.Amount-(i.MiAmount+i.SelfAmount)) > 0.01 {
		return fmt.Errorf("项目金额与医保、自付金额不一致")
	}
	if math.Abs(i.Amount-i.Qty*i.UnitPrice) > 0.02 {
		return fmt.Errorf("项目金额与数量、单价不一致")
	}
	return nil
}
func sumItems(items []ChargeItemRequest) float64 {
	var total float64
	for _, i := range items {
		total += i.Amount
	}
	return total
}
func nullID(id int64) any {
	if id == 0 {
		return nil
	}
	return id
}
func round2(v float64) float64 { return math.Round(v*100) / 100 }
func chargeWhere(req GetChargeRecordsRequest) (string, []any) {
	where := " WHERE 1=1"
	args := []any{}
	if req.EncounterID > 0 {
		where += " AND encounter_id=?"
		args = append(args, req.EncounterID)
	}
	if req.PatientID != "" {
		where += " AND patient_id=?"
		args = append(args, req.PatientID)
	}
	if req.ChargeStatus != "" {
		where += " AND charge_status=?"
		args = append(args, strings.ToUpper(req.ChargeStatus))
	}
	if req.StartDate != "" {
		where += " AND DATE(created_at)>=?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		where += " AND DATE(created_at)<=?"
		args = append(args, req.EndDate)
	}
	return where, args
}
