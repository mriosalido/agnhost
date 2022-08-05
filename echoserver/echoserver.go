package echoserver

// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"strconv"

// 	"google.golang.org/grpc/credentials"
// )

import (
	"net"
	"bytes"
	"encoding/json"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
 	"context"
	"fmt"
	"net/http"
	"strconv"

	"agnhost/version"

	"github.com/spf13/cobra"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
 	"google.golang.org/grpc"
 	"google.golang.org/grpc/reflection"
 	"google.golang.org/grpc/peer"
)

//go:generate protoc -I.  --go_out=plugins=grpc:. ./hello.proto

var CmdEchoServer = &cobra.Command{
	Use:   "echoserver",
	Args:  cobra.MaximumNArgs(0),
	Run:   main,
}

func StartHttpServer(port int) {
	http.HandleFunc("/", handler)
	server := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(port),
		Handler: http.HandlerFunc(handler),
	}
	addr := "0.0.0.0:" + strconv.Itoa(port)
	fmt.Printf("listening on %s, http1.0, http1.1\n", addr)
	server.ListenAndServe()
}

func StartTCPServer(port int) {
	addr := ":" + strconv.Itoa(port)
	l, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	defer l.Close()

	fmt.Printf("listening on %s, tcp\n", addr)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go handleTcpConnection(conn)
	}
}

func StartUDPServer(port int) {
	addr := ":" + strconv.Itoa(port)
	pc, err := net.ListenPacket("udp", addr)

	if err != nil {
		panic(err)
	}

	defer pc.Close()

	fmt.Printf("listening on %s, udp\n", addr)

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go handleUdpPackageConn(pc, addr, buf[:n])
	}
}

func StartHttp2CleartextServer(port int) {
	server := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(port),
		Handler: h2c.NewHandler(http.HandlerFunc(handler), &http2.Server{}),
	}

	fmt.Printf("listening on %s, support http1.0, http1.1, non-TLS HTTP/2 (aka h2c, upgrade, prior knowledge)\n", server.Addr)
	server.ListenAndServe()
}

func StartGrpcServer(port int) {
	addr := ":" + strconv.Itoa(port)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	RegisterHelloWorldServer(s, &GrpcServer{})

	fmt.Printf("listening on %s, grpc\n", addr)

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

type GrpcServer struct{
	UnimplementedHelloWorldServer
}

func (*GrpcServer) Greeting(ctx context.Context, msg *GreetingMessage) (*GreetingReply, error) {
	name, err := os.Hostname()

	if err != nil {
		panic(err)
	}

	p, ok := peer.FromContext(ctx)

	if !ok {
		panic(fmt.Errorf("get peer from context not success"))
	}

	return &GreetingReply{
		Hostname:      name,
		ClientAddress: p.Addr.String(),
		AuthInfo:      "",
	}, nil
}

func main(cmd *cobra.Command, args []string) {
	go StartHttpServer(8001)
	go StartHttp2CleartextServer(8002)
	// go StartHttp2TLSServer(8003)
	go StartGrpcServer(8004)
	// go StartGrpcWithTLSServer(8005)
	go StartTCPServer(8006)
	StartUDPServer(8007)
}

func getClientIP(req *http.Request) string {
	if req.Header.Get("X-Envoy-External-Address") != "" {
		return req.Header.Get("X-Envoy-External-Address")
	}

	ra, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ra
}

func printf(w io.Writer, msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, msg, args...)
	fmt.Printf(msg, args...)
}

func handler(w http.ResponseWriter, req *http.Request) {
	name, err := os.Hostname()

	if err != nil {
		panic(err)
	}

	printf(w, "Version: %s\n", version.Version)
	printf(w, "HostName: %s\n", name)

	printf(w, "\nRequest Info:\n")
	printf(w, "    content-length: %d\n", req.ContentLength)
	printf(w, "    remote address: %s\n", req.RemoteAddr)
	printf(w, "    realIP: %s\n", getClientIP(req))
	printf(w, "    method: %s\n", req.Method)
	printf(w, "    path: %s\n", req.URL.Path)
	printf(w, "    query: %s\n", req.URL.RawQuery)
	printf(w, "    request_version: %s\n", req.Proto)
	printf(w, "    uri: %s\n", req.RequestURI)
	printf(w, "    tls: %t\n", req.TLS != nil)

	printf(w, "\nHeaders:\n")

	for name, headers := range req.Header {
		for _, h := range headers {
			printf(w, "    %v: %v\n", name, h)
		}
	}

	if req.Header.Get("Kalm-Sso-Userinfo") != "" {
		printf(w, "\nKalm SSO:\n")
		claims, err := base64.RawStdEncoding.DecodeString(req.Header.Get("Kalm-Sso-Userinfo"))

		if err != nil {
			printf(w, "Base64 decode error: %s\n", err.Error())
		} else {
			var out bytes.Buffer
			prefix := "  "
			if err := json.Indent(&out, claims, prefix, "  "); err != nil {
				printf(w, "json indent error: %s\n", err.Error())
			} else {
				printf(w, "%s%s\n", prefix, string(out.Bytes()))
			}
		}
	}

	printf(w, "\nBody:\n")
	if req.Body != nil && req.ContentLength > 0 {
		defer req.Body.Close()
		bs, err := ioutil.ReadAll(req.Body)

		if err != nil {
			printf(w, "Read body error: %s\n", err.Error())
		} else {
			printf(w, "%s", bs)
		}
	} else {
		printf(w, "No Body\n")
	}

	printf(w, "\nEnvironment Variables:\n")
	for _, env := range os.Environ() {
		printf(w, "%s\n", env)
	}

	printf(w, "\n")
}

func handleTcpConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	defer conn.Close()

	for {
		reqLen, err := conn.Read(buf)

		if err != nil {

			// client exit
			if err == io.EOF {
				break
			}

			fmt.Println("Error reading:", err.Error())
			break
		}

		_, _ = conn.Write([]byte(fmt.Sprintf("Message received: %s\n", buf[:reqLen])))
	}
}

func handleUdpPackageConn(pc net.PacketConn, addr net.Addr, buf []byte) {
	pc.WriteTo([]byte(fmt.Sprintf("Message received: %s\n", string(buf))), addr)
}
