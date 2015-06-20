FROM armv7/armhf-ubuntu

RUN apt-get update &&\
    apt-get install -y build-essential pkg-config libzmq3-dev git bzr wget

RUN cd /tmp &&\
    wget http://dave.cheney.net/paste/go1.4.2.linux-arm~multiarch-armv7-1.tar.gz &&\
    tar xf go1.4.2.linux-arm~multiarch-armv7-1.tar.gz &&\
    mv go /

RUN cd /tmp &&\
    export GOROOT=/go &&\
    export GOPATH=/opt/gogadgets &&\
    wget https://github.com/cswank/gogadgets/tarball/master &&\
    mkdir -p /opt/gogadgets/src/github.com/cswank/gogadgets &&\
    cd /opt/gogadgets/src/github.com/cswank &&\
    tar xf /tmp/master -C gogadgets --strip-components 1 &&\
    cd gogadgets/gogadgets &&\
    /go/bin/go get -tags zmq_3_x github.com/alecthomas/gozmq &&\
    /go/bin/go get &&\
    /go/bin/go install

EXPOSE 6111 6112

CMD /opt/gogadgets/bin/gadgets
