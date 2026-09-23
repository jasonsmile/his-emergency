<template>
  <div class="charge-workbench" v-loading="charging">
    <el-row :gutter="20">
      <!-- 左侧：就诊查询 -->
      <el-col :span="10">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>就诊查询</span>
            </div>
          </template>
          
          <el-input
            v-model="searchKeyword"
            placeholder="患者姓名/就诊号"
            clearable
            @keyup.enter="startEncounterSearch"
          >
            <template #append>
              <el-button :icon="Search" aria-label="搜索就诊" @click="startEncounterSearch" />
            </template>
          </el-input>
          
          <el-table
            :data="encounterList"
            v-loading="encounterLoading"
            highlight-current-row
            @current-change="selectEncounter"
            style="margin-top: 15px"
          >
            <el-table-column prop="visit_no" label="就诊号" width="140" />
            <el-table-column prop="patient_name" label="姓名" width="80" />
            <el-table-column prop="dept_name" label="科室" />
            <el-table-column prop="payment_status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.payment_status)">
                  {{ getStatusText(row.payment_status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
          
          <el-pagination
            v-model:current-page="encounterPage"
            :total="encounterTotal"
            :page-size="20"
            layout="total, prev, pager, next"
            @current-change="searchEncounters"
            style="margin-top: 10px"
          />
        </el-card>
      </el-col>
      
      <!-- 右侧：收费操作 -->
      <el-col :span="14">
        <!-- 待收费项目 -->
        <el-card v-if="selectedEncounter">
          <template #header>
            <div class="card-header">
              <span>待收费项目 - {{ selectedEncounter.patient_name }}</span>
              <el-tag>就诊号: {{ selectedEncounter.visit_no }}</el-tag>
            </div>
          </template>
          
          <el-table :data="pendingItems" v-loading="pendingLoading" @selection-change="selectedItems = $event">
            <el-table-column type="selection" width="50" />
            <el-table-column prop="item_name" label="项目名称" />
            <el-table-column prop="item_spec" label="规格" width="120" />
            <el-table-column prop="qty" label="数量" width="80" />
            <el-table-column prop="unit_price" label="单价" width="100">
              <template #default="{ row }">
                ¥{{ row.unit_price }}
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">
                ¥{{ row.amount }}
              </template>
            </el-table-column>
            <el-table-column prop="source_type" label="来源" width="80" />
          </el-table>
          
          <div class="charge-summary">
            <div class="amount-row">
              <span>总金额：</span>
              <span class="amount">¥{{ totalAmount }}</span>
            </div>
            <div class="amount-row">
              <span>医保金额：</span>
              <span class="amount">¥{{ miAmount }}</span>
            </div>
            <div class="amount-row">
              <span>自付金额：</span>
              <span class="amount highlight">¥{{ selfAmount }}</span>
            </div>
          </div>
          
          <el-divider />
          
          <el-form inline>
            <el-form-item label="支付方式">
              <el-radio-group v-model="payMethod">
                <el-radio-button value="CASH">现金</el-radio-button>
                <el-radio-button value="WECHAT">微信</el-radio-button>
                <el-radio-button value="ALIPAY">支付宝</el-radio-button>
                <el-radio-button value="CARD">刷卡</el-radio-button>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :disabled="selectedItems.length === 0 || pendingLoading || charging"
                :loading="charging"
                @click="submitCharge"
              >
                确认收费
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
        
        <!-- 收费结果 -->
        <el-card v-if="chargeResult" style="margin-top: 15px">
          <template #header>
            <div class="card-header">
              <span>收费成功</span>
            </div>
          </template>
          
          <el-result
            icon="success"
            title="收费成功"
            :sub-title="`收费单号：${chargeResult.charge_no} | 发票号：${chargeResult.invoice_no}`"
          >
            <template #extra>
              <el-button type="primary" @click="printReceipt">打印收据</el-button>
              <el-button @click="resetCharge">继续收费</el-button>
            </template>
          </el-result>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getEncounterList } from '@/api/registration'
import { getPendingCharges, createCharge } from '@/api/charge'

// 就诊查询
const searchKeyword = ref('')
const encounterList = ref([])
const encounterLoading = ref(false)
const encounterPage = ref(1)
const encounterTotal = ref(0)
const selectedEncounter = ref(null)

// 待收费项目
const pendingItems = ref([])
const pendingLoading = ref(false)
const selectedItems = ref([])
let pendingRequest = 0
let encounterRequest = 0

// 支付
const payMethod = ref('CASH')
const charging = ref(false)
const chargeResult = ref(null)

// 计算金额
const totalAmount = computed(() => {
  return selectedItems.value.reduce((sum, item) => sum + Number(item.amount || 0), 0).toFixed(2)
})

const miAmount = computed(() => {
  return selectedItems.value.reduce((sum, item) => sum + Number(item.mi_amount || 0), 0).toFixed(2)
})

const selfAmount = computed(() => {
  return selectedItems.value.reduce((sum, item) => sum + Number(item.self_amount || 0), 0).toFixed(2)
})

// 状态显示
const getStatusType = (status) => {
  const map = { UNPAID: 'danger', PARTIAL: 'warning', PAID: 'success', REFUNDED: 'info' }
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = { UNPAID: '未支付', PARTIAL: '部分支付', PAID: '已支付', REFUNDED: '已退费' }
  return map[status] || status
}

// 搜索就诊
const searchEncounters = async () => {
  if (charging.value) return
  const requestId = ++encounterRequest
  resetCharge()
  encounterList.value = []
  encounterTotal.value = 0
  encounterLoading.value = true
  try {
    const res = await getEncounterList({
      keyword: searchKeyword.value,
      page: encounterPage.value,
      pageSize: 20
    })
    if (requestId !== encounterRequest) return
    encounterList.value = res.data || []
    encounterTotal.value = res.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    if (requestId === encounterRequest) encounterLoading.value = false
  }
}

const startEncounterSearch = () => {
  encounterPage.value = 1
  searchEncounters()
}

// 选择就诊
const selectEncounter = async (row) => {
  if (charging.value) return
  ++pendingRequest
  pendingItems.value = []
  selectedItems.value = []
  pendingLoading.value = false
  selectedEncounter.value = row
  chargeResult.value = null
  
  if (row) {
    await loadPendingCharges(row.id)
  } else {
    pendingItems.value = []
  }
}

// 加载待收费项目
const loadPendingCharges = async (encounterId) => {
  const requestId = ++pendingRequest
  pendingItems.value = []
  selectedItems.value = []
  pendingLoading.value = true
  try {
    const res = await getPendingCharges(encounterId)
    if (requestId === pendingRequest && selectedEncounter.value?.id === encounterId) {
      pendingItems.value = res.data.items || []
    }
  } catch (error) {
    console.error(error)
    if (requestId === pendingRequest) pendingItems.value = []
  } finally {
    if (requestId === pendingRequest) pendingLoading.value = false
  }
}

// 提交收费
const submitCharge = async () => {
  if (!selectedEncounter.value || charging.value || pendingLoading.value || !selectedItems.value.length) return
  const encounterId = selectedEncounter.value.id
  
  charging.value = true
  try {
    const res = await createCharge({
      encounter_id: encounterId,
      items: selectedItems.value.map(item => ({
        item_type: item.item_type,
        item_id: item.item_id,
        item_code: item.item_code,
        item_name: item.item_name,
        item_spec: item.item_spec,
        qty: item.qty,
        unit: item.unit,
        unit_price: item.unit_price,
        amount: item.amount,
        mi_amount: item.mi_amount || 0,
        self_amount: item.self_amount
      })),
      pay_method: payMethod.value
    })
    
    chargeResult.value = res.data
    ElMessage.success('收费成功')
    
    // 刷新待收费项目
    await loadPendingCharges(encounterId)
    const updated = await getEncounterList({ keyword: searchKeyword.value, page: encounterPage.value, pageSize: 20 })
    encounterList.value = updated.data || []
    encounterTotal.value = updated.total || 0
    
  } catch (error) {
    console.error(error)
  } finally {
    charging.value = false
  }
}

// 打印收据
const printReceipt = () => {
  ElMessage.info('打印功能开发中')
}

// 重置
const resetCharge = () => {
  if (charging.value) return
  ++pendingRequest
  selectedItems.value = []
  pendingLoading.value = false
  chargeResult.value = null
  selectedEncounter.value = null
  pendingItems.value = []
}

onMounted(() => {
  searchEncounters()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.charge-summary {
  margin-top: 15px;
  text-align: right;
}

.amount-row {
  margin: 5px 0;
  font-size: 14px;
}

.amount {
  font-weight: bold;
  font-size: 16px;
}

.amount.highlight {
  color: #f56c6c;
  font-size: 20px;
}
</style>
