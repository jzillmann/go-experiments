# Go Experiments

Go server with various tech setups:
1. Svelte
  - BE: [huma](https://huma.rocks) with [chi](https://github.com/go-chi/chi)
  - FE: [svelte](https://svelte.dev)

# Dev Env

- `go install github.com/bokwoon95/wgo@latest` For live reloading BE code

## Commands

- `make prep` run after doing changes and before commiting (cleaning up the server and checking the UI code)
- `make dev` run backend & frontend in dev mode with live/hot reloading
  - http://localhost:8080/docs for inspecting and using the API
