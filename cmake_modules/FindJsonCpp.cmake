# - Try to find JsonCpp
#
# The following variables are optionally searched for defaults
#  JsonCpp_ROOT_DIR:            Base directory where all JsonCpp components are found
#
# The following are set after configuration is done: 
#  JsonCpp_FOUND
#  JsonCpp_INCLUDE_DIRS
#  JsonCpp_LIBS
#  JsonCpp_MODEL_DIRS
#  JsonCpp_LIBRARY_DIRS
cmake_minimum_required(VERSION 2.6)
cmake_policy(SET CMP0011 OLD)

if (WIN32)
    FIND_PATH(JsonCpp_ROOT_DIR
        json/features.h
        HINTS
        $ENV{JSONCPP_ROOT})
else (WIN32)
    FIND_PATH(JsonCpp_ROOT_DIR
        libjsoncpp.dylib
        HINTS
        /opt/local/lib
        /usr/local/lib          
        /usr/lib
    )
endif (WIN32)

IF(JsonCpp_ROOT_DIR)
    # We are testing only a couple of files in the include directories
    if (WIN32)
        FIND_PATH(JsonCpp_INCLUDE_DIRS
            json/features.h
            HINTS
            ${JsonCpp_ROOT_DIR}/src/windows
        )
    else (WIN32)
        FIND_PATH(JsonCpp_INCLUDE_DIRS
            json/features.h
            HINTS
            /usr/include
            /usr/include/jsoncpp
            /usr/local/include
            /opt/local/include
            ${JsonCpp_ROOT_DIR}/src
        )
    endif (WIN32)

    # Find the libraries
    SET(JsonCpp_LIBRARY_DIRS ${JsonCpp_ROOT_DIR})


    # TODO: This can use some per-component linking
    if(MSVC)
        SET(_jsoncpp_libpath_suffixes /Release /Debug)
        FIND_LIBRARY(JsonCpp_lib_release 
                    NAMES libjsoncpp
                    HINTS
                        ${JsonCpp_LIBRARY_DIRS}
                    PATH_SUFFIXES ${_jsoncpp_libpath_suffixes})
        FIND_LIBRARY(JsonCpp_lib_debug
                    NAMES libjsoncpp-debug
                    HINTS
                        ${JsonCpp_LIBRARY_DIRS}
                    PATH_SUFFIXES ${_jsoncpp_libpath_suffixes})
        SET(JsonCpp_lib optimized ${JsonCpp_lib_release} debug ${JsonCpp_lib_debug})
    else()
        FIND_LIBRARY(JsonCpp_lib jsoncpp ${JsonCpp_LIBRARY_DIRS})
    endif()

    # set up include and link directory
    include_directories(${JsonCpp_INCLUDE_DIRS})
    link_directories(${JsonCpp_LIBRARY_DIRS})

    # MESSAGE("JsonCpp_LIBRARY_DIRS" ${JsonCpp_LIBRARY_DIRS})
    # MESSAGE("JsonCpp_lib" ${JsonCpp_lib})
    SET(JsonCpp_LIBS ${JsonCpp_lib})

    SET(JsonCpp_FOUND true)

    MARK_AS_ADVANCED(
        JsonCpp_INCLUDE_DIRS
        )
ELSE(JsonCpp_ROOT_DIR)
    FIND_PATH(JsonCpp_ROOT_DIR
        src
    )
    MARK_AS_ADVANCED(JsonCpp_ROOT_DIR)
    MESSAGE(SEND_ERROR "Cannot find Root directory of JsonCpp")
    SET(JsonCpp_FOUND false)
ENDIF(JsonCpp_ROOT_DIR)
