<template>
  <el-form
    ref="form"
    :model="form"
    :rules="rules"
    label-width="140px"
    class="writer-config-form"
  >
    <el-form-item label="Writer Data Source" prop="writerDataSourceId">
      <el-select
        v-model="form.writerDataSourceId"
        placeholder="Select writer data source"
        style="width: 100%"
        filterable
        @change="handleWriterSourceChange"
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

    <el-form-item label="Target Table" prop="writerTable">
      <el-input
        v-model="form.writerTable"
        placeholder="Enter target table name"
        :disabled="!form.writerDataSourceId"
      />
    </el-form-item>

    <el-form-item label="Write Mode" prop="writeMode">
      <el-select
        v-model="form.writeMode"
        placeholder="Select write mode"
        style="width: 100%"
        :disabled="!form.writerDataSourceId"
      >
        <el-option label="Insert" value="insert" />
        <el-option label="Replace" value="replace" />
        <el-option label="Upsert" value="upsert" />
      </el-select>
    </el-form-item>

    <el-alert
      v-if="sameSourceWarning"
      title="Warning"
      type="warning"
      :closable="false"
      show-icon
      style="margin-bottom: 20px;"
    >
      Reader and Writer are using the same data source. Make sure the target table is different from the source table.
    </el-alert>
  </el-form>
</template>

<script>
export default {
  name: 'WriterConfigForm',
  props: {
    dataSources: {
      type: Array,
      default: () => []
    },
    readerDataSourceId: {
      type: [Number, String],
      default: null
    },
    value: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      form: {
        writerDataSourceId: '',
        writerTable: '',
        writeMode: 'insert'
      },
      rules: {
        writerDataSourceId: [
          { required: true, message: 'Please select writer data source', trigger: 'change' }
        ],
        writerTable: [
          { required: true, message: 'Please enter target table name', trigger: 'blur' }
        ],
        writeMode: [
          { required: true, message: 'Please select write mode', trigger: 'change' }
        ]
      }
    }
  },
  computed: {
    sameSourceWarning() {
      return this.readerDataSourceId && 
             this.form.writerDataSourceId && 
             String(this.readerDataSourceId) === String(this.form.writerDataSourceId)
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
    handleWriterSourceChange() {
      this.form.writerTable = ''
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
.writer-config-form {
  max-width: 600px;
  margin: 0 auto;
}
</style>
