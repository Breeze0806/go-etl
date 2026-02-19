<template>
  <div class="home">
    <el-container>
      <el-header>
        <div class="logo">go-etl</div>
      </el-header>
      <el-container>
        <el-aside width="200px">
          <el-menu
            :default-active="$route.path"
            router
            class="el-menu-vertical">
            <el-menu-item index="/">
              <i class="el-icon-s-home"></i>
              <span>Home</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main>
          <h1>Welcome to go-etl</h1>
          <p>Data synchronization platform</p>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'Home'
}
</script>

<style scoped>
.home {
  height: 100vh;
}

.el-header {
  background-color: #409EFF;
  color: white;
  display: flex;
  align-items: center;
}

.logo {
  font-size: 20px;
  font-weight: bold;
}

.el-aside {
  background-color: #F5F7FA;
}

.el-main {
  background-color: #FFFFFF;
}
</style>
