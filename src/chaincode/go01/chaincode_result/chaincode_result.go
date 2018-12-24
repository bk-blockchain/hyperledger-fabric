package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ChaincodeResult example simple Chaincode implementation
type ChaincodeResult struct {
}

type Result struct {
	ResultTransactionID string `json:"result_transaction_id"`
	ResultID            string `json:"result_id"`
	RecipientID         string `json:"recipient_id"`
	RecipientName       string `json:"recipient_name"`
	CourseID            string `json:"course_id"`
	CourseName          string `json:"course_name"`
	IssueID             string `json:"issue_id"`
	IssueName           string `json:"issue_name"`
	Grade               string `json:"grade"`
	Time                string `json:"time"`
	OriginalResultID    string `json:"original_result_id"`
}

func (t *ChaincodeResult) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("result infomation init")
	return shim.Success(nil)
}

func (t *ChaincodeResult) initResult(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0           1       		2           3			 4
	// "userID", "nameUser", "dateOfBrith", "sexUser", "addressUser""
	if len(args) < 10 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	ResultTransactionID := args[0]
	ResultID := args[1]
	RecipientID := args[2]
	RecipientName := strings.ToLower(args[3])
	CourseID := strings.ToLower(args[4])
	CourseName := strings.ToLower(args[5])
	IssueID := args[6]
	IssueName := strings.ToLower(args[7])
	Grade := args[8]
	Time := args[9]

	// ==== Check if user already exists ====
	userAsBytes, err := stub.GetState(ResultID)
	if err != nil {
		return shim.Error("Failed to get result: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This result already exists: " + ResultID)
		return shim.Error("This result already exists: " + ResultID)
	}

	// ==== Create user object and marshal to JSON ====
	result := &Result{ResultTransactionID, ResultID, RecipientID, RecipientName, CourseID, CourseName, IssueID, IssueName, Grade, Time, ""}
	resultJSONasBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save user to state ===
	err = stub.PutState(ResultID, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ChaincodeResult) updateResult(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0         1			 	2			3			 4
	// "userID", "nameUser", "dateOfBrith", "sexUser", "addressUser"
	if len(args) < 10 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ResultTransactionID := args[0]
	ResultID := args[1]
	RecipientID := args[2]
	RecipientName := strings.ToLower(args[3])
	CourseID := strings.ToLower(args[4])
	CourseName := strings.ToLower(args[5])
	IssueID := args[6]
	IssueName := strings.ToLower(args[7])
	Grade := args[8]
	Time := args[9]
	var OriginalResultID string

	fmt.Println("- start updateResult ", ResultID)

	resultAsBytes, err := stub.GetState(ResultID)
	if err != nil {
		return shim.Error("Failed to get result:" + err.Error())
	} else if resultAsBytes == nil {
		return shim.Error("result does not exist")
	}

	resultOld := Result{}
	err = json.Unmarshal(resultAsBytes, &resultOld)
	if err != nil {
		return shim.Error(err.Error())
	}
	if OriginalResultID == "" {
		OriginalResultID = ResultTransactionID
	}
	resultOld.ResultTransactionID = ResultTransactionID
	resultOld.RecipientID = RecipientID
	resultOld.RecipientName = RecipientName
	resultOld.CourseID = CourseID
	resultOld.CourseName = CourseName
	resultOld.IssueID = IssueID
	resultOld.IssueName = IssueName
	resultOld.Grade = Grade
	resultOld.Time = Time

	resultJSONasBytes, _ := json.Marshal(resultOld)
	err = stub.PutState(RecipientID, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end success")
	return shim.Success(nil)
}

func (t *ChaincodeResult) getResultByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	resultID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"result_id\":\"%s\"}}", resultID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (t *ChaincodeResult) getResultByRecipientID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	recipientID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"recipient_id\":\"%s\"}}", recipientID)

	queryResults, err := getValueQueryResultForQueryString(stub, queryString)
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

func getValueQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getValueQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		buffer.WriteString(string(queryResponse.Value))
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *ChaincodeResult) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Result information Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "getResultByID" {
		// get
		return t.getResultByID(stub, args)
	} else if function == "updatResult" {
		// update
		return t.updateResult(stub, args)
	} else if function == "initResult" {
		// create
		fmt.Println("initResult")
		return t.initResult(stub, args)
	} else if function == "getResultByRecipientID" {
		// get
		return t.getResultByRecipientID(stub, args)
	}

	return shim.Error("Invalid invoke function name")
}

func main() {
	err := shim.Start(new(ChaincodeResult))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
