# Spec First

Spec First is a schema-first event guest-list application. It supports user
registration and JWT login, event management, and guest/RSVP management.

The repository is a monorepo built around
[`docs/api_schema.yaml`](docs/api_schema.yaml), the OpenAPI contract and single
source of truth for the HTTP API:

- `backend/` — Go API using Chi, PostgreSQL, sqlc, goose, and a hexagonal
  (ports-and-adapters) architecture.
- `frontend/` — Vue 3, TypeScript, Pinia, Vue Router, and Vite.
- `docs/api_schema.yaml` — generates the Go server interface/models and the
  TypeScript API client/models.
- `docs/requests.http` — a manual request collection (register, log in,
  list/create an event, add/list a guest) for the [VS Code REST Client
  extension](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).
  Open it, run requests top to bottom with "Send Request" — later requests
  reuse the JWT and event id captured from earlier responses.
- `compose.yml` — local PostgreSQL only. The API and frontend run on the
  host during development. Works with Podman or Docker Compose.

Generated files are committed, but must not be edited by hand:

- `backend/internal/api/generated.go`
- `backend/internal/store/postgres/db/`
- `frontend/src/api/generated/`

## Prerequisites

- Go 1.27.1 (as declared in `backend/go.mod`)
- Node.js and npm
- Podman or Docker, with Compose support (the Makefile prefers Podman and
  falls back to Docker automatically if Podman isn't on `PATH`). On
  macOS/Windows with Podman, a running Podman machine is required.
- `make`

Go CLI tools are installed at pinned versions in the ignored root `.bin/`
directory. They do not need global installations:

| Tool | Pinned version | Used for |
| --- | ---: | --- |
| golangci-lint | v2.13.2 | Backend linting |
| oapi-codegen | v2.8.0 | Go API model and server generation |
| sqlc | v1.31.1 | Typed PostgreSQL query generation |
| goose | v3.28.0 | Creating migration files |

The corresponding Make targets install a missing or newly pinned tool
automatically. To install all tools eagerly or inspect their versions, run
`make install-go-tools` or `make tool-versions`.

## Start the local environment

Run these setup steps once from the repository root:

```sh
cp backend/.env.example backend/.env
make install-go-tools
cd frontend && npm ci && cd ..
```

Review `backend/.env` if local ports or credentials need to change. Its default
values match `compose.yml`. The Makefile loads this file automatically
for migration and backend commands.

Start each part in this order:

1. Start PostgreSQL and wait for it to become healthy:

   ```sh
   make dev-db
   ```

2. Apply all pending migrations:

   ```sh
   make migrate-up
   ```

3. Start the backend in one terminal:

   ```sh
   make run-backend
   ```

4. Start the frontend in another terminal:

   ```sh
   make run-frontend
   ```

The frontend is available at <http://localhost:5173> and calls the API at
<http://localhost:8080>. The API does not run migrations automatically.

To stop PostgreSQL without deleting its named volume:

```sh
make dev-db-down
```

## Database migrations

Migrations live in
`backend/internal/adapters/postgres/migrations/`. They are embedded into the
backend binary and use goose's numbered SQL format.

Create the next sequential migration from the repository root:

```sh
make migration-create NAME=add_event_owner
```

Give the file a descriptive snake-case name and implement both directions:

```sql
-- +goose Up
ALTER TABLE events ADD COLUMN owner_id uuid REFERENCES users (id);

-- +goose Down
ALTER TABLE events DROP COLUMN owner_id;
```

Then:

1. Update SQL in `backend/internal/adapters/postgres/queries/` if the schema
   change affects reads or writes.
2. Run `make generate-backend` from the repository root so sqlc regenerates
   `backend/internal/store/postgres/db/`.
3. Run `make migrate-up` to apply the migration locally.
4. Use `make migrate-down` to roll back the most recently applied migration
   while testing the `Down` section, then run `make migrate-up` again.
5. Update the PostgreSQL adapter and its integration tests for the new schema.

Do not edit an already-shared migration to change an existing database. Add a
new migration instead. Keep destructive `Down` operations in reverse dependency
order.

## Adding or changing an API endpoint

The project follows a contract-first workflow:

1. **Change the OpenAPI contract.** Add the path, HTTP operation, unique
   `operationId`, request/response schemas, errors, and authentication rules to
   `docs/api_schema.yaml`. Endpoints are authenticated by default; explicitly
   add `security: []` only for a public operation.
2. **Regenerate both sides.** From the repository root run:

   ```sh
   make generate
   ```

   This regenerates the Go API types/server interface, sqlc output, and the
   frontend TypeScript SDK/types. Never patch generated files directly.
3. **Implement backend behavior from the inside out.** Put entities and domain
   errors in `backend/internal/domain/`; orchestration and required repository
   interfaces in `backend/internal/app/`; then implement infrastructure in
   `backend/internal/adapters/`.
4. **Add persistence when needed.** Add a migration and named sqlc queries,
   regenerate, then implement the repository port in
   `backend/internal/adapters/postgres/`. Database and sqlc types must not leak
   into the domain or application packages.
5. **Implement the HTTP method.** Add the generated `ServerInterface` method to
   the appropriate file under `backend/internal/adapters/http/`. Decode generated
   `api` request types, call an application service, convert domain values in
   `convert.go`, and map domain errors to documented HTTP responses.
6. **Wire new resources.** If the endpoint introduces a new service or adapter,
   construct it in `backend/cmd/api/main.go` and inject it into
   `backend/internal/adapters/http/handlers.go`.
7. **Keep public routes in sync.** Because the generated Chi middleware is
   applied uniformly, a new operation marked `security: []` must also be added
   to `publicPaths` in `backend/internal/adapters/http/middleware.go`.
8. **Use the generated frontend client.** Import SDK functions and types through
   `frontend/src/api/apiClient.ts`, which supplies the base URL and bearer token,
   rather than importing from or editing `src/api/generated/` directly.
9. **Add tests and verify the change:**

   ```sh
   make build
   make test-unit
   make test-frontend
   make lint
   cd frontend && npm run build
   ```

   Use `make test-integration` for PostgreSQL adapter changes. Those tests start
   a real PostgreSQL container with testcontainers-go.

The backend dependency direction is `domain` <- `app` <- `adapters`, with
`cmd/api` as the composition root. HTTP handlers call application services; they
must not call the PostgreSQL adapter directly.

## Common commands

| Command | Purpose |
| --- | --- |
| `make dev-db` | Start local PostgreSQL (Podman Compose, or Docker Compose if Podman isn't installed) |
| `make dev-db-down` | Stop local PostgreSQL |
| `make migrate-up` | Apply pending migrations |
| `make migrate-down` | Roll back one migration |
| `make migration-create NAME=...` | Create a sequential goose migration |
| `make run-backend` | Run the API on the configured `PORT` |
| `make run-frontend` | Run the Vite development server |
| `make generate` | Regenerate Go API/sqlc code and the TypeScript client |
| `make build` | Build all backend packages |
| `make lint` | Lint backend and frontend |
| `make test` | Run backend unit/integration and frontend tests |
| `make install-go-tools` | Install all pinned Go CLI tools into `.bin/` |
| `make tool-versions` | Install tools if needed and print their versions |

### Podman/Docker integration tests

`make test-integration` runs testcontainers-go tests against whichever
container runtime the Makefile detected (Podman preferred, Docker as
fallback — see `CONTAINER_ENGINE` in the Makefile).

- **Docker:** no extra setup — `make test-integration` just works.
- **Podman:** the Makefile automatically points testcontainers at Podman's
  rootless machine socket and disables Ryuk (testcontainers' cleanup sidecar,
  which doesn't work against rootless Podman) for this target. If you invoke
  `go test -tags=integration ./...` directly instead of through `make`,
  export both yourself first:

  ```sh
  export DOCKER_HOST="unix://$(podman machine inspect --format '{{.ConnectionInfo.PodmanSocket.Path}}')"
  export TESTCONTAINERS_RYUK_DISABLED=true
  ```

## Choosing an oapi-codegen server target

The backend currently selects `chi-server` in
`backend/oapi-codegen.config.yaml`, but `oapi-codegen` v2.8.0 can generate
server boilerplate for Chi, Echo v4/v5, Fiber v2/v3, Gin, gorilla/mux, Iris,
or the Go standard library's `net/http` router. The corresponding generation
flags are `chi-server`, `echo-server`, `echo5-server`, `fiber-server`,
`fiber-v3-server`, `gin-server`, `gorilla-server`, `iris-server`, and
`std-http-server`.

Only one server target may be enabled at a time. `strict-server` can optionally
be enabled alongside that target to generate the stricter request/response
interface; it is not a router itself. Switching targets requires regenerating
the backend and adapting framework-specific routing, middleware, and handler
signatures where necessary. Do not edit the generated code directly.
