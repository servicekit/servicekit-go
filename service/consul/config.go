package consul

type ConfigConsul struct {
	Addr      string
	Scheme    string
	Token     string
	KVPath    string
	TagPrefix string
}
