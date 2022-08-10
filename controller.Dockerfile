FROM golang:1.19
MAINTAINER vbha.mmk@gmail.com
WORKDIR /dbuild
COPY . .
RUN go build -mod=vendor -ldflags="-w -s" -o bin/controller controller/main.go
RUN chmod +x tools/hamctl

RUN wget -O distcc.tar.gz https://github.com/distcc/distcc/releases/download/v3.4/distcc-3.4.tar.gz; \
    tar -xf distcc.tar.gz
RUN apt-get update; \
    apt-get install -y python3-pip python3-dev libiberty-dev git fakeroot build-essential ncurses-dev xz-utils libssl-dev bc flex libelf-dev bison; \
    current_dir=$PWD; \
    cd /usr/local/bin; \
    ln -s /usr/bin/python3 python; \
    pip3 --no-cache-dir install --upgrade pip; \
    rm -rf /var/lib/apt/lists/*; \
    cd $current_dir;
RUN (cd distcc-3.4 ; ./configure && make && make install && update-distcc-symlinks)

RUN ln -s /usr/lib/distcc /usr/local/lib/distcc

CMD ["/dbuild/bin/controller"]
