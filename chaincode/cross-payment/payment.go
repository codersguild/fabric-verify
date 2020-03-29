package main

import (
	"encoding/json"
	"time"
	"strings"
	"strconv"
	"fmt"
	"crypto/x509"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	// "github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	peer "github.com/hyperledger/fabric-protos-go/peer"
)

type crossPaymentContract struct {

}

type bank struct {
	ObjectType 			 string 		`json:"docType"`
	BankId               string 		`json:"bank_id"`
	BankName             string 		`json:"bank_name"`
	BankType             string 		`json:"bank_type"`
	SuppNonMemberBanks	 string 		`json:"supp_non_member_banks"`
	Certificate      	 x509.Certificate `json:"certificate"`
}

type fbank_addnl_curr struct {
	ObjectType 		string  `json:"docType"`
	Bank_name		string  `json:"bank_name"`
	Currency 		string  `json:"currency"`
	Exchange_rate 	float64 `json:exchange_rate`
	Balance 		float64 `json:balance`
}

type Transaction struct {
	ObjectType 			string 	`json:"docType"`
	Origin_timestamp 	string 	`json:"origin_time"`
	Trans_id  			string  `json:"trans_id"`
	Src_bank  			string  `json:"src_bank"`
	Dest_bank 			string  `json:"dest_bank"`
	Amount    			float64 `json:"amount"`
	Src_curr  			string  `json:"src_curr"`
	Dest_curr 			string  `json:"dest_curr"`
	Fbank 				string  `json:"fbank"`
	Rbank 				string  `json:"rbank"`
	Last_approved 		string 	`json:"last_approved"`
	Assigned_to 		string 	`json:"assigned_to"`
	Trans_status 		string 	`json:"trans_status"`
	Update_timestamp 	string 	`json:"update_time"`
}

func (s *crossPaymentContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (s *crossPaymentContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "initChaincodePayment" {
		return s.initChaincodePayment(APIstub, args)
	} else if function == "add_forex_currency" {
		return s.add_forex_currency(APIstub, args)
	} else if function == "allocate_funds" {
		return s.allocate_funds(APIstub, args)
	} else if function == "create_bank" {
		return s.create_bank(APIstub, args)
	} else if function == "read_bank" {
		return s.read_bank(APIstub, args)
	} else if function == "approve_transaction" {
		return s.approve_transaction(APIstub, args)
	} else if function == "get_completed_transaction" {
		return s.get_completed_transaction(APIstub, args)
	} else if function == "get_pending_transaction" {
		return s.get_pending_transaction(APIstub, args)
	} else if function == "get_supported_currencies" {
		return s.get_supported_currencies(APIstub, args)
	} else if function == "get_supported_non_member_banks" {
		return s.get_supported_non_member_banks(APIstub, args)
	} else if function == "list_fbanks" {
		return s.list_fbanks(APIstub)
	} else if function == "list_mbanks" {
		return s.list_mbanks(APIstub)
	} else if function == "list_rbanks" {
		return s.list_rbanks(APIstub)
	} else if function == "show_bank_details" {
		return s.show_bank_details(APIstub)
	} else if function == "query_balance" {
		return s.query_balance(APIstub, args)
	} else if function == "set_exchange_rate" {
		return s.set_exchange_rate(APIstub, args)
	} else if function == "transfer_money" {
		return s.transfer_money(APIstub, args)
	} else if function == "say_hello" {
		return s.say_hello(APIstub, args)
	} else if function == "Init" {
		return s.Init(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name : " + function)
}

func createResult(APIstub shim.ChaincodeStubInterface, payload []byte, code string, msg string) []byte {
	txnID := APIstub.GetTxID()
	timestamp, _ := APIstub.GetTxTimestamp()
	resultResponse := callResponse{code, msg, payload, txnID, time.Unix(timestamp.GetSeconds(), 0).String()}
	return resultResponse.JSONformatResponse()
}

func (s *crossPaymentContract) initChaincodePayment (APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

 	return shim.Success(nil)
}

func (s *crossPaymentContract) say_hello(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

 	return shim.Success(createResult(APIstub, []byte("Chaincode Installed and Initaited"), CODESUCCESS, "say_hello() invoked."))
}

// ======================================================
// Create Bank Data - create bank data in chaincode state
// ======================================================
func (s *crossPaymentContract) create_bank(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 2 {
		return shim.Error("Expecting 2 args. bank_name, bank_type")
	}

	bankName := strings.ToLower(args[0])
	bankType := strings.ToLower(args[1])
	bankId := getMD5Hash(bankName + bankType)
	suppNonMemberBanks := "None"
	certificate := getDummyCertificate()  //  need to import real certificate later.

	bankAsBytes, err := APIstub.GetState(bankName)
	if err != nil {
		return shim.Error("Failed to Query : " + err.Error())
	} else if bankAsBytes != nil {
		return shim.Error("This bank already exists : " + bankName)
	}

	objectType := "bankDoc"
	bank := &bank{objectType, bankId, bankName, bankType, suppNonMemberBanks, certificate}
	bankJSONasBytes, err := json.Marshal(bank)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = APIstub.PutState(bankName, bankJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	byteToJSON(bankJSONasBytes, 2)
	return shim.Success(createResult(APIstub, bankJSONasBytes, CODESUCCESS, "create_bank() invoked."))
}

// ======================================================
// Read Bank Data - read a bank data from chaincode state
// ======================================================
func (s *crossPaymentContract) read_bank(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting only name of the bank to query")
	}

	name := strings.ToLower(args[0])
	valAsbytes, err := APIstub.GetState(name)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"Bank does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	byteToJSON(valAsbytes, 2)
	return shim.Success(createResult(APIstub, valAsbytes, CODESUCCESS, "read_bank() invoked."))
}

func (s *crossPaymentContract) add_forex_currency(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 3 {
		return shim.Error("Expecting atleast 3 args. bank_name, currency, exchange_rate, optional(balance)")
	}

	objectType := "fbank_addnl_curr"
	bank_name := strings.ToLower(args[0])
	currency := strings.ToLower(args[1])
	balance := 0.000

	exchange_rate, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("Parse Error in " + args[2])
	}

	if len(args) == 4 {
		balance, err = strconv.ParseFloat(args[3], 64)
		if err != nil {
			return shim.Error("Parse Error in " + args[3])
		}
	}

	fbankObj := &fbank_addnl_curr{objectType, bank_name, currency, exchange_rate, balance}
	fbankObjJSONasBytes, err := json.Marshal(fbankObj)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Add forex currency support and update fbank_addnl_curr table.
	fbankObjName := bank_name + currency + "_forex"
	addfbankObj := APIstub.PutState(fbankObjName, fbankObjJSONasBytes)
	if addfbankObj != nil {
		return shim.Error(err.Error())
	}

	byteToJSON(fbankObjJSONasBytes, 2)
	// TODO : To add support for a currency and set exchange rate of the same against INR
	// Will need a couchdb database store and set_exchange_rate().
	return shim.Success(createResult(APIstub, fbankObjJSONasBytes, CODESUCCESS, "read_bank() invoked."))
}

func (s *crossPaymentContract) allocate_funds(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Expecting 2 args currency, amount")
	}

	// currency := strings.ToLower(args[0])
	// amount, err := strconv.ParseFloat(args[1], 64)
	// if err != nil {
	// 	return shim.Error("Parse Float Error. " + args[1])
	// }

	// // Check if bank exists and forex entry available.
	// bank_name := "rbi" // Hard Coded Sponsor Bank, later we change to getSponsorBank()
	// bankIndex := bank_name + currency + "_forex"
	// bankData := &fbank_addnl_curr{}

	// // Get bank data to see if it already exists.
	// bankDataJSONasBytes, err := APIstub.GetState(bankIndex)
	// if err != nil {
	// 	fmt.Println("Failed to fetch fbank entry. " + bank_name)
	// 	return shim.Error(err.Error())
	// } else if bankDataJSONasBytes == nil {
	// 	fmt.Println("Adding new fbank Entry. " + bank_name)
	// 	s.add_forex_currency(APIstub, []string{"rbi", currency, "1.00", args[1]}) // initial exchange_rate = 1.00 default
	// }

	// // Bank data already exists
	// fmt.Println("Allocating funds to Sponsor Bank. " + bank_name)
	// // Get the fresh bank data.
	// bankDataJSONasBytes, err = APIstub.GetState(bankIndex)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// bankErr := json.Unmarshal(bankDataJSONasBytes, bankData)
	// if bankErr != nil {
	// 	return shim.Error(err.Error())
	// }

	// // Modify amount.
	// exchange_rate := bankData.exchange_rate
	// // old_amount := bankData.balance
	// objectType := "fbank_addnl_curr"

	// bankData = &fbank_addnl_curr{objectType, bank_name, currency, exchange_rate, amount}
	// bankDataJSONasBytes, _ = json.Marshal(bankData)

	// // Add back (rewrite) the data to fbank_addnl_curr table.
	// err = APIstub.PutState(bankIndex, bankDataJSONasBytes)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// // Check if object really added.
	// bankDataJSONasBytes, err = APIstub.GetState(bankIndex)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// bankErr = json.Unmarshal(bankDataJSONasBytes, bankData)
	// if bankErr != nil {
	// 	return shim.Error(err.Error())
	// }

	// fmt.Println(bankData)

	// TODO : To add currency to Sponsor Bank
	// Will need a couchdb database store.
	return shim.Success(nil)
}

func (s *crossPaymentContract) approve_transaction(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Expecting 2 args, bank_name, trans_id.")
	}

	// // Get the transaction by index
	// transcIndex := strings.ToLower(args[0]) + args[1] + "_trans"
	// transcJSONasBytes, err := APIstub.GetState(transcIndex)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// } else if transcJSONasBytes == nil {
	// 	return shim.Error("Transaction not found. " + err.Error())
	// }

	// // Modify the state
	// transcObj := Transaction{}
	// transcErr := json.Unmarshal(transcJSONasBytes, transcObj)
	// if transcErr != nil {
	// 	return shim.Error(err.Error())
	// }

	// // Update the new timestamp and approval status
	// transcObj.update_timestamp = time.Now().String()
	// if transcObj.dest_bank == args[0] {
	// 	transcObj.trans_status = "completed"
	// }  else {
	// 	transcObj.trans_status = "pending"
	// }
	// transcObj.last_approved = args[0] + "_approved"

	// // Rewrite to the state.
	// transcJSONasBytes, _ = json.Marshal(transcObj)
	// putErr := APIstub.PutState(transcIndex, transcJSONasBytes)
	// if putErr != nil {
	// 	return shim.Error(err.Error())
	// }
	// TODO : To add transactions table with new transaction.
	// Will need a couchdb database store.
	return shim.Success(nil)
}

func (s *crossPaymentContract) get_completed_transaction(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Expecting 1 args, bank_name")
	}

	bank_name := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"transaction\",\"src_bank\":\"%s\",\"trans_status\":\"%s\"}, \"fields\":[\"trans_id\"]}", bank_name, "completed")

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (s *crossPaymentContract) get_pending_transaction(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Expecting 1 args, bank_name")
	}

	bank_name := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"transaction\",\"assigned_to\":\"%s\",\"trans_status\":\"%s\", \"fields\":[\"trans_id\"]}", bank_name, "pending")

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (s *crossPaymentContract) get_supported_currencies(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Expecting 1 args, bank_name")
	}

	bank_name := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"fbank_addnl_curr\",\"bank_name\":\"%s\", \"fields\":[\"currency\"]}", bank_name)

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (s *crossPaymentContract) get_supported_non_member_banks(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Expecting 1 args, bank_name")
	}

	bank_name := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"banks\",\"bank_name\":\"%s\", \"fields\":[\"supp_non_member_banks\"]}", bank_name)

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (s *crossPaymentContract) list_fbanks(APIstub shim.ChaincodeStubInterface) peer.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"banks\",\"bank_type\":\"fbank\", \"fields\":[\"bank_name\"]}")

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (s *crossPaymentContract) list_mbanks(APIstub shim.ChaincodeStubInterface) peer.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"banks\",\"bank_type\":\"mbank\", \"fields\":[\"bank_name\"]}")

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (s *crossPaymentContract) list_rbanks(APIstub shim.ChaincodeStubInterface) peer.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"banks\",\"bank_type\":\"rbank\", \"fields\":[\"bank_name\", \"bank_id\"]}")

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (s *crossPaymentContract) show_bank_details(APIstub shim.ChaincodeStubInterface) peer.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"banks\", \"fields\":[\"bank_name\", \"bank_id\"]}")

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

func (s *crossPaymentContract) query_balance(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Expecting 2 args, bank_name, currency.")
	}

	// bank_name := strings.ToLower(args[0])
	// currency := strings.ToLower(args[1])

	// bankIndex := bank_name + currency + "_forex"
	// bankData := &fbank_addnl_curr{}

	// // Get bank data to see if it already exists.
	// bankDataJSONasBytes, err := APIstub.GetState(bankIndex)
	// if err != nil {
	// 	return shim.Error("Failed to fetch fbank entry. " + err.Error())
	// } else if bankDataJSONasBytes == nil {
	// 	return shim.Error("Bank not added. " + err.Error())
	// }

	// // Bank data already exists
	// fmt.Println("Getting Balance for " + bank_name)

	// // fresh details
	// bankDataJSONasBytes, err = APIstub.GetState(bankIndex)
	// bankErr := json.Unmarshal(bankDataJSONasBytes, bankData)
	// if bankErr != nil {
	// 	return shim.Error(err.Error())
	// }

	// fmt.Println(bankData.balance)
	return shim.Success(nil)
}

func (s *crossPaymentContract) set_exchange_rate(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("Expecting 3 args, bank_name, currency, exchange_rate.")
	}

	// bank_name := strings.ToLower(args[0])
	// currency := strings.ToLower(args[1])
	// exchange_rate, parseErr := strconv.ParseFloat(args[2], 64)
	// if parseErr != nil {
	// 	return shim.Error("Parse Error Occured : " + parseErr.Error())
	// }

	// bankIndex := bank_name + currency + "_forex"
	// bankData := &fbank_addnl_curr{}

	// // Get bank data to see if it already exists.
	// bankDataJSONasBytes, err := APIstub.GetState(bankIndex)
	// if err != nil {
	// 	return shim.Error("Failed to fetch bank details. " + err.Error())
	// } else if bankDataJSONasBytes == nil {
	// 	return shim.Error("Bank not added. " + err.Error())
	// }

	// // Bank data already exists
	// fmt.Println("Setting Exchange Rate. " + bank_name + ", " + currency)

	// // Get the fresh bank data.
	// bankDataJSONasBytes, err = APIstub.GetState(bankIndex)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// bankErr := json.Unmarshal(bankDataJSONasBytes, bankData)
	// if bankErr != nil {
	// 	return shim.Error(err.Error())
	// }

	// // Modify exchange_rate
	// bankData.exchange_rate = exchange_rate
	// bankDataJSONasBytes, _ = json.Marshal(bankData)

	// // Add back (rewrite) the data to fbank_addnl_curr table.
	// err = APIstub.PutState(bankIndex, bankDataJSONasBytes)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	return shim.Success(nil)
}

func (s *crossPaymentContract) transfer_money(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	// if len(args) < 5 {
	// 	return shim.Error("Expecting atlest 5 args. src_bank, dest_bank, amount, src_curr, dest_curr. ")
	// }

	// // Check needs to be added on what previous balances were. Deduct the amount and then continue transaction.
	// src_bank := strings.ToLower(args[0])
	// dest_bank := strings.ToLower(args[1])
	// src_curr := strings.ToLower(args[3])
	// dest_curr := strings.ToLower(args[4])
	// origin_timestamp := time.Now().String()
	// trans_id := getMD5Hash(src_bank + dest_bank + src_curr + dest_curr + origin_timestamp)
	// objectType := "transaction"
	// amount, parseErr := strconv.ParseFloat(args[2], 64)
	// if parseErr != nil {
	// 	return shim.Error("Parse Error Occured : " + parseErr.Error())
	// }


	// bankIndex := src_bank + dest_curr + "_forex"
	// bankData := &fbank_addnl_curr{}

	// bankDataJSONasBytes, err := APIstub.GetState(bankIndex)
	// if err != nil {
	// 	return shim.Error("Failed to fetch bank details. " + err.Error())
	// } else if bankDataJSONasBytes == nil {
	// 	return shim.Error("Bank not added. " + err.Error())
	// }

	// bankErr := json.Unmarshal(bankDataJSONasBytes, bankData)
	// if bankErr != nil {
	// 	return shim.Error(err.Error())
	// }

	// current_balance := bankData.balance
	// new_balance := current_balance - amount
	// bankData.balance = new_balance

	// if new_balance > 0 {

	// 	// Add back updated balance entry.
	// 	bankDataJSONasBytes, _ = json.Marshal(bankData)

	// 	// Add back (rewrite) the data to fbank_addnl_curr table.
	// 	err = APIstub.PutState(bankIndex, bankDataJSONasBytes)
	// 	if err != nil {
	// 		return shim.Error(err.Error())
	// 	}
	// 	// Add a new transactions. But need to check if transaction already exists. // TODO : Check if transaction id already exists or got completed.
	// 	// ObjectType 			string 	`json:"docType"`
	// 	// origin_timestamp 	string 	`json:"origin_time"`
	// 	// trans_id  			string  `json:"trans_id"`
	// 	// src_bank  			string  `json:"src_bank"`
	// 	// dest_bank 			string  `json:"dest_bank"`
	// 	// amount    			float64 `json:"amount"`
	// 	// src_curr  			string  `json:"src_curr"`
	// 	// dest_curr 			string  `json:"dest_curr"`
	// 	transcObj := &Transaction{objectType, origin_timestamp, trans_id, src_bank, dest_bank, amount, src_curr, dest_curr, "", "", "", "", "", ""}
	// 	transcJSONasBytes, err := json.Marshal(transcObj)
	// 	if err != nil {
	// 		return shim.Error(err.Error())
	// 	}

	// 	err = APIstub.PutState(trans_id, transcJSONasBytes)
	// 	if err != nil {
	// 		return shim.Error(err.Error())
	// 	}
	// } else {
	// 	return shim.Error("Insufficient Balance Error.")
	// }

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(crossPaymentContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}