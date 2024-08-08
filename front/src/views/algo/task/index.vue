<template>
  <el-card class="box-card" style="margin: 10px 0">
    <el-button type="primary" size="default" icon="Plus" @click="addTaskItem" :round="true">添加任务</el-button>
    <el-table style="margin: 10px 0" border :data="taskItems">
      <!-- <el-table-column label="序号" width="80px" align="center" type="index"></el-table-column> -->
      <el-table-column label="任务名称" prop="taskId" width="120px"></el-table-column>
      <el-table-column label="视频源" prop="deviceName" width="100px">
        <template v-slot="scope">
          <div class="text-overflow" :title="scope.row.deviceName">
            {{ scope.row.deviceName }}
          </div>
        </template>
      </el-table-column>
      <el-table-column label="编码格式" prop="codeName" width="100px"></el-table-column>
      <el-table-column label="分辨率" width="100px">
        <template v-slot="scope">
          <div class="text-overflow" :title="scope.row.deviceName">{{ scope.row.width }}*{{ scope.row.height }}</div>
        </template>
      </el-table-column>
      <el-table-column label="算法配置信息" prop="abilities" width="250px">
        <template v-slot="scope">
          <span v-for="(ability, index) in scope.row.abilities" :key="index" class="ability-tag">
            {{ ability }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120px">
        <template v-slot="scope">
          <div class="text-overflow">
            <el-tag :type="scope.row.status === 0 ? 'danger' : 'success'">
              {{ scope.row.status === 0 ? '已停止' : '运行中' }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template v-slot:default="{ row }">
          <el-tooltip class="item" effect="dark" content="启动任务" placement="top" :show-after="300">
            <el-button v-if="!Boolean(row.status)" type="primary" size="small" @click="start(row.taskId)" icon="TurnOff" color="rgb(244, 180, 180)"></el-button>
          </el-tooltip>
          <el-tooltip class="item" effect="dark" content="停止任务" placement="top" :show-after="300">
            <el-button v-if="Boolean(row.status)" type="primary" size="small" @click="stop(row.taskId)" icon="Open" color="rgb(93, 233, 105)"></el-button>
          </el-tooltip>
          <el-tooltip class="item" effect="dark" content="预览图片" placement="top" :show-after="300">
            <el-button type="primary" size="small" icon="PictureFilled" @click="showPic(row)"></el-button>
          </el-tooltip>

          <el-tooltip class="item" effect="dark" content="编辑" placement="top" :show-after="300">
            <el-button type="primary" size="small" icon="Edit" @click="updateTask(row)"></el-button>
          </el-tooltip>
          <el-popconfirm :title="`您确定删除任务'${row.taskId}'`" width="250px" icon="delete" @confirm="removeTask(row.taskId)">
            <template #reference>
              <el-button type="danger" size="small" icon="Delete"></el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页器组件
                pagination
                   v-model:current-page:设置分页器当前页码
                   v-model:page-size:设置每一个展示数据条数
                   page-sizes:用于设置下拉菜单数据
                   background:设置分页器按钮的背景颜色
                   layout:可以设置分页器六个子组件布局调整
            -->
    <el-pagination
      v-model:current-page="pageNo"
      v-model:page-size="limit"
      :page-sizes="[10, 20, 30]"
      :background="true"
      layout="prev, pager, next, ->, sizes, total"
      :total="total"
      @current-change="getTasklists"
      @size-change="sizeChange"
    />
  </el-card>

  <el-dialog v-model="dialogFormVisibleAdd" title="添加算法任务">
    <el-form style="width: 90%" :model="addTaskParams" :rules="rules" ref="formRef">
      <el-form-item label="任务名称" label-width="100px" prop="taskId">
        <el-input placeholder="请您输入任务名称" v-model="addTaskParams.taskId"></el-input>
      </el-form-item>
      <el-form-item label="视频源名称" label-width="100px" prop="deviceName">
        <el-input placeholder="请您输入视频源名称" v-model="addTaskParams.deviceName"></el-input>
      </el-form-item>
      <el-form-item label="视频源地址" label-width="100px" prop="deviceName">
        <el-input placeholder="请您输入视频地址(如rtsp地址或视频文件绝对路径)" v-model="addTaskParams.url"></el-input>
      </el-form-item>
      <el-form-item label="算法配置" label-width="100px" prop="types">
        <el-checkbox-group v-model="addTaskParams.types" @change="handleCheckboxChange">
          <el-checkbox v-for="ab in abys" :key="ab.type" :label="ab.type" :border="true">
            {{ ab.name }}
          </el-checkbox>
        </el-checkbox-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button type="primary" size="default" @click="cancel">取消</el-button>
      <el-button type="primary" size="default" @click="confirmAdd">确定</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="dialogFormVisibleUpdate" title="修改算法任务">
    <el-form style="width: 90%" :model="addTaskParams" :rules="rules" ref="formRef">
      <el-form-item label="任务名称" label-width="100px" prop="taskId">
        <el-input placeholder="请您输入任务名称" v-model="addTaskParams.taskId" :disabled="true"></el-input>
      </el-form-item>
      <el-form-item label="视频源名称" label-width="100px" prop="deviceName">
        <el-input placeholder="请您输入视频源名称" v-model="addTaskParams.deviceName"></el-input>
      </el-form-item>
      <el-form-item label="视频源地址" label-width="100px" prop="deviceName">
        <el-input placeholder="请您输入视频地址(如rtsp地址或视频文件绝对路径)" v-model="addTaskParams.url"></el-input>
      </el-form-item>
      <el-form-item label="算法配置" label-width="100px" prop="types">
        <el-checkbox-group v-model="addTaskParams.types">
          <el-checkbox v-for="ab in abys" :key="ab.type" :label="ab.type" :border="true">
            {{ ab.name }}
          </el-checkbox>
        </el-checkbox-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button type="primary" size="default" @click="cancel">取消</el-button>
      <el-button type="primary" size="default" @click="confirmUpdate">确定</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="dialogFormVisiblePic" title="视频图片预览" width="50%" :top="'3%'">
    <el-button type="primary" size="large" :loading="imageLoading" @click="updateImage">刷新图片</el-button>
    <br />
    <br />
    <div class="image-container">
      <img ref="image" :src="imgSrc" alt="" />
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, nextTick } from 'vue'

import { reqTaskList, getAbilites, addTask, deleteTask, modTask, startTask, stopTask, getPic } from '@/api/algo/task'

import type { TaskListResponseData, Item, ResAbilityList, AbilityList } from '@/api/algo/task/type'

let pageNo = ref<number>(1)
let limit = ref<number>(10)
let total = ref<number>(0)

let dialogFormVisibleAdd = ref<boolean>(false)
let dialogFormVisibleUpdate = ref<boolean>(false)
let dialogFormVisiblePic = ref<boolean>(false)
let imageLoading = ref<boolean>(false)

let imgSrc = ref<string>('')

let taskItems = ref<Item[]>([])
let addTaskParams = reactive<Item>({
  taskId: '',
  deviceName: '',
  url: '',
  types: [],
})
let abys = ref<AbilityList[]>([])
let formRef = ref()

const getTasklists = async (pager = 1) => {
  const req = {
    pageNo: pager,
    pageSize: limit.value,
  }
  try {
    let res: TaskListResponseData = await reqTaskList(req)

    if (res.code === 0) {
      total.value = res.result.total
      taskItems.value = res.result.items
    } else {
      ElMessage({
        type: 'error',
        message: '请求失败',
      })
    }
  } catch (error) {
    // 在这里处理超时情况
    ElMessage({
      type: 'error',
      message: '服务器异常',
    })
  }
}

onMounted(async () => {
  getTasklists()
  let res: ResAbilityList = await getAbilites()
  if (res.code === 0) {
    abys.value = res.result.map((item) => ({ name: item.name, type: item.type })).sort((a, b) => a.type - b.type) || []
  }
})

const sizeChange = () => {
  getTasklists()
}

const addTaskItem = () => {
  dialogFormVisibleAdd.value = true
  addTaskParams.taskId = ''
  addTaskParams.deviceName = ''
  addTaskParams.url = ''
  addTaskParams.types = []

  nextTick(() => {
    formRef.value.clearValidate('taskId')
    formRef.value.clearValidate('deviceName')
    formRef.value.clearValidate('types')
    formRef.value.clearValidate('url')
  })
}

const updateTask = (row: Item) => {
  nextTick(() => {
    formRef.value.clearValidate('taskId')
    formRef.value.clearValidate('deviceName')
    formRef.value.clearValidate('types')
    formRef.value.clearValidate('url')
  })
  dialogFormVisibleUpdate.value = true
  Object.assign(addTaskParams, row)
}

const start = async (id: string) => {
  let date = {
    taskId: id,
  }

  let res = await startTask(date)
  if (res.code === 0) {
    dialogFormVisibleAdd.value = false
    ElMessage({
      type: 'success',
      message: res.msg,
    })
    getTasklists(pageNo.value)
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
}

const stop = async (id: string) => {
  let date = {
    taskId: id,
  }

  let res = await stopTask(date)
  if (res.code === 0) {
    dialogFormVisibleAdd.value = false
    ElMessage({
      type: 'success',
      message: res.msg,
    })
    getTasklists(pageNo.value)
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
}

const handleCheckboxChange = (value: any) => {
  // 如果选择的数量超过1，则取消选中除当前选中项之外的所有选项
  if (value.length > 1) {
    addTaskParams.types = [value[value.length - 1]]
  }
}
const showPic = (row: Item) => {
  dialogFormVisiblePic.value = true
  Object.assign(addTaskParams, row)
  imgSrc.value = `/api/task/image?taskId=${addTaskParams.taskId}&timestamp=${new Date().getTime()}`
}

const cancel = () => {
  dialogFormVisibleAdd.value = false
  dialogFormVisibleUpdate.value = false
}

const updateImage = async () => {
  imageLoading.value = true
  let date = {
    taskId: addTaskParams.taskId,
    url: addTaskParams.url,
  }

  let res = await getPic(date)
  if (res.code === 0) {
    imgSrc.value = `/api/task/image?taskId=${addTaskParams.taskId}&timestamp=${new Date().getTime()}`
    imageLoading.value = false
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
    imageLoading.value = false
  }
  imageLoading.value = false
}

const confirmAdd = async () => {
  await formRef.value.validate()
  // console.log(addTaskParams)

  let res = await addTask(addTaskParams)
  if (res.code === 0) {
    dialogFormVisibleAdd.value = false
    ElMessage({
      type: 'success',
      message: res.msg,
    })
    getTasklists(pageNo.value)
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
}

const confirmUpdate = async () => {
  await formRef.value.validate()
  // console.log(addTaskParams)

  let res = await modTask(addTaskParams)
  if (res.code === 0) {
    dialogFormVisibleUpdate.value = false
    ElMessage({
      type: 'success',
      message: res.msg,
    })
    getTasklists(pageNo.value)
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
}

const removeTask = async (id: string) => {
  let date = {
    taskId: id,
  }
  let res = await deleteTask(date)

  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除任务成功',
    })

    getTasklists(pageNo.value)
  } else {
    ElMessage({
      type: 'error',
      message: '删除失败',
    })
  }
}

const validatorName = (rule: any, value: any, callBack: any) => {
  if (value.trim().length >= 2) {
    callBack()
  } else {
    callBack(new Error('名称位数大于等于两位'))
  }
}

const validatorTypes = (rule: any, value: any, callBack: any) => {
  if (value.length > 0) {
    callBack()
  } else {
    callBack(new Error('请选择算法配置'))
  }
}

const rules = {
  taskId: [
    {
      required: true,
      trigger: 'blur',
      validator: validatorName,
    },
  ],
  deviceName: [
    {
      required: true,
      trigger: 'blur',
      validator: validatorName,
    },
  ],
  url: [
    {
      required: true,
      trigger: 'blur',
      validator: validatorName,
    },
  ],
  types: [
    {
      required: true,
      trigger: 'change',
      validator: validatorTypes,
    },
  ],
}
</script>

<style>
.el-checkbox {
  margin: 5px;
  display: flex;
  flex-wrap: wrap;
  align-items: right;
}

.text-overflow {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: block;
  /* 确保div是块级元素 */
}

.test {
  color: rgb(93, 233, 105);
}

.image-container img {
  flex: 0 0 70%; /* flex-grow, flex-shrink, flex-basis */
  max-width: 100%; /* 最大宽度为容器的70% */
  margin-right: 20px; /* 右侧边距 */
}

.ability-tag {
  background-color: #409eff; /* Element UI的主题蓝色 */
  color: white;
  padding: 3px 6px;
  border-radius: 4px;
  margin-right: 4px;
  font-size: 0.8em;
  display: inline-block; /* 使其显示为行内块 */
  margin: 2px;
}
</style>
