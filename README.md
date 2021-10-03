# OAuth 2.0 written in Go
An example of OAuth2.0 written in Golang (Main repo can be found at [Example](https://github.com/go-oauth2/oauth2/blob/master/example/README.md)).

## OAuth 2.0 Explained
OAuth 2.0 is an authorization service that enables applications to obtain limited access to a user's account. It works by providing user authentication to the service that hosts a user account and authorizing third-party applications to access that user account.

There are three kinds of roles in an oauth service.

1. Resource Owner: A user who authorizes an application to access their account which is limited to the scope of the authorization granted.

2. Client: An application that wants to access the user’s account. The user must provide the authorization.

3. Resource/Authorization Server: Hosts of the user accounts. The server verifies the identity of the user then issues access tokens to the application.

## OAuth 2.0 Flow
The flow of the OAuth service is as

1. The application (client) requests authorization from the user (Resource Owner) to access some resources.

2. If the user authorized the request, the application receives an authorization permission.

3. The application requests an access token from the server (Resource/Authorization Server) by presenting authentication of its own identity and the authorization permission.

4. If the application identity is authenticated and the authorization permission is valid, the server issues an access token to the application.

5. The application requests the resource from the server and presents the access token for authentication.

6. If the access token is valid, the server serves the resource to the application.


## How to run the application locally

### Start the Server

``` bash
>> go run server.go

Server is running at 9096 port.
OAuth client Auth endpoint to http://localhost:9096/oauth/authorize
OAuth client Token endpoint to http://localhost:9096/oauth/token
```

## Start the Client

```bash
>> go run client.go

Client is running at 9094 port.
Please open http://localhost:9094/upload-photos
```

## Login page

Enter admin for user name and password (Optional)

![login](https://github.com/StephenDsouza90/oauth-golang/blob/main/server/static/login.png)
![auth](https://github.com/StephenDsouza90/oauth-golang/blob/main/server/static/auth.png)
![token](https://github.com/StephenDsouza90/oauth-golang/blob/main/server/static/token.png)
