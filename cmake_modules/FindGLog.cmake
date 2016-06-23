# - Try to find GLog
#
# The following variables are optionally searched for defaults
#  GLog_ROOT_DIR:            Base directory where all GLog components are found
#
# The following are set after configuration is done: 
#  GLog_FOUND
#  GLog_INCLUDE_DIRS
#  GLog_LIBS
#  GLog_MODEL_DIRS
#  GLog_LIBRARY_DIRS
cmake_minimum_required(VERSION 2.6)
cmake_policy(SET CMP0011 OLD)

if (WIN32)
    FIND_PATH(GLog_ROOT_DIR
        glog/logging.h
        HINTS
        $ENV{GFLAGS_ROOT})
else (WIN32)
    FIND_PATH(GLog_ROOT_DIR
        libglog.dylib
        HINTS
        /usr/local/lib          
        /opt/local/lib
    )
endif (WIN32)

IF(GLog_ROOT_DIR)
    # We are testing only a couple of files in the include directories
    if (WIN32)
        FIND_PATH(GLog_INCLUDE_DIRS
            glog/loggin.h
            HINTS
            ${GLog_ROOT_DIR}/src/windows
        )
    else (WIN32)
        FIND_PATH(GLog_INCLUDE_DIRS
            glog/logging.h
            HINTS
            /usr/local/include
            /opt/local/include
            ${GLog_ROOT_DIR}/src
        )
    endif (WIN32)


    # Find the libraries
    SET(GLog_LIBRARY_DIRS ${GLog_ROOT_DIR})
    
    # TODO: This can use some per-component linking
    if(MSVC)
        SET(_libglog_libpath_suffixes /Release /Debug)
        FIND_LIBRARY(GLog_lib_release 
                    NAMES libglog
                    HINTS
                        ${GLog_LIBRARY_DIRS}
                    PATH_SUFFIXES ${_libglog_libpath_suffixes})
        FIND_LIBRARY(GLog_lib_debug
                    NAMES libglog-debug
                    HINTS
                        ${GLog_LIBRARY_DIRS}
                    PATH_SUFFIXES ${_libglog_libpath_suffixes})
        SET(GLog_lib optimized ${GLog_lib_release} debug ${GLog_lib_debug})
    else()
        FIND_LIBRARY(GLog_lib glog ${GLog_LIBRARY_DIRS})
    endif()

    # set up include and link directory
    include_directories(${GLog_INCLUDE_DIRS})
    link_directories(${GLog_LIBRARY_DIRS})

    SET(GLog_LIBS ${GLog_lib})
    SET(GLog_FOUND true)

    MARK_AS_ADVANCED(GLog_INCLUDE_DIRS)
ELSE(GLog_ROOT_DIR)
    FIND_PATH(GLog_ROOT_DIR src)
    MARK_AS_ADVANCED(GLog_ROOT_DIR)
    MESSAGE(SEND_ERROR "Cannot find Root directory of Google Logger")
    SET(GLog_FOUND false)
ENDIF(GLog_ROOT_DIR)
