package internal

import (
	"encoding/json"
	"github.com/SherClockHolmes/webpush-go"
	"log"
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

func SendToAllSubscribers(request NotificationRequest) {
	var subscriptionsFromDb []WebPushSubscription
	Connection.Find(&subscriptionsFromDb)

	body, err := json.Marshal(request)

	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return
	}

	var notificationsSent int

	for _, subscriptionFromDb := range subscriptionsFromDb {
		notificationsSent++
		go sendNotificationToSubscription(subscriptionFromDb.Subscription, body)
	}

	log.Printf("Sent %v notification(s)", notificationsSent)
}

func sendNotificationToSubscription(subscription string, body []byte) {

	var webPushSubscription webpush.Subscription
	err := json.Unmarshal([]byte(subscription), &webPushSubscription)
	if err != nil {
		log.Printf("Error unmarshalling subscriptionFromDb: %v", err)
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
	}

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		log.Printf("Error sending notification, status code: %v, response: %v", resp.StatusCode, resp)
		log.Printf("Subscription was: %v", subscription)
	}

	defer resp.Body.Close()
}
