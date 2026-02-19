<template>
  <div class="synctasks">
    <el-card>
      <div slot="header" class="clearfix header-row">
        <div class="header-left">
          <span class="page-title">Sync Tasks</span>
          <el-button
            style="float: right;"
            type="primary"
            size="small"
            icon="el-icon-plus"
            @click="$router.push('/synctasks/new')"
          >
            New Task
          </el-button>
        </div>
        <div class="header-right">
          <el-button
            size="small"
            icon="el-icon-refresh"
            :loading="loading"
            @click="fetchTasks"
          >
            Refresh
          </el-button>
        </div>
      </div>

      <!-- Task Table -->
      <el-table
        v-loading="loading"
        :data="tasks"
        stripe
        border
        style="width: 100%"
        empty-text="No sync tasks created yet"
      >
        <el-table-column
          prop="id"
          label="ID"
          width="80"
          show-overflow-tooltip
        />

        <el-table-column
          prop="name"
          label="Task Name"
          min-width="150"
          show-overflow-tooltip
        />

        <el-table-column
          prop="status"
          label="Status"
          width="120"
          align="center"
        >
          <template slot-scope="scope">
            <task-status :status="scope.row.status" />
          </template>
        </el-table-column>

        <el-table-column
          label="Source"
          min-width="150"
          show-overflow-tooltip
        >
          <template slot-scope="scope">
            <div v-if="scope.row.reader">
              {{ scope.row.reader.dataSourceName || scope.row.reader.dataSourceId }}
              <span class="table-info">{{ scope.row.reader.table }}</span>
            </div>
            <span v-else class="table-empty">-</span>
          </template>
        </el-table-column>

        <el-table-column
          label="Target"
          min-width="150"
          show-overflow-tooltip
        >
          <template slot-scope="scope">
            <div v-if="scope.row.writer">
              {{ scope.row.writer.dataSourceName || scope.row.writer.dataSourceId }}
              <span class="table-info">{{ scope.row.writer.table }}</span>
            </div>
            <span v-else class="table-empty">-</span>
          </template>
        </el-table-column>

        <el-table-column
          prop="createdAt"
          label="Created"
          width="160"
          align="center"
        >
          <template slot-scope="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column
          label="Actions"
          width="200"
          align="center"
          fixed="right"
        >
          <template slot-scope="scope">
            <el-button
              v-if="scope.row.status === 'ready' || scope.row.status === 'paused'"
              type="success"
              size="mini"
              icon="el-icon-video-play"
              :loading="scope.row.starting"
              @click="startTask(scope.row)"
            >
              Start
            </el-button>
            <el-button
              v-if="scope.row.status === 'running'"
              type="warning"
              size="mini"
              icon="el-icon-video-pause"
              :loading="scope.row.stopping"
              @click="stopTask(scope.row)"
            >
              Stop
            </el-button>
            <el-button
              type="info"
              size="mini"
              icon="el-icon-view"
              @click="viewDetails(scope.row)"
            >
              Details
            </el-button>
            <el-button
              v-if="scope.row.status === 'draft' || scope.row.status === 'ready' || scope.row.status === 'completed' || scope.row.status === 'failed'"
              type="danger"
              size="mini"
              icon="el-icon-delete"
              :loading="scope.row.deleting"
              @click="deleteTask(scope.row)"
            >
              Delete
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Task Details Drawer -->
    <el-drawer
      v-model="detailsDrawer"
      :title="selectedTask ? selectedTask.name : 'Task Details'"
      size="500px"
    >
      <div v-if="selectedTask" class="task-details">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="Task ID">
            {{ selectedTask.id }}
          </el-descriptions-item>
          <el-descriptions-item label="Status">
            <task-status :status="selectedTask.status" />
          </el-descriptions-item>
          <el-descriptions-item label="Created">
            {{ formatDate(selectedTask.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="Updated">
            {{ formatDate(selectedTask.updatedAt) }}
          </el-descriptions-item>
        </el-descriptions>

        <el-divider content-position="left">Reader Configuration</el-divider>
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="Data Source">
            {{ selectedTask.reader?.dataSourceName || selectedTask.reader?.dataSourceId || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="Table">
            {{ selectedTask.reader?.table || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="Query">
            <code v-if="selectedTask.reader?.query">{{ selectedTask.reader.query }}</code>
            <span v-else class="table-empty">-</span>
          </el-descriptions-item>
        </el-descriptions>

        <el-divider content-position="left">Writer Configuration</el-divider>
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="Data Source">
            {{ selectedTask.writer?.dataSourceName || selectedTask.writer?.dataSourceId || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="Table">
            {{ selectedTask.writer?.table || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="Write Mode">
            {{ selectedTask.writer?.mode || 'insert' }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="selectedTask.stats" class="task-stats">
          <el-divider content-position="left">Statistics</el-divider>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="stat-item">
                <div class="stat-value">{{ selectedTask.stats.totalRows || 0 }}</div>
                <div class="stat-label">Total Rows</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="stat-item">
                <div class="stat-value">{{ selectedTask.stats.successRows || 0 }}</div>
                <div class="stat-label">Success</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="stat-item">
                <div class="stat-value">{{ selectedTask.stats.failedRows || 0 }}</div>
                <div class="stat-label">Failed</div>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script>
import TaskStatus from '@/components/TaskStatus.vue'

export default {
  name: 'SyncTasks',
  components: {
    TaskStatus
  },
  data() {
    return {
      tasks: [],
      loading: false,
      detailsDrawer: false,
      selectedTask: null,
      refreshInterval: null,
      REFRESH_INTERVAL_MS: 5000
    }
  },
  computed: {
    hasRunningTasks() {
      return this.tasks.some(task => task.status === 'running')
    }
  },
  created() {
    this.fetchTasks()
    this.startAutoRefresh()
  },
  beforeDestroy() {
    this.stopAutoRefresh()
  },
  methods: {
    async fetchTasks() {
      this.loading = true
      try {
        const response = await this.$http.get('/api/synctasks')
        this.tasks = (response.data || []).map(task => ({
          ...task,
          starting: false,
          stopping: false,
          deleting: false
        }))
      } catch (error) {
        this.$message.error(error.response?.data?.message || 'Failed to fetch sync tasks')
      } finally {
        this.loading = false
      }
    },
    startAutoRefresh() {
      // Auto-refresh every 5 seconds if there are running tasks
      this.refreshInterval = setInterval(() => {
        if (this.hasRunningTasks) {
          this.fetchTasks()
        }
      }, this.REFRESH_INTERVAL_MS)
    },
    stopAutoRefresh() {
      if (this.refreshInterval) {
        clearInterval(this.refreshInterval)
        this.refreshInterval = null
      }
    },
    async startTask(task) {
      this.$set(task, 'starting', true)
      try {
        await this.$http.post(`/api/synctasks/${task.id}/start`)
        this.$message.success(`Task "${task.name}" started successfully`)
        await this.fetchTasks()
      } catch (error) {
        this.$message.error(error.response?.data?.message || 'Failed to start task')
      } finally {
        this.$set(task, 'starting', false)
      }
    },
    async stopTask(task) {
      this.$set(task, 'stopping', true)
      try {
        // Using DELETE on /api/syncjobs/:id to stop as per requirements
        await this.$http.delete(`/api/syncjobs/${task.id}`)
        this.$message.success(`Task "${task.name}" stopped successfully`)
        await this.fetchTasks()
      } catch (error) {
        this.$message.error(error.response?.data?.message || 'Failed to stop task')
      } finally {
        this.$set(task, 'stopping', false)
      }
    },
    async deleteTask(task) {
      try {
        await this.$confirm(
          `Are you sure you want to delete task "${task.name}"? This action cannot be undone.`,
          'Delete Task',
          {
            confirmButtonText: 'Delete',
            cancelButtonText: 'Cancel',
            type: 'warning'
          }
        )

        this.$set(task, 'deleting', true)
        await this.$http.delete(`/api/synctasks/${task.id}`)
        this.$message.success(`Task "${task.name}" deleted successfully`)
        await this.fetchTasks()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error(error.response?.data?.message || 'Failed to delete task')
        }
      } finally {
        this.$set(task, 'deleting', false)
      }
    },
    viewDetails(task) {
      this.selectedTask = task
      this.detailsDrawer = true
    },
    formatDate(dateString) {
      if (!dateString) return '-'
      const date = new Date(dateString)
      if (isNaN(date.getTime())) return '-'
      return date.toLocaleString('en-US', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }
  }
}
</script>

<style scoped>
.synctasks {
  height: 100%;
  padding: 20px;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.table-info {
  color: #909399;
  font-size: 12px;
  margin-left: 4px;
}

.table-empty {
  color: #c0c4cc;
}

.task-details {
  padding: 0 20px;
}

.task-stats {
  margin-top: 20px;
}

.stat-item {
  text-align: center;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

code {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
  word-break: break-all;
}

/* Override element-ui drawer styles */
::v-deep .el-drawer__header {
  margin-bottom: 0;
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
}

::v-deep .el-drawer__body {
  padding: 20px;
  overflow-y: auto;
}
</style>
