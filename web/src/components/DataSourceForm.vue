<template>
  <el-dialog
    :title="isEdit ? 'Edit Data Source' : 'Add Data Source'"
    :visible.sync="dialogVisible"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="form"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="Name" prop="name">
        <el-input
          v-model="form.name"
          placeholder="Enter data source name"
        />
      </el-form-item>
      
      <el-form-item label="Type" prop="type">
        <el-select
          v-model="form.type"
          placeholder="Select data source type"
          style="width: 100%"
          :disabled="isEdit"
        >
          <el-option
            v-for="item in sourceTypes"
            :key="item"
            :label="item"
            :value="item"
          />
        </el-select>
      </el-form-item>
      
      <!-- Database Connection Fields -->
      <template v-if="isDatabaseType">
        <el-form-item label="Host" prop="host">
          <el-input
            v-model="form.host"
            placeholder="e.g., localhost"
          />
        </el-form-item>
        
        <el-form-item label="Port" prop="port">
          <el-input-number
            v-model="form.port"
            :min="1"
            :max="65535"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="Username" prop="username">
          <el-input
            v-model="form.username"
            placeholder="Database username"
          />
        </el-form-item>
        
        <el-form-item label="Password" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="Database password"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="Database" prop="database">
          <el-input
            v-model="form.database"
            placeholder="Database name"
          />
        </el-form-item>
      </template>
      
      <!-- File-based Fields -->
      <template v-if="isFileType">
        <el-form-item label="File Path" prop="path">
          <el-input
            v-model="form.path"
            placeholder="Absolute path to file"
          />
        </el-form-item>
      </template>
      
      <!-- CSV/XLSX specific -->
      <template v-if="form.type === 'CSV' || form.type === 'XLSX'">
        <el-form-item label="Encoding" prop="encoding">
          <el-select
            v-model="form.encoding"
            placeholder="Select encoding"
            style="width: 100%"
          >
            <el-option label="UTF-8" value="utf-8" />
            <el-option label="GBK" value="gbk" />
            <el-option label="GB2312" value="gb2312" />
            <el-option label="ASCII" value="ascii" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="Has Header" prop="hasHeader">
          <el-switch v-model="form.hasHeader" />
        </el-form-item>
      </template>
    </el-form>
    
    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">Cancel</el-button>
      <el-button
        type="primary"
        :loading="submitting"
        @click="handleSubmit"
      >
        {{ isEdit ? 'Update' : 'Create' }}
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
export default {
  name: 'DataSourceForm',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    dataSource: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      form: {
        name: '',
        type: '',
        host: 'localhost',
        port: 3306,
        username: '',
        password: '',
        database: '',
        path: '',
        encoding: 'utf-8',
        hasHeader: true
      },
      sourceTypes: [
        'MySQL',
        'PostgreSQL',
        'Oracle',
        'DB2',
        'SQLite3',
        'SQL Server',
        'Dameng',
        'CSV',
        'XLSX'
      ],
      defaultPorts: {
        MySQL: 3306,
        PostgreSQL: 5432,
        Oracle: 1521,
        'DB2': 50000,
        SQLite3: 0,
        'SQL Server': 1433,
        Dameng: 5236
      },
      submitting: false
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.visible
      },
      set(val) {
        this.$emit('update:visible', val)
      }
    },
    isEdit() {
      return !!this.dataSource && !!this.dataSource.id
    },
    isDatabaseType() {
      return ['MySQL', 'PostgreSQL', 'Oracle', 'DB2', 'SQLite3', 'SQL Server', 'Dameng'].includes(this.form.type)
    },
    isFileType() {
      return ['CSV', 'XLSX'].includes(this.form.type)
    },
    rules() {
      const basicRules = {
        name: [
          { required: true, message: 'Please enter data source name', trigger: 'blur' }
        ],
        type: [
          { required: true, message: 'Please select data source type', trigger: 'change' }
        ]
      }
      
      if (this.isDatabaseType) {
        return {
          ...basicRules,
          host: [
            { required: true, message: 'Please enter host', trigger: 'blur' }
          ],
          port: [
            { required: true, message: 'Please enter port', trigger: 'blur' }
          ],
          username: [
            { required: true, message: 'Please enter username', trigger: 'blur' }
          ]
        }
      }
      
      if (this.isFileType) {
        return {
          ...basicRules,
          path: [
            { required: true, message: 'Please enter file path', trigger: 'blur' }
          ]
        }
      }
      
      return basicRules
    }
  },
  watch: {
    'form.type'(newType) {
      if (newType && this.defaultPorts[newType]) {
        this.form.port = this.defaultPorts[newType]
      }
    },
    dataSource: {
      immediate: true,
      handler(val) {
        if (val) {
          this.form = { ...this.form, ...val }
        } else {
          this.resetForm()
        }
      }
    }
  },
  methods: {
    handleClose() {
      this.$emit('update:visible', false)
      this.$emit('close')
      this.resetForm()
    },
    handleSubmit() {
      this.$refs.form.validate(valid => {
        if (valid) {
          this.submitting = true
          const submitData = { ...this.form }
          
          // Don't send password if empty (for editing)
          if (this.isEdit && !submitData.password) {
            delete submitData.password
          }
          
          if (this.isEdit) {
            this.$emit('update', submitData)
          } else {
            this.$emit('create', submitData)
          }
          
          setTimeout(() => {
            this.submitting = false
          }, 500)
        }
      })
    },
    resetForm() {
      this.form = {
        name: '',
        type: '',
        host: 'localhost',
        port: 3306,
        username: '',
        password: '',
        database: '',
        path: '',
        encoding: 'utf-8',
        hasHeader: true
      }
      if (this.$refs.form) {
        this.$refs.form.clearValidate()
      }
    }
  }
}
</script>

<style scoped>
.el-form-item {
  margin-bottom: 20px;
}
</style>
