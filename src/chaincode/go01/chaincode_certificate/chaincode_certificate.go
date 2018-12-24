package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type CertificateChaincode struct{}

// Certificate --
type Certificate struct {
	CertificateTransactionID string   `json:"certificate_transaction_id"`
	CertificateID            string   `json:"certificate_id"`
	CertificateName          string   `json:"certificate_name"`
	IssuerID                 string   `json:"issuer_id"`
	IssuerName               string   `json:"issuer_name"`
	RecipientID              string   `json:"recipient_id"`
	RecipientName            string   `json:"recipient_name"`
	Grade                    string   `json:"grade"`
	Time                     string   `json:"time"`
	ResultTransactionIDs     []string `json:"result_transaction_ids"`
}

type ResultTransaction struct {
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

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

// Init CertificateChaincode
func (t *CertificateChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("certificate Init")
	return shim.Success(nil)
}

func (t *CertificateChaincode) initCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var channelName string
	var queryArgs [][]byte

	channelName = ""

	//
	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}

	chaincodeName := args[0]
	certificateTransactionID := args[1]
	certificateID := args[2]
	certificateName := args[3]
	issuerID := args[4]
	issuerName := args[5]
	recipientID := args[6]
	recipientName := args[7]
	grade := args[8]
	time := args[9]

	queryArgs = toChaincodeArgs("getResultByRecipientID", recipientID)

	response := stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}
	listResultTransaction := []ResultTransaction{}

	err = json.Unmarshal(response.Payload, &listResultTransaction)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + recipientID + "\"}"
		return shim.Error(jsonResp)
	}

	// ==== Check if user already exists ====
	certificateAsBytes, err := stub.GetState(certificateTransactionID)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if certificateAsBytes != nil {
		fmt.Println("This certificate already exists: " + certificateTransactionID)
		return shim.Error("This certificate already exists: " + certificateTransactionID)
	}

	var listResultTransactionIDs []string
	for _, value := range listResultTransaction {
		listResultTransactionIDs = append(listResultTransactionIDs, value.ResultTransactionID)
	}

	certificate := &Certificate{certificateTransactionID, certificateID,
		certificateName, issuerID, issuerName, recipientID, recipientName,
		grade, time, listResultTransactionIDs}

	certificateJSONasBytes, err := json.Marshal(certificate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save user to state ===
	err = stub.PutState(certificateID, certificateJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *CertificateChaincode) updateCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var channelName string
	var queryArgs [][]byte

	channelName = ""
	//
	if len(args) < 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}

	chaincodeName := args[0]
	certificateTransactionID := args[1]
	certificateID := args[2]
	certificateName := args[3]
	issuerID := args[4]
	issuerName := args[5]
	recipientID := args[6]
	recipientName := args[7]
	grade := args[8]
	time := args[9]

	fmt.Println("- start update ", certificateTransactionID)

	queryArgs = toChaincodeArgs("getResultByRecipientID", recipientID)

	response := stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}
	listResultTransaction := []ResultTransaction{}

	err := json.Unmarshal(response.Payload, &listResultTransaction)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to decode JSON of: " + recipientID + "\"}"
		return shim.Error(jsonResp)
	}

	certificateAsBytes, err := stub.GetState(certificateID)
	if err != nil {
		return shim.Error("Failed to get certificate:" + err.Error())
	} else if certificateAsBytes == nil {
		return shim.Error("Certificate does not exist")
	}

	certificateOld := &Certificate{}

	err = json.Unmarshal(certificateAsBytes, &certificateOld)
	if err != nil {
		return shim.Error(err.Error())
	}

	if certificateName == "" {
		certificateName = certificateOld.CertificateName
	}
	if issuerName == "" {
		issuerName = certificateOld.IssuerName
	}
	if certificateName == "" {
		recipientName = certificateOld.RecipientName
	}
	if grade == "" {
		grade = certificateOld.Grade
	}
	if time == "" {
		time = certificateOld.Time
	}

	var listResultTransactionIDs []string
	for _, value := range listResultTransaction {
		listResultTransactionIDs = append(listResultTransactionIDs, value.ResultTransactionID)
	}

	certificate := &Certificate{certificateTransactionID, certificateID,
		certificateName, issuerID, issuerName, recipientID, recipientName,
		grade, time, listResultTransactionIDs}

	certificateJSONasBytes, _ := json.Marshal(certificate)
	err = stub.PutState(certificateID, certificateJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end success")
	return shim.Success(nil)
}

func (t *CertificateChaincode) getCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	certificateID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"certificate_id\":\"%s\"}}", certificateID)

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

// Invoke --
func (t *CertificateChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("certificate Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "getCertificate" {
		// get
		return t.getCertificate(stub, args)
	} else if function == "updateCertificate" {
		// update
		return t.updateCertificate(stub, args)
	} else if function == "initCertificate" {
		// create
		return t.initCertificate(stub, args)
	}
	return shim.Error("Invalid invoke function name")
}

func main() {
	err := shim.Start(new(CertificateChaincode))
	if err != nil {
		fmt.Printf("Error starting profile chaincode: %s", err)
	}
}
