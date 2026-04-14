# marshall-labs-sdk

Platform SDK for services in the marshall-labs ecosystem. Provides shared types, HTTP clients, and middleware used by downstream application services to interact with the marshall-labs platform layer (IDAM, app registry).

## Packages

### `idam`

Shared Identity and Access Management types: users, error codes, and request/response structs for registration, login, email verification, and password reset flows. Also includes `UserAuthClient` for calling the IDAM user account endpoints.

### `appreg`

App registry types: `Application`, `Feature`, `Route`, and `Service`. Used by services that need to reason about registered applications and their routes.

### `labs`

Inter-service communication utilities:

- **`UserAuthGuard`** — Gin middleware that validates the `X-Idam-User-Auth` HMAC-signed header on internal routes. Use this to protect routes that should only be accessible from authenticated platform services.
- **`UserAuthContext`** — Holds the identity context extracted from a validated auth header. Use `ApplyRequestHeaders` to attach it to outbound requests.
- **`IdamSecurityClient`** — Client for the IDAM security endpoints (`/verify-token`, `/verify-route`).
- **`IdamUserClient`** — Client for the IDAM user endpoint (`GET /api/idam/users/:usrId`). Intended for internal service-to-service calls.
- **`IdamPublicRoutesClient`** — Client to fetch the public route list from the IDAM app registry endpoint.
- Event types (`UserCreatedEvent`, `UserLoggedOutEvent`, etc.) for consuming IDAM user events from the message bus.

## Installation

```bash
go get github.com/dmars8047/marshall-labs-sdk
```

## Usage

### Protecting an internal route with `UserAuthGuard`

```go
import "github.com/dmars8047/marshall-labs-sdk/labs"

router.GET("/internal/resource", labs.UserAuthGuard(hmacKey), func(c *gin.Context) {
    ctx := labs.GetAuthContext(c)
    // ctx.UserId, ctx.ApplicationId, ctx.TokenId are verified and safe to use
})
```

The `hmacKey` must match the key used by the caller when invoking `UserAuthContext.ApplyRequestHeaders`.

### Verifying a token

```go
client := labs.NewIdamSecurityClient(&httpClient, idamBaseUrl, hmacKey)

resp, err := client.VerifyToken(bearerToken)
if err != nil {
    // handle unauthorized / service error
}

authCtx := resp.AsNewIdentityContext()
```

### Fetching a user

```go
userClient := labs.NewIdamUserClient(&httpClient, idamBaseUrl)

user, err := userClient.GetUser(userId)
```

## License

MIT — see [LICENSE](LICENSE).
