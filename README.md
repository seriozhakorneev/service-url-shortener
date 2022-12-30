# service-url-shortener
URL Shortener Service with clean architecture based on
[Go-clean-template](https://github.com/evrone/go-clean-template).

The service can create shortened URLs, return original URLs (using GRPC).
Redirect to the original URL when follow the short one (using HTTP).


- Shortener usecase, Digitiser usecase, Mocks & Unit-Tests & Fuzz Testing.
- Shortener GRPC service, Redirect HTTP & Integration-Tests.
- Postgres repository, Migrations.
- GRPC server, request Logger.

### Local development:
```sh
# Postgres, Redis
$ make compose-up

# Run app with migrations
$ make run
```

Integration tests (can be run in CI):
```sh
# DB, app + migrations, integration tests
$ make compose-up-integration-test
```

### Usage:
GRPC endpoints described in [proto](https://github.com/seriozhakorneev/service-url-shortener/blob/main/internal/entrypoint/grpc/shortener_proto/shortener.proto).

HTTP redirect can be used by following the short URL provided by GRPC Shortener.Create response.
