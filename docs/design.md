# Design Considerations

## 1.Wireframe

- Figma

## 2.REST API

- [OpenAPI](https://learn.openapis.org/)
- [editor](https://editor.swagger.io/)
- [visualization](https://github.com/swagger-api/swagger-ui)

## 3.Auth

- [Password](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
  - Argon2id
  - PHC string format
- [Session](https://learn.openapis.org/specification/security.html)
  - [JWT](https://www.jwt.io/introduction)
  - [HttpOnly cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Cookies#security)
  - secret key (HMAC) or [public key](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_Cheat_Sheet.html#public-key-signatures)

## 3.DB

- [ER model](https://en.wikipedia.org/wiki/Entity%E2%80%93relationship_model)
- PostgreSQL
  - Docker
  - pgAdmin 4
  - [UUID](https://www.boot.dev/blog/backend/what-are-uuids-and-should-you-use-them)

![database schema](./DB%20schema.jpg)

## 4.Backend

Find a style guide? [Conventional Commits](https://www.conventionalcommits.org/), [Semantic Versioning](https://semver.org/)? Check out [secret mode](https://antonz.org/accepted/runtime-secret/)?

- Go 1.22+
  - net/http
  - [goose](https://github.com/pressly/goose)
  - pgx

Schema -> Migration -> Model -> Port -> Repository -> Service -> Controller -> Route

## 5.Frontend
