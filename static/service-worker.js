self.addEventListener('push', event => {
    const webPushData = event.data.json()
    const options = {
        body: webPushData.body,
    };
    event.waitUntil(self.registration.showNotification(webPushData.subject, options));
});