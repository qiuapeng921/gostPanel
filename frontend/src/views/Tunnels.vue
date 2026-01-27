<template>
  <div class="page-container">
    <div class="page-header">
      <h3>隧道管理</h3>
    </div>
    <el-card shadow="hover">
      <!-- 搜索栏 -->
      <div class="search-bar">
        <div class="filters">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索隧道名称"
            :prefix-icon="Search"
            clearable
            style="width: 250px"
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
          <el-select v-model="searchNodeId" placeholder="选择节点" clearable style="width: 180px" @change="handleSearch">
            <el-option v-for="node in nodeList" :key="node.id" :label="node.name" :value="node.id" />
          </el-select>
          <el-select v-model="searchStatus" placeholder="状态" clearable style="width: 120px" @change="handleSearch">
            <el-option label="运行中" value="running" />
            <el-option label="已停止" value="stopped" />
          </el-select>
          <el-button :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="fetchData">刷新</el-button>
        </div>
        <el-button type="primary" :icon="Plus" @click="openDialog()">添加隧道</el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="tunnelList" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="隧道名称" min-width="150" align="center" show-overflow-tooltip />
        <el-table-column label="入口节点" width="140" align="center">
          <template #default="{ row }">
            <el-tag size="small" type="primary">{{ row.entry_node?.name || '-' }}</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="出口节点" width="140" align="center">
          <template #default="{ row }">
            <el-tag size="small" type="success">{{ row.exit_node?.name || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="protocol" label="协议" width="80" align="center">
          <template #default="{ row }">
            <el-tag size="small">{{ row.protocol?.toUpperCase() }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="relay_port" label="Relay端口" width="100" align="center" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button 
              v-if="row.status !== 'running'" 
              type="success" link size="small" 
              @click="handleStart(row)"
            >启动</el-button>
            <el-button 
              v-else 
              type="warning" link size="small" 
              @click="handleStop(row)"
            >停止</el-button>
            <el-button type="primary" link size="small" @click="openDialog(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑隧道' : '添加隧道'"
      width="550px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item label="隧道名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入隧道名称" :prefix-icon="EditPen" />
        </el-form-item>
        <el-divider content-position="left">链路配置</el-divider>
        <el-form-item label="入口节点" prop="entry_node_id">
          <el-select v-model="form.entry_node_id" placeholder="选择入口节点" style="width: 100%" :disabled="isEdit">
            <el-option 
              v-for="node in nodeList" 
              :key="node.id"
              :label="node.name" 
              :value="node.id" 
              :disabled="node.id === form.exit_node_id"
            />
          </el-select>
          <div class="form-hint">客户端连接的节点</div>
        </el-form-item>
        <el-form-item label="出口节点" prop="exit_node_id">
          <el-select v-model="form.exit_node_id" placeholder="选择出口节点" style="width: 100%" :disabled="isEdit">
            <el-option 
              v-for="node in nodeList" 
              :key="node.id" 
              :label="node.name" 
              :value="node.id" 
              :disabled="node.id === form.entry_node_id"
            />
          </el-select>
          <div class="form-hint">流量出口节点，启动时会在该节点创建 Relay 服务</div>
        </el-form-item>
        <el-divider content-position="left">协议配置</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="协议" prop="protocol">
              <el-select v-model="form.protocol" style="width: 100%">
                <el-option label="TCP" value="tcp" />
                <el-option label="UDP" value="udp" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Relay端口" prop="relay_port">
              <el-input-number v-model="form.relay_port" :min="1" :max="65535" controls-position="right" style="width: 100%" />
              <div class="form-hint">出口节点 Relay 服务端口</div>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" :rows="2" placeholder="备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh, Search, EditPen, Connection } from '@element-plus/icons-vue'
import { getTunnelList, createTunnel, updateTunnel, deleteTunnel, startTunnel, stopTunnel } from '@/api/tunnel'
import { getNodeList } from '@/api/node'

// 节点列表
const nodeList = ref([])

// 列表数据
const tunnelList = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 搜索
const searchKeyword = ref('')
const searchNodeId = ref('')
const searchStatus = ref('')

// 对话框
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const submitLoading = ref(false)
const formRef = ref(null)

const form = reactive({
  name: '',
  entry_node_id: '',
  exit_node_id: '',
  protocol: 'tcp',
  relay_port: 8443,
  remark: ''
})

const formRules = {
  name: [{ required: true, message: '请输入隧道名称', trigger: 'blur' }],
  entry_node_id: [{ required: true, message: '请选择入口节点', trigger: 'change' }],
  exit_node_id: [{ required: true, message: '请选择出口节点', trigger: 'change' }],
  protocol: [{ required: true, message: '请选择协议', trigger: 'change' }],
  relay_port: [{ required: true, message: '请输入 Relay 端口', trigger: 'blur' }]
}

// 状态处理
const getStatusType = (status) => {
  const map = { running: 'success', stopped: 'info', error: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = { running: '运行中', stopped: '已停止', error: '错误' }
  return map[status] || status
}

// 获取节点列表
const fetchNodes = async () => {
  try {
    const res = await getNodeList({ pageSize: 100 })
    nodeList.value = res.data.list || []
  } catch (error) {
    console.error('获取节点列表失败:', error)
  }
}

// 获取数据
const fetchData = async (isSilent = false) => {
  if (!isSilent) loading.value = true
  try {
    const res = await getTunnelList({
      page: page.value,
      pageSize: pageSize.value,
      node_id: searchNodeId.value,
      status: searchStatus.value,
      keyword: searchKeyword.value
    })
    tunnelList.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error) {
    console.error('获取隧道列表失败:', error)
  } finally {
    if (!isSilent) loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  page.value = 1
  fetchData()
}

// 打开对话框
const openDialog = (row = null) => {
  isEdit.value = !!row
  editId.value = row?.id || null
  
  if (row) {
    Object.assign(form, {
      name: row.name,
      entry_node_id: row.entry_node_id,
      exit_node_id: row.exit_node_id,
      protocol: row.protocol || 'tcp',
      relay_port: row.relay_port || 8443,
      remark: row.remark || ''
    })
  } else {
    Object.assign(form, {
      name: '',
      entry_node_id: '',
      exit_node_id: '',
      protocol: 'tcp',
      relay_port: 8443,
      remark: ''
    })
  }
  
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitLoading.value = true
    try {
      const submitData = {
        name: form.name,
        entry_node_id: form.entry_node_id,
        exit_node_id: form.exit_node_id,
        protocol: form.protocol,
        relay_port: form.relay_port,
        remark: form.remark
      }
      
      if (isEdit.value) {
        await updateTunnel(editId.value, submitData)
        ElMessage.success('更新成功')
      } else {
        await createTunnel(submitData)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      console.error('操作失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

// 删除隧道
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除隧道 "${row.name}" 吗？如果有规则正在使用此隧道，将无法删除。`, 
      '提示', 
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await deleteTunnel(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 启动隧道
const handleStart = async (row) => {
  try {
    await startTunnel(row.id)
    ElMessage.success('启动成功')
    fetchData()
  } catch (error) {
    console.error('启动失败:', error)
  }
}

// 停止隧道
const handleStop = async (row) => {
  try {
    await stopTunnel(row.id)
    ElMessage.success('停止成功')
    fetchData()
  } catch (error) {
    console.error('停止失败:', error)
  }
}

// 定时刷新
let refreshTimer = null

onMounted(() => {
  fetchNodes()
  fetchData()
  
  // 每 5 秒刷新一次 (静默刷新)
  refreshTimer = setInterval(() => {
    fetchData(true)
  }, 5000)
})

onBeforeUnmount(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.search-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filters {
  display: flex;
  gap: 12px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.form-hint {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}

:deep(.el-table .el-table__cell) {
  padding: 12px 0;
}
</style>
