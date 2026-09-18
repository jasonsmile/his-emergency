<script setup lang="ts">
import PatientList from './views/PatientList.vue'
</script>

<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="brand"><span class="brand-mark">H</span><div>Emergency HIS<small>医院业务扩展平台</small></div></div>
      <div class="nav-label">业务工作台</div>
      <nav aria-label="主导航"><a class="nav-item active" href="#patients" aria-current="page"><span aria-hidden="true">⌕</span> 查询患者</a></nav>
      <div class="sidebar-footer"><span class="footer-dot"></span> 患者信息查询</div>
    </aside>
    <div class="main-shell">
      <header class="topbar"><div class="breadcrumb">业务工作台 <span>/</span> <strong>查询患者</strong></div><span class="workspace-label">医院业务扩展与应急保障平台</span></header>
      <div class="tabbar"><span class="page-tab"><span class="tab-dot"></span> 查询患者</span></div>
      <main id="patients"><PatientList /></main>
    </div>
  </div>
</template>
