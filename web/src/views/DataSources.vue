<template>
  <div class="datasources">
    <el-card>
      <div slot="header" class="clearfix">
        <span class="page-title">Data Sources</span>
        <el-button 
          style="float: right;" 
          type="primary" 
          size="small" 
          icon="el-icon-plus"
          @click="handleAdd"
        >
          Add Data Source
        </el-button>
      </div>
      
      <data-source-table
        v-if="dataSources.length > 0"
        :data-sources="dataSources"
        :loading="loading"
        @test="handleTest"
        @edit="handleEdit"
        @delete="handleDelete"
      />
      
      <el-empty 
        v-else 
        description="No data sources configured. Click 'Add Data Source' to create one."
      >
        <el-button type="primary" @click="handleAdd">
          Add Data Source
        </el-button>
      </el-empty>
    </el-card>
    
    <!-- Add/Edit Form Modal -->
    <data-source-form
      :visible.sync="formVisible"
      :data-source="currentDataSource"
      @create="handleCreate"
      @update="handleUpdate"
      @close="handleFormClose"
    />
    
    <!-- Test Connection Result Dialog -->
    <el-dialog
      title="Connection Test"
      :visible.sync="testDialogVisible"
      width="400px"
    >
      <div class="test-result">
        <el-result
          :icon="testSuccess ? 'success' : 'warning'"
          :title="testSuccess ? 'Connection Successful' : 'Connection Failed'"
          :sub-title="testMessage"
        />
      </div>
      <span slot="footer">
        <el-button type="primary" @click="testDialogVisible = false">
          OK
        </el-button>
      </span>
    </el-dialog>
    
    <!-- Delete Confirmation Dialog -->
    <el-dialog
      title="Delete Data Source"
      :visible.sync="deleteDialogVisible"
      width="400px"
    >
      <p>Are you sure you want to delete <strong>{{ currentDataSource?.name }}</strong>?</p>
      <p class="warning-text">This action cannot be undone.</p>
      <span slot="footer">
        <el-button @click="deleteDialogVisible = false">Cancel</el-button>
        <el-button 
          type="danger" 
          :loading="deleting"
          @click="confirmDelete"
        >
          Delete
        </el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import DataSourceTable from '@/components/DataSourceTable.vue'
import DataSourceForm from '@/components/DataSourceForm.vue'

export default {
  name: 'DataSources',
  components: {
    DataSourceTable,
    DataSourceForm
  },
  data() {
    return {
      dataSources: [],
      loading: false,
      formVisible: false,
      currentDataSource: null,
      testDialogVisible: false,
      testSuccess: false,
      testMessage: '',
      deleteDialogVisible: false,
      deleting: false
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
    handleAdd() {
      this.currentDataSource = null
      this.formVisible = true
    },
    handleEdit(row) {
      this.currentDataSource = { ...row }
      this.formVisible = true
    },
    async handleCreate(formData) {
      try {
        await this.$http.post('/api/datasources', formData)
        this.$message.success('Data source created successfully')
        this.formVisible = false
        this.fetchDataSources()
      } catch (error) {
        this.$message.error(error.response?.data?.message || 'Failed to create data source')
      }
    },
    async handleUpdate(formData) {
      try {
        await this.$http.put(`/api/datasources/${formData.id}`, formData)
        this.$message.success('Data source updated successfully')
        this.formVisible = false
        this.fetchDataSources()
      } catch (error) {
        this.$message.error(error.response?.data?.message || 'Failed to update data source')
      }
    },
    handleDelete(row) {
      this.currentDataSource = row
      this.deleteDialogVisible = true
    },
    async confirmDelete() {
      if (!this.currentDataSource) return
      
      this.deleting = true
      try {
        await this.$http.delete(`/api/datasources/${this.currentDataSource.id}`)
        this.$message.success('Data source deleted successfully')
        this.deleteDialogVisible = false
        this.fetchDataSources()
      } catch (error) {
        this.$message.error(error.response?.data?.message || 'Failed to delete data source')
      } finally {
        this.deleting = false
      }
    },
    async handleTest(row) {
      try {
        const response = await this.$http.post(`/api/datasources/${row.id}/test`)
        this.testSuccess = true
        this.testMessage = response.data?.message || 'Connection successful!'
        this.testDialogVisible = true
      } catch (error) {
        this.testSuccess = false
        this.testMessage = error.response?.data?.message || 'Connection failed. Please check your settings.'
        this.testDialogVisible = true
      }
    },
    handleFormClose() {
      this.currentDataSource = null
    }
  }
}
</script>

<style scoped>
.datasources {
  height: 100%;
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.warning-text {
  color: #f56c6c;
  font-size: 14px;
}

.test-result {
  padding: 20px 0;
}
</style>
