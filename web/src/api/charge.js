import request from './request'

// 后端使用 POST，JSON 请求体为 { encounter_id: number }。
export function getPendingCharges(encounterId) {
  return request({
    url: '/charge/getPendingCharges',
    method: 'post',
    data: { encounter_id: Number(encounterId) },
  })
}

/**
 * @typedef {Object} ChargeItem
 * @property {string} item_type
 * @property {number} [item_id]
 * @property {string} item_code
 * @property {string} item_name
 * @property {string} [item_spec]
 * @property {number} qty
 * @property {string} [unit]
 * @property {number} unit_price
 * @property {number} amount
 * @property {number} [mi_amount]
 * @property {number} [self_amount]
 */

/**
 * 创建收费，字段与后端 JSON 协议一致。
 * @param {{ encounter_id: number, items: ChargeItem[],
 * pay_method: 'CASH'|'WECHAT'|'ALIPAY'|'CARD',
 * discount_amount?: number, remark?: string }} data
 */
export function createCharge(data) {
  return request({
    url: '/charge/createCharge',
    method: 'post',
    data,
  })
}

/**
 * charge_id 为收费记录的数字 id；不传 refund_items 时全额退该条记录。
 * @param {{ charge_id: number, refund_reason: string,
 * refund_items?: Array<{ item_id?: number, item_code?: string, refund_qty: number }> }} data
 */
export function refundCharge(data) {
  return request({
    url: '/charge/refundCharge',
    method: 'post',
    data,
  })
}

// params: encounter_id、patient_id、charge_status、start_date、end_date、page、page_size。
// 注意：收费分页查询使用 page_size，与挂号接口的 pageSize 不同。
export function getChargeRecords(params) {
  return request({
    url: '/charge/getChargeRecords',
    method: 'get',
    params,
  })
}

// params: date（YYYY-MM-DD，默认当天）、cashier_id。
// 响应 data 为 { summary, records }。
export function getDailyReport(params) {
  return request({
    url: '/charge/getDailyReport',
    method: 'get',
    params,
  })
}
