# service-url-shortener
Fault tolerant, Self-sufficient, Customizable, Сache accelerated, TTL supported.<br/>
Service URL Shortener with clean architecture based on
[Go-clean-template](https://github.com/evrone/go-clean-template).

Service can create shortened URLs, return original URLs (using GRPC).Redirect to the original URL when follow the short one (using HTTP).

- Shortener usecase, Digitiser usecase, Mocks & Unit-Tests & Fuzz Testing.
- Shortener GRPC routes, Redirect HTTP & Integration-Tests.
- Postgres repository, Migrations, Redis Cache & Integration-Tests
- GRPC server, GRPC request Logger, Redis client.

### Local development:
```sh
# Postgres, Redis
$ make compose-up

# Run app with migrations
$ make run
```

Integration tests:
```sh
# DB, app + migrations, integration tests
$ make compose-up-integration-test
```

### Usage:

Visit [Config](https://github.com/seriozhakorneev/service-url-shortener/blob/main/config/config.yml) before use.

GRPC endpoints described in [Proto](https://github.com/seriozhakorneev/service-url-shortener/blob/main/internal/entrypoint/grpc/shortener_proto/shortener.proto).

HTTP redirect can be used by following the short URL provided by GRPC Shortener.Create response.
