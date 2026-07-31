Your backend structure is already quite feature-oriented, which is good news because your frontend can mirror that. Since you already have:

- Go API ✅
- OpenAPI spec ✅
- generated TypeScript schema ✅
- SvelteKit initialized ✅

you are at the point where you should **design the frontend architecture before writing pages**.

For your app (a recipe platform), I would not start by making a dashboard or authentication. Start with the **public browsing experience** because it will force you to learn the core SvelteKit concepts without adding auth complexity.

A good development order:

```
1. Frontend foundation
2. Public recipe browsing
3. Recipe details
4. Authentication
5. User profile
6. Recipe management dashboard
7. Polish
```

---

# Step 1 — Decide your frontend architecture

Your backend features:

```
users
recipes
comments
images
auth
```

Your frontend features should mirror that:

```
src
├── lib
│   ├── api
│   │   ├── client.ts
│   │   ├── schema.ts
│   │   └── errors.ts
│   │
│   ├── components
│   │   ├── Navbar.svelte
│   │   ├── RecipeCard.svelte
│   │   ├── RecipeCarousel.svelte
│   │   ├── CommentList.svelte
│   │   └── SearchBar.svelte
│   │
│   ├── features
│   │   ├── recipes
│   │   │   ├── api.ts
│   │   │   └── types.ts
│   │   │
│   │   ├── auth
│   │   │   ├── api.ts
│   │   │   └── store.ts
│   │   │
│   │   └── users
│   │       └── api.ts
│   │
│   └── stores
│
└── routes
    ├── +layout.svelte
    ├── +page.svelte
    │
    ├── recipes
    │   └── [id]
    │       └── +page.svelte
    │
    ├── login
    │   └── +page.svelte
    │
    ├── register
    │   └── +page.svelte
    │
    ├── dashboard
    │   └── +page.svelte
    │
    └── profile
        └── edit
            └── +page.svelte
```

The important idea:

- `routes` = pages
- `components` = reusable UI
- `features` = business logic

---

# Step 2 — Add Tailwind + component library

I would install:

## Tailwind

```bash
npm install tailwindcss @tailwindcss/vite
```

Configure it.

Then add:

## shadcn-svelte

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

For a recipe app it fits very well.

---

# Step 3 — Create your API layer

You already have:

```
src/lib/api/schema.ts
```

Now create wrappers.

Example:

```
src/lib/features/recipes/api.ts
```

Something like:

```ts
import { api } from "$lib/api/client";

export async function getRecipes() {
  const { data, error } = await api.GET("/recipes");

  if (error) {
    throw error;
  }

  return data;
}

export async function getRecipe(id: string) {
  const { data, error } = await api.GET("/recipes/{id}", {
    params: {
      path: {
        id,
      },
    },
  });

  if (error) {
    throw error;
  }

  return data;
}
```

Your pages should never directly call:

```ts
api.GET(...)
```

Keep HTTP details away from UI.

---

# Step 4 — Build the global layout

First file:

```
src/routes/+layout.svelte
```

This is your application shell.

It should contain:

```
--------------------------------
Navbar
--------------------------------

Page content

--------------------------------
Footer
--------------------------------
```

Example structure:

```svelte
<script>
 import Navbar from '$lib/components/Navbar.svelte';
</script>

<Navbar />

<main>
    {@render children()}
</main>
```

Every page gets inserted there.

---

# Step 5 — Build the homepage first

Your homepage:

```
/
```

should become:

```
Navbar

Hero section

Search bar

Recipe cards grid


RecipeCard
RecipeCard
RecipeCard
```

Your first real Svelte component:

```
src/lib/components/RecipeCard.svelte
```

Input:

```ts
let { recipe } = $props();
```

Output:

```
+----------------+
| image          |
|                |
| Recipe name    |
| author         |
| tags           |
+----------------+
```

Click:

```
/recipes/123
```

---

# Step 6 — Recipe page

Create:

```
routes/recipes/[id]/+page.svelte
```

The URL:

```
/recipes/42
```

gives you:

```ts
params.id;
```

The page contains:

```
Recipe title

Author

Image carousel

Ingredients

Instructions

Comments
```

Components:

```
RecipeHeader.svelte

RecipeCarousel.svelte

IngredientsList.svelte

Instructions.svelte

CommentList.svelte

CommentForm.svelte
```

---

# Step 7 — Add authentication

Only after browsing works.

You already have:

```
internal/app/adapters/jwt.go
middleware/auth.go
```

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

# Step 8 — Dashboard

After login:

```
/dashboard
```

Show:

```
My recipes

+ Create recipe


Recipe A   Edit Delete

Recipe B   Edit Delete
```

Components:

```
DashboardCard

RecipeEditor
```

---

# Step 9 — Recipe editor

This will be your biggest page.

Split it:

```
RecipeForm.svelte

TitleInput

DescriptionInput

IngredientsEditor

InstructionsEditor

ImageUploader
```

Do not create one giant:

```
EditRecipe.svelte (800 lines)
```

---

# Step 10 — Add state only when needed

Svelte beginners often add stores too early.

You probably need only:

```
auth store
```

Example:

```
src/lib/features/auth/store.ts
```

Maybe:

```ts
export const user = writable(null);
```

Do not create stores for recipes, comments, etc. Initially.

SvelteKit already handles page state well.

---

# Suggested milestone plan

## Milestone 1: "I can browse recipes"

Finish:

✅ navbar
✅ homepage
✅ search bar UI
✅ recipe cards
✅ recipe details page
✅ images

No login yet.

---

## Milestone 2: "Users exist"

Finish:

✅ register
✅ login
✅ logout
✅ profile page

---

## Milestone 3: "Users create content"

Finish:

✅ dashboard
✅ create recipe
✅ edit recipe
✅ delete recipe

---

## Milestone 4: "Community"

Finish:

✅ comments
✅ likes/favorites (if desired)
✅ ratings

---

# One thing I would check in your OpenAPI spec before coding

Your frontend experience depends heavily on these API shapes:

### Recipes

Do you have:

```
GET /recipes
```

returning:

```json
[
  {
    "id": 1,
    "title": "Banitsa",
    "thumbnail": "..."
  }
]
```

or does it return the whole recipe including:

- instructions
- images
- comments?

Ideally:

```
GET /recipes
```

returns lightweight cards.

```
GET /recipes/{id}
```

returns the full detail.

---

### Images

Since you have:

```
images.go
recipe_images.go
```

you probably want:

```
GET /recipes/{id}/images
```

returning:

```json
[
  {
    "url": "..."
  }
]
```

so the carousel is easy.

---

### Authentication

Your OpenAPI should define:

```
POST /auth/login

POST /auth/register

GET /users/me
```

because your dashboard will need "who am I?"

---

If I were building this project from your current state, my **next three concrete actions** would be:

1. Install Tailwind + shadcn-svelte.
2. Create `Navbar.svelte` and the global layout.
3. Implement `/` homepage calling `GET /recipes` and rendering `RecipeCard`.

That will give you your first complete frontend-backend integration and teach you the SvelteKit workflow.

---

---

Awesome. I think this is where the project starts becoming fun.

One thing I want to avoid is the common beginner trap of jumping straight into pages with lots of HTML. Instead, we'll build a **small design system** for _your_ app. After that, creating pages is mostly composing components.

---

# Roadmap

Here's the roadmap I'd follow from this point:

```
Foundation ✅
├── SvelteKit
├── Tailwind
├── shadcn
├── OpenAPI
└── API client

↓

Application Shell

↓

Homepage (fake data)

↓

Connect homepage to backend

↓

Recipe page

↓

Authentication

↓

Dashboard

↓

Recipe editor
```

---

# Today's goal

Let's build something that already feels like a recipe website.

We'll end up with this:

```
+-------------------------------------------------------------+
| 🍲 RecipeHub        Search recipes...      Login Register    |
+-------------------------------------------------------------+

            Discover Bulgarian Recipes

Find traditional and modern recipes

+--------------------------------------------+
| Search...                                  |
+--------------------------------------------+

Recipes

┌──────────────┐
│   image      │
│              │
├──────────────┤
│ Banitsa      │
│ 45 min       │
│ Bulgarian    │
└──────────────┘
```

Notice something:

There is **no backend yet.**

We're only learning Svelte.

---

# Step 1

Create folders.

Inside

```
src/lib/components
```

make

```
components
│
├── layout
│
├── recipe
│
└── shared
```

Eventually it'll become

```
components
│
├── layout
│   ├── Navbar.svelte
│   └── Footer.svelte
│
├── recipe
│   ├── RecipeCard.svelte
│   ├── RecipeGrid.svelte
│   ├── RecipeCarousel.svelte
│   └── RecipeSearch.svelte
│
└── shared
    ├── PageTitle.svelte
    ├── EmptyState.svelte
    └── LoadingSpinner.svelte
```

Much cleaner than throwing everything into one folder.

---

# Step 2

Let's create our **first real component**.

```
src/lib/components/layout/Navbar.svelte
```

We'll intentionally keep it tiny.

```svelte
<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { Input } from "$lib/components/ui/input";
</script>

<nav class="border-b bg-background">
	<div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
		<h1 class="text-xl font-bold">🍲 RecipeHub</h1>

		<div class="w-96">
			<Input placeholder="Search recipes..." />
		</div>

		<div class="flex gap-2">
			<Button variant="ghost">Login</Button>
			<Button>Register</Button>
		</div>
	</div>
</nav>
```

That's your first reusable component.

---

# Step 3

Open

```
src/routes/+layout.svelte
```

Instead of only rendering children:

```svelte
<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
```

Import the navbar.

```svelte
<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";

	import Navbar from "$lib/components/layout/Navbar.svelte";

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Navbar />

{@render children()}
```

Immediately your whole application gets a navbar.

That's one of the nicest things about SvelteKit.

---

# Step 4

Replace

```
src/routes/+page.svelte
```

with something tiny.

```svelte
<script lang="ts">
</script>

<div class="mx-auto max-w-7xl px-6 py-16">
	<h1 class="text-5xl font-bold">
		Discover Bulgarian Recipes
	</h1>

	<p class="mt-4 max-w-xl text-muted-foreground">
		Traditional recipes shared by the community.
	</p>
</div>
```

Now your website already has:

- navigation
- title
- search box
- buttons

---

# Why we're doing it this way

Notice we haven't touched:

- stores
- authentication
- OpenAPI
- fetch
- routing

Those are all **harder concepts**.

Instead, we're learning the basics:

> A page is composed of components.

That's the core of frontend development.

---

# The next component

Once you're happy with the navbar, the next thing we'll build is `RecipeCard.svelte`.

It's the perfect next component because it introduces almost every important Svelte concept you'll use:

- receiving data through props
- displaying images and text
- using shadcn's `Card`
- navigating to another page
- responsive layouts

After we have a `RecipeCard`, building the homepage becomes as simple as rendering a grid of cards. Once that works with fake data, swapping in data from your Go backend will be a small change rather than a major rewrite.

So for now, create the `layout`, `recipe`, and `shared` component folders, add the `Navbar.svelte`, wire it into `+layout.svelte`, and replace the home page with the simple hero section. Run the app and see how it feels. If anything looks off—or if you run into import or styling issues—show me the code or the error, and we'll fix it before moving on to the recipe cards.

---

---

I actually like those ideas a lot. The one thing I'd change is the **order of information** on the card.

When people browse recipes, they usually decide in this order:

1. **Does it look delicious?** (image)
2. **What is it?** (title)
3. **Is it worth making?** (rating)
4. **Can I make it?** (time, difficulty)
5. **Is it enough food?** (servings)

So I'd let the image do more of the work.

## I'd structure the card like this

```
╭──────────────────────────────╮
│                              │
│         Recipe photo         │
│                              │
│                    ⭐ 4.8     │
├──────────────────────────────┤
│ Banitsa                      │
│ Traditional Bulgarian pastry │
│                              │
│ ⏱ 45 min                     │
│ 👨‍🍳 Medium                   │
│ 🍽 6 servings                │
╰──────────────────────────────╯
```

A few reasons:

- The image immediately catches attention.
- The title is the most important text.
- The metadata is easy to scan.
- A short subtitle or description makes different recipes distinguishable.

---

## Should metadata be icons or text?

I'd definitely use icons.

Something like:

```
⏱ 45 min

🔥 Medium

🍽 6
```

is much easier to scan than

```
Preparation time: 45 minutes

Difficulty: Medium

Servings: 6
```

Lucide already includes excellent icons for this.

---

## Rating

I'd avoid showing five filled stars.

Instead:

```
⭐ 4.8 (124)
```

or

```
★★★★★ 4.8
```

The first takes less space.

---

## Difficulty

Instead of colors only:

```
Easy
Medium
Hard
```

I'd add a badge.

Example:

```
🟢 Easy

🟠 Medium

🔴 Hard
```

or a shadcn Badge component.

---

## Grid layout

I agree with three cards on desktop.

I'd make it responsive:

```
Desktop
──────────────

□ □ □

□ □ □

Tablet
──────────────

□ □

□ □

□ □

Phone
──────────────

□

□

□
```

Tailwind makes this almost trivial:

```html
<div class="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-3"></div>
```

---

## Search

I'd actually separate the search from the navbar.

Navbar:

```
Logo

Browse

Dashboard

Login
```

Homepage:

```
Discover Bulgarian Recipes

[ Search recipes........ ]

Categories

Popular recipes
```

Reason:

A wide search bar belongs to the page, not the navigation. It gives you room later for filters like cuisine, difficulty, and preparation time.

---

## Homepage sections

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

Initially, only implement the first "Recipes" section. The others can come later.

---

## One recommendation for your API

When we connect the backend, I'd expose **recipe cards** as lightweight data.

Instead of returning the entire recipe list with comments and instructions, I'd have your list endpoint return only what's needed for the cards, for example:

```json
{
  "id": 42,
  "title": "Banitsa",
  "thumbnailUrl": "/images/banitsa.webp",
  "averageRating": 4.8,
  "preparationMinutes": 45,
  "difficulty": "Medium",
  "servings": 6,
  "author": "Ivan"
}
```

Then `GET /recipes/{id}` can return the full details including ingredients, instructions, comments, and the full image gallery. That keeps the homepage fast and avoids downloading data the user doesn't need yet.

## What I'd build next

I think the next milestone should be **one polished `RecipeCard.svelte`**.

Not a grid.

Not fake data.

Just **one reusable card** that looks exactly how you want.

Once we're happy with that single component, putting nine of them into a responsive grid is only a few lines of code. More importantly, you'll have a reusable building block that you'll keep using throughout the project. I find this "perfect one component, then repeat it" approach leads to cleaner code and makes the UI feel more consistent.
