FROM opensuse/tumbleweed:latest

RUN set -x                        \
  && zypper in --no-recommends -y \
       make                       \
       automake                   \
       autoconf                   \
       libtool                    \
       libuuid-devel              \
       pkgconfig                  \
       glib2-devel                \
       libxml2-devel              \
       libxslt-devel              \
       libbz2-devel               \
       libqb-devel                \
       python3                    \
       gettext-tools              \
       libcorosync-devel          \
       corosync                   \
       findutils                  \
       gcc12                      \
       git                        \
       util-linux                 \
       perf                       \
                                  \
  && zypper clean --all           \
                                  \
  && update-alternatives          \
       --install                  \
       /usr/bin/gcc               \
       gcc                        \
       /usr/bin/gcc-12 100        \
                                  \
  && curl -Ss https://dl.google.com/go/go1.19.linux-amd64.tar.gz \
       | tar -C /usr/local -xzf - \
                                  \
  && update-alternatives          \
       --install                  \
       /usr/bin/go                \
       golang                     \
       /usr/local/go/bin/go       \
       100                        \
                                  \
  && mkdir -p /etc/init.d         \
                                  \
  && git clone https://github.com/brendangregg/FlameGraph /FlameGraph

RUN echo '[ ! -z "$TERM" -a -r /etc/motd ] && cat /etc/issue && cat /etc/motd' >> /etc/bash.bashrc

COPY etc/* /etc/



