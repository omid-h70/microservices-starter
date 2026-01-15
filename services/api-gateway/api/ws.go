package api

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/proto/driver"
)

var (
	userIDVar   = "userID"
	packageSlug = "packageSlug"
)

/*
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
*/

var (
	connManger = messaging.NewConnManger()
)

func (api *HttpApi) handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := connManger.Upgrade(w, r)
	if err != nil {
		log.Printf("websocket upgrade failed %v", err)
		return
	}

	//for a case that we didn't have any userID
	defer conn.Close()

	userID := r.URL.Query().Get(userIDVar)
	if len(userID) == 0 {
		log.Println("userID is not provided")
	}

	connManger.Add(userID, conn)
	//for a case that we have userID
	defer connManger.Remove(userID)

	//what queues we want to consume
	queues := []string{
		messaging.NotifyDriverNoDriverFoundQueue,
		messaging.NotifyDriverAssingQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(api.rabbitmq, connManger, q)

		if err := consumer.Start(); err != nil {
			log.Printf("failed to start the consumer %v", err)
		}
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read ws message failed %v", err)
			return
		}
		//conn.WriteJSON("ok")
		log.Printf("got ws message %s", string(msg))
	}
}

func (api *HttpApi) handleDriverWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := connManger.Upgrade(w, r)
	if err != nil {
		log.Printf("websocket upgrade failed %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get(userIDVar)
	if len(userID) == 0 {
		log.Println(userIDVar + "is not provided")
	}

	packageSlug := r.URL.Query().Get(packageSlug)
	if len(packageSlug) == 0 {
		log.Println(packageSlug + "is not provided")
	}

	connManger.Add(userID, conn)
	driverService, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Printf("make grpc NewDriverServiceClient failed %v", err)
		return
	}
	ctx := r.Context()
	defer func() {
		connManger.Remove(userID)

		driverService.Client.UnregisterDriver(ctx, &driver.RegisterDriverRequest{
			DriverId:    userID,
			PackageSlug: packageSlug,
		})

		driverService.Close()
		log.Printf("driver %s unregistered", userID)
	}()

	driverData, err := driverService.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverId:    userID,
		PackageSlug: packageSlug,
	})

	if err != nil {
		log.Printf("RegisterDriver failed %v", err)
		return
	}

	//TODO refactor here - Done !
	/*
		type Driver struct {
			UserID         string `json:"userID"`
			Name           string `json:"name"`
			ProfilePicture string `json:"profilePicture"`
			CarPlate       string `json:"carPlate"`
			PackageSlug    string `json:"packageSlug"`
		}*/

	msg := contracts.WSMessage{
		Type: contracts.DriverCmdRegister,
		Data: driverData.Driver,
		/*
			Data: Driver{
				UserID:         userID,
				Name:           "dodo",
				ProfilePicture: util.GetRandomAvatar(1),
				CarPlate:       "ABCDEFG123",
				PackageSlug:    packageSlug,
			},
		*/
	}

	if err := connManger.SendMessage(userID, msg); err != nil {
		log.Printf("write ws response failed %v", err)
		return
	}

	//what queues we want to consume - test
	queues := []string{
		messaging.FindAvailableDriversQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(api.rabbitmq, connManger, q)

		if err := consumer.Start(); err != nil {
			log.Printf("failed to start the consumer %v", err)
		}
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read ws message failed %v", err)
			return
		}
		log.Printf("got ws message %v", msg)

		type driverMessage struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}

		var driverMsg driverMessage
		if err := json.Unmarshal(msg, &driverMsg); err != nil {
			log.Printf("error unmarshal driver message failed %v", err)
			continue
		}

		switch driverMsg.Type {
		case contracts.DriverCmdLocation:
			//TODO: add it later
		case contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline:
			//Forward the message to rabbit
			if err := api.rabbitmq.PublishMessage(ctx, driverMsg.Type, contracts.AmqpMessage{
				OwnerID: userID,
				Data:    driverMsg.Data,
			}); err != nil {
				log.Printf("error publishing message to rabbit %v", err)
			}
		default:
			log.Printf("unknown message type %s", driverMsg.Type)
		}

	}
}
