shadcn-svelte

This gives you:

- Button
- Card
- Dialog
- Input
- Form
- Dropdown
- Sheet
- Avatar
- Tabs

You will probably use:

```
Button
Card
Input
Textarea
Dialog
DropdownMenu
Avatar
Carousel
Form
Skeleton
```

---

https://en.wikipedia.org/wiki/JSON_Web_Token#Use

So your frontend needs:

```
login
    |
    POST /login
    |
    receive JWT
    |
    store token
```

I would avoid localStorage if possible.

Prefer:

```
HTTP-only cookie
```

because your backend already controls authentication.

Flow:

```
Browser

   |
   |
login request

   |

Go backend

   |
sets cookie

   |

Browser automatically sends cookie
```

Then your SvelteKit app doesn't manually manage tokens.

---

I'd organize the homepage like this:

```
Navbar

Hero

Search

Popular recipes

Newest recipes

Most rated recipes

Footer
```

---

A common stack in SvelteKit is:

Zod for defining the schema.
sveltekit-superforms if you want rich form handling, validation, and error display.

With a Zod schema you can share validation between client and server if you later move more logic into SvelteKit actions. Even if you continue posting to your Go REST API, Zod is still useful for client-side validation.

---

https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS
https://en.wikipedia.org/wiki/Cross-origin_resource_sharing
Hexagonal Architecture (Ports and Adapters)
