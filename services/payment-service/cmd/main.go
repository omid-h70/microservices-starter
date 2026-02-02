package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"ride-sharing/services/payment-service/internal/infrastructure/events"
	"ride-sharing/services/payment-service/internal/infrastructure/stripe"
	"ride-sharing/services/payment-service/internal/service"
	"ride-sharing/services/payment-service/pkg/types"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"syscall"
)

var (
	grpcServeAdrr = env.GetString("GRPC_ADDR", ":9004")
	rabbitmqURI   = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672")
	appURL        = env.GetString("APP_URL", "http://localhost:3000")
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	stripeCfg := &types.PaymentConfig{
		StripeSecretKey: env.GetString("STRIP_SECRET_KEY", ""),
		SuccessURL:      env.GetString("STRIP_SUCCESS_URL", appURL+"?payment=success"),
		CancelURL:       env.GetString("STRIP_CANCEL_URL", appURL+"?payment=cancel"),
	}

	if stripeCfg.StripeSecretKey == "" {
		log.Fatal("stripe secret key is empty")
	}

	paymentProcessor := stripe.NewStripeClient(stripeCfg)
	svc := service.NewPaymentService(paymentProcessor)

	rabbitmq, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatalf("rabbit is down %v", err)
	}
	defer rabbitmq.Close()
	log.Printf("Rabbitmq started on %s ", rabbitmqURI)

	// Start Driver Consumer
	tripConsumer := events.NewTripConsumer(rabbitmq, svc)
	go tripConsumer.Listen(context.Background(), messaging.FindAvailableDriversQueue)

	//wait for shutdown signal
	<-ctx.Done()
	log.Println("i'm done here !")
}
