
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ChainCodeVegetable example simple Chaincode implementation
type ChainCodeVegetable struct {
}

//Vegetable
type Vegetable struct {
	TransactionID string        `json:"transaction_id"`
	Code          string        `json:"code"`
	Product       Product       `json:"product"`
	Import        Import        `json:"import"`
	ContainerFrom ContainerFrom `json:"container_from"`
	ContainerTo   ContainerTo   `json:"container_to"`
	Export        Export        `json:"export"`
	Type          string        `json:"type"`
}

//Product
type Product struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	Manufacture string `json:"manufacture"`
	Expiry      string `json:"expiry"`
}

type Import struct {
	CodeOrder       string `json:"code_order"`
	DateOrder       string `json:"date_order"`
	ProviderCode    string `json:"provider_code"`
	ProviderName    string `json:"provider_name"`
	ProviderAddress string `json:"provider_address"`
	ProviderEmail   string `json:"provider_email"`
	ProviderMobile  string `json:"provider_mobile"`
	ProviderImage   string `json:"provider_image"`
}

type ContainerFrom struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	Image   string `json:"image"`
	Code    string `json:"code"`
}

type ContainerTo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	Image   string `json:"image"`
	Code    string `json:"code"`
}

type Export struct {
	CodeOrder          string `json:"code_order"`
	DateOrder          string `json:"date_order"`
	SupermarketCode    string `json:"supermarket_code"`
	SupermarketName    string `json:"supermarket_name"`
	SupermarketAddress string `json:"supermarket_address"`
	SupermarketEmail   string `json:"supermarket_email"`
	SupermarketMobile  string `json:"supermarket_mobile"`
	SupermarketImage   string `json:"supermarket_image"`
}

//Init ChainCodeVegetable
func (t *ChainCodeVegetable) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Vegetable infomation init")
	return shim.Success(nil)
}

func (t *ChainCodeVegetable) createProductLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//
	if len(args) != 22 {
		return shim.Error("Incorrect number of arguments. Expecting 20")
	}

	var TransactionID string
	var ContainerFrom ContainerFrom
	var Export Export

	CodeVegetable := args[0]
	CodeProduct := args[1]
	NameProduct := args[2]
	CategoryProduct := args[3]
	ImageProduct := args[4]
	Manufacture := args[5]
	Expiry := args[6]
	CodeImport := args[7]
	DateImport := args[8]
	ProviderCodeImport := args[9]
	ProviderNameImport := args[10]
	ProviderAddressImport := args[11]
	EmailImport := args[12]
	MobileImport := args[13]
	ImageImport := args[14]
	NameContainerTo := args[15]
	AddressContainerTo := args[16]
	MobileContainerTo := args[17]
	EmailContainerTo := args[18]
	ImageContainerTo := args[19]
	CodeContainerTo := args[20]
	Type := args[21]

	Product := Product{CodeProduct, NameProduct, CategoryProduct, ImageProduct, Manufacture, Expiry}

	Import := Import{CodeImport, DateImport, ProviderCodeImport, ProviderNameImport,
		ProviderAddressImport, EmailImport, MobileImport, ImageImport}

	ContainerTo := ContainerTo{NameContainerTo, AddressContainerTo, MobileContainerTo,
		EmailContainerTo, ImageContainerTo, CodeContainerTo}

	// ==== Check if product already exists ====
	valueAsBytes, err := stub.GetState(CodeVegetable)
	if err != nil {
		return shim.Error("Failed to get result: " + err.Error())
	} else if valueAsBytes != nil {
		fmt.Println("This result already exists: " + CodeVegetable)
		return shim.Error("This result already exists: " + CodeVegetable)
	}

	// ==== Create product object and marshal to JSON ====
	result := &Vegetable{TransactionID, CodeVegetable, Product, Import, ContainerFrom, ContainerTo,
		Export, Type}

	resultJSONasBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	err = stub.PutState(CodeVegetable, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ChainCodeVegetable) updateProductFromTo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//
	if len(args) != 24 {
		return shim.Error("Incorrect number of arguments. Expecting 19")
	}

	var Import Import
	var Export Export

	TransactionID := args[0]
	CodeVegetable := args[1]
	CodeProduct := args[2]
	NameProduct := args[3]
	CategoryProduct := args[4]
	ImageProduct := args[5]
	Manufacture := args[6]
	Expiry := args[7]
	NameContainerFrom := args[8]
	AddressContainerFrom := args[9]
	MobileContainerFrom := args[10]
	EmailContainerFrom := args[11]
	ImageContainerFrom := args[12]
	CodeContainerFrom := args[13]
	NameContainerTo := args[14]
	AddressContainerTo := args[15]
	MobileContainerTo := args[16]
	EmailContainerTo := args[17]
	ImageContainerTo := args[18]
	CodeContainerTo := args[19]
	Type := args[20]

	Product := Product{CodeProduct, NameProduct, CategoryProduct, ImageProduct, Manufacture, Expiry}

	ContainerFrom := ContainerFrom{NameContainerFrom, AddressContainerFrom, MobileContainerFrom,
		EmailContainerFrom, ImageContainerFrom, CodeContainerFrom}

	ContainerTo := ContainerTo{NameContainerTo, AddressContainerTo, MobileContainerTo,
		EmailContainerTo, ImageContainerTo, CodeContainerTo}

	// ==== Check if product already exists ====
	valueAsBytes, err := stub.GetState(CodeVegetable)
	if err != nil {
		return shim.Error("Failed to get result: " + err.Error())
	} else if valueAsBytes != nil {
		fmt.Println("This result already exists: " + CodeVegetable)
		return shim.Error("This result already exists: " + CodeVegetable)
	}

	// ==== Create product object and marshal to JSON ====
	result := &Vegetable{TransactionID, CodeVegetable, Product, Import, ContainerFrom, ContainerTo,
		Export, Type}

	resultJSONasBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	err = stub.PutState(CodeVegetable, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ChainCodeVegetable) updateProductExport(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//
	if len(args) != 23 {
		return shim.Error("Incorrect number of arguments. Expecting 21")
	}

	var Import Import
	var ContainerTo ContainerTo

	TransactionID := args[0]
	CodeVegetable := args[1]
	CodeProduct := args[2]
	NameProduct := args[3]
	CategoryProduct := args[4]
	ImageProduct := args[5]
	Manufacture := args[6]
	Expiry := args[7]
	NameContainerFrom := args[8]
	AddressContainerFrom := args[9]
	MobileContainerFrom := args[10]
	EmailContainerFrom := args[11]
	ImageContainerFrom := args[12]
	CodeContainerFrom := args[13]
	CodeExport := args[14]
	DateExport := args[15]
	SupermarketCodeExport := args[16]
	SupermarketNameExport := args[17]
	SupermarketAddressExport := args[18]
	SupermarketEmailExport := args[19]
	SupermarketMobileExport := args[20]
	SupermarketImageExport := args[21]
	Type := args[22]

	Product := Product{CodeProduct, NameProduct, CategoryProduct, ImageProduct, Manufacture, Expiry}

	ContainerFrom := ContainerFrom{NameContainerFrom, AddressContainerFrom, MobileContainerFrom,
		EmailContainerFrom, ImageContainerFrom, CodeContainerFrom}

	Export := Export{CodeExport, DateExport, SupermarketCodeExport, SupermarketNameExport,
		SupermarketAddressExport, SupermarketEmailExport, SupermarketMobileExport, SupermarketImageExport}

	// ==== Check if product already exists ====
	valueAsBytes, err := stub.GetState(CodeVegetable)
	if err != nil {
		return shim.Error("Failed to get result: " + err.Error())
	} else if valueAsBytes != nil {
		fmt.Println("This result already exists: " + CodeVegetable)
		return shim.Error("This result already exists: " + CodeVegetable)
	}

	// ==== Create product object and marshal to JSON ====
	result := &Vegetable{TransactionID, CodeVegetable, Product, Import, ContainerFrom, ContainerTo,
		Export, Type}

	resultJSONasBytes, err := json.Marshal(result)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	err = stub.PutState(CodeVegetable, resultJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ChainCodeVegetable) getProductLogByCode(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	Code := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\"}}", Code)

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

func (t *ChainCodeVegetable) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Result Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "createProductLog" {
		// create
		return t.createProductLog(stub, args)
	} else if function == "updateProductFromTo" {
		// update
		fmt.Println("updateProductFromTo")
		return t.updateProductFromTo(stub, args)
	} else if function == "updateProductExport" {
		// get all
		fmt.Println("updateProductExport")
		return t.updateProductExport(stub, args)
	} else if function == "getProductLogByCode" {
		// create
		return t.getProductLogByCode(stub, args)
	}

	return shim.Error("Invalid invoke function name")
}

func main() {
	err := shim.Start(new(ChainCodeVegetable))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

