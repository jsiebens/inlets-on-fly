FROM ghcr.io/inlets/inlets-pro:{{.InletsVersion}}
CMD ["{{.Mode}}", "server", "--auto-tls-san={{.Name}}.fly.dev", "--token-env=TOKEN"]