<template>
  <el-dialog v-model="centerDialogVisible" :title="dialogStatus" width="500" align-center>
    <span>{{ dialogMsg }}</span>
    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="centerDialogVisible = false">确认</el-button>
      </div>
    </template>
  </el-dialog>
  <el-card style="margin: 10px 10px">
    <div>
      <h1>人脸1：1比对</h1>
    </div>
  </el-card>
  <el-card style="margin: 10px 10px; min-height: calc(100vh - 35px)">
    <el-row :gutter="20">
      <el-col class="upload-col" v-for="(item, index) in uploadConfigs" :key="index" :push="2" :span="10">
        <el-upload
          class="avatar-uploader"
          drag
          action="/file/upload"
          multiple
          v-model:file-list="item.fileList"
          :limit="1"
          :disabled="!item.disabled"
          :on-success="(response, uploadFile, uploadFiles) => handleSuccess(response, uploadFile, uploadFiles, index)"
          :on-remove="handleRemove"
          :before-upload="beforeAvatarUpload"
        >
          <img v-if="!item.disabled" :src="item.fileList[0].url" class="avatar" />
          <el-icon v-else class="avatar-uploader-icon"><upload-filled /></el-icon>
          <div class="el-upload__text" v-if="item.disabled">
            拖曳文件到此处或
            <em>点击</em>
            以上传文件
          </div>
        </el-upload>
        <div v-if="!item.disabled" class="avatar-delete-icon">
          <el-button type="primary" :icon="Delete" @click="deleteFile(index)" plain>删除图片</el-button>
        </div>
      </el-col>
    </el-row>
    <div class="flex-container">
      <el-button class="compare-button" type="success" :icon="Search" @click="compareFace" :disabled="isComparable">人脸比对</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { UploadConfig, Features, CompareResult } from '@/api/face/compare/type'
import { reqImageFeature, reqCompareFeature } from '@/api/face/compare'
import { UploadFile, UploadFiles, UploadProps, UploadUserFile } from 'element-plus'
import { Search, Delete } from '@element-plus/icons-vue'

const centerDialogVisible = ref(false)
const dialogMsg = ref('' as string)
const uploadConfigs = ref<UploadConfig[]>([
  {
    fileList: [] as UploadUserFile[],
    action: '/file/upload',
    listType: 'picture',
    limit: 1,
    disabled: true,
    fileName: ''
  },
  {
    fileList: [] as UploadUserFile[],
    action: '/file/upload',
    listType: 'picture',
    limit: 1,
    disabled: true,
    fileName: ''
  },
])

const dialogStatus = ref(' ')
const deleteFile = (index: number) => {
  uploadConfigs.value[index].fileList = [] as unknown as UploadUserFile[]
  uploadConfigs.value[index].disabled = true
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
const handleSuccess = (response: any, file: UploadFile, fileList: UploadFiles, configId: number) => {
  uploadConfigs.value[configId].fileList[0].url = '/api/face/alarm/image?fileName=' + file.name
  console.log(uploadConfigs.value[configId].fileList[0].url)
  let featureNum = getFileFeature(configId)
  console.log()
  console.log(featureNum)
  // if fileList
}

const isComparable = computed(() => {
  return uploadConfigs.value[0].fileList.length === 0 || uploadConfigs.value[1].fileList.length === 0
})

const getFileFeature = async (configId: number) => {
  let res: Features = await reqImageFeature({ FileID: "/data/face/upload/" + uploadConfigs.value[configId].fileList[0].name })
  console.log(res)
  if (res.code === 0) {
    if (res.result.length === 0) {
      dialogStatus.value = ' '

      centerDialogVisible.value = true

      dialogMsg.value = '未检测到人脸'
      uploadConfigs.value[configId].fileList = [] as unknown as UploadUserFile[]

      return
    } else if (res.result.length > 1) {
      dialogStatus.value = ' '

      centerDialogVisible.value = true
      dialogMsg.value = '请上传仅包含一张人脸的图片'
      uploadConfigs.value[configId].fileList = [] as unknown as UploadUserFile[]

      return
    }
    uploadConfigs.value[configId].disabled = false
  } else {
    dialogStatus.value = ' '
    centerDialogVisible.value = true

    dialogMsg.value = '未检测到人脸'
    uploadConfigs.value[configId].fileList = [] as unknown as UploadUserFile[]
  }
}

const handleRemove = () => {
  console.log('hello')
}
const compareFace = async () => {
  if (uploadConfigs.value[0].fileList.length === 0 || uploadConfigs.value[1].fileList.length === 0) {
    dialogStatus.value = ''

    centerDialogVisible.value = true
    dialogMsg.value = '请上传两张人脸图片'
  } else {
    let res: CompareResult = await reqCompareFeature({ FileID: [uploadConfigs.value[0].fileList[0].name, uploadConfigs.value[1].fileList[0].name] as string[] })
    console.log(res)

    if (res.code != 0) {
      console.log(res)
      dialogStatus.value = ' '
      centerDialogVisible.value = true
      dialogMsg.value = '人脸比对失败'
    } else {
      dialogStatus.value = '比对成功'
      centerDialogVisible.value = true
      dialogMsg.value = '相似度：' + (res.result.Similarity * 100).toFixed(2) + '%'
    }
  }
}
</script>

<style lang="scss" scoped>
.el-upload-dragger {
  height: 1000px;
}

.show-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

<style>
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
  width: 25vw;
  height: 25vw;
  text-align: center;
}

.el-col.upload-col {
  height: 60vh;
}

.el-icon.avatar-delete-icon {
  font-size: 5vw;
  color: #8c939d;
  width: 1vw;
  height: 1vw;
  text-align: center;
}

.flex-container {
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  height: 10vh; /* 或者你需要的任何高度 */
}

.flex-container .compare-button {
  font-size: 2vh;
  height: 5vh;
  width: 10vw;
}
</style>
