package internal

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"sync/atomic"

	"github.com/SherClockHolmes/webpush-go"
)

func SaveSubscription(subscriptionStr []byte) error {
	if Connection == nil {
		log.Fatalf("Database connection not initialized")
	}

	var subscriptionFromRequest webpush.Subscription
	err := json.Unmarshal(subscriptionStr, &subscriptionFromRequest)

	if err != nil {
		log.Printf("Can't unmarshall subscription string into webpush subscription type: %v", err)
		return nil
	}

	dbSubscription := WebPushSubscription{
		Endpoint:     subscriptionFromRequest.Endpoint,
		Subscription: string(subscriptionStr),
	}

	result := Connection.Create(&dbSubscription)

	if result.Error != nil {
		log.Printf("Error saving subscription: %v", result.Error)
		return result.Error
	}
	return nil
}

func RemoveSubscription(unsubscriptionRequest WebPushUnsubscriptionRequest) error {
	if Connection == nil {
		log.Fatalf("Database connection not initialized")
	}

	// Delete the subscription
	result := Connection.Where("endpoint = ?", unsubscriptionRequest.Endpoint).Delete(&WebPushSubscription{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SendToAllSubscribers(request NotificationRequest) error {
	var subscriptionsFromDb []WebPushSubscription
	Connection.Find(&subscriptionsFromDb)

	body, err := json.Marshal(request)

	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return err
	}

	var notificationsSent int
	var waitingGroup sync.WaitGroup
	var errorOccurred atomic.Bool

	for _, subscriptionFromDb := range subscriptionsFromDb {
		notificationsSent++
		waitingGroup.Add(1)
		go sendNotificationToSubscription(subscriptionFromDb.Subscription, body, &waitingGroup, &errorOccurred)
	}

	waitingGroup.Wait()
	if errorOccurred.Load() {
		return errors.New("error sending notifications")
	}

	return nil
}

func sendNotificationToSubscription(subscription string, body []byte, waitingGroup *sync.WaitGroup, errorOccurred *atomic.Bool) {
	defer waitingGroup.Done()

	var webPushSubscription webpush.Subscription
	err := json.Unmarshal([]byte(subscription), &webPushSubscription)
	if err != nil {
		log.Printf("Error unmarshalling subscriptionFromDb: %v", err)
		errorOccurred.Store(true)
		return
	}

	resp, err := webpush.SendNotification(body, &webPushSubscription, &webpush.Options{
		Subscriber:      "notify@jskweb.de",
		VAPIDPublicKey:  VapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30,
	})

	if err != nil {
		log.Printf("Error sending notification: %v", err)
		errorOccurred.Store(true)
		return
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		if resp.StatusCode == 410 { // endpoint does not exist anymore
			unsubscriptionRequest := WebPushUnsubscriptionRequest{
				Endpoint: webPushSubscription.Endpoint,
			}

			err = RemoveSubscription(unsubscriptionRequest)
			if err != nil {
				log.Printf("Error deleting gone subscription: %v", err)
				errorOccurred.Store(true)
				return
			}

			log.Printf("Successfully removed subscription after receiving 410 status code: %v", webPushSubscription.Endpoint)
			return
		}

		log.Printf("Error sending notification, status code: %v, response: %v", resp.StatusCode, resp)
		log.Printf("Subscription was: %v", subscription)
		errorOccurred.Store(true)
	}

	defer resp.Body.Close()
}
