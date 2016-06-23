mkdir -p vendor
cd vendor
(wget -c 'https://googletest.googlecode.com/files/gtest-1.7.0.zip' \
    && unzip -n gtest-1.7.0.zip && ln -s gtest-1.7.0 gtest)

(cd gtest && cmake . && make)

(wget -c -O gflags.zip https://github.com/gflags/gflags/archive/v2.1.2.zip \
 && unzip -n gflags.zip && ln -s gflags-2.1.2 gflags)
(cd gflags && cmake . && make)

(wget -c -O glog.zip https://github.com/google/glog/archive/v0.3.4.zip \
  && unzip -n glog.zip && ln -s glog-0.3.4 glog)
(cd glog && ./configure && make && sudo make install)
