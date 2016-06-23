# cmake file examples
# https://code.google.com/p/opencv-feature-tracker/source/browse/CMakeLists.txt?r=f804b03e704147e65183c19a50f57abedb22f45c
cmake_minimum_required(VERSION 2.8)
project(c6)

# set(CMAKE_CXX_COMPILER "/Users/c9s/src/clang-openmp/llvm/build/Release+Asserts/bin/clang")

# -DCMAKE_BUILD_TYPE=DEBUG
# -DCMAKE_BUILD_TYPE=RELEASE
# if(NOT CMAKE_BUILD_TYPE )
#     set( CMAKE_BUILD_TYPE Debug CACHE STRING "Choose the type of build, options are: None Debug Release RelWithDebInfo MinSizeRel." FORCE)
# endif()

SET(CMAKE_MODULE_PATH
    ${PROJECT_SOURCE_DIR}/cmake_modules
    ${CMAKE_MODULE_PATH}
    )

find_package(OpenCV REQUIRED)
find_package(GooglePerfTools REQUIRED)
# find_package(opencv)
# link_directories(${LINK_DIRECTORIES} ${opencv_LIBRARY_DIRS})

set(Boost_USE_STATIC_LIBS OFF)
set(Boost_USE_MULTITHREADED ON)
set(Boost_USE_STATIC_RUNTIME OFF)

# find_package(Boost 1.54.0 COMPONENTS program_options filesystem regex)
find_package(Boost 1.54.0 COMPONENTS program_options)
find_package(GFlags REQUIRED)
# find_package(OpenMP)
find_package(GLog REQUIRED)

link_directories(${LINK_DIRECTORIES} /opt/local/lib /usr/lib)
include_directories(. ${PROJECT_SOURCE_DIR}/include ${PROJECT_SOURCE_DIR}/t3/src ${PROJECT_SOURCE_DIR}/src /opt/local/include /usr/include)

include_directories(${INCLUDE_DIRECTORIES} ${OPENCV_INCLUDE_DIR})
# set (CMAKE_CXX_FLAGS "-std=c++0x -arch x86_64 -stdlib=libc++ -g3 -Wall -O0")

# debug
# -pg for profile
# set (CMAKE_CXX_FLAGS "-std=c++0x -arch x86_64 -stdlib=libc++ -Wall -Wno-long-long -Wno-unused-value -Wno-overloaded-virtual -O3")
# set (CMAKE_CXX_FLAGS "-std=c++0x -arch x86_64 -stdlib=libc++ -g3 -Wall -Wno-long-long -Wno-unused-value -Wno-overloaded-virtual -O0")

# http://www.openguru.com/2009/04/cmake-detecting-platformoperating.html
if(${CMAKE_SYSTEM_NAME} MATCHES "Darwin")
    set(CMAKE_CXX_FLAGS "-std=c++0x -arch x86_64 -stdlib=libc++ -Wall -Wno-long-long -Wno-unused-value -Wno-overloaded-virtual")
elseif(${CMAKE_SYSTEM_NAME} MATCHES "Linux")
    set(CMAKE_CXX_FLAGS "-std=c++11 -Wall -Wno-long-long -Wno-unused-value -Wno-overloaded-virtual")
endif()

if(${CMAKE_BUILD_TYPE} MATCHES "Release")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS}")
# elseif(${CMAKE_BUILD_TYPE} MATCHES "Debug")
else()
    # add_definitions(-DThread)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -pg -O0 -g3")
endif()


# IF(${CMAKE_SYSTEM_NAME} MATCHES "Linux")
#     # Linux specific code
#     SET(OperatingSystem "Linux")
# ENDIF(${CMAKE_SYSTEM_NAME} MATCHES "Linux")

enable_testing()

# set(gtest_SOURCE_DIR gtest-1.7.0)
# add_subdirectory(${gtest_SOURCE_DIR})
# include_directories(${gtest_SOURCE_DIR}/include ${gtest_SOURCE_DIR})
set(ENV{GTEST_ROOT} ${PROJECT_SOURCE_DIR}/gtest-1.7.0)
# add_subdirectory(gtest-1.7.0)
add_subdirectory(t3)
add_subdirectory(src)
add_subdirectory(tests)
