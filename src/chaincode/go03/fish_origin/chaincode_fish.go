package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ChaincodeFish example simple Chaincode implementation
type ChaincodeFish struct {
}

type Fish struct {
	FishID         string        `json:"fish_id"`
	Name           string        `json:"name"`
	WeightPackage  string        `json:"weight_package"`
	ImagePackage   string        `json:"image_package"`
	TimeFishing    string        `json:"time_fishing"`
	AddressFishing string        `json:"address_fishing"`
	IDTransaction  string        `json:"id_transaction"`
	Certificates   []Certificate `json:"certifications"`
}

type Certificate struct {
	CodeCertificate    string `json:"code_certification"`
	NameOrgCertificate string `json:"name_org_certification"`
	DateCertificate    string `json:"date_certification"`
}

func (t *ChaincodeFish) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("fish infomation init")
	return shim.Success(nil)
}

func (t *ChaincodeFish) initProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0           1       		2           3			 4
	// "FishID", "NameFish", "TimeFishing", "AddressFishing", "CodeCertificate", "NameOrgCertificate", "DateCertificate"
	if len(args) < 7 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	IDTransaction := args[0]
	FishID := args[1]
	Name := args[2]
	WeightPackage := args[3]
	ImagePackage := strings.ToLower(args[4])
	TimeFishing := strings.ToLower(args[5])
	AddressFishing := strings.ToLower(args[6])
	var Certificates []Certificate

	// ==== Check if product already exists ====
	userAsBytes, err := stub.GetState(FishID)
	if err != nil {
		return shim.Error("Failed to get result: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This result already exists: " + FishID)
		return shim.Error("This result already exists: " + FishID)
	}

	// ==== Create product object and marshal to JSON ====
	result := &Fish{FishID, Name, WeightPackage, ImagePackage, TimeFishing, AddressFishing, IDTransaction, Certificates}
	resultJSONasBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	err = stub.PutState(FishID, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ChaincodeFish) updateCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	IDTransaction := args[0]
	FishID := args[1]
	CodeCertificate := args[2]
	NameOrgCertificate := args[3]
	DateCertificate := args[4]

	resultAsBytes, err := stub.GetState(FishID)
	if err != nil {
		return shim.Error("Failed to get result:" + err.Error())
	} else if resultAsBytes == nil {
		return shim.Error("result does not exist")
	}

	resultOld := Fish{}
	err = json.Unmarshal(resultAsBytes, &resultOld)
	if err != nil {
		return shim.Error(err.Error())
	}

	resultOld.IDTransaction = IDTransaction

	Certificates := Certificate{}

	Certificates.CodeCertificate = CodeCertificate
	Certificates.NameOrgCertificate = NameOrgCertificate
	Certificates.DateCertificate = DateCertificate

	resultOld.Certificates = append(resultOld.Certificates, Certificates)

	resultJSONasBytes, _ := json.Marshal(resultOld)
	err = stub.PutState(FishID, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ChaincodeFish) deleteProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	FishID := args[1]

	resultAsBytes, err := stub.GetState(FishID)
	if err != nil {
		return shim.Error("Failed to get result:" + err.Error())
	} else if resultAsBytes == nil {
		return shim.Error("result does not exist")
	}

	err = stub.DelState(FishID) //remove the product from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}

func (t *ChaincodeFish) getResultByFishID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	FishID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"fish_id\":\"%s\"}}", FishID)

	queryResults, err := getValueQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (t *ChaincodeFish) getAllProductNotCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"certifications\":null}}")

	queryResults, err := getValueQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (t *ChaincodeFish) getAllProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"fish_id\":{\"$gt\":null}}}")

	queryResults, err := getValueQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getValueQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getValueQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *ChaincodeFish) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Result information Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "getResultByFishID" {
		// get
		return t.getResultByFishID(stub, args)
	} else if function == "getAllProductNotCertificate" {
		// get all
		fmt.Println("getAllProductNotCertificate")
		return t.getAllProductNotCertificate(stub, args)
	} else if function == "getAllProduct" {
		// get all
		fmt.Println("getAllProduct")
		return t.getAllProduct(stub, args)
	} else if function == "initProduct" {
		// create
		return t.initProduct(stub, args)
	} else if function == "updateCertificate" {
		// update certificate
		return t.updateCertificate(stub, args)
	} else if function == "deleteProduct" {
		// delete
		return t.deleteProduct(stub, args)
	}

	return shim.Error("Invalid invoke function name")
}

func main() {
	err := shim.Start(new(ChaincodeFish))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

