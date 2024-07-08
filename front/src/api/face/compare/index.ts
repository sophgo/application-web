import request from '@/utils/request'
import type { ResConfig, FileIDs, FileID, ModAlgorithm, ResponseData } from './type'
import type { UploadProps, UploadUserFile } from 'element-plus'

enum API {
  MODCONFIG_URL = 'api/config/mod',
  GETCONFIG_URL = 'api/config/get',
  FEATURE_URL = 'api/face/feature',
  GALLERY_CREATE_URL = 'api/face/gallery/create',
  COMPARE_URL = 'api/face/compare',
}


export const ModTaskConfig = (data: ModAlgorithm) => request.post<any, ResponseData>(API.MODCONFIG_URL, data)

export const reqImageFeature = (data: FileID) => request.post<any, any>(API.FEATURE_URL, data)

export const reqCompareFeature = (data: FileIDs) => request.post<any, any>(API.COMPARE_URL, data)

export const reqGalleryCreate = (data: FileID) => request.post<any, any>(API.GALLERY_CREATE_URL, data)