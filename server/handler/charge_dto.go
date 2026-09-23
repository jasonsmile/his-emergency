package handler

import "time"

type GetPendingChargesRequest struct {
	EncounterID int64 `json:"encounter_id" binding:"required"`
}

type GetChargeRecordsRequest struct {
	EncounterID  int64  `form:"encounter_id"`
	PatientID    string `form:"patient_id"`
	ChargeStatus string `form:"charge_status"`
	StartDate    string `form:"start_date"`
	EndDate      string `form:"end_date"`
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
}

type CreateChargeRequest struct {
	EncounterID    int64               `json:"encounter_id" binding:"required"`
	Items          []ChargeItemRequest `json:"items" binding:"required,min=1"`
	PayMethod      string              `json:"pay_method" binding:"required"`
	DiscountAmount float64             `json:"discount_amount"`
	Remark         string              `json:"remark"`
}

type ChargeItemRequest struct {
	ItemType   string  `json:"item_type" binding:"required"`
	ItemID     int64   `json:"item_id"`
	ItemCode   string  `json:"item_code" binding:"required"`
	ItemName   string  `json:"item_name" binding:"required"`
	ItemSpec   string  `json:"item_spec"`
	Qty        float64 `json:"qty" binding:"required,min=0.01"`
	Unit       string  `json:"unit"`
	UnitPrice  float64 `json:"unit_price" binding:"min=0"`
	Amount     float64 `json:"amount" binding:"min=0"`
	MiAmount   float64 `json:"mi_amount" binding:"min=0"`
	SelfAmount float64 `json:"self_amount" binding:"min=0"`
}

type RefundChargeRequest struct {
	ChargeID     int64               `json:"charge_id" binding:"required"`
	RefundReason string              `json:"refund_reason" binding:"required"`
	RefundItems  []RefundItemRequest `json:"refund_items"`
}

type RefundItemRequest struct {
	ItemID    int64   `json:"item_id"`
	ItemCode  string  `json:"item_code"`
	RefundQty float64 `json:"refund_qty"`
}
type GetDailyReportRequest struct {
	Date      string `form:"date"`
	CashierID string `form:"cashier_id"`
}

type PendingChargeItem struct {
	ItemType     string  `json:"item_type"`
	ItemID       int64   `json:"item_id"`
	ItemCode     string  `json:"item_code"`
	ItemName     string  `json:"item_name"`
	ItemSpec     string  `json:"item_spec"`
	Qty          float64 `json:"qty"`
	Unit         string  `json:"unit"`
	UnitPrice    float64 `json:"unit_price"`
	Amount       float64 `json:"amount"`
	MiAmount     float64 `json:"mi_amount"`
	SelfAmount   float64 `json:"self_amount"`
	SourceType   string  `json:"source_type"`
	SourceNo     string  `json:"source_no"`
	SourceStatus string  `json:"source_status"`
}
type PendingChargesResponse struct {
	EncounterID int64                `json:"encounter_id"`
	VisitNo     string               `json:"visit_no"`
	PatientID   string               `json:"patient_id"`
	PatientName string               `json:"patient_name"`
	TotalAmount float64              `json:"total_amount"`
	MiAmount    float64              `json:"mi_amount"`
	SelfAmount  float64              `json:"self_amount"`
	Items       []*PendingChargeItem `json:"items"`
}
type ChargeResult struct {
	ChargeNo    string  `json:"charge_no"`
	EncounterID int64   `json:"encounter_id"`
	TotalAmount float64 `json:"total_amount"`
	MiAmount    float64 `json:"mi_amount"`
	SelfAmount  float64 `json:"self_amount"`
	PayMethod   string  `json:"pay_method"`
	PayTime     string  `json:"pay_time"`
	InvoiceNo   string  `json:"invoice_no"`
	CashierName string  `json:"cashier_name"`
}
type RefundResult struct {
	RefundChargeNo   string  `json:"refund_charge_no"`
	OriginalChargeNo string  `json:"original_charge_no"`
	RefundAmount     float64 `json:"refund_amount"`
	RefundTime       string  `json:"refund_time"`
	CashierName      string  `json:"cashier_name"`
}
type ChargeRecord struct {
	ID               int64      `json:"id"`
	ChargeNo         string     `json:"charge_no"`
	EncounterID      int64      `json:"encounter_id"`
	VisitNo          string     `json:"visit_no"`
	PatientID        string     `json:"patient_id"`
	PatientName      string     `json:"patient_name"`
	ItemType         string     `json:"item_type"`
	ItemCode         string     `json:"item_code"`
	ItemName         string     `json:"item_name"`
	ItemSpec         string     `json:"item_spec"`
	Qty              float64    `json:"qty"`
	Unit             string     `json:"unit"`
	UnitPrice        float64    `json:"unit_price"`
	Amount           float64    `json:"amount"`
	DiscountAmount   float64    `json:"discount_amount"`
	MiAmount         float64    `json:"mi_amount"`
	SelfAmount       float64    `json:"self_amount"`
	ChargeStatus     string     `json:"charge_status"`
	PayMethod        string     `json:"pay_method"`
	PayTime          *time.Time `json:"pay_time"`
	CashierID        string     `json:"cashier_id"`
	CashierName      string     `json:"cashier_name"`
	InvoiceNo        string     `json:"invoice_no"`
	IsRefund         int        `json:"is_refund"`
	OriginalChargeID *int64     `json:"original_charge_id"`
	Remark           string     `json:"remark"`
	CreatedAt        time.Time  `json:"created_at"`
}
type DailyReport struct {
	Date           string  `json:"date"`
	CashierID      string  `json:"cashier_id"`
	TotalCount     int64   `json:"total_count"`
	TotalAmount    float64 `json:"total_amount"`
	CashAmount     float64 `json:"cash_amount"`
	WechatAmount   float64 `json:"wechat_amount"`
	AlipayAmount   float64 `json:"alipay_amount"`
	CardAmount     float64 `json:"card_amount"`
	MiAmount       float64 `json:"mi_amount"`
	SelfAmount     float64 `json:"self_amount"`
	RefundCount    int64   `json:"refund_count"`
	RefundAmount   float64 `json:"refund_amount"`
	NetAmount      float64 `json:"net_amount"`
	InvoiceStartNo string  `json:"invoice_start_no"`
	InvoiceEndNo   string  `json:"invoice_end_no"`
}
type DailyReportDetail struct {
	Summary *DailyReport    `json:"summary"`
	Records []*ChargeRecord `json:"records"`
}
