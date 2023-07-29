FROM ghcr.io/ggerganov/llama.cpp:light as bin
FROM golang
WORKDIR /workspace
COPY --from=bin /main /workspace
COPY llama-2-7b-chat.ggmlv3.q8_0.bin /workspace
COPY server.go client.go .
# RUN rm $GOPATH/go.mod
RUN go mod init lol
RUN go get github.com/golang/glog && go get golang.org/x/net/websocket

# RUN git clone https://github.com/ggerganov/llama.cpp.git && cd llama.cpp && make -j && cp ./main /workspace
CMD ["go", "run", "server.go"]
