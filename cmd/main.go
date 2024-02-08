package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	jaegar "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/vishnusunil243/UserService/db"
	"github.com/vishnusunil243/UserService/initializer"
	"github.com/vishnusunil243/UserService/service"
	servicediscovery "github.com/vishnusunil243/UserService/service_discovery"
	"github.com/vishnusunil243/proto-files/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf(err.Error())
	}
	addr := os.Getenv("DB_KEY")

	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	services := initializer.Initializer(DB)
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, services)
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen in port 8082")
	}
	log.Printf("user server listening on port 8082")
	go func() {
		time.Sleep(5 * time.Second)
		servicediscovery.RegisterService()
	}()
	healthService := &service.HealthChecker{}

	grpc_health_v1.RegisterHealthServer(server, healthService)
	tracer, closer := initTracer()
	defer closer.Close()
	service.RetrieveTracer(tracer)
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to listen on port 8082")
	}

}
func initTracer() (tracer opentracing.Tracer, closer io.Closer) {
	jargarEndpoint := "http://localhost:14268/api/traces"
	cfg := &config.Configuration{
		ServiceName: "user-service",
		Sampler: &config.SamplerConfig{
			Type:  jaegar.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: jargarEndpoint,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("updated")
	return
}
