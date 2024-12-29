## About GoNotify

GoNotify is a notification service that provides a simple JSON REST API for sending both web push and email notifications. It is designed to be lightweight, easy to deploy, and can be run using Docker Compose for an easy setup process.

## Getting started

GoNotify can be run using Docker Compose. Make sure to provide the necessary environment variables to docker, e.g. via an `.env` file.

## Required Environment Variables

The following environment variables are required for GoNotify to function properly:

- **`SMTP_HOST`**: The hostname or IP address of the SMTP server.
- **`SMTP_PORT`**: The port used to connect to the SMTP server.
- **`EMAIL`**: The email address used as the sender for outgoing notifications.
- **`EMAIL_PASSWORD`**: The password or app-specific key for the sender's email account.
- **`NETWORK_NAME`**: The name of the docker network the container should run in.

Once your environment variables are set, you can start the service using:

```bash
docker-compose up -d
```

With this setup, GoNotify will be up and running, ready to send web push and email notifications via its REST API.


## VAPID Keys

VAPID keys are required for sending Web Push notifications. By default, GoNotify automatically generates a pair of VAPID keys on the first run. These keys are stored in the `keys` directory as:

- **`keys/private.key`** — The private key used to sign push requests.
- **`keys/public.key`** — The public key shared with push services.

### Using Custom VAPID Keys

If you prefer to use your own VAPID keys, simply replace the `private.key` and `public.key` files in the `keys` directory with your own keys, ensuring they have the same filenames. GoNotify will then use your custom keys for all Web Push notifications.


## API Endpoints

GoNotify provides two main API endpoints for sending notifications: **Web Push** and **Email**. Below is a guide on how to use each endpoint.


### 📢 **Web Push Endpoint** (`/api/webpush`)

#### **How it Works**
GoNotify serves a simple webpage at the root path (`/`) where users can subscribe to Web Push notifications. This page allows users to register their devices to receive push notifications.

#### **Sending a Web Push Notification**
To send a Web Push notification to subscribed users, make a `POST` request to the following endpoint:  

```
POST /api/webpush
```

#### **Request Body**
The body of the request must be in JSON format and match the following structure:

```json
{
  "subject": "Your notification title",
  "body": "The message content you want to send"
}
```

#### **Field Descriptions**
- **`subject`**: The title of the notification.  
- **`body`**: The message content displayed in the push notification.  

#### **Example Request (cURL)**
```bash
curl -X POST http://localhost:8080/api/webpush \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "New Update Available!",
    "body": "Check out the latest features we have added."
  }'
```

---

### ✉️ **Mail Endpoint** (`/api/mail`)

Mail notifications are sent **from the provided email address to itself**. This is designed to notify **yourself** about important events or updates.  

If you want multiple people to receive these notifications, it is recommended at the moment to create a dedicated email account specifically for this purpose and share access with the intended recipients.  

🚀 **Planned Update:** In future versions, GoNotify will support sending notifications to multiple email addresses directly, making it even easier to keep a group informed.

#### **Sending an Email Notification**
To send an email, make a `POST` request to the following endpoint:  

```
POST /api/mail
```

#### **Request Body**
The request body must be in JSON format and follow this structure:

```json
{
  "subject": "The subject of the email",
  "content_type": "Content Type of the mail",
  "body": "The content of the email"
}
```

#### **Field Descriptions**
- **`subject`**: The subject line of the email.  
- **`content_type`**: The format of the email body, e.g. `text/plain` or `text/html`.  
- **`body`**: The main content of the email message.  

#### **Example Request (cURL)**
```bash
curl -X POST http://localhost:8080/api/mail \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "Welcome to GoNotify!",
    "content_type": "text/html",
    "body": "<h1>Welcome!</h1><p>We are excited to have you on board.</p>"
  }'
```

---

These endpoints make it easy to send both Web Push and Email notifications directly via simple JSON requests. The provided examples should help you get started quickly.


## Gopher Icon 
This project uses a Gopher icon as its logo from the lovely [Free Gophers Pack](https://github.com/MariaLetta/free-gophers-pack).


## Test Status
![Tests](https://github.com/JaKu01/GoNotify/actions/workflows/go-test.yml/badge.svg)
