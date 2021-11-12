package cmd

type App struct {
	Mode          string
	InletsVersion string
	Name          string
	Org           string
	Region        string
	Token         string
	Ports         []Ports
}

type Ports struct {
	InternalPort  int
	ExternalPorts []int
}
