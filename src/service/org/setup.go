package org

type Org struct {
	OrgName  string `json:"org_name"`
	Peers    []PeerNode `json:"peers"`
	Orderers []OrdererNode `json:"orderers"`
}

func NewOrg(orgName string, peers []PeerNode, orderers []OrdererNode) *Org {
	return &Org{
		OrgName:  orgName,
		Peers:    peers,
		Orderers: orderers,
	}
}

func (org *Org) Setup() {

}
