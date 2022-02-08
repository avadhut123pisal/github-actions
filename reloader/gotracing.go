package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

func detectResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("product-service"),
			semconv.ServiceNamespaceKey.String("US-West-1"),
			semconv.HostNameKey.String("prodsvc.us-west-1.example.com"),
			semconv.HostTypeKey.String("system"),
		),
	)
	return res, err
}

func spanExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	// Steps for mTLS
	// Step 1: Create a CA certificate pool and add ca-certificate to it
	caCert, err := ioutil.ReadFile("/home/ubuntu/tls-certificates/ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Step 2: Read the client application TLS key pair to create certificate
	cert, err := tls.LoadX509KeyPair("/home/ubuntu/tls-certificates/client.pem", "/home/ubuntu/tls-certificates/client-key.pem")
	if err != nil {
		return nil, err
	}
	// Step 3:
	// Use credentials.NewTLS for mTLS
	cred := credentials.NewTLS(&tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	})

	// create exporter
	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithTLSCredentials(cred),
		otlptracegrpc.WithEndpoint(os.Getenv("OTLP_ENDPOINT")),
		// otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
}

func main() {

	ctx := context.Background()
	traceExporter, err := spanExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize stdouttrace export pipeline: %v", err)
	}

	res, err := detectResource(ctx)
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	tracer := otel.Tracer("example.com/basictracer")
	priority := attribute.Key("business.priority")
	appEnv := attribute.Key("prod.env")

	func(ctx context.Context) {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "HTTP GET /products/{id}")
		defer span.End()
		span.AddEvent("Authentication", trace.WithAttributes(attribute.String("Username", "TestUser")))
		span.AddEvent("Products", trace.WithAttributes(attribute.Int("ID", 100)))
		span.SetAttributes(appEnv.String("UAT"))

		func(ctx context.Context) {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "SELECT * from Products where pID={id}")
			defer span.End()
			span.SetAttributes(priority.String("CRITICAL"))
			span.AddEvent("Datababse", trace.WithAttributes(attribute.Int("Count", 20)))
		}(ctx)
	}(ctx)

	http.ListenAndServe(":3030", nil)
}
