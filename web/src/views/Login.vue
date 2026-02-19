<template>
  <div class="login-container">
    <div class="login-background">
      <div class="background-shape shape-1"></div>
      <div class="background-shape shape-2"></div>
      <div class="background-shape shape-3"></div>
    </div>
    <div class="login-wrapper">
      <div class="login-brand">
        <div class="brand-icon">
          <i class="el-icon-data-line"></i>
        </div>
        <h1 class="brand-name">go-etl</h1>
        <p class="brand-tagline">Data Synchronization Platform</p>
      </div>
      <login-form />
    </div>
  </div>
</template>

<script>
import LoginForm from '@/components/LoginForm'

export default {
  name: 'Login',
  components: {
    LoginForm
  },
  mounted() {
    // Check if already logged in
    const token = localStorage.getItem('jwt_token')
    if (token) {
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.login-container {
  position: relative;
  width: 100%;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 1;
}

.background-shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  animation: float 20s ease-in-out infinite;
}

.shape-1 {
  width: 600px;
  height: 600px;
  background: #409eff;
  top: -200px;
  right: -200px;
  animation-delay: 0s;
}

.shape-2 {
  width: 400px;
  height: 400px;
  background: #e6a23c;
  bottom: -100px;
  left: -100px;
  animation-delay: -5s;
}

.shape-3 {
  width: 300px;
  height: 300px;
  background: #67c23a;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -10s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(20px, -20px) rotate(5deg);
  }
  50% {
    transform: translate(-10px, 20px) rotate(-5deg);
  }
  75% {
    transform: translate(-20px, -10px) rotate(3deg);
  }
}

.login-wrapper {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 420px;
  padding: 20px;
}

.login-brand {
  text-align: center;
  margin-bottom: 30px;
}

.brand-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  border-radius: 20px;
  font-size: 40px;
  color: white;
  box-shadow: 0 10px 30px rgba(64, 158, 255, 0.3);
}

.brand-name {
  font-size: 32px;
  font-weight: 700;
  color: #ffffff;
  margin: 0 0 8px 0;
  letter-spacing: 2px;
}

.brand-tagline {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
  letter-spacing: 1px;
}
</style>
