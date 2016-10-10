package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChainCode struct {
}

func (cc *SimpleChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var A string
	var Aval int
	var err error
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	fmt.Printf("Aval= %d\n", Aval)
	err = stub.PutState(A, []byte(args[1]))
	return nil, err
}

func (cc *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "invoke" {
		return nil, errors.New("invalid function name. Expecting invoke")
	}
	if len(args) != 2 {
		return nil, errors.New("incorrect number of arguments. Expecting 2")
	}
	var A string
	var Aval int
	var err error
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	fmt.Printf("Aval=%d\n", Aval)
	err = stub.PutState(A, []byte(args[1]))
	return nil, errors.New("simulate fail")
	//return nil, err
}

func (cc *SimpleChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("invalid query function name. Expecting query")
	}
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments. Expecting 1")
	}
	A := args[0]
	if A == "" {
		return nil, errors.New("invalid key name. key is empty")
	}
	val, err := stub.GetState(A)
	fmt.Printf("Aval = %d\n", val)
	return val, err
}

func main() {
	err := shim.Start(new(SimpleChainCode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
