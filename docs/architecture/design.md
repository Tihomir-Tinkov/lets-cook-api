# Design Considerations

## 1.Wireframe

- Figma

## 2.Git

[License](https://choosealicense.com/), [Conventional Commits](https://www.conventionalcommits.org/), [Semantic Versioning](https://semver.org/)?

## 3.REST API

- [OpenAPI](https://learn.openapis.org/)
- [editor](https://editor.swagger.io/)
- [visualization](https://github.com/swagger-api/swagger-ui)

## 4.Auth

- [Password](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
  - [Argon2id](https://www.boot.dev/lessons/ccca05e8-fde6-4e02-b99f-5353069baa2b)
  - PHC string format
- [Session](https://learn.openapis.org/specification/security.html)
  - [JWT](https://www.jwt.io/introduction)
  - [HttpOnly cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Cookies#security)
  - secret key (HMAC) or [public key](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_Cheat_Sheet.html#public-key-signatures)

## 5.DB

- [ER model](https://en.wikipedia.org/wiki/Entity%E2%80%93relationship_model)
- PostgreSQL
  - Docker
  - pgAdmin 4
  - [UUID](https://www.boot.dev/blog/backend/what-are-uuids-and-should-you-use-them)

![database schema](./DB%20schema.jpg)

## 6.Backend

Find a style guide? Check out [secret mode](https://antonz.org/accepted/runtime-secret/)?

- Go 1.22+
  - net/http
  - [goose](https://github.com/pressly/goose)
  - pgx

Schema -> Migration -> Model -> Port -> Repository -> Service -> Controller -> Route

## 7.Frontend

https://openapi-ts.dev/introduction

- TypeScript
- Svelte/Kit
  - https://flowbite-svelte.com/
  - https://melt-ui.com/
  - https://sveltematerialui.com/
  - https://www.bits-ui.com/
  - https://www.shadcn-svelte.com/
  - https://www.skeleton.dev/
- (React)
  - https://ant.design/
  - https://chakra-ui.com/
  - https://mantine.dev/
  - https://mui.com/material-ui/
  - https://ui.shadcn.com/
  - https://www.radix-ui.com/
- framework-agnostic
  - https://daisyui.com/
  - https://lucide.dev/icons/
- Tailwind CSS


At this point you have:

✅ SvelteKit
✅ TypeScript
✅ Tailwind CSS v4
✅ shadcn-svelte
✅ Luma design preset
✅ OpenAPI-generated types

That's a solid modern frontend stack.

Next things I'd do:

Add a RecipeGrid.svelte component and move the {#each} there.
Add a homepage structure:
hero/title
search bar
recipe grid
Replace the emoji metadata with Lucide icons and shadcn Badge.

But first, let's get the current cards rendering. Let me know if the three cards appear.