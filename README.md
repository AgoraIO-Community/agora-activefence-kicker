# agora-activefence-kicker

This Go package provides a library for kicking users from Agora channels via the REST API. If used directly with the cmd/main.go file it will expect the same POST request that the ActiveFence backend sends at the "/kick" endpoint.

## One-Click Deployments

| Railway | Render | Heroku |
|:-:|:-:|:-:|
| [![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/piHt63?referralCode=waRWUT) | [![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/AgoraIO-Community/agora-activefence-kicker) | [![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://www.heroku.com/deploy/?template=https://github.com/AgoraIO-Community/agora-activefence-kicker) |

## Usage

To use this package in your own server, you could set the following environment variables:

- `APP_ID` - Your Agora app ID  
- `CUSTOMER_KEY` - Your Agora customer key
- `CUSTOMER_SECRET` - Your Agora customer secret

You can optionally add these in a `.env` file at the root directory of the project. 

Then import the kickService and start serving with default settings:

```go
import kickService "github.com/AgoraIO-Community/agora-activefence-kicker/service"

s := kickService.NewService()
go s.Stop()
s.Start()
```

> This will add the endpoint `/kick` and expect the default post request shown below. It will kick a user for 5 minutes.

Alternatively, if you want to send different types of POST requests, you can call KickUser directly:

```go
kickService.KickUser(
  "app-id",
  "channel",
  999, // user ID
  300, // Seconds to ban user for
  "rest-token"
)
```

This will call the Agora API to ban the user from joining the specified channel for the duration. 

## Running 

To run the example server with default implementation:

```
go run cmd/main.go
```

This will start a basic web server that exposes a `/kick` endpoint for kicking users.

## ActiveFence Webhooks

An example of the POST data from ActiveFence includes:

```json
{
  "contentUrl": "s3://activefence-content/temp/rehosted/agora-demo/agore-demo/ba5be0b6-6867-411b-b9c2-4d4e01b62b66",
  "status": "61c97b8b9abfed64af087e22",
  "userId": "556544384",
  "addedAt": "2023-07-11T17:46:47.000Z",
  "contentId": "D73F290CE6EC4E038E84E2BE52687109",
  "reason": [
    "Unauthorised Sales"
  ],
  "metadata": "{\"callbackData\":\"\",\"cname\":\"test3\",\"requestId\":\"D73F290CE6EC4E038E84E2BE52687109\",\"sid\":\"7921EE8B14674DA48B1C690BAB3758B8\",\"source\":\"agora\",\"timestamp\":20230711174617885,\"uid\":556544384}"
}
```

The minimum required from ActiveFence would be userId, and metadata.

## License

This project is licensed under the MIT license. See [LICENSE](/LICENSE) for details.