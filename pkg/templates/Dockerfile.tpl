FROM ghcr.io/inlets/inlets-pro:{{.InletsVersion}}
CMD ["tcp", "server", "--auto-tls-san={{.Name}}.fly.dev", "--token-env=TOKEN"]