cmake_minimum_required(VERSION 2.8 FATAL_ERROR)
project(get_frame)


set(FLAG_CROSS_COMPILE true)
set(CMAKE_C_COMPILER "aarch64-linux-gnu-gcc-9")
set(CMAKE_CXX_COMPILER "aarch64-linux-gnu-g++-9")


set(SOPHON_LIB_ROOT_DIR /opt/sophon)
set(CROSS_LIB_ROOT_DIR /opt/sophon)

set(LIBSOPHON_ROOT_DIR ${SOPHON_LIB_ROOT_DIR}/libsophon-0.4.9)
set(OpenCV_ROOT_DIR ${SOPHON_LIB_ROOT_DIR}/sophon-opencv-latest)
set(FFMPEG_ROOT_DIR ${SOPHON_LIB_ROOT_DIR}/sophon-ffmpeg-latest)

#opencv
set(OpenCV_INCLUDE_DIRS /home/linaro/sophon_media/sophon-opencv_1.5.0/include/opencv4)
set(OpenCV_LIB_DIR ${OpenCV_ROOT_DIR}/lib)
set(OpenCV_LIBRARIES ${OpenCV_LIB_DIR}/libopencv_core.so;${OpenCV_LIB_DIR}/libopencv_imgproc.so;
  ${OpenCV_LIB_DIR}/libopencv_highgui.so;${OpenCV_LIB_DIR}/libopencv_imgcodecs.so;
  ${OpenCV_LIB_DIR}/libopencv_videoio.so;${OpenCV_LIB_DIR}/libopencv_calib3d.so;
  ${OpenCV_LIB_DIR}/libopencv_features2d.so;${OpenCV_LIB_DIR}/libopencv_flann.so)

#ffmpeg
set(FFMPEG_INCLUDE_DIRS /home/linaro/sophon_media/sophon-ffmpeg_1.5.0/include)
set(FFMPEG_LIB_DIR ${FFMPEG_ROOT_DIR}/lib)
set(FFMPEG_LIBRARIES ${FFMPEG_LIB_DIR}/libavcodec.so;
    ${FFMPEG_LIB_DIR}/libavformat.so;
    ${FFMPEG_LIB_DIR}/libavutil.so;
    ${FFMPEG_LIB_DIR}/libswscale.so;
    ${FFMPEG_LIB_DIR}/libswresample.so;
    ${FFMPEG_LIB_DIR}/libavdevice.so;
    ${FFMPEG_LIB_DIR}/libavfilter.so;)

#libsophon
set(LIBSOPHON_INCLUDE_DIRS ${LIBSOPHON_ROOT_DIR}/include)
set(LIBSOPHON_LIB_DIR ${LIBSOPHON_ROOT_DIR}/lib)
set(LIBSOPHON_LIBRARIES ${LIBSOPHON_LIB_DIR}/libbmlib.so;
    ${LIBSOPHON_LIB_DIR}/libbmrt.so;
    ${LIBSOPHON_LIB_DIR}/libbmcv.so)

include_directories(${CMAKE_SOURCE_DIR})
include_directories(${OpenCV_INCLUDE_DIRS})
include_directories(${FFMPEG_INCLUDE_DIRS})
include_directories(${LIBSOPHON_INCLUDE_DIRS})


link_directories(${OpenCV_LIB_DIR})
link_directories(${FFMPEG_LIB_DIR})
link_directories(${LIBSOPHON_LIB_DIR})

add_executable(get_frame src/get_frame.cpp)
target_link_libraries(get_frame ${OpenCV_LIBRARIES})
