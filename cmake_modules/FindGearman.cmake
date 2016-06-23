cmake_minimum_required(VERSION 2.6)
cmake_policy(SET CMP0011 OLD)

if (WIN32)
    FIND_PATH(Gearman_ROOT_DIR
        libgearman/gearman.h
        HINTS
        $ENV{GFLAGS_ROOT})
else (WIN32)
    FIND_PATH(Gearman_ROOT_DIR
        libgearman.dylib
        HINTS
        /usr/local/lib          
        /opt/local/lib
    )
endif (WIN32)

IF(Gearman_ROOT_DIR)
    # We are testing only a couple of files in the include directories
    if (WIN32)
        FIND_PATH(Gearman_INCLUDE_DIRS
            libgearman/gearman.h
            HINTS
            ${Gearman_ROOT_DIR}/src/windows
        )
    else (WIN32)
        FIND_PATH(Gearman_INCLUDE_DIRS
            libgearman/gearman.h
            HINTS
            /usr/local/include
            /opt/local/include
            ${Gearman_ROOT_DIR}/src
        )
    endif (WIN32)


    # Find the libraries
    SET(Gearman_LIBRARY_DIRS ${Gearman_ROOT_DIR})
    
    # TODO: This can use some per-component linking
    if(MSVC)
        SET(_libglog_libpath_suffixes /Release /Debug)
        FIND_LIBRARY(Gearman_lib_release 
                    NAMES libglog
                    HINTS
                        ${Gearman_LIBRARY_DIRS}
                    PATH_SUFFIXES ${_libglog_libpath_suffixes})
        FIND_LIBRARY(Gearman_lib_debug
                    NAMES libglog-debug
                    HINTS
                        ${Gearman_LIBRARY_DIRS}
                    PATH_SUFFIXES ${_libglog_libpath_suffixes})
        SET(Gearman_lib optimized ${Gearman_lib_release} debug ${Gearman_lib_debug})
    else()
        FIND_LIBRARY(Gearman_lib gearman ${Gearman_LIBRARY_DIRS})
    endif()

    # set up include and link directory
    include_directories(${Gearman_INCLUDE_DIRS})
    link_directories(${Gearman_LIBRARY_DIRS})

    SET(Gearman_LIBS ${Gearman_lib})
    SET(Gearman_FOUND true)

    MARK_AS_ADVANCED(Gearman_INCLUDE_DIRS)
ELSE(Gearman_ROOT_DIR)
    FIND_PATH(Gearman_ROOT_DIR src)
    MARK_AS_ADVANCED(Gearman_ROOT_DIR)
    MESSAGE(SEND_ERROR "Cannot find Root directory of Libgearman")
    SET(Gearman_FOUND false)
ENDIF(Gearman_ROOT_DIR)
