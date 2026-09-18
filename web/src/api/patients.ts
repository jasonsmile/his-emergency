import axios from 'axios'

export interface Patient {
  patientId: string
  name: string
  gender?: string
  birthDate?: string
  phone?: string
  cardNo?: string
}

interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}

export async function listPatients(keyword: string): Promise<Patient[]> {
  const { data } = await axios.get<ApiResponse<Patient[]>>('/api/v1/patients', {
    params: { keyword: keyword.trim() },
    timeout: 15000,
  })
  if (data.code !== 0) throw new Error(data.message || '查询患者失败')
  if (!Array.isArray(data.data)) throw new Error('接口返回的数据格式不正确')
  return data.data
}

export function getRequestError(error: unknown): string {
  if (axios.isAxiosError(error)) {
    if (error.code === 'ECONNABORTED') return '请求超时，请稍后重试。'
    if (!error.response) return '无法连接服务，请检查网络或后端服务是否已启动。'
    const message = error.response.data?.message
    return typeof message === 'string' ? message : `查询失败（HTTP ${error.response.status}），请稍后重试。`
  }
  return error instanceof Error ? error.message : '查询失败，请稍后重试。'
}
