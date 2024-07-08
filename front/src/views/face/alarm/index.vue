<template>
  <el-card style="margin: 10px 10px">
    <el-form :inline="true" class="form">
      <el-form-item label="任务">
        <el-select class="selector" v-model="reqItem.taskId" clearable filterable placeholder="请选择">
          <el-option v-for="item in taskIds" :key="item" :label="item" :value="item"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="视频通道">
        <el-select class="selector" v-model="reqItem.srcId" clearable filterable placeholder="请选择">
          <el-option v-for="item in srcIds" :key="item" :label="item" :value="item"></el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="开始时间">
        <el-date-picker v-model="beginDate" type="datetime" format="YYYY-MM-DD HH:mm" value-format="X" placeholder="选择日期时间"></el-date-picker>
      </el-form-item>
      <el-form-item label="结束时间">
        <el-date-picker v-model="endDate" type="datetime" format="YYYY-MM-DD HH:mm" value-format="X" placeholder="选择日期时间"></el-date-picker>
      </el-form-item>
      <el-form-item style="float: right; text-align: right">
        <el-button type="primary" size="default" @click="search">搜索</el-button>
        <el-button size="default" @click="reset">重置</el-button>
      </el-form-item>
    </el-form>
  </el-card>

  <el-card style="margin: 10px 10px">
    <div>
      <span>删除告警图片：</span>
      <el-input-number :min="1" :max="resItems?.total" v-model="deleteNum" />
      <el-button style="margin: 1px 10px" type="primary" size="default" @click="deleteAlarms">删除</el-button>
      <span>告警图片总数：{{ resItems?.total }}</span>
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
      @current-change="getAlarmLists"
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
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import { ref, onMounted, nextTick, computed, CSSProperties, watch } from 'vue'
import { reqAlarmList, getInfo, reqDeleteAlarm } from '@/api/face/alarm'
import type { ResAlarmList, AlarmRes, ResInfo, ReqData, AlarmItem, Box } from '@/api/face/alarm/type'
import type { ResAbilityList, AbilityList } from '@/api/face/task/type'
import { getAbilites } from '@/api/face/task'

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
  getAlarmLists()
  getInfos()
  let res: ResAbilityList = await getAbilites()
  if (res.code === 0) {
    abys.value = res.result.map((item) => ({ name: item.name, type: item.type })).sort((a, b) => a.type - b.type) || []
  }
})

const getAlarmLists = async (pager = 1) => {
  reqItem.value.pageNo = pager
  reqItem.value.pageSize = limit.value

  let res: ResAlarmList = await reqAlarmList(reqItem.value)
  if (res.code === 0) {
    resItems.value = res.result
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
    getAlarmLists()
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

const search = () => {
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

  getAlarmLists()
  console.log(reqItem)
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
  getAlarmLists()
}

const sizeChange = () => {
  getAlarmLists()
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
</style>
