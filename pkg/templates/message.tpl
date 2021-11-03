==================================================================
inlets PRO TCP ({{.Name}}) server summary:

  URL: wss://{{.Name}}.fly.dev:10023/connect
  Auth-token: {{.Token}}

Command:

# Obtain a license at https://inlets.dev
# Store it at $HOME/.inlets/LICENSE or use --help for more options

export PORTS="{{range .Ports}}{{.InternalPort}},{{end}}"
export UPSTREAM="localhost"

inlets-pro tcp client \
  --url wss://{{.Name}}.fly.dev:10023/connect \
  --token {{.Token}} \
  --upstream $UPSTREAM \
  --ports $PORTS

To delete:
  flyctl destroy {{.Name}}
