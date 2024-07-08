<template>

  <el-dialog v-model="dialogFormVisibleAdd" title="添加人脸">
    <el-upload
      class="avatar-uploader"
      drag
      action="http://192.168.0.66:8089/file/upload"
      multiple
      v-model:file-list="uploadConfigs.fileList"
      :limit="1"
      :disabled="!uploadConfigs.disabled"
      :on-success="(response, uploadFile, uploadFiles) => handleSuccess(response, uploadFile, uploadFiles)"
      :on-remove="handleRemove"
      :before-upload="beforeAvatarUpload"
    >
      <img v-if="!uploadConfigs.disabled" :src="uploadConfigs.fileList[0].url" class="avatar" />
      <el-icon v-else class="avatar-uploader-icon"><upload-filled /></el-icon>
      <div class="el-upload__text" v-if="uploadConfigs.disabled">
        拖曳文件到此处或
        <em>点击</em>
        以上传文件
      </div>
    </el-upload>
    <el-form style="width: 90%" :model="pinfo" :rules="rules" ref="formRef">
      <el-form-item label="姓名" label-width="100px" prop="Name">
        <el-input placeholder="请输入姓名" v-model="pinfo.Name"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button type="primary" size="default" @click="cancel">取消</el-button>
      <el-button type="primary" size="default" @click="confirmAdd">确定</el-button>
    </template>
  </el-dialog>
  <el-dialog v-model="centerDialogVisible" :title="dialogStatus" width="500" align-center>
    <span>{{ dialogMsg }}</span>
    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="centerDialogVisible = false">确认</el-button>
      </div>
    </template>
  </el-dialog>

  <el-row style="min-height: calc(100vh - 325px); max-height: max-content" :gutter="20">
    <el-col style="min-height: calc(100vh - 325px); max-height: max-content" :span="6">
      <el-card style="margin: 10px 10px">
        <el-button type="primary" size="default" icon="Plus" @click="addPersonInfo" :round="true">添加人脸</el-button>
        <el-table style="margin: 10px 0" border :data="personInfo">
          <el-table-column label="人员名称" prop="Name" width="120px">
          </el-table-column>

          <el-table-column label="操作">
            <template v-slot:default="{ row }">
              <el-tooltip class="item" effect="dark" content="以图搜图" placement="top" :show-after="300">
                <el-button type="primary" size="small" icon="VideoPlay" @click="addSearchTask(row)"></el-button>
              </el-tooltip>
              <el-popconfirm :title="`您确定删除人脸'${row}'`" width="250px" icon="delete" @confirm="removePerson(row)">
                <template #reference>
                  <el-button type="danger" size="small" icon="Delete"></el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
      <el-card style="margin: 10px 10px">

        <el-form :inline="true" class="form">
          <el-form-item label="视频通道">
            <el-select class="selector" v-model="addTaskItem.SrcID" clearable filterable placeholder="请选择">
              <el-option v-for="item in srcIds" :key="item" :label="item" :value="item"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="开始时间">
            <el-date-picker v-model="beginDate" type="datetime" format="YYYY-MM-DD HH:mm" value-format="X" placeholder="选择日期时间"></el-date-picker>
          </el-form-item>
          <el-form-item label="结束时间">
            <el-date-picker v-model="endDate" type="datetime" format="YYYY-MM-DD HH:mm" value-format="X" placeholder="选择日期时间"></el-date-picker>
          </el-form-item>
        </el-form>
      </el-card>
      
    </el-col>
    <el-col :span="18">
      <el-card style="margin: 10px 10px">
        <div>
          <span>搜图结果：</span>
        </div>
        <div style="width: 100%">
          <el-row style="min-height: calc(100vh - 325px); max-height: max-content" :gutter="20">
            <el-col :xs="{ span: 24 }" :lg="{ span: 6 }" :md="{ span: 12 }" :xl="{ span: 4 }" v-for="(item, index) in resItems?.items" :key="index">
              <div class="myCard" :body-style="{ padding: '10px' }" @click="alarmDetail(item)">
                <img :src="item.smallImage" class="image" />
                <div style="padding: 0px">
                  <div class="des_card">
                    <p style="color: blue">
                      <el-icon>
                        <clock />
                      </el-icon>
                      {{ dayjs(item.time * 1000).format('YYYY-MM-DD HH:mm:ss') }}
                    </p>
                    <p style="color: blue">
                      <el-icon>
                        <location />
                      </el-icon>
                      任务：{{ item.taskId }}
                    </p>
                  </div>
                  <div class="des_card">
                    <p style="color: blue">
                      <el-icon><video-camera /></el-icon>
                      视频源：{{ item.srcId }}
                    </p>
                    <p style="color: blue">
                      <el-icon>
                        <tickets />
                      </el-icon>
                      {{ getNameByType(item.type) }}
                    </p>
                  </div>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
        <el-pagination
          v-model:current-page="reqItem.pageNo"
          v-model:page-size="limit"
          :page-sizes="[12, 24, 48]"
          :background="true"
          layout="prev, pager, next,  jumper, ->, sizes, total"
          :total="total"
          @current-change="getSearchResults"
          @size-change="sizeChange"
        />
        <el-dialog v-model="dialogFormVisibleAlarm" title="告警详情" width="50%">
          <div class="dialog-content">
            <!-- 图片部分 -->
            <div class="image-container" style="position: relative">
              <img ref="image" :src="alarmItem?.bigImage" alt="告警图片" @load="onImageLoad(alarmItem?.box as Box)" />
              <!-- 检测框 -->
              <div class="detection-box" :style="detectionBoxStyle"></div>
            </div>

            <!-- 信息部分 -->
            <div class="info-container">
              <div class="info-item">
                <span class="info-title">任务名称：</span>
                <span class="info-content">{{ alarmItem?.taskId }}</span>
              </div>
              <div class="info-item">
                <span class="info-title">告警类型：</span>
                <span class="info-content">{{ getNameByType(alarmItem?.type as number) }}</span>
              </div>
              <div class="info-item">
                <span class="info-title">视频源：</span>
                <span class="info-content">{{ alarmItem?.srcId }}</span>
              </div>
              <div class="info-item">
                <span class="info-title">抓拍时间：</span>
                <span class="info-content">{{ dayjs((alarmItem?.time as number) * 1000).format('YYYY-MM-DD HH:mm:ss') }}</span>
              </div>
              <div class="info-item" v-for="(value, key) in parsedJson(alarmItem?.Extend as string)" :key="key">
                <span class="info-label">{{ key }}</span>
                :
                <span class="info-content">{{ value }}</span>
              </div>
            </div>
          </div>
        </el-dialog>
      </el-card>
    </el-col>
  </el-row>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import { ref, onMounted, nextTick, computed, CSSProperties, watch } from 'vue'
import { reqAlarmList, getInfo, reqDeleteAlarm, reqSearchResList, reqAddSearchTask, reqAddPerson, reqDeletePerson, reqListPerson } from '@/api/face/library'
import type { ResAlarmList, AlarmRes, ResInfo, ReqData, AlarmItem, Box, SearchReqData, SearchResReq, reqPersonInfo, ResPersonList} from '@/api/face/library/type'
import type { ResAbilityList, AbilityList, } from '@/api/face/task/type'
import { getAbilites } from '@/api/face/task'
import { UploadFile, UploadFiles, UploadProps, UploadUserFile } from 'element-plus'
import { reqImageFeature, reqCompareFeature, reqGalleryCreate } from '@/api/face/compare'
import type {PersonInfo} from '@/api/face/library/type'


import { UploadConfig, Features, CompareResult, ResponseData } from '@/api/face/compare/type'
import { reactive } from 'vue'

const centerDialogVisible = ref(false)
const dialogStatus = ref(' ')
const dialogMsg = ref('' as string)

const uploadConfigs = ref<UploadConfig>(
  {
    fileList: [] as UploadUserFile[],
    action: '/file/upload',
    listType: 'picture',
    limit: 1,
    disabled: true,
    fileName: ''
  },
)

const handleSuccess = (response: any, file: UploadFile, fileList: UploadFiles) => {
  uploadConfigs.value.fileList[0].url = '/api/face/alarm/image?fileName=' + file.name
  uploadConfigs.value.fileName = file.name
  console.log(uploadConfigs.value.fileList[0].url)
  let featureNum = getFileFeature()
  console.log()
  console.log(featureNum)
}

const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
  const allowedTypes = ['image/jpeg', 'image/png']
  if (!allowedTypes.includes(rawFile.type)) {
    dialogStatus.value = ' '
    centerDialogVisible.value = true
    dialogMsg.value = '请上传jpg/png格式的图片'
    return false
  } else if (rawFile.size / 1024 / 1024 > 2) {
    dialogStatus.value = ' '
    centerDialogVisible.value = true
    dialogMsg.value = '请上传不超过2MB的图片'
    return false
  }
  return true
}


const getFileFeature = async () => {

  let res: Features = await reqImageFeature({ FileID: "/data/face/upload/" + uploadConfigs.value.fileList[0].name })
  console.log(res)
  if (res.code === 0) {
    if (res.result.length === 0) {
      dialogStatus.value = ' '

      centerDialogVisible.value = true

      dialogMsg.value = '未检测到人脸'
      uploadConfigs.value.fileList = [] as unknown as UploadUserFile[]

      return
    } else if (res.result.length > 1) {
      dialogStatus.value = ' '

      centerDialogVisible.value = true
      dialogMsg.value = '请上传仅包含一张人脸的图片'
      uploadConfigs.value.fileList = [] as unknown as UploadUserFile[]

      return
    }
    uploadConfigs.value.disabled = false
  } else {
    dialogStatus.value = ' '
    centerDialogVisible.value = true

    dialogMsg.value = '未检测到人脸'
    uploadConfigs.value.fileList = [] as unknown as UploadUserFile[]
  }
}

const handleRemove = () => {
  console.log('hello')
}

let limit = ref<number>(12)
let total = ref<number>(0)
let deleteNum = ref<number>(1)

let resItems = ref<AlarmRes>()
let types = ref<number[]>()
let taskIds = ref<string[]>()
let srcIds = ref<string[]>()
let dialogFormVisibleAlarm = ref<boolean>(false)

let alarmItem = ref<AlarmItem>()
let reqItem = ref<ReqData>({
  pageNo: 1,
  pageSize: 10,
  beginTime: 0,
  endTime: 0,
  srcId: '',
  taskId: '',
  trackId: 0,
  types: [],
})

let reqSearchResItem = ref<SearchResReq> ({
  pageNo: 1,
  pageSize: 1,
  ComparisonTaskID: ''
})

let addTaskItem = ref<SearchReqData>({
  ImageFile: '',
  ComparisonTaskID: '',
  SrcID: '',
  Threshold: 0.1,
  BeginTime: 0,
  EndTime: 0,
  Top: 10,
})

let beginDate = ref<string>()
let endDate = ref<string>()

let displayedWidth = ref(0)
let displayedHeight = ref(0)
let image = ref<HTMLImageElement | null>(null)
let position_xy = ref({
  top: 0,
  left: 0,
  width: 0,
  height: 0,
})

let abys = ref<AbilityList[]>([])

onMounted(async () => {
  getSearchResults()
  getInfos()
  getPersonList()
  let res: ResAbilityList = await getAbilites()
  if (res.code === 0) {
    abys.value = res.result.map((item) => ({ name: item.name, type: item.type })).sort((a, b) => a.type - b.type) || []
  }
})

const getSearchResults = async (pager = 1) => {
  reqSearchResItem.value.pageNo = pager
  reqSearchResItem.value.pageSize = limit.value
  let res: ResAlarmList = await reqSearchResList(reqSearchResItem.value)
  if (res.code === 0) {
    resItems.value = res.result
    resItems.value.items = [resItems.value.items[0]]

    total.value = res.result.total
  }
}

const getInfos = async (pager = 1) => {
  try {
    let res: ResInfo = await getInfo()
    if (res.code === 0) {
      types.value = res.result.types
      taskIds.value = res.result.taskIds
      srcIds.value = res.result.srcIds
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

const deleteAlarms = async () => {
  let res: ResInfo = await reqDeleteAlarm({ number: deleteNum.value })
  if (res.code === 0) {
    getSearchResults()
    ElMessage({
      type: 'success',
      message: '删除成功',
    })
    deleteNum.value = 0
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
}


const search = async () => {
  if (beginDate.value) {
    reqItem.value.beginTime = Number(beginDate.value)
  } else {
    reqItem.value.beginTime = 0
  }

  if (endDate.value) {
    reqItem.value.endTime = Number(endDate.value)
  } else {
    reqItem.value.endTime = 0
  }

  addTaskItem.value.ImageFile = '/data/face/upload/' + uploadConfigs.value.fileName
  
  getSearchResults()
  console.log(reqItem)
}

const addSearchTask = async (row : reqPersonInfo) => {
  if (beginDate.value) {
    addTaskItem.value.BeginTime = Number(beginDate.value)
  } else {
    addTaskItem.value.BeginTime = 0
  }

  if (endDate.value) {
    addTaskItem.value.EndTime = Number(endDate.value)
  } else {
    addTaskItem.value.EndTime = 0
  }

  addTaskItem.value.ImageFile = row.ImageFile
  let res: ResponseData = await reqAddSearchTask(addTaskItem.value)

  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg,
    })
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
  getSearchResults()
  console.log(reqItem)
}

function generateRandomString(): string {
  const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
  let result = '';
  for (let i = 0; i < 10; i++) {
    result += letters.charAt(Math.floor(Math.random() * letters.length));
  }
  return result;
}

const reset = () => {
  reqItem.value.pageNo = 1
  reqItem.value.pageSize = 10
  reqItem.value.srcId = ''
  reqItem.value.taskId = ''
  reqItem.value.types = []
  reqItem.value.beginTime = 0
  reqItem.value.endTime = 0

  beginDate.value = ''
  endDate.value = ''
  getSearchResults()
}

const sizeChange = () => {
  getSearchResults()
}

const parsedJson = (str: string) => {
  try {
    return JSON.parse(str)
  } catch (e) {
    console.error('Error parsing JSON string:', e)
    return {}
  }
}
// 图片加载时调用的方法
const onImageLoad = (box: Box) => {
  if (image.value) {
    displayedWidth.value = image.value.offsetWidth
    displayedHeight.value = image.value.offsetHeight

    const img = image.value // 使用ref引用获取图片元素
    const scale = img.clientWidth / img.naturalWidth // 计算缩放比例
    position_xy.value.top = box.LeftTopY * scale
    position_xy.value.left = box.LeftTopX * scale
    position_xy.value.width = (box.RightBtmX - box.LeftTopX) * scale
    position_xy.value.height = (box.RightBtmY - box.LeftTopY) * scale
    console.log(position_xy.value)
    console.log(displayedHeight.value)
  }
}

// 当对话框变为可见时，使用 nextTick 确保 DOM 已经更新
watch(dialogFormVisibleAlarm, async (newVal) => {
  if (newVal) {
    await nextTick()
    // 确保图片已经加载，如果已经加载，执行 onImageLoad
    if (image.value?.complete && alarmItem.value) {
      onImageLoad(alarmItem.value.box)
    }
    // 如果图片还没有加载，onImageLoad 将通过 @load 事件触发
  }
})

// 当 alarmItem 发生变化时（例如，用户选择了一个新的告警项），确保更新图片的 src
watch(
  () => alarmItem.value,
  (newVal) => {
    if (newVal && newVal.bigImage) {
      // 如果有新的 bigImage URL，重置图片加载状态
      if (image.value) {
        image.value.src = '' // 先清空 src 强制重新加载
        image.value.src = newVal.bigImage // 重新设置 src
      }
    }
  },
  { deep: true },
)

// 打开对话框的函数，现在不需要直接调用 onImageLoad
const alarmDetail = (item: AlarmItem) => {
  dialogFormVisibleAlarm.value = true // 这将触发对话框的 watch
  alarmItem.value = item // 这将触发 alarmItem 的 watch
}
const detectionBoxStyle = computed(
  (): CSSProperties => ({
    top: position_xy.value.top + 'px',
    left: position_xy.value.left + 'px',
    width: position_xy.value.width + 'px',
    height: position_xy.value.height + 'px',
    position: 'absolute',
    border: '2px solid red',
  }),
)

// 函数根据type查找name
function getNameByType(type: number): string | undefined {
  const ability = abys.value.find((item) => item.type === type)
  return ability ? ability.name : undefined
}

const validatorName = (rule: any, value: any, callBack: any) => {
  if (value.trim().length >= 2) {
    callBack()
  } else {
    callBack(new Error('名称位数大于等于两位'))
  }
}

const rules = {
  'Name': [
    {
      required: true,
      trigger: 'blur',
      validator: validatorName,
    },
  ],
}

const dialogFormVisibleAdd = ref(false)
let personInfo = ref<reqPersonInfo[]> ([])
let pinfo = ref<reqPersonInfo>({
  Name: '',
  ImageFile: '',
})
let formRef = ref()

const addPersonInfo = () => {
  dialogFormVisibleAdd.value = true

  nextTick(() => {
    formRef.value.clearValidate('Name')
  })
}

const cancel = () => {
  dialogFormVisibleAdd.value = false
}

const confirmAdd = async () => {
  await formRef.value.validate()
  // console.log(addTaskParams)

  pinfo.value.ImageFile = "/data/face/upload/" + uploadConfigs.value.fileName
  let res = await reqAddPerson(pinfo.value)
  if (res.code === 0) {
    dialogFormVisibleAdd.value = false
    ElMessage({
      type: 'success',
      message: res.msg,
    })
    getPersonList()
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
}

let personList = ref<string[]>([] as string[])

const getPersonList = async () => {

  try {
    let res : ResPersonList = await reqListPerson()
    if (res.code === 0) {
      personInfo.value = res.result
      console.log(personInfo.value)
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


const removePerson = async (id: string) => {
  let data : PersonInfo = {
    Name: id
  }

  let res = await reqDeletePerson(data)

  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除人脸成功',
    })

    getPersonList()
  } else {
    ElMessage({
      type: 'error',
      message: '删除失败',
    })
  }
}

</script>

<style lang="scss" scoped>
.form {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  /* 允许按钮换行 */
}

.form-item-buttons {
  display: flex;
  align-items: center;
  /* 垂直居中 */
}

.myCard {
  cursor: pointer;
  // overflow: hidden;
  border: 1px dashed #b4b1b1e7; /* 2px是边框粗细，#000是颜色（黑色），dashed是虚线样式 */
}

.image {
  width: 100%;
  height: 280px;
  object-fit: contain;
  margin-right: 10px;
}

.des_card {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}

.des_card p {
  margin: 1px 1px;
  display: flex;
  align-items: center;
}

.des_card i {
  margin-right: 5px;
}

.el-pagination {
  margin: 10px 0;
}

.dialog-content {
  display: flex;
  // align-items: flex-start; /* 如果你想要顶部对齐 */
  flex-wrap: wrap;
  /* 允许内容换行 */
}

.image-container img {
  flex: 0 0 70%;
  /* flex-grow, flex-shrink, flex-basis */
  max-width: 100%;
  /* 最大宽度为容器的70% */
  margin-right: 20px;
  /* 右侧边距 */
}

.info-container {
  // display: flex;
  flex-direction: column;
}

.info-item {
  margin-bottom: 20px;
  /* 信息条目之间的间距 */
}

.info-title {
  font-weight: bold;
}

.info-content {
  margin-left: 5px;
  /* 标题和内容之间的间距 */
}

.el-select.selector {
  width:10vw
}

.avatar-uploader .el-upload {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
  font-size: 5vw;
  color: #8c939d;
  width: 15vw;
  height: 15vw;
  text-align: center;
}

.el-col.upload-col {
  height: 60vh;
}

</style>


<style scoped>
.avatar-uploader .avatar {
  width: 25vw;
  height: calc(25vw / var(--aspect-ratio));
  --aspect-ratio: 1; /* 定义宽高比，1代表正方形 */
  background-color: lightgray; /* 用于无图像时的背景色 */
  justify-content: center; /* 或者 display: flex; align-items: center; justify-content: center; 用于居中内容 */
  overflow: hidden; /* 遮盖超出边界的图像内容 */
}
</style>