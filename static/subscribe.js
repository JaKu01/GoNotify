window.addEventListener('load', () => {
    if ('serviceWorker' in navigator) {
        navigator.serviceWorker
            .register('./service-worker.js', { scope: './' })
            .then(function (registration) {
                console.log("Service Worker Registered");
                // Handle the subscribe button click
                const subscribeButton = document.getElementById('subscribeButton');
                subscribeButton.addEventListener('click', () => {
                    // Register for push notifications when the button is clicked
                    navigator.serviceWorker.getRegistration()
                        .then(registration => {
                            if (registration) {
                                console.log('Registration found')
                                subscribeUserToPush(registration);
                            } else {
                                console.error('Service Worker registration not found.');
                            }
                        });
                });

                const unsubscribeButton = document.getElementById('unsubscribeButton');
                unsubscribeButton.addEventListener('click', () => {
                    // Register for push notifications when the button is clicked
                    navigator.serviceWorker.getRegistration()
                        .then(registration => {
                            if (registration) {
                                console.log('Registration found')
                                unsubscribeUserFromPush(registration)
                            } else {
                                console.error('Service Worker registration not found.');
                            }
                        });
                });
            })
            .catch(function (err) {
                console.log("Service Worker Failed to Register", err);
            })

    } else {
        console.log('service worker not ready');
    }
});

 // Function to handle the push subscription
 function subscribeUserToPush(registration) {
    // Request permission to show push notifications
    Notification.requestPermission().then(permission => {
        if (permission === 'granted') {
            // Create a push subscription
            const convertedVapidKey = urlBase64ToUint8Array(key);
            registration.pushManager.subscribe({
                userVisibleOnly: true, // Required for push notifications
                applicationServerKey: convertedVapidKey
            })
                .then(subscription => {
                    console.log('User subscribed to push notifications:');
                    console.log(JSON.stringify(subscription));
                    // Send the subscription to the server
                    sendSubscriptionToServer(subscription);
                })
                .catch(err => {
                    console.error('Failed to subscribe user:', err);
                });
        } else {
            console.log('permission not granted');
        }
    });
}

// Function to send the subscription to the server
function sendSubscriptionToServer(subscription) {
    fetch('/api/subscribe', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(subscription)
    });
}

// Convert the VAPID public key to Uint8Array
function urlBase64ToUint8Array(base64String) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding)
        .replace(/\-/g, '+')
        .replace(/_/g, '/');

    const rawData = atob(base64);
    const outputArray = new Uint8Array(rawData.length);

    for (let i = 0; i < rawData.length; i++) {
        outputArray[i] = rawData.charCodeAt(i);
    }

    return outputArray;
}

function unsubscribeUserFromPush(registration) {
    // TODO: Remove subscription in backend first
    registration.pushManager.getSubscription()
        .then(subscription => {
            if (subscription) {
                // Unsubscribe the user from push notifications
                console.log(subscription);
                
                subscription.unsubscribe()
                    .then(() => {
                        console.log('User unsubscribed from push notifications');
                        sendUnsubscribeRequestToServer(subscription.endpoint)
                    })
                    .catch(err => {
                        console.error('Failed to unsubscribe user:', err);
                    });
            } else {
                console.log('No subscription found.');
            }
        })
        .catch(err => {
            console.error('Error fetching subscription:', err);
        });
}

function sendUnsubscribeRequestToServer(endpoint) {
    fetch('/api/subscribe', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            endpoint: endpoint
        })
    }).catch((err) => {
        console.log("Failed to send unsubscription request to server.", err);
    });
}