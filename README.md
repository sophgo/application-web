
## sophon-web是基于SOPHON SDK的算法可视化应用平台

- **快速搭建**：前后端一键部署，快速对接sophon-stream算法应用；
- **任务管理**：通过 Web 端轻松管理算法任务等；
- **告警展示**：展示算法告警图片，显示详细信息；

## 快速开始
安装包位置：
``` bash
ls back/release/
application-web-linux_arm64.tgz
```

将安装包拷贝到算能SE9设备
``` bash
tar -xzvf application-web.tgz
cd application_web
sudo ./install.sh
```
安装完成后，打开8089端口即可接入页面

## 前端开发
[前端开发环境搭建](./front/README.md)  
#### 前端项目打包
```bash
cd front/
pnpm i
pnpm run build
```
得到dist包
## 后端开发
[后端开发环境搭建](./back/README.md)  
#### 后端项目打包
```bash
cd back/build/
cp -r  ../../front/dist/*  ../dist/
./build_test.sh
```
在release目录得到`application-web-linux_arm64.tgz`包
## tool工具说明
get_frame是基于bmcv获取视频流或视频文件一帧图片工具，供后端代码调用