<template>
  <el-card style="margin: 10px 0" shadow="never">
    <span>选择算法任务：</span>
    <el-select v-model="taskId" clearable filterable placeholder="请选择任务" @change="SelectTask">
      <el-option v-for="item in taskItems" :key="item.taskId" :label="item.taskId" :value="item.taskId"></el-option>
    </el-select>
  </el-card>
  <el-card style="margin: 10px 0" shadow="never">
    <el-descriptions title="视频源详情" class="channel" :column="4" ref="descRef" style="margin: 10px 0px">
      <el-descriptions-item label="名称：">
        {{ taskConfig?.device.name }}
      </el-descriptions-item>
      <el-descriptions-item label="视频源地址：">
        {{ taskConfig?.device.url }}
      </el-descriptions-item>
      <el-descriptions-item label="分辨率：">
        {{ taskConfig?.device.resolution }}
      </el-descriptions-item>
      <el-descriptions-item label="编码格式：">
        {{ taskConfig?.device.codeName }}
      </el-descriptions-item>
    </el-descriptions>

    <h2 class="h2">算法功能列表</h2>
    <br />
    <el-tabs v-model="activeKey" @tab-click="handleTabChange" type="card">
      <el-tab-pane v-for="item in taskConfig?.algorithms" :key="item.Type" :label="getAbName(Number(item.Type))" :name="item.Type" />
    </el-tabs>
    <div style="display: flex; justify-content: flex-end; width: 820px">
      <el-button size="large" type="primary" @click="submit">保存</el-button>
      <el-button size="large" type="primary" @click="reset">重置</el-button>
    </div>

    <h2 class="h2">参数设置</h2>
    <br />
    <el-form :model="taskConfVal" :inline="true" :rules="rules" label-width="auto">
      <el-form-item label="检测帧间隔" style="width: 250px" prop="DetectInterval">
        <el-input v-model.number="taskConfVal.DetectInterval" size="small"><template #append>帧</template></el-input>
      </el-form-item>
      <el-form-item label="连续跟踪帧数" style="width: 250px" prop="TrackInterval">
        <el-input v-model.number="taskConfVal.TrackInterval" size="small"><template #append>帧</template></el-input>
      </el-form-item>
      <el-form-item label="告警间隔" style="width: 250px" prop="AlarmInterval">
        <el-input v-model.number="taskConfVal.AlarmInterval" size="small"><template #append>帧</template></el-input>
      </el-form-item>
      <br />
      <el-form-item label="检测阈值" style="width: 250px" prop="threshold">
        <el-input v-model.number="taskConfVal.threshold" size="small"><template #append>1~100</template></el-input>
      </el-form-item>
      <el-form-item label="目标最大像素" style="width: 250px" prop="MaxDetect">
        <el-input v-model.number="taskConfVal.MaxDetect" size="small"><template #append>px</template></el-input>
      </el-form-item>
      <el-form-item label="目标最小像素" style="width: 250px" prop="MinDetect">
        <el-input v-model.number="taskConfVal.MinDetect" size="small"><template #append>px</template></el-input>
      </el-form-item>
    </el-form>

    <el-button size="small" :type="drawRegion ? 'primary' : 'default'" @click="drawPolygon">绘制检测区域</el-button>
    <el-button size="small" :type="drawLine ? 'primary' : 'default'" @click="drawLineDerect">绘制检测线</el-button>

    <el-button size="small" type="default" @click="clearDraw">清除绘制</el-button>

    <br />
    <br />
    <el-button size="small" :type="rectMode ? 'primary' : 'default'" v-if="drawRegion" @click="drawRect">绘制矩形检测区域</el-button>
    <el-button size="small" :type="drawMode ? 'primary' : 'default'" v-if="drawRegion || drawLine" @click="changeDrawMode">{{ drawText }}</el-button>
    <br />
    <br />
    <div class="video" :style="{ width: `${800}px`, height: `${450}px` }" style="position: relative; border: 1px solid #f7faff">
      <img style="width: 100%; height: 100%; object-fit: cover; position: absolute; z-index: 1" :src="imgSrc" alt="" />
      <canvas id="canvasId" style="position: absolute; top: 0; left: 0; z-index: 1"></canvas>
    </div>
  </el-card>
</template>

<script lang="ts" setup>
import { ref, onMounted, nextTick, watch, reactive } from 'vue'
import { fabric } from 'fabric'
import type { FormRules } from 'element-plus'
import { reqTaskList, getAbilites } from '@/api/algo/task'
import type { TaskListResponseData, Item, ResAbilityList, AbilityList } from '@/api/algo/task/type'

import { reqTaskConfig, ModTaskConfig } from '@/api/algo/config'
import type { ResConfig, ConfigResult, ModAlgorithm, Algorithm, ResponseData } from '@/api/algo/config/type'

const taskItems = ref<Item[]>([])
const taskId = ref<string>('')
const taskConfig = ref<ConfigResult>()
const taskConfVal = reactive({
  TrackInterval: 0,
  DetectInterval: 0,
  MinDetect: 0,
  MaxDetect: 0,
  threshold: 0,
  AlarmInterval: 0,
})

const imgSrc = ref<string>('')

let abys = ref<AbilityList[]>([])
// 绘制
const canvas = ref()
const drawRegion = ref(false) // 绘制区域，打开在视频上显示从接口获取的区域
const drawLine = ref(false) // 绘制检测线，打开在视频上显示从接口获取的检测线
const drawMode = ref(false) // 打开绘制区域或者绘制检测线之后，打开绘制，可在视频上绘制对应图形
const rectMode = ref(false) // 快速绘制矩形模式
const pylogonPoints = ref<any>([]) //当前绘制了的所有点数据
const lastPoint = ref() // 上一个点
const activeLine = ref() //绘制草稿，移动时跟随鼠标移动的线对象
const withArrow = ref(false)
const draftLines = ref([{}]) //绘制草稿，线对象
const draftPoints = ref([{}]) //绘制草稿，点对象
const pylogonObject = ref() //当前区域对象
const pylogonSubmitPoint = ref<any>([]) //当前区域的点数组
const lineObject = ref() //当前检测线对象
const lineSubmitPoint = ref() //当前检测线的点
const originPaintData = ref()
const rectDownPoint = ref() //绘制矩形-按下的点
const rectUpPoint = ref() //绘制矩形-抬起的点
const activeRect = ref() //按下之后移动的临时矩形对象
const drawText = ref('')

// video 宽高计算所需值
const pixRatioX = ref() // 当前canvas宽高与标准宽高的缩放比，兼容不同大小的canvas提交的绘制
const pixRatioY = ref() // 当前canvas宽高与标准宽高的缩放比，兼容不同大小的canvas提交的绘制

const activeKey = ref(1) //当前算法任务值

interface Point {
  x: number
  y: number
}

// 初始化 Fabric.js 的 Canvas
onMounted(async () => {
  getTasklists()
  initCanvas()
  let res: ResAbilityList = await getAbilites()
  if (res.code === 0) {
    abys.value = res.result.map((item) => ({ name: item.name, type: item.type })).sort((a, b) => a.type - b.type) || []
  }
  // console.log(abys.value)
})

const rules: any = reactive({
  TrackInterval: [
    { required: true, message: '请输入', trigger: 'blur' },
    { type: 'number', min: 1, message: '大于0整数', trigger: ['blur', 'change'] },
  ],
  MinDetect: [
    { required: true, message: '请输入', trigger: 'blur' },
    { type: 'number', min: 1, message: '大于0整数', trigger: ['blur', 'change'] },
  ],
  MaxDetect: [
    { required: true, message: '请输入', trigger: 'blur' },
    { type: 'number', min: 1, message: '大于0整数', trigger: ['blur', 'change'] },
  ],
  threshold: [
    { required: true, message: '请输入', trigger: 'blur' },
    { type: 'number', min: 1, max: 100, message: '1~100之间整数', trigger: ['blur', 'change'] },
  ],
  DetectInterval: [
    { required: true, message: '请输入', trigger: 'blur' },
    { type: 'number', min: 1, message: '大于0整数', trigger: ['blur', 'change'] },
  ],
  AlarmInterval: [
    { required: true, message: '请输入', trigger: 'blur' },
    { type: 'number', min: 1, message: '大于0整数', trigger: ['blur', 'change'] },
  ],
})

watch([drawRegion, drawLine], () => {
  if (drawRegion.value) {
    drawText.value = '绘制多边形区域'
  } else if (drawLine.value) {
    drawText.value = '绘制检测线'
  }
})

const getTasklists = async (pager = 1) => {
  const req = {
    pageNo: -1,
    pageSize: 20,
  }
  let res: TaskListResponseData = await reqTaskList(req)

  if (res.code === 0) {
    taskItems.value = res.result.items
    if (taskItems.value.length > 0) {
      taskId.value = taskItems.value[0].taskId
      SelectTask()
    }
  } else {
    ElMessage({
      type: 'error',
      message: '请求失败',
    })
  }
}

const SelectTask = async () => {
  if (!taskId.value) return

  const req = {
    taskId: taskId.value,
  }
  let res: ResConfig = await reqTaskConfig(req)

  if (res.code === 0) {
    taskConfig.value = res.result
    imgSrc.value = `/api/task/image?taskId=${taskId.value}&timestamp=${new Date().getTime()}`
    pixRatioX.value = taskConfig.value.device.width / 800
    pixRatioY.value = taskConfig.value.device.height / 450
    // 绘制的点坐标适配当前canvas
    const points = taskConfig.value.algorithms[0].DetectInfos[0].HotArea
    const line = taskConfig.value.algorithms[0].DetectInfos[0].TripWire
    originPaintData.value = [points, line]
    initObjectByApi(points, line)

    activeKey.value = taskConfig.value.algorithms[0].Type

    handleTabChange()
    drawLine.value = false
    drawRegion.value = false
    clearCanvas()
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
}

function handleTabChange() {
  clearCanvas()
  nextTick(() => {
    const algo = taskConfig.value?.algorithms.filter((v) => v.Type === activeKey.value)[0]
    // console.log(activeKey.value)
    initDrawSelect()
    clearDraftDraw()
    // 绘制的点坐标适配当前canvas
    if (algo) {
      const points = algo.DetectInfos[0].HotArea
      const line = algo.DetectInfos[0].TripWire
      originPaintData.value = [points, line]
      initObjectByApi(points, line)

      taskConfVal.TrackInterval = algo.TrackInterval
      taskConfVal.DetectInterval = algo.DetectInterval
      taskConfVal.threshold = algo.threshold
      taskConfVal.MinDetect = algo.TargetSize.MinDetect
      taskConfVal.MaxDetect = algo.TargetSize.MaxDetect
      taskConfVal.AlarmInterval = algo.AlarmInterval
      // console.log(taskConfVal)
    }
  })
}

function initCanvas() {
  const width = 800
  const height = 450
  canvas.value = new fabric.Canvas('canvasId', {
    width: width,
    height: height,
    hoverCursor: 'auto',
  })
}

/**
 * 开启关闭区域
 * 先清除删除当前canvas，设置drawmode=false、设置多边形可选、更新darwregion值
 * 若drawregion=true:
 *    置drawline绘制检测线为false，只能绘制一个东西
 *    初始化canvas
 *    如果有多边形object，就把它绘制到canvas上
 */
function drawPolygon() {
  // pylogonObject.value?.set({ selectable: true });
  clearCanvas()
  drawRegion.value = !drawRegion.value
  if (drawRegion.value) {
    drawLine.value = false
    nextTick(() => {
      initCanvas()
      if (pylogonObject.value && canvas.value) {
        canvas.value.add(pylogonObject.value)
      }
    })
  }
}

// 快捷绘制-矩形
const drawRect = () => {
  rectMode.value = !rectMode.value
  if (rectMode.value) {
    clearDraftDraw()
    drawLine.value = false
    drawMode.value && changeDrawMode()
    //监听鼠标down和up
    canvas.value.on('mouse:down', onMouseDownRect).on('mouse:up', onMouseUpRect).on('mouse:move', onMouseMoveRect)
    canvas.value.defaultCursor = 'crosshair' // 画布光标样式设置为十字
  } else {
    canvas.value.off('mouse:down', onMouseDownRect).off('mouse:up', onMouseUpRect).off('mouse:move', onMouseMoveRect)
    canvas.value.defaultCursor = 'auto' // 画布光标样式设置为十字
  }
}
function onMouseDownRect(event: fabric.IEvent) {
  if (rectMode.value) {
    rectDownPoint.value = event.pointer
  }
}

function onMouseMoveRect(event: fabric.IEvent) {
  if (rectMode.value) {
    if (rectDownPoint.value) {
      getRectPoints(rectDownPoint.value, event.pointer as Point)
      if (!activeRect.value) {
        activeRect.value = makePylogon(pylogonPoints.value)
        canvas.value.add(activeRect.value)
      } else {
        // 更新现有的矩形对象的点
        activeRect.value.set({ points: pylogonPoints.value })
        activeRect.value.setCoords() // 更新对象的坐标
      }
      canvas.value.renderAll()
    }
  }
}

function getRectPoints(downPoint: Point, upPoint: Point) {
  const ponit2 = { x: downPoint.x, y: upPoint.y }
  const point4 = { x: upPoint.x, y: downPoint.y }
  pylogonPoints.value = [downPoint, ponit2, upPoint, point4]
}

// up的时候绘制矩形
function onMouseUpRect(event: fabric.IEvent) {
  // 避免单击
  if (JSON.stringify(rectDownPoint.value) === JSON.stringify(event.pointer)) {
    clearDraftDraw()
    return
  }
  if (rectMode.value) {
    rectUpPoint.value = event.pointer
    getRectPoints(rectDownPoint.value as Point, rectUpPoint.value as Point)
    pylogonSubmitPoint.value = pylogonPoints.value
    pylogonObject.value && canvas.value.remove(pylogonObject.value) //多边形清除
    const thisPylogon = makePylogon(pylogonPoints.value)
    pylogonObject.value = thisPylogon
    canvas.value.add(thisPylogon)
    clearDraftDraw()
    // console.log(pylogonObject.value)
  }
}

// 绘制检测线
function drawLineDerect() {
  // lineObject.value?.set({ selectable: true });
  clearCanvas()
  drawLine.value = !drawLine.value
  // console.log(drawLine.value)
  if (drawLine.value) {
    drawRegion.value = false
    nextTick(() => {
      initCanvas()
      if (lineObject.value && canvas.value) {
        canvas.value.add(lineObject.value)
      }
    })
  }
}

function onMouseDown(event: fabric.IEvent) {
  if (drawMode.value) {
    let opts = {}
    if (pylogonPoints.value.length && lastPoint.value) {
      // 绘制直线
      const thisLine = makeLine([lastPoint.value, event.pointer])
      canvas.value.add(thisLine)
      draftLines.value.push(thisLine)
    }
    if (drawRegion.value && pylogonPoints.value.length === 0) opts = { radius: 15, fill: 'red', stroke: 'red' }
    pylogonPoints.value.push(event.pointer)
    const thisPoint = makePoint(event.pointer, opts)
    canvas.value.add(thisPoint)
    draftPoints.value.push(thisPoint)
    // 完成多边形绘制
    if (pylogonPoints.value.length > 3 && checkPointIntersect(pylogonPoints.value[0], event.pointer)) {
      pylogonPoints.value.pop()
      const thisPylogon = makePylogon(pylogonPoints.value)
      canvas.value.add(thisPylogon)
      pylogonObject.value && canvas.value.remove(pylogonObject.value)
      pylogonObject.value = thisPylogon
      pylogonSubmitPoint.value = pylogonPoints.value
      clearDraftDraw()
      changeDrawMode()
      return
    }
    // 完成检测线绘制
    if (drawLine.value && withArrow.value) {
      const lines = makeGroup([...draftLines.value])
      lineObject.value && canvas.value.remove(lineObject.value)
      lineObject.value = lines
      lineSubmitPoint.value = [...pylogonPoints.value, lastPoint.value]
      canvas.value.add(lines)
      clearDraftDraw()
      changeDrawMode()
      return
    }
    lastPoint.value = event.pointer
  }
}

// 开启/关闭绘制 多边形
function changeDrawMode() {
  if (!canvas.value) return
  drawMode.value = !drawMode.value
  if (drawMode.value) {
    rectMode.value && drawRect()
    // 开启绘制
    clearDraftDraw() //初始化草稿
    canvas.value.defaultCursor = 'crosshair' // 画布光标样式设置为十字
    canvas.value.on('mouse:down', onMouseDown).on('mouse:move', onMousemove).on('mouse:dblclick', onDbclick) //开启鼠标监听
  } else {
    canvas.value.off('mouse:down', onMouseDown).off('mouse:move', onMousemove).off('mouse:dblclick', onDbclick)
    canvas.value.defaultCursor = 'auto'
  }
}
function onMousemove(event: fabric.IEvent) {
  if (drawMode.value) {
    const movePoint = event.pointer
    if (drawLine.value && pylogonPoints.value.length === 2) {
      lastPoint.value = getCenterPoint(pylogonPoints.value[0], pylogonPoints.value[1])
      withArrow.value = true
    }
    if (activeLine.value) {
      canvas.value.remove(activeLine.value)
    }
    if (lastPoint.value) {
      activeLine.value = makeLine([lastPoint.value, movePoint])
      canvas.value.add(activeLine.value)
    }
  }
}

function onDbclick() {
  clearDraftDraw()
}
// 清除绘制
function clearDraw() {
  if (canvas.value) {
    canvas.value.clear()
    initDraft()
    if (drawLine.value) {
      lineObject.value = null
      lineSubmitPoint.value = []
    }
    if (drawRegion.value) {
      pylogonObject.value = null
      pylogonSubmitPoint.value = []
    }
  }
}

// 重置，重置表单和绘制
function reset() {
  clearCanvas()
  initDrawSelect()
  clearDraftDraw()
  const [points, line] = originPaintData.value
  initObjectByApi(points, line)
  handleTabChange()
}

function initDrawSelect() {
  drawRegion.value = false
  drawLine.value = false
  drawMode.value = false
  rectMode.value = false
}

function makeLine(points: any, opts = {}, needArrow = false) {
  const pointArray: number[] = []
  points.forEach((v: any) => {
    pointArray.push(v.x, v.y)
  })
  let color = 'green'
  let arrow
  if (withArrow.value || needArrow) {
    color = 'red'
    arrow = makeTriangle(points[0], points[1])
  }
  const obj = new fabric.Line(pointArray, {
    fill: color,
    stroke: color,
    strokeWidth: 2,
    evented: false,
    selectable: false,
    hasControls: false,
    hasBorders: false,
    strokeUniform: true,
    ...opts,
  })
  // 返回带箭头的线
  if (arrow) return makeGroup([obj, arrow])
  return obj
}
function makePoint(point: any, opts = {}) {
  const obj = new fabric.Circle({
    originY: 'center',
    originX: 'center',
    left: point.x,
    top: point.y,
    strokeWidth: 1,
    radius: 4,
    fill: 'green',
    stroke: 'green',
    selectable: false,
    hasControls: false,
    hasBorders: false,
    strokeUniform: true,
    ...opts,
  })
  return obj
}

function makeGroup(objs: any, opts = {}) {
  const group = new fabric.Group(objs, {
    selectable: false,
    hasBorders: true,
    strokeUniform: true,
    ...opts,
  })
  return group
}

// 绘制箭头
function makeTriangle(point1: any, point2: any, opts = {}) {
  let angle = 0
  const X = point2.x - point1.x
  const Y = point2.y - point1.y
  if (X === 0 && Y === 0) return null
  else if (X === 0 && Y > 0) {
    angle = 180
  } else if (Y === 0) {
    if (X > 0) angle = 90
    else angle = -90
  } else if (X > 0 && Y > 0) {
    angle = getDegree(Math.atan(Y / X)) + 90
  } else if (X > 0 && Y < 0) {
    angle = getDegree(Math.atan(-X / Y))
  } else if (X < 0 && Y > 0) {
    angle = getDegree(Math.atan(Y / X)) - 90
  } else if (X < 0 && Y < 0) {
    angle = getDegree(Math.atan(-X / Y))
  }
  const obj = new fabric.Triangle({
    originY: 'center',
    originX: 'center',
    left: point2.x + 1,
    top: point2.y + 1,
    width: 16, // 底边长度
    height: 16, // 底边到对角的距离（三角形的高）
    fill: 'red',
    angle: angle,
    ...opts,
  })
  return obj
}

function initObjectByApi(points: any, line: any) {
  if (points.length) {
    const convertPoints = points.map((item: any) => ({
      x: item.X / pixRatioX.value,
      y: item.Y / pixRatioY.value,
    }))
    pylogonObject.value = makePylogon(convertPoints)
    pylogonSubmitPoint.value = convertPoints
  }
  if (line) {
    const p1 = { x: line.LineStart.X / pixRatioX.value, y: line.LineStart.Y / pixRatioY.value }
    const p2 = { x: line.LineEnd.X / pixRatioX.value, y: line.LineEnd.Y / pixRatioY.value }
    const p3 = { x: line.DirectStart.X / pixRatioX.value, y: line.DirectStart.Y / pixRatioY.value }
    const p4 = { x: line.DirectEnd.X / pixRatioX.value, y: line.DirectEnd.Y / pixRatioY.value }
    const line1 = makeLine([p1, p2])
    const line2 = makeLine([p3, p4], {}, true)
    lineObject.value = makeGroup([line1, line2])
    lineSubmitPoint.value = [p1, p2, p4, p3]
  }
}

function getDegree(angle: number) {
  return angle * (180 / Math.PI)
}

// 检测两个点是否重叠，
function checkPointIntersect(p1: any, p2: any) {
  return Math.abs(p1.x - p2.x) <= 15 && Math.abs(p1.y - p2.y) <= 15
}

// 获取两个点的中点
function getCenterPoint(p1: any, p2: any) {
  return { x: (p1.x + p2.x) / 2, y: (p1.y + p2.y) / 2 }
}

const getAbName = (val: number) => {
  const tem = abys.value.filter((ability) => ability.type === val).map((ability) => ability.name)
  return tem[0]
}

function clearDraftDraw() {
  if (canvas.value) {
    canvas.value.remove(...draftLines.value)
    canvas.value.remove(...draftPoints.value)
    canvas.value.remove(activeLine.value)
    canvas.value.remove(activeRect.value)
  }
  initDraft()
  // changeDrawMode();
}

function makePylogon(points: Point[], opts = {}) {
  const obj = new fabric.Polygon(points, {
    fill: 'transparent',
    strokeWidth: 2,
    stroke: 'green',
    objectCaching: false,
    transparentCorners: false,
    selectable: false,
    strokeUniform: true,
    ...opts,
  })
  return obj
}

// 清空画布，退出绘画模式
function clearCanvas() {
  if (canvas.value) {
    canvas.value.dispose()
    canvas.value = null
  }
  drawMode.value = false
  // rectMode.value = false
  // drawLine.value = false
}

function initDraft() {
  draftLines.value = []
  draftPoints.value = []
  pylogonPoints.value = []
  activeLine.value = null
  lastPoint.value = null
  withArrow.value = false
  activeRect.value = null
  rectDownPoint.value = null
  rectUpPoint.value = null
}
async function submit() {
  // console.log(pylogonSubmitPoint)
  // console.log(lineSubmitPoint.value)
  const params: any = {
    TaskID: taskId.value,
    Algorithm: {
      Type: activeKey.value,
      TrackInterval: taskConfVal.TrackInterval,
      DetectInterval: taskConfVal.DetectInterval,
      AlarmInterval: taskConfVal.AlarmInterval,
      TargetSize: {
        MinDetect: taskConfVal.MinDetect,
        MaxDetect: taskConfVal.MaxDetect,
      },
      threshold: taskConfVal.threshold,
      DetectInfos: [
        {
          TripWire: {
            LineStart: {
              X: lineSubmitPoint.value && lineSubmitPoint.value.length > 0 && pixRatioX.value != null ? Math.round(lineSubmitPoint.value[0].x * pixRatioX.value) : 0,
              Y: lineSubmitPoint.value && lineSubmitPoint.value.length > 0 && pixRatioX.value != null ? Math.round(lineSubmitPoint.value[0].y * pixRatioY.value) : 0,
            },
            LineEnd: {
              X: lineSubmitPoint.value && lineSubmitPoint.value.length > 0 && pixRatioX.value != null ? Math.round(lineSubmitPoint.value[1].x * pixRatioX.value) : 0,
              Y: lineSubmitPoint.value && lineSubmitPoint.value.length > 0 && pixRatioX.value != null ? Math.round(lineSubmitPoint.value[1].y * pixRatioY.value) : 0,
            },
            DirectStart: {
              X: lineSubmitPoint.value && lineSubmitPoint.value.length > 0 && pixRatioX.value != null ? Math.round(lineSubmitPoint.value[3].x * pixRatioX.value) : 0,
              Y: lineSubmitPoint.value && lineSubmitPoint.value.length > 0 && pixRatioX.value != null ? Math.round(lineSubmitPoint.value[3].y * pixRatioY.value) : 0,
            },
            DirectEnd: {
              X: lineSubmitPoint.value && lineSubmitPoint.value.length > 0 && pixRatioX.value != null ? Math.round(lineSubmitPoint.value[2].x * pixRatioX.value) : 0,
              Y: lineSubmitPoint.value && lineSubmitPoint.value.length > 0 && pixRatioX.value != null ? Math.round(lineSubmitPoint.value[2].y * pixRatioY.value) : 0,
            },
          },
          HotArea: pylogonSubmitPoint.value.map((item: any) => ({
            X: Math.round(item.x * pixRatioX.value),
            Y: Math.round(item.y * pixRatioY.value),
          })),
        },
      ],
    },
  }
  // console.log(params)
  let res: ResponseData = await ModTaskConfig(params)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg,
    })
    // getTasklists()
    const tmp = activeKey.value

    await SelectTask()
    nextTick(() => {
      activeKey.value = tmp
    })
    handleTabChange()
  } else {
    ElMessage({
      type: 'error',
      message: res.msg,
    })
  }
}
</script>

<style>
.h2 {
  font-weight: bold; /* 加粗 */
  font-size: 15px; /* 字号15px */
  font-family: 'Arial Black', Gadget, sans-serif; /* 黑体字 */
}
.h3 {
  font-weight: bold; /* 加粗 */
  font-size: 12px; /* 字号15px */
  font-family: 'Arial Black', Gadget, sans-serif; /* 黑体字 */
}
#myCanvas {
  border: 1px solid black;
}
.canvas-container {
  position: absolute !important;
  bottom: 0;
}

.channel .el-descriptions-item {
  background-color: #f0f0f0; /* 浅灰色背景 */
}
</style>
