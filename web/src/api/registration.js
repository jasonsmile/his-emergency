import request from './request'

// params: keyword、page、pageSize（后端分页查询使用驼峰 pageSize）。
export function searchPatients(params) {
  return request({
    url: '/registration/searchPatients',
    method: 'get',
    params,
  })
}

// id 传患者的 patient_id，而非数据库记录的数字 id。
export function getPatientDetail(id) {
  return request({
    url: '/registration/getPatientDetail',
    method: 'get',
    params: { id },
  })
}

// params: outpatient_only（布尔值）。
export function getDepts(params) {
  return request({
    url: '/registration/getDepts',
    method: 'get',
    params,
  })
}

// params: dept_code。
export function getDoctors(params) {
  return request({
    url: '/registration/getDoctors',
    method: 'get',
    params,
  })
}

// params: clinic_date、dept_code、doctor_id、time_desc。
export function getSchedules(params) {
  return request({
    url: '/registration/getSchedules',
    method: 'get',
    params,
  })
}

/**
 * @param {{ patient_id: string, dept_code: string, doctor_id?: string,
 * schedule_id?: number, charge_type?: string, is_emergency?: boolean,
 * is_green_channel?: boolean, pay_method?: string, remark?: string }} data
 */
export function createRegistration(data) {
  return request({
    url: '/registration/createRegistration',
    method: 'post',
    data,
  })
}

// params: status、dept_code、doctor_id、start_date、end_date、keyword、page、pageSize。
export function getEncounterList(params) {
  return request({
    url: '/registration/getEncounterList',
    method: 'get',
    params,
  })
}

// id 传就诊记录的数字 id。
export function getEncounterDetail(id) {
  return request({
    url: '/registration/getEncounterDetail',
    method: 'get',
    params: { id },
  })
}

// 保留 cancelEncounter(id, { reason }) 调用方式；后端接收 JSON { id, reason }。
export function cancelEncounter(id, data = {}) {
  return request({
    url: '/registration/cancelEncounter',
    method: 'post',
    data: { id: Number(id), reason: data.reason ?? '' },
  })
}
