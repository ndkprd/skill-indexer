# Examples

Two ways to get a Skillstore site running, plus the sample data both of them
(and the test suite) use.

- `skills/` — 33 real Claude Code skills, used as fixture data by `go test`
  and as the sample input for the quickstarts below.
- `.gitlab-ci.yml` — a copy-pasteable pipeline that builds the CLI and
  publishes the generated site to GitLab Pages.

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

## Building a container image instead

The repository root also has a `Dockerfile` that packages the CLI (not the
generated site) into a minimal image:

```bash
docker build -t skillstore .
docker run --rm \
  -v "$(pwd)/examples/skills:/skills:ro" \
  -v "$(pwd)/public:/public" \
  skillstore --skill-dir /skills --output-dir /public
```

Useful if you'd rather run generation inside your own infrastructure than
install Go locally.
