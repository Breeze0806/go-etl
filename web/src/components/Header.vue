<template>
  <el-header class="app-header">
    <div class="header-left">
      <h2 class="page-title">{{ pageTitle }}</h2>
    </div>
    <div class="header-right">
      <el-dropdown @command="handleCommand" trigger="click">
        <div class="user-info">
          <el-avatar :size="36" class="user-avatar">
            <i class="el-icon-user-solid"></i>
          </el-avatar>
          <span class="user-name">{{ username }}</span>
          <i class="el-icon-arrow-down el-icon--right"></i>
        </div>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="profile">
            <i class="el-icon-user"></i>
            Profile
          </el-dropdown-item>
          <el-dropdown-item command="settings">
            <i class="el-icon-setting"></i>
            Settings
          </el-dropdown-item>
          <el-dropdown-item divided command="logout">
            <i class="el-icon-switch-button"></i>
            Logout
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </el-header>
</template>

<script>
export default {
  name: 'Header',
  computed: {
    username() {
      const user = localStorage.getItem('user')
      if (user) {
        try {
          return JSON.parse(user).username || 'Admin'
        } catch (e) {
          return 'Admin'
        }
      }
      return 'Admin'
    },
    pageTitle() {
      const route = this.$route
      const titles = {
        '/': 'Dashboard',
        '/datasources': 'Data Sources',
        '/synctasks': 'Sync Tasks',
        '/synctasks/new': 'New Sync Task'
      }
      return titles[route.path] || 'go-etl'
    }
  },
  methods: {
    handleCommand(command) {
      switch (command) {
        case 'logout':
          this.logout()
          break
        case 'profile':
          this.$message.info('Profile page coming soon')
          break
        case 'settings':
          this.$message.info('Settings page coming soon')
          break
      }
    },
    logout() {
      localStorage.removeItem('jwt_token')
      localStorage.removeItem('user')
      this.$message.success('Logged out successfully')
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.app-header {
  background-color: #ffffff;
  color: #333;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  height: 60px;
  line-height: 60px;
}

.header-left {
  flex: 1;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a1a2e;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.user-avatar {
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  margin-right: 10px;
}

.user-name {
  font-size: 14px;
  color: #606266;
  margin-right: 6px;
}

.el-dropdown-menu {
  margin-top: 8px;
}

.el-dropdown-menu__item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
}

.el-dropdown-menu__item i {
  margin-right: 8px;
}
</style>
