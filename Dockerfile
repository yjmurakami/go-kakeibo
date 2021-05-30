FROM registry.access.redhat.com/ubi7/ubi:7.8

ARG GOFILE="go1.16.4.linux-amd64.tar.gz"
ARG GOPATH="/go"

RUN yum -y update
RUN yum -y install wget

# Go
RUN wget https://golang.org/dl/${GOFILE}
RUN tar -C /usr/local -xzf ${GOFILE}
RUN rm ${GOFILE}
ENV GOPATH ${GOPATH}
ENV PATH $PATH:${GOPATH}/bin
ENV PATH $PATH:/usr/local/go/bin

WORKDIR ${GOPATH}/src

# Gitのインストールに必要
RUN yum -y install make autoconf gcc gettext zlib-devel curl-devel

# Git 最新バージョンのインストール
# https://www.digitalocean.com/community/tutorials/how-to-install-git-on-centos-7
RUN wget https://github.com/git/git/archive/v2.30.0.tar.gz
RUN tar -xzf v2.30.0.tar.gz
WORKDIR ${GOPATH}/src/git-2.30.0
RUN make configure
RUN ./configure --prefix=/usr/local
RUN make install

# openapi-generator
RUN yum -y install java-1.8.0-openjdk
RUN mkdir ${GOPATH}/bin
RUN wget https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/5.0.0/openapi-generator-cli-5.0.0.jar -O ${GOPATH}/bin/openapi-generator-cli.jar

WORKDIR ${GOPATH}/src

# go install : ${GOPATH}/bin にコンパイル
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/davidrjenni/reftools/cmd/fillstruct@latest
RUN go install github.com/cweill/gotests/...@latest
RUN go install github.com/josharian/impl@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/ramya-rao-a/go-outline@latest
RUN go install github.com/vektra/mockery/v2/...@latest
RUN go install golang.org/x/tools/gopls@latest

# フォークされたGitリポジトリは go get でコンパイルできないため、git clone & go install を使用
# 実際のパス(github.com/yjmurakami/xo)とソース内のインポートパス(github.com/xo/xo)が異なり、エラーになる
WORKDIR ${GOPATH}/src/github.com/yjmurakami
RUN git clone https://github.com/yjmurakami/xo.git
WORKDIR ${GOPATH}/src/github.com/yjmurakami/xo
RUN go install
WORKDIR ${GOPATH}/src

# コンパイル後は不要になったソースを削除
RUN rm -rf *

CMD ["sh"]
