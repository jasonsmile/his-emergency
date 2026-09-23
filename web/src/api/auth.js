import request from './request'

/**
 * POST /api/login
 * @param {{ login_name: string, password?: string }} data
 * @returns {Promise<{code: number, message: string, data: {
 *   token: string, user_id: string, user_name: string, dept_code: string
 * }}>}
 */
export function login({ login_name, password = '' }) {
  return request({
    url: '/login',
    method: 'post',
    data: { login_name, password },
    skipAuth: true,
  })
}
