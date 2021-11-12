==================================================================
inlets PRO TCP ({{.Name}}) server summary:

  URL: wss://{{.Name}}.fly.dev:10023/connect
  Auth-token: {{.Token}}

Command:

# Obtain a license at https://inlets.dev
# Store it at $HOME/.inlets/LICENSE or use --help for more options

{{if eq .Mode "tcp"}}
export PORTS="{{range .Ports}}{{.InternalPort}},{{end}}"
export UPSTREAM="localhost"

inlets-pro tcp client \
  --url wss://{{.Name}}.fly.dev:10023 \
  --token {{.Token}} \
  --upstream $UPSTREAM \
  --ports $PORTS
{{else}}
export UPSTREAM="http://localhost:8080"

inlets-pro http client \
  --url wss://{{.Name}}.fly.dev:10023 \
  --token {{.Token}} \
  --upstream $UPSTREAM
{{end}}

To delete:
  flyctl destroy {{.Name}}
