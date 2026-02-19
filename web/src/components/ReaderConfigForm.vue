<template>
  <el-form
    ref="form"
    :model="form"
    :rules="rules"
    label-width="140px"
    class="reader-config-form"
  >
    <el-form-item label="Task Name" prop="taskName">
      <el-input
        v-model="form.taskName"
        placeholder="Enter sync task name"
      />
    </el-form-item>

    <el-form-item label="Reader Data Source" prop="readerDataSourceId">
      <el-select
        v-model="form.readerDataSourceId"
        placeholder="Select reader data source"
        style="width: 100%"
        filterable
        @change="handleReaderSourceChange"
      >
        <el-option
          v-for="ds in dataSources"
          :key="ds.id"
          :label="ds.name"
          :value="ds.id"
        >
          <span>{{ ds.name }}</span>
          <span style="float: right; color: #8492a6; font-size: 12px">
            {{ ds.type }}
          </span>
        </el-option>
      </el-select>
    </el-form-item>

    <el-form-item label="Source Table" prop="readerTable">
      <el-input
        v-model="form.readerTable"
        placeholder="Enter source table name"
        :disabled="!form.readerDataSourceId"
      />
    </el-form-item>

    <el-form-item label="Query (Optional)" prop="readerQuery">
      <el-input
        v-model="form.readerQuery"
        type="textarea"
        :rows="4"
        placeholder="Enter custom query (e.g., SELECT * FROM table WHERE condition)"
        :disabled="!form.readerDataSourceId"
      />
      <div class="form-tip">
        Leave empty to read all data from the table. Custom query takes precedence over table name.
      </div>
    </el-form-item>
  </el-form>
</template>

<script>
export default {
  name: 'ReaderConfigForm',
  props: {
    dataSources: {
      type: Array,
      default: () => []
    },
    value: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      form: {
        taskName: '',
        readerDataSourceId: '',
        readerTable: '',
        readerQuery: ''
      },
      rules: {
        taskName: [
          { required: true, message: 'Please enter task name', trigger: 'blur' },
          { min: 2, max: 100, message: 'Task name must be 2-100 characters', trigger: 'blur' }
        ],
        readerDataSourceId: [
          { required: true, message: 'Please select reader data source', trigger: 'change' }
        ],
        readerTable: [
          { required: true, message: 'Please enter source table name', trigger: 'blur' }
        ],
        readerQuery: [
          { validator: this.validateQuery, trigger: 'blur' }
        ]
      }
    }
  },
  watch: {
    value: {
      immediate: true,
      handler(val) {
        if (val) {
          this.form = { ...this.form, ...val }
        }
      }
    },
    form: {
      deep: true,
      handler(val) {
        this.$emit('input', val)
        this.$emit('update', val)
      }
    }
  },
  methods: {
    handleReaderSourceChange() {
      this.form.readerTable = ''
      this.form.readerQuery = ''
    },
    validateQuery(rule, value, callback) {
      if (value && value.trim()) {
        const trimmed = value.trim().toLowerCase()
        if (!trimmed.startsWith('select')) {
          callback(new Error('Query must start with SELECT'))
        } else {
          callback()
        }
      } else {
        callback()
      }
    },
    validate() {
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(valid => {
          if (valid) {
            resolve(true)
          } else {
            reject(new Error('Please fill in all required fields'))
          }
        })
      })
    },
    resetFields() {
      this.$refs.form.resetFields()
    }
  }
}
</script>

<style scoped>
.reader-config-form {
  max-width: 600px;
  margin: 0 auto;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}
</style>
