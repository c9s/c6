#!/bin/bash
cmake -DGFlags_ROOT_DIR=$(pwd)vendor/gflags -DGTEST_ROOT=$(pwd)vendor/gtest -DGLog_ROOT_DIR=$(pwd)/vendor/glog .
