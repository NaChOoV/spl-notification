package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"spl-notification/internal/config"
	"spl-notification/internal/errors"
	"spl-notification/internal/model"
	"time"

	"cloud.google.com/go/pubsub"
)

type notificationServiceImpl struct {
	enviromentConfig   *config.EnvironmentConfig
	whatsappClient     *http.Client
	pubsubClient       *pubsub.Client
	pubsubTopic        *pubsub.Topic
	pubsubSubscription *pubsub.Subscription
}

func NewNotificationServiceImpl(
	enviromentConfig *config.EnvironmentConfig,
) NotificationService {
	ctx := context.Background()

	// Inicializar cliente de Pub/Sub
	pubsubClient, err := pubsub.NewClient(ctx, enviromentConfig.PubSubProjectID)
	if err != nil {
		panic(fmt.Sprintf("Error al crear cliente de Pub/Sub: %v", err))
	}

	// Obtener el topic
	topic := pubsubClient.Topic(enviromentConfig.PubSubTopicID)

	// Obtener la suscripción
	subscription := pubsubClient.Subscription(enviromentConfig.PubSubSubscriptionID)

	return &notificationServiceImpl{
		enviromentConfig: enviromentConfig,
		whatsappClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		pubsubClient:       pubsubClient,
		pubsubTopic:        topic,
		pubsubSubscription: subscription,
	}
}

func (n *notificationServiceImpl) SendNotification(requests []*model.NotificationRequest) *errors.AppError {
	ctx := context.Background()

	for _, request := range requests {
		messageData, err := json.Marshal(request)
		if err != nil {
			return errors.NewAppError("NotificationService",
				fmt.Errorf("error serializing notification: %w", err))
		}

		msg := &pubsub.Message{
			Data: messageData,
			Attributes: map[string]string{
				"type":     request.Type.String(),
				"chatId":   request.ChatID,
				"run":      request.Run,
				"location": fmt.Sprintf("%d", request.Location),
			},
		}

		result := n.pubsubTopic.Publish(ctx, msg)

		// Esperar a que se complete la publicación
		_, err = result.Get(ctx)
		if err != nil {
			return errors.NewAppError("NotificationService",
				fmt.Errorf("error publishing message to Pub/Sub: %w", err))
		}
	}

	return nil
}

func (n *notificationServiceImpl) HandleNotification() {
	ctx := context.Background()

	log.Println("[NotificationService] Starting Pub/Sub notification consumer...")

	err := n.pubsubSubscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var notificationRequest model.NotificationRequest
		if err := json.Unmarshal(msg.Data, &notificationRequest); err != nil {
			log.Printf("Error deserializing message: %v\n", err)
			msg.Nack() // Reject the message to retry
			return
		}

		log.Printf("[NotificationService] Message received: Type=%s, ChatID=%s, Run=%s, Location=%d\n",
			notificationRequest.Type.String(),
			notificationRequest.ChatID,
			notificationRequest.Run,
			notificationRequest.Location,
		)

		err := n.notifyTemplate(&notificationRequest)
		if err != nil {
			log.Printf("%v\n", err)
			msg.Nack()
			return
		}

		// Message sended with success
		msg.Ack()
	})

	if err != nil {
		log.Printf("Error in Pub/Sub subscription: %v\n", err)
	}
}

func (n *notificationServiceImpl) notifyTemplate(request *model.NotificationRequest) *errors.AppError {
	body := map[string]string{
		"chatId":   request.ChatID,
		"fullName": request.FullName,
		"location": request.LocationName(),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return n.error(err)
	}

	finalPath := "notify-entry"
	if request.Type == model.NotificationTypeExit {
		finalPath = "notify-exit"
	}

	req, err := http.NewRequest("POST", n.enviromentConfig.NotificationBaseUrl+"webhook/whatsapp/"+finalPath, bytes.NewBuffer(jsonBody))
	if err != nil {
		return n.error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(n.enviromentConfig.NotificationUsername, n.enviromentConfig.NotificationPassword)

	resp, err := n.whatsappClient.Do(req)
	if err != nil {
		return n.error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return n.error(fmt.Errorf("error sending notification: %s", resp.Status))
	}

	return nil
}

func (n *notificationServiceImpl) SendMessage(chatID string, message string) *errors.AppError {
	body := map[string]string{
		"chatId":  chatID,
		"message": message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return n.error(err)
	}

	req, err := http.NewRequest("POST", n.enviromentConfig.NotificationBaseUrl+"webhook/whatsapp", bytes.NewBuffer(jsonBody))
	if err != nil {
		return n.error(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(n.enviromentConfig.NotificationUsername, n.enviromentConfig.NotificationPassword)

	resp, err := n.whatsappClient.Do(req)
	if err != nil {
		return n.error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return n.error(fmt.Errorf("error sending WhatsApp message: %s", resp.Status))
	}

	return nil
}

func (n *notificationServiceImpl) SendTracks(chatId string, tracks []*model.Track) *errors.AppError {
	message := "No tienes seguimientos."
	if len(tracks) > 0 {
		var trackList string
		for _, track := range tracks {
			name := track.FullName
			if track.Alias != nil {
				name = *track.Alias
			}
			trackList += fmt.Sprintf("- %s %s\n", track.Run, name)
		}
		message = fmt.Sprintf("📋 Listado:\n%s", trackList)
	}

	err := n.SendMessage(chatId, message)
	if err != nil {
		return n.error(err)
	}

	return nil
}

func (n *notificationServiceImpl) Close() error {
	// Detener el topic para que no acepte más publicaciones
	n.pubsubTopic.Stop()

	// Cerrar el cliente de Pub/Sub
	err := n.pubsubClient.Close()
	if err != nil {
		return fmt.Errorf("error al cerrar cliente de Pub/Sub: %w", err)
	}

	return nil
}

func (n *notificationServiceImpl) error(err error) *errors.AppError {
	return errors.NewAppError("NotificationService", err)
}
