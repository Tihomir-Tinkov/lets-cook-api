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

---

**Backend mindset:**

```
Request → Handler → Database → Response
```

**Frontend mindset:**

```
User action
    ↓
Component state changes
    ↓
UI updates automatically
    ↓
Maybe call API
    ↓
Update state
    ↓
UI updates again
```

The frontend is a long-running program in the browser. It reacts to events.

For example:

- User clicks "favorite recipe"
- Svelte runs a function
- Function calls your Go API
- Response comes back
- Svelte updates the button

No page refresh required.

---

## What would look good in a portfolio?

Beyond the choice of framework, employers often notice the quality of the implementation. Consider including:

Responsive design (mobile and desktop)
Dark/light mode
Search and filtering
Pagination or infinite scrolling
Optimistic UI updates for favorites
Loading skeletons
Error handling and empty states
Authentication (JWT or session-based)
Clean component architecture
Accessible forms and keyboard navigation

Those features demonstrate practical frontend skills regardless of the framework.

## Current trend (2025–2026)

The ecosystem has shifted toward:

1. **Headless components**
   - Radix UI
   - Headless UI
   - Ark UI

2. **Component ownership instead of dependency-heavy kits**
   - shadcn/ui-style copy/paste components

3. **Tailwind-based design systems**
   - Tailwind + component primitives has become a dominant pattern

4. **Web Components for cross-framework reuse**
   - Especially in large organizations maintaining multiple frontend stacks

---

---

## 3. Introduce loading and error states

This is where projects begin to feel polished.

Every API call should have

- loading UI
- empty state
- error state

Example:

```
Recipes

Loading...

No recipes found.

Unable to load recipes.
```

Instead of assuming data always exists.

---

## 9. Authentication improvements

Currently you have login/logout.

Eventually you'll want

```
hooks.server.ts

↓

read JWT

↓

load user

↓

protect routes
```

For example

```
/dashboard

if not logged in

↓

redirect("/login")
```

instead of checking inside components.

---

## 10. Route protection

Dashboard and edit pages should require authentication.

```
login

↓

dashboard

↓

edit recipe
```

Anonymous users shouldn't even reach those pages.

---

## 11. Notifications

Users appreciate feedback.

Instead of

```
click Save

nothing happens
```

show

```
Recipe saved.

Review posted.

Unable to save recipe.
```

A toast system goes a long way.

---

## 15. Accessibility

Check

- keyboard navigation
- labels
- focus states
- alt text
- semantic HTML

Luma already helps here.

---

## 16. Responsive layout

Test only after major features exist.

Use

- phone
- tablet
- desktop

Don't leave this until the end.

---

# Svelte-specific improvements

Since you're learning Svelte 5, I'd also spend some time adopting its idioms rather than writing React-style code.

### Use `+page.ts` for data

Avoid fetching inside components when the data belongs to a route.

Good:

```
+page.ts

↓

returns data

↓

+page.svelte
```

instead of

```
onMount()

↓

fetch()
```

---

### Keep components mostly presentational

For example

```
RecipeCard
```

should receive

```ts
recipe;
```

not fetch its own data.

---

### Keep business logic inside `features/`

Your structure already looks good.

```
features/

recipe/

review/

auth/
```

I would continue that approach.

---

### Avoid giant stores

A common beginner mistake is putting everything in stores.

Stores should mostly be for

- current user
- theme
- notifications

Recipes should come from route loads.
