# ===================================================================================
#  The OpenCV CMake configuration file
#
#             ** File generated automatically, do not modify **
#
#  Usage from an external project:
#    In your CMakeLists.txt, add these lines:
#
#    FIND_PACKAGE(OpenCV REQUIRED)
#    TARGET_LINK_LIBRARIES(MY_TARGET_NAME ${OpenCV_LIBS})
#
#    Or you can search for specific OpenCV modules:
#
#    FIND_PACKAGE(OpenCV REQUIRED core highgui)
#
#    If the module is found then OPENCV_<MODULE>_FOUND is set to TRUE.
#
#    This file will define the following variables:
#      - OpenCV_LIBS                     : The list of all imported targets for OpenCV modules.
#      - OpenCV_INCLUDE_DIRS             : The OpenCV include directories.
#      - OpenCV_COMPUTE_CAPABILITIES     : The version of compute capability
#      - OpenCV_ANDROID_NATIVE_API_LEVEL : Minimum required level of Android API
#      - OpenCV_VERSION                  : The version of this OpenCV build: "2.4.8"
#      - OpenCV_VERSION_MAJOR            : Major version part of OpenCV_VERSION: "2"
#      - OpenCV_VERSION_MINOR            : Minor version part of OpenCV_VERSION: "4"
#      - OpenCV_VERSION_PATCH            : Patch version part of OpenCV_VERSION: "8"
#      - OpenCV_VERSION_TWEAK            : Tweak version part of OpenCV_VERSION: "0"
#
#    Advanced variables:
#      - OpenCV_SHARED
#      - OpenCV_CONFIG_PATH
#      - OpenCV_INSTALL_PATH  (not set on Windows)
#      - OpenCV_LIB_COMPONENTS
#      - OpenCV_USE_MANGLED_PATHS
#      - OpenCV_HAVE_ANDROID_CAMERA
#
# ===================================================================================

set(modules_file_suffix "")
if(ANDROID)
  string(REPLACE - _ modules_file_suffix "_${ANDROID_NDK_ABI_NAME}")
endif()

if(NOT TARGET opencv_core)
  include(${CMAKE_CURRENT_LIST_DIR}/OpenCVModules${modules_file_suffix}.cmake)
endif()

# TODO All things below should be reviewed. What is about of moving this code into related modules (special vars/hooks/files)

# Version Compute Capability from which OpenCV has been compiled is remembered
set(OpenCV_COMPUTE_CAPABILITIES "")

set(OpenCV_CUDA_VERSION )
set(OpenCV_USE_CUBLAS   )
set(OpenCV_USE_CUFFT    )
set(OpenCV_USE_NVCUVID  )

# Android API level from which OpenCV has been compiled is remembered
set(OpenCV_ANDROID_NATIVE_API_LEVEL 0)

# Some additional settings are required if OpenCV is built as static libs
set(OpenCV_SHARED OFF)

# Enables mangled install paths, that help with side by side installs
set(OpenCV_USE_MANGLED_PATHS FALSE)

# Extract the directory where *this* file has been installed (determined at cmake run-time)
get_filename_component(OpenCV_CONFIG_PATH "${CMAKE_CURRENT_LIST_FILE}" PATH CACHE)

if(NOT WIN32 OR OpenCV_ANDROID_NATIVE_API_LEVEL GREATER 0)
  if(OpenCV_ANDROID_NATIVE_API_LEVEL GREATER 0)
    set(OpenCV_INSTALL_PATH "${OpenCV_CONFIG_PATH}/../../..")
  else()
    set(OpenCV_INSTALL_PATH "${OpenCV_CONFIG_PATH}/../..")
  endif()
  # Get the absolute path with no ../.. relative marks, to eliminate implicit linker warnings
  if(${CMAKE_MAJOR_VERSION}.${CMAKE_MINOR_VERSION} VERSION_LESS 2.8)
    get_filename_component(OpenCV_INSTALL_PATH "${OpenCV_INSTALL_PATH}" ABSOLUTE)
  else()
    get_filename_component(OpenCV_INSTALL_PATH "${OpenCV_INSTALL_PATH}" REALPATH)
  endif()
endif()

# Presence of Android native camera wrappers
set(OpenCV_HAVE_ANDROID_CAMERA OFF)

# ======================================================
# Include directories to add to the user project:
# ======================================================

# Provide the include directories to the caller
set(OpenCV_INCLUDE_DIRS "${OpenCV_INSTALL_PATH}/include/opencv2;${OpenCV_INSTALL_PATH}/include")
include_directories(${OpenCV_INCLUDE_DIRS})

# ======================================================
# Link directories to add to the user project:
# ======================================================

# Provide the libs directories to the caller
set(OpenCV_LIB_DIR_OPT  CACHE PATH "Path where release OpenCV libraries are located")
set(OpenCV_LIB_DIR_DBG  CACHE PATH "Path where debug OpenCV libraries are located")
set(OpenCV_3RDPARTY_LIB_DIR_OPT  CACHE PATH "Path where release 3rdpaty OpenCV dependencies are located")
set(OpenCV_3RDPARTY_LIB_DIR_DBG  CACHE PATH "Path where debug 3rdpaty OpenCV dependencies are located")
mark_as_advanced(FORCE OpenCV_LIB_DIR_OPT OpenCV_LIB_DIR_DBG OpenCV_3RDPARTY_LIB_DIR_OPT OpenCV_3RDPARTY_LIB_DIR_DBG OpenCV_CONFIG_PATH)

# ======================================================
#  Version variables:
# ======================================================
SET(OpenCV_VERSION 2.4.8)
SET(OpenCV_VERSION_MAJOR  2)
SET(OpenCV_VERSION_MINOR  4)
SET(OpenCV_VERSION_PATCH  8)
SET(OpenCV_VERSION_TWEAK  0)

# ====================================================================
# Link libraries: e.g. opencv_core;opencv_imgproc; etc...
# ====================================================================

SET(OpenCV_LIB_COMPONENTS opencv_videostab;opencv_video;opencv_ts;opencv_superres;opencv_stitching;opencv_photo;opencv_objdetect;opencv_nonfree;opencv_ml;opencv_legacy;opencv_imgproc;opencv_highgui;opencv_gpu;opencv_flann;opencv_features2d;opencv_core;opencv_contrib;opencv_calib3d)

# ==============================================================
#  Extra include directories, needed by OpenCV 2 new structure
# ==============================================================
SET(OpenCV2_INCLUDE_DIRS "")
if(OpenCV2_INCLUDE_DIRS)
  include_directories(${OpenCV2_INCLUDE_DIRS})
  list(APPEND OpenCV_INCLUDE_DIRS ${OpenCV2_INCLUDE_DIRS})

  set(OpenCV_ADD_DEBUG_RELEASE FALSE)
  if(OpenCV_ADD_DEBUG_RELEASE)
    set(OpenCV_LIB_DIR_OPT "${OpenCV_LIB_DIR_OPT}/Release")
    set(OpenCV_LIB_DIR_DBG "${OpenCV_LIB_DIR_DBG}/Debug")
    set(OpenCV_3RDPARTY_LIB_DIR_OPT "${OpenCV_3RDPARTY_LIB_DIR_OPT}/Release")
    set(OpenCV_3RDPARTY_LIB_DIR_DBG "${OpenCV_3RDPARTY_LIB_DIR_DBG}/Debug")
  endif()
endif()

# ==============================================================
#  Check OpenCV availability
# ==============================================================
if(ANDROID AND OpenCV_ANDROID_NATIVE_API_LEVEL GREATER ANDROID_NATIVE_API_LEVEL)
  message(FATAL_ERROR "Minimum required by OpenCV API level is android-${OpenCV_ANDROID_NATIVE_API_LEVEL}")
  #always FATAL_ERROR because we can't say to the caller that OpenCV is not found
  #http://www.mail-archive.com/cmake@cmake.org/msg37831.html
  if(OpenCV_FIND_REQUIRED)
    message(FATAL_ERROR "Minimum required by OpenCV API level is android-${OpenCV_ANDROID_NATIVE_API_LEVEL}")
  elseif(NOT OpenCV_FIND_QUIETLY)
    message(WARNING "Minimum required by OpenCV API level is android-${OpenCV_ANDROID_NATIVE_API_LEVEL}")
  endif()
  set(OpenCV_FOUND "OpenCV_FOUND-NOTFOUND")
  return()#Android toolchain requires CMake > 2.6
endif()

# ==============================================================
#  Form list of modules (components) to find
# ==============================================================
if(NOT OpenCV_FIND_COMPONENTS)
  set(OpenCV_FIND_COMPONENTS ${OpenCV_LIB_COMPONENTS})
  list(REMOVE_ITEM OpenCV_FIND_COMPONENTS opencv_java)
  if(GTest_FOUND OR GTEST_FOUND)
    list(REMOVE_ITEM OpenCV_FIND_COMPONENTS opencv_ts)
  endif()
endif()

# expand short module names and see if requested components exist
set(OpenCV_FIND_COMPONENTS_ "")
foreach(__cvcomponent ${OpenCV_FIND_COMPONENTS})
  if(NOT __cvcomponent MATCHES "^opencv_")
    set(__cvcomponent opencv_${__cvcomponent})
  endif()
  list(FIND OpenCV_LIB_COMPONENTS ${__cvcomponent} __cvcomponentIdx)
  if(__cvcomponentIdx LESS 0)
    #requested component is not found...
    if(OpenCV_FIND_REQUIRED)
      message(FATAL_ERROR "${__cvcomponent} is required but was not found")
    elseif(NOT OpenCV_FIND_QUIETLY)
      message(WARNING "${__cvcomponent} is required but was not found")
    endif()
    #indicate that module is NOT found
    string(TOUPPER "${__cvcomponent}" __cvcomponent)
    set(${__cvcomponent}_FOUND "${__cvcomponent}_FOUND-NOTFOUND")
  else()
    list(APPEND OpenCV_FIND_COMPONENTS_ ${__cvcomponent})
    # Not using list(APPEND) here, because OpenCV_LIBS may not exist yet.
    # Also not clearing OpenCV_LIBS anywhere, so that multiple calls
    # to find_package(OpenCV) with different component lists add up.
    set(OpenCV_LIBS ${OpenCV_LIBS} "${__cvcomponent}")
    #indicate that module is found
    string(TOUPPER "${__cvcomponent}" __cvcomponent)
    set(${__cvcomponent}_FOUND 1)
  endif()
endforeach()
set(OpenCV_FIND_COMPONENTS ${OpenCV_FIND_COMPONENTS_})

# ==============================================================
#  Resolve dependencies
# ==============================================================
if(OpenCV_USE_MANGLED_PATHS)
  set(OpenCV_LIB_SUFFIX ".${OpenCV_VERSION_MAJOR}.${OpenCV_VERSION_MINOR}.${OpenCV_VERSION_PATCH}")
else()
  set(OpenCV_LIB_SUFFIX "")
endif()

foreach(__opttype OPT DBG)
  SET(OpenCV_LIBS_${__opttype} "${OpenCV_LIBS}")
  SET(OpenCV_EXTRA_LIBS_${__opttype} "")

  # CUDA
  if(OpenCV_CUDA_VERSION AND (CMAKE_CROSSCOMPILING OR (WIN32 AND NOT OpenCV_SHARED)))
    if(NOT CUDA_FOUND)
      find_package(CUDA ${OpenCV_CUDA_VERSION} EXACT REQUIRED)
    else()
      if(NOT CUDA_VERSION_STRING VERSION_EQUAL OpenCV_CUDA_VERSION)
        message(FATAL_ERROR "OpenCV static library was compiled with CUDA ${OpenCV_CUDA_VERSION} support. Please, use the same version or rebuild OpenCV with CUDA ${CUDA_VERSION_STRING}")
      endif()
    endif()

    list(APPEND OpenCV_EXTRA_LIBS_${__opttype} ${CUDA_LIBRARIES})

    if(${CUDA_VERSION} VERSION_LESS "5.5")
      list(APPEND OpenCV_EXTRA_LIBS_${__opttype} ${CUDA_npp_LIBRARY})
    else()
      find_cuda_helper_libs(nppc)
      find_cuda_helper_libs(nppi)
      find_cuda_helper_libs(npps)
      list(APPEND OpenCV_EXTRA_LIBS_${__opttype} ${CUDA_nppc_LIBRARY} ${CUDA_nppi_LIBRARY} ${CUDA_npps_LIBRARY})
    endif()

    if(OpenCV_USE_CUBLAS)
      list(APPEND OpenCV_EXTRA_LIBS_${__opttype} ${CUDA_CUBLAS_LIBRARIES})
    endif()

    if(OpenCV_USE_CUFFT)
      list(APPEND OpenCV_EXTRA_LIBS_${__opttype} ${CUDA_CUFFT_LIBRARIES})
    endif()

    if(OpenCV_USE_NVCUVID)
      list(APPEND OpenCV_EXTRA_LIBS_${__opttype} ${CUDA_nvcuvid_LIBRARIES})
    endif()

    if(WIN32)
      list(APPEND OpenCV_EXTRA_LIBS_${__opttype} ${CUDA_nvcuvenc_LIBRARIES})
    endif()
  endif()
endforeach()

# ==============================================================
#  Android camera helper macro
# ==============================================================
if(OpenCV_HAVE_ANDROID_CAMERA)
  macro(COPY_NATIVE_CAMERA_LIBS target)
    get_target_property(target_location ${target} LOCATION)
    get_filename_component(target_location "${target_location}" PATH)
    file(GLOB camera_wrappers "${OpenCV_LIB_DIR_OPT}/libnative_camera_r*.so")
    foreach(wrapper ${camera_wrappers})
      add_custom_command(
        TARGET ${target}
        POST_BUILD
        COMMAND ${CMAKE_COMMAND} -E copy "${wrapper}" "${target_location}"
      )
    endforeach()
  endmacro()
endif()

# ==============================================================
# Compatibility stuff
# ==============================================================
if(CMAKE_BUILD_TYPE MATCHES "Debug")
  SET(OpenCV_LIB_DIR ${OpenCV_LIB_DIR_DBG} ${OpenCV_3RDPARTY_LIB_DIR_DBG})
else()
  SET(OpenCV_LIB_DIR ${OpenCV_LIB_DIR_OPT} ${OpenCV_3RDPARTY_LIB_DIR_OPT})
endif()
set(OpenCV_LIBRARIES ${OpenCV_LIBS})

if(CMAKE_CROSSCOMPILING AND OpenCV_SHARED AND (CMAKE_SYSTEM_NAME MATCHES "Linux"))
  foreach(dir ${OpenCV_LIB_DIR})
    set(CMAKE_EXE_LINKER_FLAGS    "${CMAKE_EXE_LINKER_FLAGS}    -Wl,-rpath-link,${dir}")
    set(CMAKE_SHARED_LINKER_FLAGS "${CMAKE_SHARED_LINKER_FLAGS} -Wl,-rpath-link,${dir}")
    set(CMAKE_MODULE_LINKER_FLAGS "${CMAKE_MODULE_LINKER_FLAGS} -Wl,-rpath-link,${dir}")
  endforeach()
endif()



#
# Some macroses for samples
#
macro(ocv_check_dependencies)
  set(OCV_DEPENDENCIES_FOUND TRUE)
  foreach(d ${ARGN})
    if(NOT TARGET ${d})
      set(OCV_DEPENDENCIES_FOUND FALSE)
      break()
    endif()
  endforeach()
endmacro()

# adds include directories in such way that directories from the OpenCV source tree go first
function(ocv_include_directories)
  set(__add_before "")
  file(TO_CMAKE_PATH "${OpenCV_DIR}" __baseDir)
  foreach(dir ${ARGN})
    get_filename_component(__abs_dir "${dir}" ABSOLUTE)
    if("${__abs_dir}" MATCHES "^${__baseDir}")
      list(APPEND __add_before "${dir}")
    else()
      include_directories(AFTER SYSTEM "${dir}")
    endif()
  endforeach()
  include_directories(BEFORE ${__add_before})
endfunction()

macro(ocv_include_modules)
  include_directories(BEFORE "${OpenCV_INCLUDE_DIRS}")
endmacro()

# remove all matching elements from the list
macro(ocv_list_filterout lst regex)
  foreach(item ${${lst}})
    if(item MATCHES "${regex}")
      list(REMOVE_ITEM ${lst} "${item}")
    endif()
  endforeach()
endmacro()
