package spec

type Service struct {
    ID          string
    Service     string
    Tags        []string
    Version     string
    Address     string
    Port        int
    CreateIndex uint64
    ModifyIndex uint64
    NodeID      string
    NodeAddress string
    Node        string
    Datacenter  string
}
