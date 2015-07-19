FROM armv7/armhf-ubuntu

RUN apt-get update &&\
    apt-get install -y build-essential pkg-config libzmq3-dev git bzr wget

RUN cd /tmp &&\
    wget http://dave.cheney.net/paste/go1.4.2.linux-arm~multiarch-armv7-1.tar.gz &&\
    tar xf go1.4.2.linux-arm~multiarch-armv7-1.tar.gz &&\
    mv go /

ADD . /opt/gogadgets/src/github.com/cswank/gogadgets

RUN cd /opt/gogadgets/src/github.com/cswank/gogadgets/gogadgets &&\
    export GOROOT=/go &&\
    export GOPATH=/opt/gogadgets &&\
    /go/bin/go get -tags zmq_3_x github.com/alecthomas/gozmq &&\
    /go/bin/go get &&\
    /go/bin/go install

EXPOSE 6111 6112

CMD /opt/gogadgets/bin/gadgets
