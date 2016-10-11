package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Campany struct {
	CID          int    `json:"cid"`
	Name         string `json:"name"`
	CreateTime   int    `json:"create_time"` //nanoseconds
	ContractHash string `json:"contract_hash"`
	FundInitial  int    `json:"fund_initial"`
	FundBalance  int    `json:"fund_balance"`
	CType        int    `json:"ctype"`
}

type CampanyContract struct {
	UID             int `json:"uid"`
	CID             int `json:"cid"`
	SignatureStatus int `json:"signature_status"`
	SignatureTime   int `json:"signature_time"`
}

type User struct {
	UID        int    `json:"uid"`
	Name       string `json:"name"`
	PWD        string `json:"pwd"`
	CreateTime int    `json:"create_time"`
}

type UserRole struct {
	UID   int `json:"uid"`
	CID   int `json:"cid"`
	CType int `json:"ctype"`
	Role  int `json:"role"`
}

type UserFund struct {
	FID        int `json:"fid"`
	UID        int `json:"uid"`
	CID        int `json:"cid"`
	CreateTime int `json:"create_time"`
	ModifyTime int `json:"modify_time"`
	Fund       int `json:"fund"`
}

type Transaction struct {
	TXID               int    `json:"txid"`
	CID                int    `json:"cid"`
	FromUID            int    `json:"from_uid"`
	ToUID              int    `json:"to_uid"`
	FundAmount         int    `json:"fund_amount"`
	CreateTime         int    `json:"create_time"`
	LawyerAuditUID     int    `json:"lawyer_audit_uid"`
	LawyerAuditResult  int    `json:"lawyer_audit_result"`
	LawyerAuditRemark  string `json:"lawyer_audit_remark"`
	LawyerAuditTime    int    `json:"lawyer_audit_time"`
	AuditorAuditUID    int    `json:"auditor_audit_uid"`
	AuditorAuditResult int    `json:"auditor_audit_result"`
	AuditorAuditRemark string `json:"auditor_audit_remark"`
	AuditorAuditTime   int    `json:"auditor_audit_time"`
	Status             int    `json:"status"`
}

type SimpleChainCode struct {
}

func (cc *SimpleChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	cid := nextCID(stub)
	c1 := Campany{CID: cid, Name: "test_c1", CreateTime: time.Now().Nanosecond(), FundInitial: 10000, FundBalance: 2000, CType: 1}
	cByte, _ := json.Marshal(c1)
	stub.PutState("c_"+strconv.Itoa(cid), cByte)
	uid := nextUID(stub)
	u1 := User{UID: uid, Name: "test_u1", PWD: "111111", CreateTime: time.Now().Nanosecond()}
	uByte, _ := json.Marshal(u1)
	stub.PutState("u_"+strconv.Itoa(uid), uByte)
	return nil, nil
}

func (cc *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "createCompany":
		return cc.createCompany(stub, args)
	case "createUser":
		return cc.createUser(stub, args)
	case "registerUserRole":
		return cc.registerUserRole(stub, args)
	case "transfer":
		return cc.transfer(stub, args)
	case "lawyerAuditTransaction":
		return cc.lawyerAuditTransaction(stub, args)
	case "auditorAuditTransaction":
		return cc.auditorAuditTransaction(stub, args)
	default:
		return nil, fmt.Errorf("unexpected function:%s", function)
	}
}

//Query implements the
func (cc *SimpleChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {
	case "getCompanyByID":
		return cc.getCompanyByID(stub, args)
	case "getUserByID":
		return cc.getUserByID(stub, args)
	case "getUserCompanyFund":
		return cc.getUserCompanyFund(stub, args)
	case "getUserFund":
		return cc.getUserFund(stub, args)
	case "getUserRoleByCID":
		return cc.getUserRoleByCID(stub, args)
	case "getUserRole":
		return cc.getUserRole(stub, args)
	case "getTransactionByID":
		return cc.getTransactionByID(stub, args)
	default:
		return nil, fmt.Errorf("unexpected function:%s", function)
	}
}

func (cc *SimpleChainCode) createCompany(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("incorrect number of arguments. expecting 4")
	}
	fundInitial, _ := strconv.Atoi(args[2])
	cType, _ := strconv.Atoi(args[3])
	cid := nextCID(stub)
	cp := Campany{CID: cid, Name: args[0], ContractHash: args[1], FundInitial: fundInitial, FundBalance: fundInitial, CType: cType, CreateTime: time.Now().Nanosecond()}
	cpByte, err := json.Marshal(&cp)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal:%v", err)
	}
	key := "c_" + strconv.Itoa(cid)
	err = stub.PutState(key, cpByte)
	if err != nil {
		return nil, fmt.Errorf("write error:%v", err)
	}
	return cpByte, nil
}

func (cc *SimpleChainCode) getCompanyByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments, expecting 1")
	}
	cid := args[0]
	if cid == "" {
		return nil, errors.New("expecting non-empty args")
	}
	key := "c_" + cid
	cByte, err := stub.GetState(key)
	if err != nil {
		return nil, fmt.Errorf("fail to get company state:%v", err)
	}
	return cByte, nil
}

func (cc *SimpleChainCode) createUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("incorrect number of arguments. expecting 2")
	}
	uid := nextUID(stub)
	user := User{UID: uid, Name: args[0], PWD: args[1], CreateTime: time.Now().Nanosecond()}
	userByte, err := json.Marshal(&user)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal:%v", err)
	}
	key := "u_" + strconv.Itoa(uid)
	err = stub.PutState(key, userByte)
	if err != nil {
		return nil, fmt.Errorf("write error:%v", err)
	}
	return userByte, nil
}

func (cc *SimpleChainCode) getUserByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments. expecting 1")
	}
	uid := args[0]
	if uid == "" {
		return nil, errors.New("expecting non-empty args")
	}
	key := "u_" + uid
	uByte, err := stub.GetState(key)
	if err != nil {
		return nil, fmt.Errorf("fail to get user state:%v", err)
	}
	return uByte, nil
}

func (cc *SimpleChainCode) getUserCompanyFund(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("incorrect number of arguments. expecting 2")
	}
	cid := args[0]
	uid := args[1]
	key := "uf_" + cid + "_" + uid
	fByte, err := stub.GetState(key)
	if err != nil {
		return nil, fmt.Errorf("fail to get userfund state:%v", err)
	}
	return fByte, nil
}

func (cc *SimpleChainCode) getUserFund(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("incorrect number of arguments. expecting 2")
	}
	cids := strings.Split(args[0], ",")
	uid := args[1]
	fundList := make([]UserFund, 0, len(cids))
	for _, cid := range cids {
		key := "uf_" + cid + "_" + uid
		fByte, err := stub.GetState(key)
		if err != nil {
			fmt.Printf("fail to get user fund state:%v\n", err)
			continue
		}
		var uf UserFund
		err = json.Unmarshal(fByte, &uf)
		if err != nil {
			fmt.Printf("fail to unmarsh fByte:%v\n", err)
			continue
		}
		fundList = append(fundList, uf)
	}
	flByte, err := json.Marshal(fundList)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal userfund list:%v", err)
	}
	return flByte, nil
}

func (cc *SimpleChainCode) registerUserRole(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("incorrect number of arguments. expecting 4")
	}
	uid, _ := strconv.Atoi(args[0])
	cid, _ := strconv.Atoi(args[1])
	cType, _ := strconv.Atoi(args[2])
	role, _ := strconv.Atoi(args[3])
	userRole := UserRole{UID: uid, CID: cid, CType: cType, Role: role}
	urByte, err := json.Marshal(&userRole)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal userRole:%v", err)
	}
	key := "ur_" + args[1] + "_" + args[0]
	err = stub.PutState(key, urByte)
	if err != nil {
		return nil, fmt.Errorf("fail to put user role state:%v", err)
	}
	return urByte, nil
}

func (cc *SimpleChainCode) getUserRoleByCID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("incorrect number of arguments. expecting 2")
	}
	key := "ur_" + args[0] + "_" + args[1]
	urByte, err := stub.GetState(key)
	if err != nil {
		return nil, fmt.Errorf("fail to get user role state:%v", err)
	}
	return urByte, nil
}

func (cc *SimpleChainCode) getUserRole(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("incorrect number of arguments. expecting 2")
	}
	cids := strings.Split(args[0], ",")
	uid := args[1]
	urList := make([]UserRole, 0, len(cids))
	for _, cid := range cids {
		key := "ur_" + cid + "_" + uid
		urByte, err := stub.GetState(key)
		if err != nil {
			fmt.Printf("fail to get user role state:%v", err)
			continue
		}
		var ur UserRole
		err = json.Unmarshal(urByte, &ur)
		if err != nil {
			fmt.Printf("fail to unmarshal user role:%v", err)
			continue
		}
		urList = append(urList, ur)
	}
	urlByte, err := json.Marshal(urList)
	if err != nil {
		return nil, err
	}
	return urlByte, nil
}

func (cc *SimpleChainCode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("incorrect number of arguments. expecting 4")
	}

	fundAmount, _ := strconv.Atoi(args[3])
	//get company fund balance
	cKey := "c_" + args[0]
	cByte, err := stub.GetState(cKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get company state:%v", err)
	}
	var cp Campany
	err = json.Unmarshal(cByte, &cp)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal company info:%v", err)
	}
	if cp.FundBalance < fundAmount {
		return nil, fmt.Errorf("not enough fund balance:%d", cp.FundBalance)
	}
	//update company's fund balance
	cp.FundBalance = cp.FundBalance - fundAmount
	newCByte, err := json.Marshal(cp)
	err = stub.PutState(cKey, newCByte)
	//create transaction record
	cid, _ := strconv.Atoi(args[0])
	fromUID, _ := strconv.Atoi(args[1])
	toUID, _ := strconv.Atoi(args[2])
	txid := nextTXID(stub)
	tx := Transaction{TXID: txid, CID: cid, FromUID: fromUID, ToUID: toUID, FundAmount: fundAmount, Status: 0, CreateTime: time.Now().Nanosecond()}
	txByte, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal transaction data:%v", err)
	}
	err = stub.PutState("tx_"+strconv.Itoa(txid), txByte)
	if err != nil {
		return nil, fmt.Errorf("fail to put transaction state:%v", err)
	}
	return txByte, nil
}

func (cc *SimpleChainCode) getTransactionByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments. expecting 1")
	}
	txKey := "tx_" + args[0]
	txByte, err := stub.GetState(txKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get transaction state:%v", err)
	}
	return txByte, nil
}

func (cc *SimpleChainCode) lawyerAuditTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("incorrect number of arguments. expecting 4")
	}
	txKey := "tx_" + args[0]
	txByte, err := stub.GetState(txKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get transaction state:%v", err)
	}
	var tx Transaction
	err = json.Unmarshal(txByte, &tx)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal transaction data:%v", err)
	}
	if tx.LawyerAuditUID > 0 {
		return nil, errors.New("the transaction has been audited by other lawyer")
	}
	//update transaction state
	lawyerUID, _ := strconv.Atoi(args[1])
	auditResult, _ := strconv.Atoi(args[2])
	tx.LawyerAuditUID = lawyerUID
	tx.LawyerAuditResult = auditResult
	tx.LawyerAuditRemark = args[3]
	tx.LawyerAuditTime = time.Now().Nanosecond()
	tx.updateStatus()
	newTXByte, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal transaction data:%v", err)
	}
	err = stub.PutState(txKey, newTXByte)
	if err != nil {
		return nil, fmt.Errorf("fail to put transaction state:%v", err)
	}

	return newTXByte, nil
}

func (cc *SimpleChainCode) auditorAuditTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("incorrect number of arguments. expecting 4")
	}
	txKey := "tx_" + args[0]
	txByte, err := stub.GetState(txKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get transaction state:%v", err)
	}
	var tx Transaction
	err = json.Unmarshal(txByte, &tx)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal transaction data:%v", err)
	}
	if tx.AuditorAuditUID > 0 {
		return nil, errors.New("the transaction has been audited by other auditor")
	}
	//update transaction state
	auditorUID, _ := strconv.Atoi(args[1])
	auditResult, _ := strconv.Atoi(args[2])
	tx.AuditorAuditUID = auditorUID
	tx.AuditorAuditResult = auditResult
	tx.AuditorAuditRemark = args[3]
	tx.AuditorAuditTime = time.Now().Nanosecond()
	tx.updateStatus()
	newTXByte, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal transaction date:%v", err)
	}
	err = stub.PutState(txKey, newTXByte)
	if err != nil {
		return nil, fmt.Errorf("fail to put transaction state:%v", err)
	}
	return newTXByte, nil
}

func (tx *Transaction) updateStatus() {
	if tx.LawyerAuditResult == 1 && tx.AuditorAuditResult == 1 {
		tx.Status = 1
		return
	}
	if tx.LawyerAuditResult == -1 || tx.AuditorAuditResult == -1 {
		tx.Status = -1
		return
	}
}

func (cc *SimpleChainCode) revertTransaction(stub shim.ChaincodeStubInterface, cid int, fundAmount int) ([]byte, error) {
	cKey := "c_" + strconv.Itoa(cid)
	cByte, err := stub.GetState(cKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get company state:%v", err)
	}
	var cp Campany
	err = json.Unmarshal(cByte, &cp)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal company data:%v", err)
	}
	cp.FundBalance = cp.FundBalance + fundAmount
	newCByte, err := json.Marshal(cp)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal company data:%v", err)
	}
	err = stub.PutState(cKey, newCByte)
	if err != nil {
		return nil, fmt.Errorf("fail to put company state:%v", err)
	}
	return newCByte, nil
}

func (cc *SimpleChainCode) executeTransaction(stub shim.ChaincodeStubInterface, cid, fundAmount, uid int) ([]byte, error) {
	ufKey := "uf_" + strconv.Itoa(cid) + "_" + strconv.Itoa(uid)
	ufByte, err := stub.GetState(ufKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get user fund state:%v", err)
	}
	var userFund UserFund
	if len(ufByte) == 0 {
		userFund = UserFund{CID: cid, UID: uid, Fund: fundAmount, CreateTime: time.Now().Nanosecond(), ModifyTime: time.Now().Nanosecond()}
	} else {
		err = json.Unmarshal(ufByte, &userFund)
		if err != nil {
			return nil, fmt.Errorf("fail to Unmarshal user fund state:%v", err)
		}
		userFund.Fund = userFund.Fund + fundAmount
		userFund.ModifyTime = time.Now().Nanosecond()
	}
	newUFByte, err := json.Marshal(userFund)
	if err != nil {
		return nil, fmt.Errorf("fail to marshal user fund state:%v", err)
	}
	err = stub.PutState(ufKey, newUFByte)
	if err != nil {
		return nil, fmt.Errorf("fail to put user fund state:%v", err)
	}
	return newUFByte, nil
}

func nextCID(stub shim.ChaincodeStubInterface) int {
	cidByte, err := stub.GetState("current_cid")
	var cid int
	if err != nil {
		cid = 0
	} else {
		cid, _ = strconv.Atoi(string(cidByte))
	}
	stub.PutState("current_cid", []byte(strconv.Itoa(cid+1)))
	return cid + 1
}

func nextUID(stub shim.ChaincodeStubInterface) int {
	uidByte, err := stub.GetState("current_uid")
	var uid int
	if err != nil {
		uid = 0
	} else {
		uid, _ = strconv.Atoi(string(uidByte))
	}
	stub.PutState("current_uid", []byte(strconv.Itoa(uid+1)))
	return uid + 1
}

func nextTXID(stub shim.ChaincodeStubInterface) int {
	txidByte, err := stub.GetState("current_txid")
	var txid int
	if err != nil {
		txid = 0
	} else {
		txid, _ = strconv.Atoi(string(txidByte))
	}
	stub.PutState("current_txid", []byte(strconv.Itoa(txid+1)))
	return txid + 1
}

func main() {
	err := shim.Start(new(SimpleChainCode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
