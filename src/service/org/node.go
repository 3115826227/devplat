package org

type Node struct {
	Address       string
	Port          []int
	EnableNodeOUs bool
	Name          string
	Domain        string
}

type PeerNode struct {
	Node
	TemplateCount int
	UserCount     int
}

type OrdererNode struct {
	Node
	Hostname string
}
