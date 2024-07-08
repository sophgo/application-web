## 编译条件
1. go >= 1.21
2. 安装gcc-aarch64-linux-gnu：`sudo apt-get install gcc-aarch64-linux-gnu`

## 构建  

1. 进入build目录，执行build脚本
``` bash
cd build
./build_test.sh 
```

3. 安装包文件  
``` bash
release/
├── application-web-linux_arm64.tgz

``` 

## 安装运行
tar -xvf application-web-linux_arm64.tgz 
cd application-web
sudo ./install.sh
