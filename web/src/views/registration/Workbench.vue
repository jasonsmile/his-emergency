<template>
  <div class="registration-workbench" v-loading="submitting">
    <el-row :gutter="20">
      <!-- 左侧：患者搜索 -->
      <el-col :span="10">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>患者查询</span>
            </div>
          </template>
          
          <el-input
            v-model="searchKeyword"
            placeholder="姓名/身份证/手机号/就诊卡号"
            clearable
            @keyup.enter="startPatientSearch"
          >
            <template #append>
              <el-button :icon="Search" aria-label="搜索患者" @click="startPatientSearch" />
            </template>
          </el-input>
          
          <el-table
            :data="patientList"
            v-loading="patientLoading"
            highlight-current-row
            @current-change="selectPatient"
            style="margin-top: 15px"
          >
            <el-table-column prop="name" label="姓名" width="80" />
            <el-table-column prop="gender" label="性别" width="60" />
            <el-table-column prop="age" label="年龄" width="70" />
            <el-table-column prop="phone" label="手机号" width="120" />
            <el-table-column prop="charge_type" label="费别" />
          </el-table>
          
          <el-pagination
            v-model:current-page="patientPage"
            :total="patientTotal"
            :page-size="20"
            layout="total, prev, pager, next"
            @current-change="searchPatients"
            style="margin-top: 10px"
          />
        </el-card>
        
        <!-- 选中患者信息 -->
        <el-card v-if="selectedPatient" style="margin-top: 15px">
          <template #header>
            <div class="card-header">
              <span>患者信息</span>
            </div>
          </template>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="姓名">{{ selectedPatient.name }}</el-descriptions-item>
            <el-descriptions-item label="性别">{{ selectedPatient.gender }}</el-descriptions-item>
            <el-descriptions-item label="年龄">{{ selectedPatient.age }}</el-descriptions-item>
            <el-descriptions-item label="费别">{{ selectedPatient.charge_type }}</el-descriptions-item>
            <el-descriptions-item label="手机号" :span="2">{{ selectedPatient.phone }}</el-descriptions-item>
            <el-descriptions-item label="身份证号" :span="2">{{ selectedPatient.id_card_no }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      
      <!-- 右侧：挂号操作 -->
      <el-col :span="14">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>挂号信息</span>
            </div>
          </template>
          
          <el-form :model="regForm" label-width="100px">
            <el-form-item label="就诊日期">
              <el-date-picker
                v-model="regForm.clinicDate"
                type="date"
                :clearable="false"
                :disabled-date="disabledDate"
                @change="loadSchedules"
              />
            </el-form-item>
            
            <el-form-item label="科室">
              <el-select
                v-model="regForm.deptCode"
                placeholder="请选择科室"
                filterable
                @change="onDeptChange"
              >
                <el-option
                  v-for="dept in deptList"
                  :key="dept.dept_code"
                  :label="dept.dept_name"
                  :value="dept.dept_code"
                />
              </el-select>
            </el-form-item>
            
            <el-form-item label="午别">
              <el-radio-group v-model="regForm.timeDesc" @change="loadSchedules">
                <el-radio-button value="">全部</el-radio-button>
                <el-radio-button value="上午">上午</el-radio-button>
                <el-radio-button value="下午">下午</el-radio-button>
                <el-radio-button value="晚间">晚间</el-radio-button>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="医生">
              <el-select
                v-model="regForm.doctorId"
                placeholder="请选择医生（可不选）"
                clearable
                filterable
                @change="loadSchedules"
              >
                <el-option
                  v-for="doctor in doctorList"
                  :key="doctor.doctor_code"
                  :label="`${doctor.name} (${doctor.title})`"
                  :value="doctor.doctor_code"
                />
              </el-select>
            </el-form-item>
          </el-form>
        </el-card>
        
        <!-- 排班列表 -->
        <el-card v-if="scheduleList.length > 0" style="margin-top: 15px">
          <template #header>
            <div class="card-header">
              <span>可选号源</span>
            </div>
          </template>
          
          <el-table
            :data="scheduleList"
            highlight-current-row
            @current-change="selectSchedule"
          >
            <el-table-column prop="doctor_name" label="医生" width="100" />
            <el-table-column prop="clinic_label" label="号别" width="100" />
            <el-table-column prop="time_desc" label="午别" width="80" />
            <el-table-column prop="regist_price" label="挂号费" width="100">
              <template #default="{ row }">
                ¥{{ row.regist_price }}
              </template>
            </el-table-column>
            <el-table-column label="剩余号源">
              <template #default="{ row }">
                <el-tag :type="row.available_num > 0 ? 'success' : 'danger'">
                  {{ row.available_num }}/{{ row.registration_limits }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
        
        <!-- 挂号确认 -->
        <el-card v-if="selectedPatient && regForm.deptCode" style="margin-top: 15px">
          <template #header>
            <div class="card-header">
              <span>确认挂号</span>
            </div>
          </template>
          
          <el-form :model="confirmForm" label-width="100px">
            <el-form-item label="费别">
              <el-select v-model="confirmForm.chargeType" placeholder="请选择费别">
                <el-option label="自费" value="自费" />
                <el-option label="医疗保险" value="医疗保险" />
                <el-option label="公费医疗" value="公费医疗" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="拟用支付方式">
              <el-radio-group v-model="confirmForm.payMethod">
                <el-radio-button value="CASH">现金</el-radio-button>
                <el-radio-button value="WECHAT">微信</el-radio-button>
                <el-radio-button value="ALIPAY">支付宝</el-radio-button>
                <el-radio-button value="CARD">刷卡</el-radio-button>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="急诊">
              <el-switch v-model="confirmForm.isEmergency" />
            </el-form-item>
            
            <el-form-item label="绿色通道">
              <el-switch v-model="confirmForm.isGreenChannel" />
            </el-form-item>
            
            <el-form-item>
              <p>挂号成功后，请前往收费工作台缴费。</p>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :disabled="!canSubmit"
                :loading="submitting"
                @click="submitRegistration"
              >
                确认挂号
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import {
  searchPatients as searchPatientsApi,
  getDepts,
  getDoctors,
  getSchedules,
  createRegistration
} from '@/api/registration'

// 患者搜索
const searchKeyword = ref('')
const patientList = ref([])
const patientLoading = ref(false)
const patientPage = ref(1)
const patientTotal = ref(0)
const selectedPatient = ref(null)

// 科室/医生
const deptList = ref([])
const doctorList = ref([])

// 排班
const scheduleList = ref([])
const selectedSchedule = ref(null)
const scheduleLoading = ref(false)
let patientRequest = 0
let scheduleRequest = 0
let doctorRequest = 0

// 挂号表单
const regForm = reactive({
  clinicDate: new Date(),
  deptCode: '',
  timeDesc: '',
  doctorId: ''
})

// 确认表单
const confirmForm = reactive({
  chargeType: '自费',
  payMethod: 'CASH',
  isEmergency: false,
  isGreenChannel: false
})

const submitting = ref(false)

// 计算是否可以提交
const canSubmit = computed(() => {
  return !!(selectedPatient.value && regForm.deptCode && !scheduleLoading.value &&
    !submitting.value && formatDate(regForm.clinicDate) === formatDate(new Date()))
})

// 后端按当天创建就诊记录，暂不支持保存未来预约日期。
const disabledDate = (date) => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return date.getTime() !== today.getTime()
}

// 搜索患者
const searchPatients = async () => {
  if (submitting.value) return
  const requestId = ++patientRequest
  selectedPatient.value = null
  patientList.value = []
  patientTotal.value = 0
  patientLoading.value = true
  try {
    const res = await searchPatientsApi({
      keyword: searchKeyword.value,
      page: patientPage.value,
      pageSize: 20
    })
    if (requestId !== patientRequest) return
    patientList.value = res.data || []
    patientTotal.value = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    if (requestId === patientRequest) patientLoading.value = false
  }
}

const startPatientSearch = () => {
  patientPage.value = 1
  searchPatients()
}

// 选择患者
const selectPatient = (row) => {
  if (submitting.value) return
  selectedPatient.value = row
  if (row) {
    confirmForm.chargeType = row.charge_type || '自费'
  }
}

// 加载科室
const loadDepts = async () => {
  try {
    const res = await getDepts({ outpatient_only: true })
    deptList.value = res.data
  } catch (error) {
    console.error(error)
  }
}

// 科室变化
const onDeptChange = async (deptCode) => {
  const requestId = ++doctorRequest
  regForm.doctorId = ''
  selectedSchedule.value = null
  doctorList.value = []
  loadSchedules()
  if (deptCode) {
    try {
      const res = await getDoctors({ dept_code: deptCode })
      if (requestId === doctorRequest) doctorList.value = res.data || []
    } catch (error) {
      console.error(error)
    }
  } else {
    doctorList.value = []
  }
}

// 加载排班
const loadSchedules = async () => {
  const requestId = ++scheduleRequest
  selectedSchedule.value = null
  scheduleList.value = []
  scheduleLoading.value = false
  if (!regForm.clinicDate || !regForm.deptCode) {
    scheduleList.value = []
    return
  }
  
  scheduleLoading.value = true
  try {
    const res = await getSchedules({
      clinic_date: formatDate(regForm.clinicDate),
      dept_code: regForm.deptCode,
      doctor_id: regForm.doctorId || undefined,
      time_desc: regForm.timeDesc
    })
    if (requestId === scheduleRequest) scheduleList.value = res.data || []
  } catch (error) {
    console.error(error)
  } finally {
    if (requestId === scheduleRequest) scheduleLoading.value = false
  }
}

// 选择排班
const selectSchedule = (row) => {
  if (submitting.value) return
  if (row && (row.available_num <= 0 || row.states !== '正常')) {
    selectedSchedule.value = null
    ElMessage.warning('该排班已满或已停诊，请选择其他号源')
    return
  }
  selectedSchedule.value = row
  if (row && row.doctor_id) {
    regForm.doctorId = row.doctor_id
  }
}

// 格式化日期
const formatDate = (date) => {
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

// 提交挂号
const submitRegistration = async () => {
  if (submitting.value) return
  if (!selectedPatient.value) {
    ElMessage.warning('请先选择患者')
    return
  }
  if (!canSubmit.value) {
    ElMessage.warning('请选择科室和当天就诊日期，并等待排班加载完成')
    return
  }
  
  submitting.value = true
  try {
    const res = await createRegistration({
      patient_id: selectedPatient.value.patient_id,
      dept_code: regForm.deptCode,
      doctor_id: regForm.doctorId || undefined,
      schedule_id: selectedSchedule.value?.id,
      charge_type: confirmForm.chargeType,
      is_emergency: confirmForm.isEmergency,
      is_green_channel: confirmForm.isGreenChannel,
      pay_method: confirmForm.payMethod
    })
    
    ElMessage.success(`挂号成功！就诊号：${res.data.visit_no}`)
    
    // 重置表单
    selectedPatient.value = null
    regForm.deptCode = ''
    regForm.doctorId = ''
    doctorList.value = []
    scheduleList.value = []
    selectedSchedule.value = null
    
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadDepts()
  searchPatients()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
