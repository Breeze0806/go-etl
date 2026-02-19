<template>
  <el-tag
    :type="tagType"
    :effect="effect"
    :closable="closable"
    :disable-transitions="false"
    :hit="false"
    :size="size"
    :close-transition="false"
    @close="handleClose"
  >
    {{ statusText }}
  </el-tag>
</template>

<script>
export default {
  name: 'TaskStatus',
  props: {
    status: {
      type: String,
      required: true,
      validator: function(value) {
        return ['draft', 'ready', 'running', 'paused', 'completed', 'failed'].indexOf(value) !== -1
      }
    },
    size: {
      type: String,
      default: 'small'
    },
    closable: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    tagType() {
      const statusMap = {
        draft: 'info',
        ready: 'primary',
        running: 'warning',
        paused: 'warning',
        completed: 'success',
        failed: 'danger'
      }
      return statusMap[this.status] || 'info'
    },
    statusText() {
      const textMap = {
        draft: 'Draft',
        ready: 'Ready',
        running: 'Running',
        paused: 'Paused',
        completed: 'Completed',
        failed: 'Failed'
      }
      return textMap[this.status] || this.status
    },
    effect() {
      return this.status === 'running' ? 'dark' : 'plain'
    }
  },
  methods: {
    handleClose() {
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
.el-tag {
  font-weight: 500;
  text-transform: capitalize;
}

/* Running status animation */
.el-tag.running {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}
</style>
