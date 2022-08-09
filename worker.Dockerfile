FROM golang:1.19
MAINTAINER vbha.mmk@gmail.com
WORKDIR /dbuild
#COPY . .
#RUN go build -mod=vendor -ldflags="-w -s" -o bin/worker worker/main.go
#RUN chmod +x tools/hamctl

ENV HOME=/home/distcc
RUN useradd -s /bin/bash distcc

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


# Define how to start distccd by default
# (see "man distccd" for more information)
ENTRYPOINT [\
  "distccd", \
  "--daemon", \
  "--no-detach", \
  "--user", "distcc", \
  "--port", "3632", \
  "--stats", \
  "--stats-port", "3633", \
  "--log-stderr", \
  "--listen", "0.0.0.0", \
  "--log-level", "debug" \
]

# By default the distcc server will accept clients from everywhere.
# Feel free to run the docker image with different values for the
# following params.
CMD [\
  "--allow", "0.0.0.0/0" \
]

# 3632 is the default distccd port
# 3633 is the default distccd port for getting statistics over HTTP
EXPOSE \
  3632/tcp \
  3633/tcp

# We check the health of the container by checking if the statistics
# are served. (See
# https://docs.docker.com/engine/reference/builder/#healthcheck)
HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://0.0.0.0:3633/ || exit 1
