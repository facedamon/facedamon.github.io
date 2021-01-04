
# 为什么要使用protobuf

> protobuf既Protocol Buffers, Google开发的一种数据描述语言，是一种轻便高效的结构化数据存储格式，与语言、平台无关，可扩展可序列化。protobug以二进制方式存储，占用空间小。

&emsp;&emsp;protobuf广泛地应用于远程过程调用地二进制传输，使用protobuf地目的非常简单，为了获得更高地性能。传输前使用protobuf编码，接收方再进行解码，可以显著地降低二进制传输的大小。另一方面，protobuf可非常适合传输结构化数据，便于通信字段的扩展。

&emsp;&emsp;使用protobuf一般分为以下2步：

- 按照protobuf的语法，在`.proto`文件中定义数据结构，并使用`protoc`生产Go代码(`.proto`文件是跨平台的，还可以生成C、Java等其它源码文件)
- 在项目代码中引用生成的Go代码

# 使用protobuf通信

    syntax = "proto3";

    package geecachepb;

    message Request {
        string group = 1;
        string key = 2;
    }

    message Response {
        bytes value = 1;
    }

    service GroupCache {
        rpc Get(Request) returns (Response);
    }

- `Request`包含2个字段，group和key，这与我们之前定义的接口`/_geec/<group>/<name>`所需的参数吻合
- `Response`包含1个字段，bytes，类型为byte数组，与之前吻合

&emsp;&emsp;生成geec.pb.go

    $ protoc --go_out=. *.proto
    $ ls
    geecpb.pb.go geecpb.proto

&emsp;&emsp;可以看到`geecpb.pb.go` 中有如下数据类型：

    type Request struct {
        Group string `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
        Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitepty"`
        ...
    }
    type Response struct {
        Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
    }

&emsp;&emsp;接下来，修改`peers.go`中`PeerGetter`接口，参数使用`geecpb.pb.go`中的数据类型

