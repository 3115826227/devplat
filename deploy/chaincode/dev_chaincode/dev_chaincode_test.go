package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func Test(t *testing.T) {
	cc := new(SimpleChaincode)
	stub := shim.NewMockStub("SimpleChaincode", cc)
	initArgs := [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}
	res := stub.MockInit("1", initArgs)
	fmt.Println(res)

	queryArgs := [][]byte{[]byte("query"), []byte("a")}
	res = stub.MockInvoke("1", queryArgs)
	fmt.Println(res)

	invokeArgs := [][]byte{[]byte("invoke"), []byte("a"), []byte("b"), []byte("10")}
	res = stub.MockInvoke("1", invokeArgs)
	fmt.Println(res)

	queryArgs = [][]byte{[]byte("query"), []byte("a")}
	res = stub.MockInvoke("1", queryArgs)
	fmt.Println(res)
}
