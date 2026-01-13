<template>
  <div class="container">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="任务管理" name="tasks">
        <el-card>
          <template #header>
            <div class="header">
              <span>日志清理任务</span>
              <el-button type="primary" @click="showDialog()">添加任务</el-button>
            </div>
          </template>
      <el-table :data="tasks" stripe>
        <el-table-column prop="name" label="任务名称" />
        <el-table-column prop="path" label="路径" show-overflow-tooltip />
        <el-table-column prop="mode" label="模式" width="100">
          <template #default="{ row }">
            {{ row.mode === 'truncate' ? '清空文件' : '删除过期' }}
          </template>
        </el-table-column>
        <el-table-column prop="retentionDays" label="保留天数" width="100">
          <template #default="{ row }">
            {{ row.mode === 'retention' ? row.retentionDays + '天' : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="cronSpec" label="执行时间" width="120" />
        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" @change="saveTask(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" @click="runNow(row.id)">执行</el-button>
              <el-button size="small" @click="showDialog(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="deleteTask(row.id)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="错误日志" name="logs">
        <el-card>
          <template #header>
            <div class="header">
              <span>错误日志</span>
              <div>
                <el-button @click="loadLogs">刷新</el-button>
                <el-button type="danger" @click="clearLogs">清空</el-button>
              </div>
            </div>
          </template>
          <div style="display: flex; gap: 10px; margin-bottom: 10px;">
            <el-input v-model="searchTask" placeholder="搜索任务名称" clearable style="width: 200px" />
            <el-date-picker v-model="searchDate" type="date" placeholder="选择日期" format="YYYY-MM-DD" value-format="YYYY-MM-DD" clearable />
          </div>
          <pre style="max-height: 500px; overflow: auto; white-space: pre-wrap; font-size: 12px;">{{ filteredLogs || '暂无错误日志' }}</pre>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑任务' : '添加任务'" width="500">
      <el-form :model="form" label-width="100px">
        <el-form-item label="任务名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="路径">
          <div style="display: flex; gap: 8px; width: 100%">
            <el-input v-model="form.path" style="flex: 1" />
            <el-button @click="selectPath">文件</el-button>
            <el-button @click="selectDir">目录</el-button>
          </div>
        </el-form-item>
        <el-form-item label="清理模式">
          <el-radio-group v-model="form.mode">
            <el-radio value="truncate">清空文件内容</el-radio>
            <el-radio value="retention">删除过期文件</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.mode === 'retention'" label="保留天数">
          <el-input-number v-model="form.retentionDays" :min="1" />
        </el-form-item>
        <el-form-item v-if="form.mode === 'retention'" label="文件类型">
          <el-select v-model="form.filePattern" style="width: 100%">
            <el-option label="所有文件" value="*" />
            <el-option label="仅 .log 文件" value=".log" />
            <el-option label="仅 .txt 文件" value=".txt" />
            <el-option label="仅 .tmp 文件" value=".tmp" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行时间">
          <el-select v-model="form.cronSpec" style="width: 100%">
            <el-option label="每天凌晨2点" value="0 2 * * *" />
            <el-option label="每天凌晨4点" value="0 4 * * *" />
            <el-option label="每周日凌晨3点" value="0 3 * * 0" />
            <el-option label="每月1号凌晨3点" value="0 3 1 * *" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTask(form)">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="errorDialogVisible" title="失败详情" width="600">
      <el-table :data="errorList" max-height="400">
        <el-table-column prop="" label="错误信息">
          <template #default="{ row }">{{ row }}</template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { GetTasks, SaveTask, DeleteTask, RunTaskNow, SelectDirectory, SelectFile, GetErrorLogs, ClearErrorLogs } from '../wailsjs/go/main/App'

const activeTab = ref('tasks')
const tasks = ref([])
const dialogVisible = ref(false)
const errorDialogVisible = ref(false)
const errorList = ref([])
const logContent = ref('')
const searchTask = ref('')
const searchDate = ref('')
const form = ref({
  id: '',
  name: '',
  path: '',
  mode: 'retention',
  retentionDays: 7,
  filePattern: '.log',
  cronSpec: '0 2 * * *',
  enabled: true
})

const loadTasks = async () => {
  tasks.value = await GetTasks()
}

const showDialog = (task = null) => {
  if (task) {
    form.value = { ...task }
  } else {
    form.value = { id: '', name: '', path: '', mode: 'retention', retentionDays: 7, filePattern: '.log', cronSpec: '0 2 * * *', enabled: true }
  }
  dialogVisible.value = true
}

const selectPath = async () => {
  const path = await SelectFile()
  if (path) form.value.path = path
}

const selectDir = async () => {
  const path = await SelectDirectory()
  if (path) form.value.path = path
}

const saveTask = async (task) => {
  await SaveTask(task)
  ElMessage.success('保存成功')
  dialogVisible.value = false
  loadTasks()
}

const deleteTask = async (id) => {
  await DeleteTask(id)
  ElMessage.success('删除成功')
  loadTasks()
}

const runNow = async (id) => {
  const result = await RunTaskNow(id)
  if (result.failed > 0) {
    ElMessage({
      type: 'warning',
      message: `执行完成: 成功 ${result.success} 个, 失败 ${result.failed} 个`,
      duration: 5000
    })
    errorList.value = result.errors
    errorDialogVisible.value = true
  } else {
    ElMessage.success(`执行完成: 成功处理 ${result.success} 个文件`)
  }
}

const loadLogs = async () => {
  logContent.value = await GetErrorLogs()
}

const filteredLogs = computed(() => {
  if (!logContent.value) return ''
  let lines = logContent.value.split('\n')
  if (searchTask.value) {
    lines = lines.filter(line => line.toLowerCase().includes(searchTask.value.toLowerCase()))
  }
  if (searchDate.value) {
    lines = lines.filter(line => line.includes(searchDate.value))
  }
  return lines.join('\n')
})

const clearLogs = async () => {
  await ClearErrorLogs()
  logContent.value = ''
  ElMessage.success('日志已清空')
}

onMounted(() => {
  loadTasks()
  loadLogs()
})
</script>

<style>
.container { padding: 20px; }
.header { display: flex; justify-content: space-between; align-items: center; }
</style>