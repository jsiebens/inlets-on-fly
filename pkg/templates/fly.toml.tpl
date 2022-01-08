app = "{{.Name}}"
[[services]]
  internal_port = 8123
  protocol = "tcp"
  [[services.ports]]
    handlers = []
    port = 8123
  [[services.tcp_checks]]
    grace_period = "1s"
    interval = "15s"
    restart_limit = 6
    timeout = "2s"

{{range .Ports}}[[services]]
  internal_port = {{.InternalPort}}
  protocol = "tcp"
  {{range .ExternalPorts}}[[services.ports]]
    handlers = [{{if eq . 80}}"http"{{else if eq . 443}}"tls","http"{{end}}]
    port = {{.}}
  {{end}}
{{end}}