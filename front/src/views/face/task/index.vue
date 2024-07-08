<template>
  <el-card class="box-card" style="margin: 10px 0">
    <el-button type="primary" size="default" icon="Plus" @click="addTaskItem" :round="true">添加任务</el-button>
    <el-table style="margin: 10px 0" border :data="taskItems">
      <el-table-column label="任务名称" prop="TaskID" width="120px"></el-table-column>
      <el-table-column label="视频源名称" prop="InputSrc.SrcID" width="100px"></el-table-column>
      <el-table-column label="视频源地址" prop="InputSrc.StreamSrc.Address" width="300px"></el-table-column>

      <el-table-column label="状态" width="120px">
        <template v-slot="scope">
          <div class="text-overflow">
            <el-tag :type="scope.row.Status === 0 ? 'danger' : 'success'">
              {{ scope.row.Status === 0 ? '已停止' : '运行中' }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template v-slot:default="{ row }">
          <el-tooltip class="item" effect="dark" content="启动任务" placement="top" :show-after="300">
            <el-button v-if="!Boolean(row.Status)" type="primary" size="small" @click="start(row.TaskID)" icon="TurnOff" color="rgb(244, 180, 180)"></el-button>
          </el-tooltip>
          <el-tooltip class="item" effect="dark" content="停止任务" placement="top" :show-after="300">
            <el-button v-if="Boolean(row.Status)" type="primary" size="small" @click="stop(row.TaskID)" icon="Open" color="rgb(93, 233, 105)"></el-button>
          </el-tooltip>
          <el-tooltip class="item" effect="dark" content="预览图片" placement="top" :show-after="300">
            <el-button type="primary" size="small" icon="PictureFilled" @click="showPic(row)"></el-button>
          </el-tooltip>
          <el-tooltip class="item" effect="dark" content="实时预览" placement="top" :show-after="300">
            <el-button v-if="Boolean(row.Status)" type="primary" size="small" icon="VideoPlay" @click="playPic(row)"></el-button>
          </el-tooltip>
          <el-popconfirm :title="`您确定删除任务'${row.TaskID}'`" width="250px" icon="delete" @confirm="removeTask(row.TaskID)">
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
        <el-input placeholder="请您输入任务名称" v-model="addTaskParams.TaskID"></el-input>
      </el-form-item>
      <el-form-item label="视频源名称" label-width="100px" prop="deviceName">
        <el-input placeholder="请您输入视频源名称" v-model="addTaskParams.InputSrc.SrcID"></el-input>
      </el-form-item>
      <el-form-item label="视频源地址" label-width="100px" prop="deviceName">
        <el-input placeholder="请您输入视频地址(如rtsp地址或视频文件绝对路径)" v-model="addTaskParams.InputSrc.StreamSrc.Address"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button type="primary" size="default" @click="cancel">取消</el-button>
      <el-button type="primary" size="default" @click="confirmAdd">确定</el-button>
    </template>
  </el-dialog>

  <!-- <el-dialog v-model="dialogFormVisibleUpdate" title="修改算法任务">
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
    </el-form>
    <template #footer>
      <el-button type="primary" size="default" @click="cancel">取消</el-button>
      <el-button type="primary" size="default" @click="confirmUpdate">确定</el-button>
    </template>
  </el-dialog> -->

  <el-dialog v-model="dialogFormVisiblePic" title="视频图片预览" width="50%" :top="'3%'">
    <el-button type="primary" size="large" :loading="imageLoading" @click="updateImage">刷新图片</el-button>
    <br />
    <br />
    <div class="image-container">
      <img ref="image" :src="imgSrc" alt="" />
    </div>
  </el-dialog>
  <el-dialog v-model="dialogFormVisibleWs" title="检测结果实时预览" width="70%">
    <div ref="imageContainer" style="position: relative" :class="{ 'blurred': isBlurred }">
      <img :src="imageSource" alt="Streamed Image" class="streamed-image" style="object-fit: cover; position: absolute; z-index: 1; width: 60vw; height: auto;" />
      <canvas id="canvasId" style="position: absolute; top: 0; left: 0; z-index: 1;"></canvas>
    </div>
    <el-dialog
      :model-value="dialogVisible"
      title="识别结果"
      class="dialog-overlay"
      @close="handleFaceClose"
    >
      <img :src="matcherFace" alt=" "/>
      <p>{{matcherName}}</p>
    </el-dialog>
    
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, nextTick, onUnmounted, watchEffect } from 'vue'
import { fabric } from 'fabric'
import { reqTaskList, getAbilites, addTask, deleteTask, startTask, stopTask, getPic} from '@/api/face/task'

import type { TaskListResponseData, RequestBody, AbilityList } from '@/api/face/task/type'
import { WebStreamImpl } from '@/api/face/task/type'

import { isBuiltin } from 'module';

let pageNo = ref<number>(1)
let limit = ref<number>(10)
let total = ref<number>(0)

let dialogFormVisibleAdd = ref<boolean>(false)
let dialogFormVisibleUpdate = ref<boolean>(false)
let dialogFormVisiblePic = ref<boolean>(false)
let dialogFormVisibleWs = ref<boolean>(false)
let showPreview = ref<boolean>(true)
let imageLoading = ref<boolean>(false)

let matcherIds = new Map<string, number>()
let isBlurred = ref<boolean>(false)
let dialogVisible = ref<boolean>(false)
let matcherFace = ref<string>("")
let matcherName = ref<string>("")
let imgSrc = ref<string>('')
let taskItems = ref<RequestBody[]>([])
let addTaskParams = reactive<RequestBody>({
  TaskID: '',
  InputSrc: {
    SrcID: '',
    StreamSrc: {
      Address: '',
    },
  },
  Algorithm: {
    TrackInterval: 1,
    DetectInterval: 5,
    AttributeInclude: true,
    FeatureInclude: false,
    TargetSize: {
      MinDetect: 30,
      MaxDetect: 200,
    },
  },
  Status: 0,
})
let abys = ref<AbilityList[]>([])
let formRef = ref()

const imageSource = ref('') // 图片的数据URL
const canvas = ref()
let width = 1
let height = 1
let ratio = 0

const globalWebStream = ref<WebStreamImpl | null>(null);

const handleFaceClose = () => {
  isBlurred.value = false
  dialogVisible.value = false
  WsPlay(globalWebStream.value.taskid)
}

const initCanvas = (width: number, height: number) => {
  canvas.value = new fabric.Canvas('canvasId', {
    width: width,
    height: height,
    hoverCursor: 'auto',
  })
}
const colors = [
  '#FF0000', // 0 - 红色
  '#00FF00', // 1 - 绿色
  '#0000FF', // 2 - 蓝色
  '#FFFF00', // 3 - 黄色
  '#FF00FF', // 4 - 紫色
  '#00FFFF', // 5 - 青色
  '#FFA500', // 6 - 橙色
  '#800080', // 7 - 紫色
  '#008000', // 8 - 绿色
  '#808080', // 9 - 灰色
]

function updateSize() {
  ratio = 0.6 * window.innerWidth / width;
  initCanvas(width * ratio, height *ratio);
}

window.addEventListener('resize', updateSize);


const WsPlay = (taskid: string) => {
  const currentDomain = "192.168.0.66"
  // const currentDomain = window.location.hostname

  // globalSocket = new WebSocket(`ws://${currentDomain}:9091/preview`)
  globalWebStream.value.taskid = taskid
  globalWebStream.value.connect()
  updateSize()

  globalWebStream.value.socket.onmessage = (event) => {
    const reader = new FileReader()

    reader.onload = () => {
      const text = reader.result
      const data = JSON.parse(text as string)
      if (data.video_id !== taskid) {
        return
      }
      const imageUrl = `data:image/png;base64,${data.frame_img.data}`
      // 释放之前的图片资源
      if (imageSource.value) {
        URL.revokeObjectURL(imageSource.value)
      }
      imageSource.value = imageUrl

      // 提取 preview_objs 数组
      const previewObjects = data.preview_objs
      const viewportWidth: number = window.innerWidth;
      const viewportHeight: number = window.innerHeight;
      // if (width === 0 || height === 0) {
      if (width === 1 || height === 1) {
        width = data.frame_img.width
        height = data.frame_img.height
        updateSize()
        // initCanvas(width * ratio, data.frame_img.height *ratio)
      }

      canvas.value.clear()

      // 对每个目标进行处理
      for (const obj of previewObjects) {

        // 获取track_id的最后一位数字
        const lastDigit = obj.track_id % 10
        const color = colors[lastDigit] // 使用取模确保不会超出数组范围
        // console.log(obj.detect_box)
        const detectBox = obj.detect_box // 获取矩形框信息
        detectBox[0] *= ratio
        detectBox[1] *= ratio
        detectBox[2] *= ratio
        detectBox[3] *= ratio
        const trackid = obj.track_id
        const matcherid = obj.match_id
        const matchername = obj.matched_name
        // if (matcherid !== undefined && matcherid != -1 && matcherIds.has(matcherid)) {
        //   const count = matcherIds.get(matcherid)!
        //   matcherIds.set(matcherid, count + 1);

        //   if (count > 50) {
        //     matcherFace.value = "/api/face/alarm/image?fileName=" + matchername
        //     matcherName.value = matchername
        //     handleMatcher(matcherid);
        //     matcherIds.set(matcherid, 0);
        //   } 
          
        // }else {
        //   matcherIds.set(matcherid, 1);
        // }
        const rect = new fabric.Rect({
          left: detectBox[0],
          top: detectBox[1],
          width: detectBox[2] - detectBox[0],
          height: detectBox[3] - detectBox[1],
          fill: 'transparent',
          stroke: color,
          strokeWidth: 2,
        })
        
        const t = (matchername === undefined) ? "未识别" : matchername.replace(/\.\w+$/, '');

        const text = new fabric.Text(`${t}`, {
          left: detectBox[0], 
          top: detectBox[1] + (detectBox[3] - detectBox[1]) + 5,
          fontSize:20,
          fill: color,
        })

        canvas.value.add(rect)
        canvas.value.add(text)
      }
    }
    reader.readAsText(event.data)
  }
}

const handleMatcher = (matcherid : number) => {
  globalWebStream.value.disconnect()
  isBlurred.value = true
  dialogVisible.value = true
}

const getTasklists = async (pager = 1) => {
  const req = {
    pageNo: pager,
    pageSize: limit.value,
  }
  try {
    let res: TaskListResponseData = await reqTaskList(req)

    if (res.code === 0) {
      total.value = res.result.total
      taskItems.value = res.result.Items
      // console.log(taskItems.value)
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
  const currentDomain = "192.168.0.66"

  // const currentDomain = window.location.hostname
  updateSize();

  // const currentDomain = window.location.hostname
  const domain = `ws://${currentDomain}:9091/preview`
  // const currentDomain = window.location.hostname
  globalWebStream.value = new WebStreamImpl();

  globalWebStream.value.init(domain)
})

const sizeChange = () => {
  getTasklists()
}

const addTaskItem = () => {
  dialogFormVisibleAdd.value = true
  addTaskParams.TaskID = ''
  addTaskParams.InputSrc.SrcID = ''
  addTaskParams.InputSrc.StreamSrc.Address = ''

  nextTick(() => {
    formRef.value.clearValidate('TaskID')
    formRef.value.clearValidate('InputSrc.SrcID')
    formRef.value.clearValidate('InputSrc.StreamSrc.Address')
  })
}


// const updateTask = (row: Item) => {
//   nextTick(() => {
//     formRef.value.clearValidate('taskId')
//     formRef.value.clearValidate('deviceName')
//     formRef.value.clearValidate('types')
//     formRef.value.clearValidate('url')
//   })
//   dialogFormVisibleUpdate.value = true
//   Object.assign(addTaskParams, row)
// }

const start = async (id: string) => {
  let date = {
    TaskID: id,
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
    TaskID: id,
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

const showPic = (row: RequestBody) => {
  dialogFormVisiblePic.value = true
  Object.assign(addTaskParams, row)
  imgSrc.value = `/api/task/image?taskId=${addTaskParams.TaskID}&timestamp=${new Date().getTime()}`
}

const playPic = (row: RequestBody) => {
  if (row.Status === 0) {
    return
  }
  dialogFormVisibleWs.value = true
  WsPlay(row.TaskID)
}

const cancel = () => {
  dialogFormVisibleAdd.value = false
  dialogFormVisibleUpdate.value = false
}

onUnmounted(() => {
  if (globalWebStream.value.socket) {
    globalWebStream.value.disconnect()
  }
  // 释放图片资源
  if (imageSource.value) {
    URL.revokeObjectURL(imageSource.value)
  }

  // 清空Canvas对象
  if (canvas.value) {
    canvas.value.dispose()
  }
})

watchEffect(() => {
  if (!dialogFormVisibleWs.value) {
    console.log('close socket')
    width = 0
    if (globalWebStream.value && globalWebStream.value.socket) {
      globalWebStream.value.disconnect()
    }
    if (imageSource.value) {
      URL.revokeObjectURL(imageSource.value)
      imageSource.value = ''
    }
    if (canvas.value) {
      canvas.value.dispose()
    }
  }
})

const updateImage = async () => {
  imageLoading.value = true
  let date = {
    TaskID: addTaskParams.TaskID,
    url: addTaskParams.InputSrc.StreamSrc.Address,
  }

  let res = await getPic(date)
  if (res.code === 0) {
    imgSrc.value = `/api/task/image?taskId=${addTaskParams.TaskID}&timestamp=${new Date().getTime()}`
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

// const confirmUpdate = async () => {
//   await formRef.value.validate()
//   // console.log(addTaskParams)

//   let res = await modTask(addTaskParams)
//   if (res.code === 0) {
//     dialogFormVisibleUpdate.value = false
//     ElMessage({
//       type: 'success',
//       message: res.msg,
//     })
//     getTasklists(pageNo.value)
//   } else {
//     ElMessage({
//       type: 'error',
//       message: res.msg,
//     })
//   }
// }

const removeTask = async (id: string) => {
  let date = {
    TaskID: id,
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

const rules = {
  TaskId: [
    {
      required: true,
      trigger: 'blur',
      validator: validatorName,
    },
  ],
  'InputSrc.SrcID': [
    {
      required: true,
      trigger: 'blur',
      validator: validatorName,
    },
  ],
  'InputSrc.StreamSrc.Address': [
    {
      required: true,
      trigger: 'blur',
      validator: validatorName,
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

.blurred {
  filter: blur(5px); /* 控制模糊度，根据需要调整 */
}

.dialog-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  /* 调整对话框大小 */
  width: 300px; /* 或根据需要调整 */
}
</style>
