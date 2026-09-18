<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getRequestError, listPatients, type Patient } from '../api/patients'

const keyword = ref('')
const patients = ref<Patient[]>([])
const loading = ref(false)
const errorMessage = ref('')
const hasLoaded = ref(false)
const updatedAt = ref('')
const selectedPatient = ref<Patient | null>(null)
const detailVisible = ref(false)

async function search() {
  if (loading.value) return
  loading.value = true
  errorMessage.value = ''
  patients.value = []
  updatedAt.value = ''
  try {
    patients.value = await listPatients(keyword.value)
    hasLoaded.value = true
    updatedAt.value = new Date().toLocaleTimeString('zh-CN', { hour12: false })
  } catch (error) {
    errorMessage.value = getRequestError(error)
    hasLoaded.value = false
  } finally {
    loading.value = false
  }
}

function reset() {
  keyword.value = ''
  void search()
}

function showDetails(patient: Patient) {
  selectedPatient.value = patient
  detailVisible.value = true
}

function display(value?: string) { return value?.trim() || '—' }

onMounted(search)
</script>

<template>
  <section class="patient-page" aria-labelledby="page-title">
    <div class="page-heading"><div><h1 id="page-title">查询患者</h1><p>快速检索患者，查看基础档案信息。</p></div><span class="page-badge">患者服务</span></div>
    <section class="surface search-card" aria-labelledby="search-title">
      <h2 id="search-title">患者检索</h2>
      <form class="search-form" @submit.prevent="search">
        <div class="search-field"><label for="patient-keyword">查询条件</label><el-input id="patient-keyword" v-model="keyword" placeholder="请输入姓名、患者编号或手机号" clearable :disabled="loading" /></div>
        <div class="search-actions"><el-button type="primary" native-type="submit" :loading="loading">查询</el-button><el-button :disabled="loading" @click="reset">重置</el-button></div>
      </form>
      <p class="search-hint">姓名支持模糊查询，患者编号和手机号需完整输入；留空可查询最近入库的患者记录。</p>
    </section>
    <section class="surface results-card" aria-labelledby="results-title" :aria-busy="loading">
      <div class="results-heading"><div class="results-title"><h2 id="results-title">患者列表</h2><span v-if="hasLoaded && !loading" class="count-badge">{{ patients.length }} 条</span></div><el-button :disabled="loading" @click="search">刷新列表</el-button></div>
      <div v-if="errorMessage" class="error-panel" role="alert"><el-alert title="未能获取患者列表" :description="errorMessage" type="error" show-icon :closable="false" /><el-button type="primary" plain @click="search">重新查询</el-button></div>
      <el-table v-loading="loading" :data="patients" row-key="patientId" class="patient-table" :height="440" stripe>
        <el-table-column type="index" label="序号" width="70" align="center" />
        <el-table-column prop="patientId" label="患者编号" min-width="160"><template #default="{ row }"><span class="patient-id">{{ row.patientId }}</span></template></el-table-column>
        <el-table-column prop="name" label="姓名" min-width="120"><template #default="{ row }"><span class="patient-name">{{ display(row.name) }}</span></template></el-table-column>
        <el-table-column prop="gender" label="性别" width="90"><template #default="{ row }">{{ display(row.gender) }}</template></el-table-column>
        <el-table-column prop="birthDate" label="出生日期" min-width="140"><template #default="{ row }">{{ display(row.birthDate) }}</template></el-table-column>
        <el-table-column prop="phone" label="手机号" min-width="155"><template #default="{ row }">{{ display(row.phone) }}</template></el-table-column>
        <el-table-column prop="cardNo" label="就诊卡号" min-width="175"><template #default="{ row }">{{ display(row.cardNo) }}</template></el-table-column>
        <el-table-column label="操作" width="100" fixed="right"><template #default="{ row }"><el-button link type="primary" :aria-label="`查看${row.name}的详情`" @click="showDetails(row)">查看详情</el-button></template></el-table-column>
        <template #empty><el-empty v-if="!loading" :image-size="88" :description="errorMessage ? '数据暂不可用，请重新查询' : '暂无符合条件的患者，请调整查询条件'" /><span v-else>正在查询患者…</span></template>
      </el-table>
      <div class="results-footer" aria-live="polite"><span>{{ loading ? '正在获取患者信息…' : hasLoaded ? `本次返回 ${patients.length} 条记录` : '尚未获取查询结果' }}<span class="footer-divider">·</span>每次最多返回 100 条，请通过条件缩小范围</span><span v-if="updatedAt">更新于 {{ updatedAt }}</span></div>
    </section>
    <el-drawer v-model="detailVisible" title="患者详情" size="min(480px, 100%)">
      <template v-if="selectedPatient"><div class="detail-intro"><div class="patient-avatar">{{ selectedPatient.name?.slice(0, 1) || '患' }}</div><div><h2>{{ display(selectedPatient.name) }}</h2><p>患者基础档案</p></div></div><el-descriptions :column="1" border><el-descriptions-item label="患者编号">{{ selectedPatient.patientId }}</el-descriptions-item><el-descriptions-item label="姓名">{{ display(selectedPatient.name) }}</el-descriptions-item><el-descriptions-item label="性别">{{ display(selectedPatient.gender) }}</el-descriptions-item><el-descriptions-item label="出生日期">{{ display(selectedPatient.birthDate) }}</el-descriptions-item><el-descriptions-item label="手机号">{{ display(selectedPatient.phone) }}</el-descriptions-item><el-descriptions-item label="就诊卡号">{{ display(selectedPatient.cardNo) }}</el-descriptions-item></el-descriptions></template>
    </el-drawer>
  </section>
</template>
