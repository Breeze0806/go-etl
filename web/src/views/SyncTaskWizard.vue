<template>
  <div class="sync-task-wizard">
    <el-card>
      <div slot="header" class="clearfix">
        <span class="page-title">Create New Sync Task</span>
      </div>
      
      <el-steps :active="activeStep" finish-status="success" align-center>
        <el-step title="Configure Reader" description="Select source data and table"></el-step>
        <el-step title="Configure Writer" description="Select destination data and table"></el-step>
      </el-steps>
      
      <div class="wizard-content">
        <!-- Step 1: Reader Configuration -->
        <div v-show="activeStep === 0" class="step-content">
          <reader-config-form
            ref="readerForm"
            :data-sources="dataSources"
            :value="taskConfig"
            @update="handleReaderUpdate"
          />
        </div>
        
        <!-- Step 2: Writer Configuration -->
        <div v-show="activeStep === 1" class="step-content">
          <writer-config-form
            ref="writerForm"
            :data-sources="dataSources"
            :reader-data-source-id="taskConfig.readerDataSourceId"
            :value="taskConfig"
            @update="handleWriterUpdate"
          />
        </div>
      </div>
      
      <div class="wizard-actions">
        <el-button v-if="activeStep > 0" @click="prevStep">
          Previous
        </el-button>
        <el-button v-if="activeStep === 0" type="primary" @click="nextStep">
          Next: Configure Writer
        </el-button>
        <el-button v-if="activeStep === 1" type="primary" @click="submitTask">
          Create Task
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import ReaderConfigForm from '@/components/ReaderConfigForm.vue'
import WriterConfigForm from '@/components/WriterConfigForm.vue'

export default {
  name: 'SyncTaskWizard',
  components: {
    ReaderConfigForm,
    WriterConfigForm
  },
  data() {
    return {
      activeStep: 0,
      dataSources: [],
      loading: false,
      taskConfig: {
        taskName: '',
        readerDataSourceId: '',
        readerTable: '',
        readerQuery: '',
        writerDataSourceId: '',
        writerTable: '',
        writeMode: 'insert'
      }
    }
  },
  created() {
    this.fetchDataSources()
  },
  methods: {
    async fetchDataSources() {
      this.loading = true
      try {
        const response = await this.$http.get('/api/datasources')
        this.dataSources = response.data || []
      } catch (error) {
        this.$message.error(error.response?.data?.message || 'Failed to fetch data sources')
      } finally {
        this.loading = false
      }
    },
    handleReaderUpdate(config) {
      this.taskConfig = { ...this.taskConfig, ...config }
    },
    handleWriterUpdate(config) {
      this.taskConfig = { ...this.taskConfig, ...config }
    },
    async nextStep() {
      try {
        await this.$refs.readerForm.validate()
        
        // Check if at least one data source is available
        if (this.dataSources.length === 0) {
          this.$message.warning('No data sources available. Please create a data source first.')
          return
        }
        
        // Check if reader and writer are the same data source
        if (this.taskConfig.readerDataSourceId === this.taskConfig.writerDataSourceId) {
          // Allow proceeding but warn user (handled in writer form)
        }
        
        this.activeStep = 1
      } catch (error) {
        this.$message.error(error.message || 'Please fill in all required fields')
      }
    },
    prevStep() {
      if (this.activeStep > 0) {
        this.activeStep--
      }
    },
    async submitTask() {
      try {
        await this.$refs.writerForm.validate()
        
        // Prepare the task configuration for API
        const submitData = {
          name: this.taskConfig.taskName,
          reader: {
            dataSourceId: this.taskConfig.readerDataSourceId,
            table: this.taskConfig.readerTable,
            query: this.taskConfig.readerQuery || undefined
          },
          writer: {
            dataSourceId: this.taskConfig.writerDataSourceId,
            table: this.taskConfig.writerTable,
            mode: this.taskConfig.writeMode
          }
        }
        
        // Remove empty query from submission
        if (!submitData.reader.query) {
          delete submitData.reader.query
        }
        
        await this.$http.post('/api/synctasks', submitData)
        this.$message.success('Sync task created successfully!')
        this.$router.push('/synctasks')
      } catch (error) {
        this.$message.error(error.response?.data?.message || 'Failed to create sync task')
      }
    }
  }
}
</script>

<style scoped>
.sync-task-wizard {
  height: 100%;
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.wizard-content {
  min-height: 350px;
  padding: 40px 20px;
}

.step-content {
  max-width: 700px;
  margin: 0 auto;
}

.wizard-actions {
  text-align: center;
  padding: 20px;
  border-top: 1px solid #ebeef5;
}

.wizard-actions .el-button {
  min-width: 120px;
  margin: 0 10px;
}
</style>
