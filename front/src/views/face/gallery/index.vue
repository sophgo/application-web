<template>
  <el-dialog v-model="centerDialogVisible" :title="dialogStatus" width="500" align-center>
    <span>{{ dialogMsg }}</span>
    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="centerDialogVisible = false">确认</el-button>
      </div>
    </template>
  </el-dialog>

  <el-card style="margin: 10px 10px; min-height: calc(100vh - 35px)">
    <el-row :gutter="20">
      <el-col class="upload-col" :push="2" :span="10">
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
            拖曳压缩包文件到此处或
            <em>点击</em>
            以上传文件
          </div>
        </el-upload>
        <div v-if="!uploadConfigs.disabled" class="avatar-delete-icon">
          <el-button type="primary" :icon="Delete" @click="deleteFile()" plain>删除人脸库</el-button>
        </div>
      </el-col>
      <el-col :span="14">
        <p>hello</p>
      </el-col>
    </el-row>
    <div class="flex-container">
      <el-button class="compare-button" type="success" :icon="Search" @click="compareFace" :disabled="isComparable">人脸比对</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { UploadConfig, Features, CompareResult, ResponseData } from '@/api/face/compare/type'
import { reqImageFeature, reqCompareFeature, reqGalleryCreate } from '@/api/face/compare'
import { UploadFile, UploadFiles, UploadProps, UploadUserFile } from 'element-plus'
import { Search, Delete } from '@element-plus/icons-vue'

const centerDialogVisible = ref(false)
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

const dialogStatus = ref(' ')
const deleteFile = () => {
  uploadConfigs.value.fileList = [] as unknown as UploadUserFile[]
  uploadConfigs.value.disabled = true
}

const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
  const allowedExtensions = ['.gz', '.tgz', '.rar', '.zip'];
  const isAllowedType = allowedExtensions.some(ext => rawFile.name.endsWith(ext));
  
  if (!isAllowedType) {
    dialogStatus.value = ' ';
    centerDialogVisible.value = true;
    dialogMsg.value = '请上传.tar.gz、.tgz、.rar或.zip格式的文件';
    return false;
  }
  return true;
}



const handleSuccess = (response: any, file: UploadFile, fileList: UploadFiles) => {

  const filepath = response.result

  processGallery(filepath)
  console.log(response)
  // uploadConfigs.value.fileList[0].url = '/api/face/alarm/image?fileName=' + file.name
  // console.log(uploadConfigs.value.fileList[0].url)
  // let featureNum = getFileFeature()
  // console.log()
  // console.log(featureNum)
  // if fileList
}

const isComparable = computed(() => {
  return true
})

const processGallery = async (filename: string) => {

  let res: ResponseData = await reqGalleryCreate({ FileID: filename })
  console.log(res)

  if (res.code != 0) {
    console.log(res)
    dialogStatus.value = ' '
    centerDialogVisible.value = true
    dialogMsg.value = '解析失败'
  } else {
    dialogStatus.value = ' '
    centerDialogVisible.value = true
    dialogMsg.value = '解析成功'
  }
  
}

const getFileFeature = async () => {
  let res: Features = await reqImageFeature({ FileID: uploadConfigs.value.fileList[0].name })
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
const compareFace = async () => {
  
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
