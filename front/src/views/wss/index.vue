<template>
  <div>
    <el-card class="controls">
      <el-input v-model="wsAddress" placeholder="请输入WebSocket地址" style="margin: 10px 10px; width: 250px">请输入WebSocket地址</el-input>
      <el-button type="primary" @click="startPlaying" style="margin: 10px 10px">播放</el-button>
      <!-- <el-button type="primary" @click="Pause" style="margin: 10px 10px">暂停</el-button> -->
      <el-button type="danger" @click="stopPlaying" style="margin: 10px 10px">停止</el-button>
      <el-button :type="debug ? 'danger' : 'primary'" @click="debugSet" style="margin: 10px 10px">{{ debug ? '关闭调试' : '调试模式' }}</el-button>
    </el-card>
    <el-card class="canvas-container">
      <canvas ref="canvas" :width="width" :height="height" class="my-canvas"></canvas>
      <el-slider v-model="playbackSpeed" :min="0.5" :max="50" :step="0.1" />
      <span>播放速度: {{ playbackSpeed }}帧/秒</span>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

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

const wsAddress = ref('ws://192.168.0.101:9002/') // 默认ws地址
const canvas = ref(null)
const context = ref(null)
const socket = ref(null)
const queue = ref([])

const pauseFalg = ref(false)

const playbackSpeed = ref(25) // 播放速度
let frameTimer = null // 定时器变量

const width = 1920 / 2 // 画布宽度
const height = 1080 / 2 // 画布高度
const ratex = ref(1)
const ratey = ref(1)

let debug = ref(false)

function debugSet() {
  debug.value = !debug.value
}

function startPlaying() {
  if (socket.value) {
    ElMessage({
      type: 'error',
      message: '请先停止当前连接',
    })
    return
  }
  initializeWebSocket()
}

function Pause() {
  // if (socket.value) {
  //   socket.value.close()
  //   socket.value = null
  //   ElMessage.success('停止播放成功。')
  // }
  pauseFalg.value = true
  if (frameTimer) {
    clearTimeout(frameTimer)
    frameTimer = null
  }
}

function stopPlaying() {
  if (socket.value) {
    socket.value.close()
    socket.value = null
    ElMessage.success('停止播放成功。')
  }
  if (frameTimer) {
    clearTimeout(frameTimer)
    frameTimer = null
  }
  queue.value = [] // 清空队列
}

function initializeWebSocket() {
  socket.value = new WebSocket(wsAddress.value)
  socket.value.onopen = () => {
    ElMessage.success('WebSocket连接成功。')
  }
  socket.value.onerror = (error) => {
    ElMessage.error('WebSocket连接失败')
    // 关闭WebSocket连接并释放资源
    socket.value.close()
    socket.value = null
  }
  socket.value.onmessage = (event) => {
    const data = JSON.parse(event.data)
    queue.value.push(data)
    if (!frameTimer) {
      drawFrame()
    }
  }
  socket.value.onclose = () => {
    ElMessage.info('WebSocket连接已关闭。')
    if (socket.value) {
      socket.value.close()
    }
    socket.value = null
  }
}

function drawFrame() {
  if (queue.value.length > 0 && pauseFalg.value === false) {
    const data = queue.value.shift()
    const image = new Image()
    image.onload = () => {
      ratex.value = data['mFrame']['mWidth'] / width === 0 ? 1 : data['mFrame']['mWidth'] / width
      ratey.value = data['mFrame']['mHeight'] / height === 0 ? 1 : data['mFrame']['mHeight'] / height
      context.value.clearRect(0, 0, width, height)
      context.value.drawImage(image, 0, 0, width, height)

      const fps = `实时处理帧率:${data['mFps'].toFixed(2)}` // 格式化文本
      context.value.font = '20px Arial'
      context.value.fillStyle = 'black' // 设置文本颜色
      context.value.fillText(fps, 10, 30)

      // 如果有检测到的对象，绘制方框
      if (data.mDetectedObjectMetadatas) {
        drawBoxObj(data.mDetectedObjectMetadatas, data.mSubObjectMetadatas, data.mTrackedObjectMetadatas)
      }
      if (data.mRecognizedObjectMetadatas) {
        drawRecognizedObj(data.mRecognizedObjectMetadatas)
      }
      if (data.mPosedObjectMetadatas) {
        drawKeyPoints(data.mPosedObjectMetadatas)
      }
      if (data.mFaceObjectMetadata) {
        drawFace(data.mFaceObjectMetadata)
      }
    }
    image.src = 'data:image/jpeg;base64,' + data['mFrame']['mSpData']
    if (debug.value) {
      console.log('data:', data)
    }
  }
  frameTimer = setTimeout(drawFrame, 1000 / playbackSpeed.value)
}

function drawRecognizedObj(objs) {
  objs.forEach((obj, index) => {
    const text = `Label:${obj['mLabelName']} Scores: ${(obj['mScores'][0] * 100).toFixed(2)} mTopKLabels:${obj['mTopKLabels']}` // 格式化文本
    // 设置文本样式
    context.value.font = '20px Arial'
    context.value.fillStyle = colors[0] // 设置文本颜色
    // 绘制文本
    context.value.fillText(text, 50, 50)
  })
}

function drawBox(boxs) {
  boxs.forEach((box) => {
    let mbox = box['mBox']
    const boxWidth = mbox['mWidth'] / ratex.value
    const boxHeight = mbox['mHeight'] / ratey.value
    const x = mbox['mX'] / ratex.value
    const y = mbox['mY'] / ratey.value
    context.value.beginPath()
    context.value.strokeStyle = colors[box['mClassify'] % 10]
    context.value.lineWidth = 1.5
    context.value.rect(x, y, boxWidth, boxHeight)
    context.value.stroke()

    const text = `Class：${box['mClassify']} Scores: ${(box['mScores'] * 100).toFixed(2)}` // 格式化文本
    // 设置文本样式
    context.value.font = '20px Arial'
    context.value.fillStyle = colors[box['mClassify'] % 10] // 设置文本颜色
    // 绘制文本
    context.value.fillText(text, x, y - 6) // 将文本放在矩形框的上方
  })
}

function drawFace(faces) {
  faces.forEach((face) => {
    const startX = face.left / ratex.value
    const startY = face.top / ratey.value
    const boxWidth = (face.right - face.left + 1) / ratex.value
    const boxHeight = (face.bottom - face.top + 1) / ratey.value

    context.value.beginPath()
    context.value.strokeStyle = colors[0]
    context.value.lineWidth = 1.5
    context.value.rect(startX, startY, boxWidth, boxHeight)
    context.value.stroke()

    const text = ` S: ${(face.score * 100).toFixed(2)}` // 格式化文本
    // 设置文本样式
    context.value.font = '16px Arial'
    context.value.fillStyle = colors[0] // 设置文本颜色
    // 绘制文本
    context.value.fillText(text, startX, startY - 6) // 将文本放在矩形框的上方
  })
}

function drawKeyPoints(points) {
  points.forEach((point) => {
    if (point.length === 0) {
      return
    }

    const keypoints = point.keypoints

    for (let index = 0; index < keypoints.length; index += 3) {
      if (keypoints[index] === 0 && keypoints[index + 1] === 0) {
        continue
      }
      if (keypoints[index + 2] < 0.1) {
        continue
      }

      const start = { x: keypoints[index] / ratex.value, y: keypoints[index + 1] / ratey.value }
      // 绘制小圆圈
      const radius = 4 // 圆圈半径
      context.value.beginPath()
      context.value.arc(start.x, start.y, radius, 0, 2 * Math.PI)
      context.value.fillStyle = colors[0]
      context.value.fill()
    }
  })
}

function drawBoxObj(boxs, subObj, tracks) {
  boxs.forEach((box, index) => {
    let mbox = box['mBox']
    const boxWidth = mbox['mWidth'] / ratex.value
    const boxHeight = mbox['mHeight'] / ratey.value
    const x = mbox['mX'] / ratex.value
    const y = mbox['mY'] / ratey.value
    context.value.beginPath()
    context.value.strokeStyle = colors[box['mClassify'] % 10]
    context.value.lineWidth = 1.5
    context.value.rect(x, y, boxWidth, boxHeight)
    context.value.stroke()

    const text = `C:${box['mClassify']},S:${(box['mScores'][0] * 100).toFixed(1)}` // 格式化文本
    // console.log(text, x, y)
    // 设置文本样式
    context.value.font = 'bold 20px Arial'
    context.value.fillStyle = colors[box['mClassify'] % 10] // 设置文本颜色
    context.value.fillText(text, x, y - 6) // 将文本放在矩形框的上方

    if (subObj && subObj.length > 0) {
      let text1 = ''
      if (subObj[index].mRecognizedObjectMetadatas[0].mLabelName !== '') {
        text1 += subObj[index].mRecognizedObjectMetadatas[0].mLabelName
      }
      if (subObj[index].mRecognizedObjectMetadatas[0].mScores.length > 0) {
        text1 += ` S:${(subObj[index].mRecognizedObjectMetadatas[0].mScores[0] * 100).toFixed(0)}`
      }
      if (subObj[index].mRecognizedObjectMetadatas[0].mTopKLabels.length > 0) {
        text1 += ` L:${subObj[index].mRecognizedObjectMetadatas[0].mTopKLabels[0]}`
      }
      context.value.fillText(text1, x + 3, y + boxHeight) // 将文本放在矩形框
    }

    if (tracks && tracks.length > 0) {
      let tId = tracks[index].mTrackId.toString()
      context.value.font = 'bold 22px Arial'
      context.value.fillText(tId, x + 10, y + 30) // 将文本放在矩形框
    }
  })
}

onUnmounted(() => {
  if (socket.value) {
    socket.value.close()
  }
  if (frameTimer) {
    clearTimeout(frameTimer)
  }
})

onMounted(() => {
  setTimeout(() => {
    context.value = canvas.value.getContext('2d')
  }, 100) // 100毫秒延迟
})
</script>

<style scoped>
.canvas-container {
  display: flex;
  flex-direction: column; /* 垂直排列子元素 */
  justify-content: center;
  align-items: center;
  margin-top: 20px;
}
.controls {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 20px; /* 增加底部间距 */
}
.my-canvas {
  background-color: #f0f0f0;
}
</style>
