self.addEventListener('push', event => {
    const webPushData = event.data.json()
    const options = {
        body: webPushData.body,
        icon: "https://notify.jskweb.de/apple-icon.png",
        image: "https://notify.jskweb.de/apple-icon.png"
    };
    event.waitUntil(self.registration.showNotification(webPushData.subject, options));
});