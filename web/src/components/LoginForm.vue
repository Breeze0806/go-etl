<template>
  <el-form
    ref="loginForm"
    :model="loginForm"
    :rules="rules"
    class="login-form"
    autocomplete="on"
    label-position="left"
  >
    <div class="title-container">
      <h3 class="title">go-etl Login</h3>
    </div>

    <el-form-item prop="username">
      <el-input
        ref="username"
        v-model="loginForm.username"
        placeholder="Username"
        name="username"
        type="text"
        tabindex="1"
        autocomplete="on"
        prefix-icon="el-icon-user"
      />
    </el-form-item>

    <el-form-item prop="password">
      <el-input
        ref="password"
        v-model="loginForm.password"
        placeholder="Password"
        name="password"
        type="password"
        tabindex="2"
        autocomplete="on"
        prefix-icon="el-icon-lock"
        show-password
        @keyup.enter.native="handleLogin"
      />
    </el-form-item>

    <el-button
      :loading="loading"
      type="primary"
      class="login-button"
      @click.native.prevent="handleLogin"
    >
      {{ loading ? 'Logging in...' : 'Login' }}
    </el-button>

      <div class="tips">
        <span>Don't have an account? <router-link to="/register">Register</router-link></span>
      </div>
  </el-form>
</template>

<script>
export default {
  name: 'LoginForm',
  data() {
    const validateUsername = (rule, value, callback) => {
      if (!value) {
        callback(new Error('Please enter username'))
      } else if (value.length < 3 || value.length > 20) {
        callback(new Error('Username must be 3-20 characters'))
      } else {
        callback()
      }
    }
    const validatePassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error('Please enter password'))
      } else if (value.length < 6) {
        callback(new Error('Password must be at least 6 characters'))
      } else {
        callback()
      }
    }
    return {
      loginForm: {
        username: '',
        password: ''
      },
      rules: {
        username: [
          { required: true, validator: validateUsername, trigger: 'blur' }
        ],
        password: [
          { required: true, validator: validatePassword, trigger: 'blur' }
        ]
      },
      loading: false
    }
  },
  methods: {
    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          this.loading = true
          this.$http.post('/api/auth/login', this.loginForm)
            .then(response => {
              const { token } = response.data
              localStorage.setItem('jwt_token', token)
              this.$message.success('Login successful')
              
              const redirect = this.$route.query.redirect || '/'
              this.$router.push(redirect)
            })
            .catch(error => {
              this.$message.error(error.response?.data?.message || 'Login failed. Please check your credentials.')
            })
            .finally(() => {
              this.loading = false
            })
        } else {
          return false
        }
      })
    }
  }
}
</script>

<style scoped>
.login-form {
  position: relative;
  width: 100%;
  max-width: 400px;
  padding: 40px;
  margin: 0 auto;
  overflow: hidden;
}

.title-container {
  text-align: center;
  margin-bottom: 30px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
  margin: 0;
  letter-spacing: 1px;
}

.login-button {
  width: 100%;
  margin-bottom: 20px;
  font-weight: 600;
}

.tips {
  text-align: center;
  font-size: 14px;
  color: #909399;
}

.tips a {
  color: #409eff;
  text-decoration: none;
}

.tips a:hover {
  text-decoration: underline;
}
</style>
