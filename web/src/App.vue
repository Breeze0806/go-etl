<template>
  <div id="app">
    <!-- Login page - full screen without layout -->
    <router-view v-if="isLoginPage" />

    <!-- Protected routes - with sidebar and header -->
    <el-container v-else class="app-layout">
      <sidebar />
      <el-container class="main-container">
        <header />
        <el-main class="app-main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import Sidebar from '@/components/Sidebar'
import Header from '@/components/Header'

export default {
  name: 'App',
  components: {
    Sidebar,
    Header
  },
  computed: {
    isLoginPage() {
      return this.$route.path === '/login'
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', Arial, sans-serif;
}

#app {
  height: 100%;
}

.app-layout {
  height: 100vh;
  overflow: hidden;
}

.main-container {
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.app-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
  height: calc(100vh - 60px);
}

/* Element UI overrides */
.el-menu {
  border-right: none;
}
</style>
