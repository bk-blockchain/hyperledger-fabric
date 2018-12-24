
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type ChaincodeCustomer struct {
}

type Customer struct {
	CustomerId string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	CustomerMobile string `json:"customer_mobile"`
	CustomerStatus string `json:"customer_status"`
	CustomerPassword string `json:"customer_password"`
	CustomerDescription string `json:"customer_description"`
	CustomerAddress string 	`json:"customer_address`
}

func (t *ChaincodeCustomer) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("result infomation init")
	return shim.Success(nil)
}

func (t *ChaincodeCustomer) createCustomer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) < 7 {
		return shim.Error("Incorrect number of arguments")
	}

	CustomerId := args[0]
	CustomerName := args[1]
	CustomerMobile := args[2]
	CustomerStatus := args[3]
	CustomerPassword := args[4]
	CustomerDescription := args[5]
	CustomerAddress := args[6]

	userAsBytes, err := stub.GetState(CustomerId)
	if err != nil {
		return shim.Error("Failed to get result: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This result already exists: " + CustomerId)
		return shim.Error("This result already exists: " + CustomerId)
	}

	customer := &Customer{CustomerId, CustomerName, CustomerMobile, CustomerStatus, CustomerPassword, CustomerDescription, CustomerAddress}
	resultJSONasBytes, err := json.Marshal(customer)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(CustomerId, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ChaincodeCustomer) updateCustomer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 7 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	CustomerId := args[0]
	CustomerName := args[1]
	CustomerMobile := args[2]
	CustomerStatus := args[3]
	CustomerPassword := args[4]
	CustomerDescription := args[5]
	CustomerAddress := args[6]

	fmt.Println("- start updateCustomer ", CustomerId)

	resultAsBytes, err := stub.GetState(CustomerId)
	if err != nil {
		return shim.Error("Failed to get result:" + err.Error())
	} else if resultAsBytes == nil {
		return shim.Error("result does not exist")
	}

	customerOld := Customer{}
	err = json.Unmarshal(resultAsBytes, &customerOld)
	if err != nil {
		return shim.Error(err.Error())
	}
	customerOld.CustomerId = CustomerId
	customerOld.CustomerName = CustomerName
	customerOld.CustomerMobile = CustomerMobile
	customerOld.CustomerStatus = CustomerStatus
	customerOld.CustomerPassword = CustomerPassword
	customerOld.CustomerDescription = CustomerDescription
	customerOld.CustomerAddress = CustomerAddress

	resultJSONasBytes, _ := json.Marshal(customerOld)
	err = stub.PutState(CustomerId, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end success")
	return shim.Success(nil)
}

func (t *ChaincodeCustomer) getCustomerById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

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
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *ChaincodeCustomer) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Result information Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "createCustomer" {
		// get
		return t.createCustomer(stub, args)
	} else if function == "updateCustomer" {
		// update
		return t.updateCustomer(stub, args)
	}

	return shim.Error("Invalid invoke function name")
}

func main() {
	err := shim.Start(new(ChaincodeCustomer))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
