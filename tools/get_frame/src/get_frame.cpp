#include <iostream>
#include <fstream>
#include <opencv2/opencv.hpp>
#include <chrono>
#include <nlohmann/json.hpp>
#include <getopt.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>

using json = nlohmann::json;


void json_failed(const char* str = "failed", ...) {
    char buffer[1024]; 
    memset(buffer, ' ', sizeof(buffer));
    va_list args;
    va_start(args, str);
    vsnprintf(buffer, sizeof(buffer), str, args);
    va_end(args);

    json j;
    j["code"] = -1;
    j["msg"] = buffer;

    printf(j.dump().c_str());
}

void json_success(json data) {
    json j;
    j["code"] = 0;
    j["msg"] = "success";
    j["data"] = data;

    printf(j.dump().c_str());
}

void Usage() {
    json_failed("Usage:--dir: Set frame save path. "
        "--name: Set frame save name. "
        "--src: Set video source."
        "--frame != 0: Set frame count."
    );
}

void processArguments(int argc, char *argv[], std::string &save_path, std::string &save_name,
                      std::string& video_src, int& frame_count) {
    struct option longOptions[] = {{"dir", required_argument, nullptr, 'd'},
                                 {"name", required_argument, nullptr, 'n'},
                                 {"src", required_argument, nullptr, 's'},
                                 {"help", no_argument, nullptr, 'h'},
                                 {"frame", required_argument, nullptr, 'f'},
                                 {nullptr, 0, nullptr, 0}};

    int optionIndex = 0;
    int option;

    while ((option = getopt_long(argc, argv, "f:d:n:s:h", longOptions,
                               &optionIndex)) != -1) {
    switch (option) {
    case 'd':
        save_path = std::string(optarg);
        break;
    case 'n':
        save_name = std::string(optarg);
        break;
    case 's':
        video_src = std::string(optarg);
        break;
    case 'f':
        frame_count = std::stoi(optarg);
        break;
    case 'h':
        Usage();
        exit(EXIT_FAILURE);
    case '?':
        Usage();
        exit(EXIT_FAILURE);
    default:
        exit(EXIT_FAILURE);
    }
    }

    if (save_path.empty() || save_name.empty() || video_src.empty() || frame_count == 0) {
        Usage();
        exit(EXIT_FAILURE);
    }   
}

class timeCount {
public:
    void begin() {
        start_time = std::chrono::high_resolution_clock::now();
        counting = true;
    }

    void end(const char* msg) {
        if (!counting) {
            std::cout << "not in counting" << std::endl;
            return;
        }
        end_time = std::chrono::high_resolution_clock::now();
        std::chrono::duration<double, std::milli> duration = end_time - start_time;
        std::cout << msg << ": " << duration.count() << " ms" << std::endl;
    };

    std::chrono::_V2::system_clock::time_point start_time;
    std::chrono::_V2::system_clock::time_point end_time;
    bool counting = false;
};

int work(const char* save_path, const char* save_name, const char* video_src, int frame_num) {

    if (mkdir(save_path, 0755) != 0 && errno != EEXIST) {
        json_failed("debug: mkdir failed");
        exit(EXIT_FAILURE);
    }
    
    cv::VideoCapture capture(video_src); 

    int frame_counter = 2;
    if (!capture.isOpened()) {
        json_failed("Error open video: %s", video_src);
        return -1;
    }

    // 获取视频流的分辨率
    int frame_width = static_cast<int>(capture.get(cv::CAP_PROP_FRAME_WIDTH));
    int frame_height = static_cast<int>(capture.get(cv::CAP_PROP_FRAME_HEIGHT));

    // 获取视频流的编码格式
    int fourcc_int = static_cast<int>(capture.get(cv::CAP_PROP_FOURCC));
    char fourcc_char[5] = {0};
    for (int i = 0; i < 4; ++i) {
        fourcc_char[i] = fourcc_int >> (8 * i) & 0xFF;
    }

    cv::Mat frame;
    if (capture.read(frame)) {
        std::string save_file = std::string(save_path) + "/" + std::string(save_name) + ".jpg";
        cv::imwrite(save_file, frame);
        if (frame_num == 1) {
            json j;
            j["codeName"] = std::string(fourcc_char);
            j["height"] = frame_height;
            j["width"] = frame_width;
            json_success(j);
            capture.release();
            return 0;
        }
    } else {
        json_failed("capture failed");
        return -1;
    }

    while (capture.read(frame)) {
        if (frame_counter++ != frame_num) {
            continue;
        }

        std::string save_file = std::string(save_path) + "/" + std::string(save_name) + ".jpg";
        cv::imwrite(save_file, frame);
        json j;

        j["codeName"] = std::string(fourcc_char);
        j["height"] = frame_height;
        j["width"] = frame_width;
        json_success(j);
        capture.release();
        return 0;
    } 

    json_failed("no enough frames, first frame saved");
    capture.release();
    return 0;
}

int main(int argc, char** argv) {
    std::string save_path;
    std::string save_name;
    std::string video_src;
    int frame_count;
    if (argc == 2 && std::string(argv[1]) == "debug") {
        if (mkdir("debug", 0755) != 0 && errno != EEXIST) {
            json_failed("debug: mkdir failed");
            exit(EXIT_FAILURE);
        }
        save_path = "debug";
        save_name = "debug_frame";
        video_src = "rtsp://172.28.8.34:26644/live/123/face_long.mp4";
        frame_count = 1;
    } else {
        processArguments(argc, argv, save_path, save_name, video_src, frame_count);
    }

    work(save_path.c_str(), save_name.c_str(), video_src.c_str(), frame_count);
    return 0;
}
