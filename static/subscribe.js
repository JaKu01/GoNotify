let subscribeButton = document.getElementById('subscribeButton');
let unsubscribeButton = document.getElementById('unsubscribeButton');
let successMessage = document.getElementById('successMessage');
let errorMessage = document.getElementById('errorMessage');

window.addEventListener('load', async () => {
    if (!('serviceWorker' in navigator)) {
        console.error('Service Worker not supported');
        showErrorMessage('Service Worker not supported');
        return;
    }
    // Register the service worker
    let registration = await registerServiceWorker();
    if (!registration) {
        console.error('Service Worker registration failed.');
        showErrorMessage('Service Worker registration failed');
        return;
    }

    registerEventListeners();

    // Check if the user is already subscribed

    if(await registration.pushManager.getSubscription()) {
        console.log('User is already subscribed to push notifications');
        subscribeButton.disabled = true;
        unsubscribeButton.disabled = false;
        return;
    }

    console.log('User is not subscribed to push notifications');
    subscribeButton.disabled = false;
    unsubscribeButton.disabled = true;
});

async function registerServiceWorker() {
    try {
        let registration = await navigator.serviceWorker
            .register('./static/service-worker.js', {scope: './static/'});
        console.log('Service Worker Registered');
        return registration;
    } catch (err) {
        console.error("Failed to register Service Worker", err);
    }
    return null;
}

function registerEventListeners() {
    subscribeButton.addEventListener('click', async () => {
        await subscribeUserToPush();
    });

    unsubscribeButton.addEventListener('click', async () => {
        await unsubscribeUserFromPush();
    });
}

// Function to handle the push subscription
async function subscribeUserToPush() {
    // Register for push notifications when the subscribe button is clicked
    let registration = await navigator.serviceWorker.getRegistration('/static/');
    if (!registration) {
        console.error('Service Worker not registered');
        showErrorMessage('Service Worker not registered');
        return;
    }
    console.log('Registration found')

    let permission = await Notification.requestPermission();
    if (permission !== 'granted') {
        console.error('Permission not granted for Notification');
        showErrorMessage('Permission not granted for Notification');
        return;
    }

    const convertedVapidKey = urlBase64ToUint8Array(key);
    const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: convertedVapidKey
    });

    if (!subscription) {
        console.error('Failed to subscribe to push notifications');
        showErrorMessage('Failed to subscribe to push notifications');
        return;
    }
    console.log('User subscribed to push notifications:', subscription);

    if (!await sendSubscriptionToServer(subscription)) {
        console.error('Failed to send subscription to server');
        showErrorMessage('Failed to send subscription to server');
        return;
    }
    console.log('Sent subscription to Server');
    subscribeButton.disabled = true;
    unsubscribeButton.disabled = false;
    showSuccessMessage('Successfully subscribed to push notifications!');
}


// Function to send the subscription to the server
async function sendSubscriptionToServer(subscription) {
    try {
        return await fetch('/api/subscribe', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(subscription)
        });
    } catch (err) {
        console.error('Failed to send subscription to server', err);
    }
    return null;
}

async function unsubscribeUserFromPush() {
    let registration = await navigator.serviceWorker.getRegistration();
    if (!registration) {
        console.error('Service Worker not registered');
        showErrorMessage('Service Worker not registered');
        return;
    }
    console.log('Registration found')

    let subscription = await registration.pushManager.getSubscription();
    if (!subscription) {
        console.error('User is not subscribed to push notifications');
        showErrorMessage('User is not subscribed to push notifications');
        return;
    }

    if (!await subscription.unsubscribe()) {
        console.error('Failed to unsubscribe from push notifications');
        showErrorMessage('Failed to unsubscribe from push notifications');
        return;
    }
    console.log('User unsubscribed from push notifications');

    if (!await sendUnsubscribeRequestToServer(subscription.endpoint)) {
        console.error('Failed to send unsubscription request to server');
        showErrorMessage('Failed to send unsubscription request to server');
        return;
    }
    console.log('Sent unsubscription request to server');
    subscribeButton.disabled = false;
    unsubscribeButton.disabled = true;
    showSuccessMessage('Successfully unsubscribed to push notifications!');
}

async function sendUnsubscribeRequestToServer(endpoint) {
    try {
        return await fetch('/api/subscribe', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                endpoint: endpoint
            })
        });
    } catch (err) {
        console.error('Failed to send unsubscription request to server', err);
    }
    return null;
}

function showErrorMessage(message) {
    errorMessage.textContent = message + ' 🥺';
    errorMessage.style.display = 'block';
    successMessage.style.display = 'none';
}

function showSuccessMessage(message) {
    successMessage.textContent = message + ' 😎';
    successMessage.style.display = 'block';
    errorMessage.style.display = 'none';
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