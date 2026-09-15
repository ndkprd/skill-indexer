# Examples

Four ways to get a Skillstore site running, plus the sample data most of
them (and the test suite) use.

- `skills/` — 33 real Claude Code skills, used as fixture data by `go test`
  and as the sample input for the quickstarts below.
- `.gitlab-ci.yml` — a copy-pasteable pipeline that builds the CLI and
  publishes the generated site to GitLab Pages.
- `docker-compose.yaml` — a one-shot generate service plus nginx, for
  running the whole thing with Docker.
- `kubernetes.yaml` — a Deployment template: init containers clone your
  skills and generate the site, then nginx serves it.

## 1. Direct generate

From the repository root:

```bash
go build -o skillstore .
./skillstore --skill-dir examples/skills --output-dir public
```

Then serve `public/` with any static file server, e.g.:

```bash
python3 -m http.server -d public 8080
```

Open `http://localhost:8080`. Point `--skill-dir` at your own skills
directory instead of `examples/skills` once you're past trying it out.

## 2. Publish to GitLab Pages via CI

1. Copy `examples/.gitlab-ci.yml` to the root of your repository as
   `.gitlab-ci.yml`.
2. Make sure your skills live in a `skills/` directory at the repo root (or
   edit the `--skill-dir` argument in the copied file to match).
3. Push to your default branch.

GitLab runs the `pages` job, which builds the CLI, generates the site into
`public/`, and publishes it. The pipeline's Pages job details will show the
live URL once it completes — typically
`https://<namespace>.gitlab.io/<project>/`.

## 3. Run with Docker Compose

```bash
docker compose -f examples/docker-compose.yaml up --build
```

This builds the CLI from the repository's `Dockerfile`, runs it once
against `examples/skills/` into a shared volume, then starts nginx serving
that volume at `http://localhost:8080`. The `generate` service exits after
one run; `web` (nginx) keeps running.

To use your own skills instead of the bundled examples, copy
`examples/docker-compose.yaml` (and the `Dockerfile`) into your project and
point the `generate` service's `./skills` volume at your own directory. See
the comments in the file for details.

## 4. Deploy to Kubernetes

`examples/kubernetes.yaml` is a Deployment template — not something to
`kubectl apply -f` unmodified. It chains two init containers (clone your
skills repo with `git`, then generate the site with the Skillstore image)
ahead of an nginx container that serves the result, all sharing per-pod
`emptyDir` volumes.

Before applying it, you need to:

1. Build and push your own Skillstore image (see the `Dockerfile`) and
   replace `registry.gitlab.com/endekasoft/skillstore:latest` in the
   manifest with wherever you pushed it.
2. Replace `https://git.example.com/skills.git` with your real skills
   repository (a private repo needs a mounted credentials Secret for the
   clone step — not included in the template).

The manifest's own comments cover both steps, plus a note on scaling
considerations (each replica clones and generates independently on
startup).

## Building a container image without Compose or Kubernetes

The repository root also has a `Dockerfile` that packages the CLI (not the
generated site) into a minimal image:

```bash
docker build -t skillstore .
docker run --rm \
  -v "$(pwd)/examples/skills:/skills:ro" \
  -v "$(pwd)/public:/public" \
  skillstore --skill-dir /skills --output-dir /public
```

Useful if you'd rather run generation inside your own infrastructure by
hand than install Go locally.
