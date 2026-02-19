<template>
  <div class="register-container">
    <el-card class="register-card">
      <div slot="header">
        <h2>Create Account</h2>
      </div>
      
      <el-form
        ref="registerForm"
        :model="registerForm"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="Username" prop="username">
          <el-input
            v-model="registerForm.username"
            placeholder="Choose a username"
            prefix-icon="el-icon-user"
          />
        </el-form-item>
        
        <el-form-item label="Email" prop="email">
          <el-input
            v-model="registerForm.email"
            placeholder="Enter your email"
            prefix-icon="el-icon-message"
          />
        </el-form-item>
        
        <el-form-item label="Password" prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="Choose a password (min 6 characters)"
            prefix-icon="el-icon-lock"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="Confirm Password" prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="Confirm your password"
            prefix-icon="el-icon-lock"
            show-password
            @keyup.enter.native="handleRegister"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            :loading="loading"
            type="primary"
            class="register-button"
            @click.native.prevent="handleRegister"
          >
            {{ loading ? 'Creating account...' : 'Register' }}
          </el-button>
        </el-form-item>
        
        <div class="tips">
          Already have an account? <router-link to="/login">Login</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Register',
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
    
    const validateEmail = (rule, value, callback) => {
      if (!value) {
        callback(new Error('Please enter email'))
      } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
        callback(new Error('Please enter a valid email'))
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
    
    const validateConfirmPassword = (rule, value, callback) => {
      if (!value) {
        callback(new Error('Please confirm password'))
      } else if (value !== this.registerForm.password) {
        callback(new Error('Passwords do not match'))
      } else {
        callback()
      }
    }
    
    return {
      registerForm: {
        username: '',
        email: '',
        password: '',
        confirmPassword: ''
      },
      rules: {
        username: [
          { required: true, validator: validateUsername, trigger: 'blur' }
        ],
        email: [
          { required: true, validator: validateEmail, trigger: 'blur' }
        ],
        password: [
          { required: true, validator: validatePassword, trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, validator: validateConfirmPassword, trigger: 'blur' }
        ]
      },
      loading: false
    }
  },
  methods: {
    handleRegister() {
      this.$refs.registerForm.validate(valid => {
        if (valid) {
          this.loading = true
          this.$http.post('/auth/register', {
            username: this.registerForm.username,
            email: this.registerForm.email,
            password: this.registerForm.password
          })
            .then(() => {
              this.$message.success('Registration successful! Please login.')
              this.$router.push('/login')
            })
            .catch(error => {
              this.$message.error(error.response?.data?.message || 'Registration failed. Please try again.')
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
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-card {
  width: 100%;
  max-width: 400px;
  margin: 20px;
}

.register-card h2 {
  text-align: center;
  margin: 0;
  color: #303133;
}

.register-button {
  width: 100%;
  margin-top: 10px;
  font-weight: 600;
}

.tips {
  text-align: center;
  font-size: 14px;
  color: #909399;
  margin-top: 10px;
}

.tips a {
  color: #409eff;
  text-decoration: none;
}

.tips a:hover {
  text-decoration: underline;
}
</style>
