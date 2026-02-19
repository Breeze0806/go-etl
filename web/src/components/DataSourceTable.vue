<template>
  <div class="datasource-table">
    <el-table
      v-loading="loading"
      :data="dataSources"
      stripe
      style="width: 100%"
    >
      <el-table-column prop="name" label="Name" min-width="150">
        <template slot-scope="scope">
          <el-tag size="medium">{{ scope.row.name }}</el-tag>
        </template>
      </el-table-column>
      
      <el-table-column prop="type" label="Type" width="120">
        <template slot-scope="scope">
          <el-tag :type="getTypeTag(scope.row.type)" size="small">
            {{ scope.row.type }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column prop="host" label="Host" min-width="180">
        <template slot-scope="scope">
          <span class="host-text">{{ formatHost(scope.row) }}</span>
        </template>
      </el-table-column>
      
      <el-table-column prop="createdAt" label="Created" width="160">
        <template slot-scope="scope">
          <i class="el-icon-time"></i>
          {{ formatDate(scope.row.createdAt) }}
        </template>
      </el-table-column>
      
      <el-table-column label="Actions" width="240" fixed="right">
        <template slot-scope="scope">
          <el-button
            size="mini"
            type="success"
            plain
            @click="handleTest(scope.row)"
            :loading="testingId === scope.row.id"
          >
            Test
          </el-button>
          <el-button
            size="mini"
            type="primary"
            plain
            @click="handleEdit(scope.row)"
          >
            Edit
          </el-button>
          <el-button
            size="mini"
            type="danger"
            plain
            @click="handleDelete(scope.row)"
          >
            Delete
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  name: 'DataSourceTable',
  props: {
    dataSources: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      testingId: null
    }
  },
  methods: {
    getTypeTag(type) {
      const typeMap = {
        'MySQL': 'primary',
        'PostgreSQL': 'success',
        'Oracle': 'warning',
        'DB2': 'danger',
        'SQLite3': 'info',
        'SQL Server': 'warning',
        'Dameng': 'primary',
        'CSV': 'success',
        'XLSX': 'success'
      }
      return typeMap[type] || 'info'
    },
    formatHost(row) {
      if (row.type === 'CSV' || row.type === 'XLSX') {
        return row.path || row.filePath || '-'
      }
      return `${row.host}:${row.port}`
    },
    formatDate(dateStr) {
      if (!dateStr) return '-'
      const date = new Date(dateStr)
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    },
    handleTest(row) {
      this.testingId = row.id
      this.$emit('test', row)
    },
    handleEdit(row) {
      this.$emit('edit', row)
    },
    handleDelete(row) {
      this.$emit('delete', row)
    }
  }
}
</script>

<style scoped>
.datasource-table {
  width: 100%;
}

.host-text {
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 13px;
  color: #606266;
}

.el-table {
  font-size: 14px;
}

.el-button {
  margin-left: 5px;
}
</style>
