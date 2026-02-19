<template>
  <el-aside width="220px" class="app-sidebar">
    <div class="sidebar-logo">
      <i class="el-icon-data-line"></i>
      <span class="logo-text">go-etl</span>
    </div>
    <el-menu
      :default-active="activeMenu"
      router
      class="sidebar-menu"
      background-color="#1a1a2e"
      text-color="#a0aec0"
      active-text-color="#409eff">
      <el-menu-item index="/">
        <i class="el-icon-s-home"></i>
        <span slot="title">Dashboard</span>
      </el-menu-item>
      <el-menu-item index="/datasources">
        <i class="el-icon-connection"></i>
        <span slot="title">Data Sources</span>
      </el-menu-item>
      <el-menu-item index="/synctasks">
        <i class="el-icon-s-order"></i>
        <span slot="title">Sync Tasks</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script>
export default {
  name: 'Sidebar',
  computed: {
    activeMenu() {
      const route = this.$route
      const { meta, path } = route
      if (meta.activeMenu) {
        return meta.activeMenu
      }
      return path
    }
  }
}
</script>

<style scoped>
.app-sidebar {
  background-color: #1a1a2e;
  height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.sidebar-logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #16213e;
  border-bottom: 1px solid #2d3748;
}

.sidebar-logo i {
  font-size: 28px;
  color: #409eff;
  margin-right: 10px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 1px;
}

.sidebar-menu {
  border-right: none;
  flex: 1;
  overflow-y: auto;
}

.sidebar-menu /deep/ .el-menu-item {
  height: 50px;
  line-height: 50px;
  font-size: 14px;
}

.sidebar-menu /deep/ .el-menu-item:hover,
.sidebar-menu /deep/ .el-menu-item.is-active {
  background-color: #2d3748 !important;
}

.sidebar-menu /deep/ .el-menu-item i {
  margin-right: 12px;
  font-size: 18px;
}
</style>
