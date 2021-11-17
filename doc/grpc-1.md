### gRPC的简介

* 进程间通信通常会采用消息传递的方式来实现，要么是同步的请求–响应风格，要么是异步的事件驱动风格。在同步通信风格中，客户端进程通过网络发送请求消息到服务器进程，并等待响应消息。在异步的事件驱动风格中，进程间会通过异步消息传递进行通信，这个过程会用到一个中介，也就是事件代理（event broker）。我们可以根据业务场景，选择希望实现的通信模式。

  当为现代云原生应用程序和微服务实现同步的请求–响应风格的通信时，最常见和最传统的方式就是将它们构建为 RESTful 服务。也就是说，将应用程序或服务建模为一组资源，这些资源可以通过 HTTP 的网络调用进行访问和状态变更。但是，对大多数使用场景来说，使用 RESTful 服务来实现进程间通信显得过于笨重、低效并且易于出错。我们通常需要扩展性强、松耦合的进程间通信技术，该技术比 RESTful 服务更高效。这也就是 gRPC 的优势所在，gRPC 是构建分布式应用程序和微服务的现代进程间通信风格（本章稍后会对比 gRPC 和 RESTful 服务）。gRPC 主要采用同步的请求–响应风格进行通信，但在建立初始连接后，它完全可以以异步模式或流模式进行操作。

![image-20210527133027474](\grpc.assets\image-20210527133027474.png)

> 服务器端：

* 通过重载服务基类，实现所生成的服务器端骨架的逻辑。
* 运行 gRPC 服务器，监听来自客户端的请求并返回服务响应。

> 客户端：

	*  当调用 gRPC 服务时，客户端的 gRPC 库会使用 protocol buffers，并将 RPC 的请求编排（marshal）为 protocol buffers 格式，然后将其通过 HTTP/2 进行发送。在服务器端，请求会被解排（unmarshal），对应的过程调用会使用 protocol buffers 来执行。

`gRPC 会使用 HTTP/2 来进行有线传输，HTTP/2 是一个高性能的二进制消息协议，支持双向的消息传递。`

### gRPC的使用

![img](\grpc.assets\c335c9631c270091.png)



```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

下载protoc 并配置环境变量

```protobuf
syntax = "proto3";

import "google/protobuf/wrappers.proto";

package gRPC;

option go_package = 'gRPC/service';

service ProductInfo {
  rpc addProduct(Product) returns (ProductID);
  rpc getProduct(ProductID) returns (Product);
}

message Product {
  string id = 1;
  string name = 2;
  string description = 3;
}

message ProductID {
  string value = 1;
}

message AAAA{
  google.protobuf.StringValue aa = 1;
}
```

```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  -I=. product.proto

-I={google/protobuf/wrappers.proto所在路径}  如果引用了 "google/protobuf/wrappers.proto" 这些内置proto
```





```go

type Server struct {
	productMap map[string]*pb.Product
	// UnimplementedProductInfoServer has implemented all mtd of service.ProductInfoServer
	pb.UnimplementedProductInfoServer
}

func NewServer() *Server {
	return &Server{productMap: nil}
}

//	AddProduct implements service.ProductInfoServer
func (s *Server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductID, error) {
	out, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Error while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}

	s.productMap[in.Id] = in
	return &pb.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

//	GetProduct implements service.ProductInfoServer
func (s *Server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
	value, exist := s.productMap[in.Value]
	if exist {
		return value, status.New(codes.OK, "").Err()
	}
	return value, status.New(codes.NotFound, "").Err()
}
```



```go
package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "gRPC/service"
	"gRPC/service/impl"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1"
	port    = ":20051"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go serverMain()

	go clintMain(&wg)

	wg.Wait()

	log.Printf("main function is exiting")
}

func serverMain() {
	// listen the tcp port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, impl.NewServer())
	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func clintMain(wg *sync.WaitGroup) {
	conn, err := grpc.Dial(address+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect :%v", err)
	}

	defer conn.Close()

	c := pb.NewProductInfoClient(conn)
	name := "Apple x"
	description := "just description"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)

	}
	log.Printf("Product ID: %s added successfully", r.Value)

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)

	}
	log.Println("Product", product.String())

	wg.Done()
	log.Printf("Client is closed")
}
```





### gRPC 的通信模式

一元 RPC、服务器端流 RPC、客户端流 RPC 以及双向流 RPC。

```protobuf
syntax = "proto3";

import "google/protobuf/wrappers.proto";

package demo;

service OrderManagement {
  //一元RPC模式
  rpc getOrder(google.protobuf.StringValue) returns (Order);

  //服务器端流RPC模式
  rpc searchOrders(google.protobuf.StringValue) returns (stream Order);

  //客户端流RPC模式
  rpc updateOrders(stream Order) returns (google.protobuf.StringValue);

  //双向流RPC模式
  rpc processOrders(stream google.protobuf.StringValue)returns (stream google.protobuf.StringValue);

}

message Order {
  string id = 1;
  repeated string items = 2;
  string description = 3;
  float price = 4;
  string destination = 5;
}
```





### gRPC 的底层原理

![img](\grpc.assets\dd7c1a48de94d091.png)1.客户端进程通过生成的存根调用 getProduct 方法。
2.客户端存根使用已编码的消息创建 HTTP POST 请求。在 gRPC 中，所有的请求都是 HTTP POST 请求，并且 content-type 前缀为 application/grpc。要调用的远程方法（/ProductInfo/getProduct）是以单独的 HTTP 头信息的形式发送的。
3.HTTP 请求消息通过网络发送到服务器端。
4.当接收到消息后，服务器端检查消息头信息，从而确定需要调用的服务方法，然后将消息传递给服务器端骨架。
5.服务器端骨架将消息字节解析成特定语言的数据结构。
6.借助解析后的消息，服务发起对 getProduct 方法的本地调用。





##### 字段索引

```protobuf
message ProductID {
  string value = 1;  // pb中字段需要指定顺序 即字段索引
  int32 price = 2 ;
  repeated string labels = 3;
}
```



##### 

##### varint编码

 编码过程：

1. 将一个uint64整数按照7位从小端开始往前进行划分	（有数值之前的段都舍弃）
2. 根据小端序规则将每段进行reverse（protobuf协议规定按小端序进行排列）
3. 末尾段补前面补0 其余段补1

```go
ex1:8(00000000 00000000 00000000 00001000)
	1.0001000
	2.0001000
	3.00001000 
	4字节-->1字节
================
ex2:188(00000000 00000000 00000000 10111100)
	1.0000001 0111100
	2.0111100 0000001
	3.10111100 00000001
	4字节-->2字节
================
ex3:2^30(1000000 00000000 00000000 00000000)
	1.0000100 0000000 0000000 0000000 0000000
	2.0000000 0000000 0000000 0000000 0000100
	3.10000000 10000000 10000000 10000000 00000100
	4字节-->5字节
```

 编码过程：

​	将编码过程reverse即可

![img](\grpc.assets\16f285fd3dd42565.png)

由此可见 对于varint对2^28以内的数有效果，数值普遍超过的话使用fixed64

##### zigzag编码

编码过程:

| 原始值 | 编码值 |
| :----: | :----: |
|   0    |   0    |
|   -1   |   1    |
|   1    |   2    |
|   -2   |   3    |
|   2    |   4    |
|  ...   |  ...   |

由于负数最高位为1，不能直接使用varint编码，采用zigzag编码后可以有效解决该问题。







##### 线路类型

| 线路类型 | 分类         | 字段类型                                                 |
| -------- | ------------ | -------------------------------------------------------- |
| 0        | Varint       | int32、int64、uint32、uint64、sint32、sint64、bool、enum |
| 1        | 64 位        | fixed64、sfixed64、double                                |
| 2        | 基于长度分隔 | string、bytes、嵌入式消息、打包的 repeated 字段          |
| 3        | 起始组       | groups（已废弃）                                         |
| 4        | 结束组       | groups（已废弃）                                         |
| 5        | 32 位        | fixed32、sfixed32、float                                 |







![image-20210531172736664](\grpc.assets\image-20210531172736664.png)

对于线路类型为0，1，5的无需length，对于线路类型位2（string） 需要length



```json
value="milk";
price=188;
labels=["delicious","drink"];
```



![image-20210531171613257](\grpc.assets\image-20210531171613257.png)



#### 基于长度前缀的消息分帧

长度前缀分帧是指在写入消息本身之前，写入长度信息，来表明每条消息的大小。如图 4-4 所示，已编码的二进制消息前面分配了 4 字节来指明消息的大小。在 gRPC 通信中，每条消息都有额外的 4 字节用来设置其大小。消息大小是一个有限的数字，为其分配 4 字节来表示消息的大小，也就意味着 gRPC 通信可以处理大小不超过 4GB 的所有消息。

![img](\grpc.assets\9ec23a8cefb96565.png)



对于客户端的请求消息，收件方是服务器；而对于响应消息，收件方则是客户端。在收件方一侧，当收到消息之后，首先要读取其第一字节，来检查该消息是否经过压缩。然后，收件方读取接下来的 4 字节，以获取编码二进制消息的大小，接着就可以从流中精确地读取确切长度的字节了。对于简单的消息，只需处理一条以长度为前缀的消息；而对于流消息，就会有多条以长度为前缀的消息要处理。

#### 基于HTTP/2的gRPC

在 HTTP/2 中，客户端和服务器端的所有通信都是通过一个 TCP 连接完成的，这个连接可以传送任意数量的双向字节流。为了理解 HTTP/2 的过程，最好熟悉下面这些重要术语。

* 数据流 `Stream`：在一个已建立的连接上的双向字节流。一个流可以携带一条或多条消息 `Message`
* 消息 `Message`:对应 `HTTP/1` 中的请求或者响应，包含一条或者多条 `Frame`，允许消息进行多路复用，客户端和服务器端能够将消息分解成独立的帧，交叉发送它们，然后在另一端进行重新组合
* 数据帧 `Frame`：最小单位，以二进制压缩格式存放 `HTTP/1` 中的内容。

远程调用中的消息以 HTTP/2 帧的形式进行发送，帧可能会携带一条 gRPC 长度前缀的消息，也可能在 gRPC 消息非常大的情况下，一条消息跨多帧



![img](\grpc.assets\03c4458d38f22386.png)

![img](\grpc.assets\d49971b9635dbc40.png)





1. 一元 RPC 模式

   ![img](\grpc.assets\701fe421d15b2bbb.png)

2. 服务器端流 RPC 模式

   ![img](\grpc.assets\bed7fa28f93865cd.png)

3. 客户端流 RPC 模式

   ![img](\grpc.assets\ca799bd46d965a63.png)

4. 双向流 RPC 模式

   ![img](\grpc.assets\cf438c762fab2b9d.png)

