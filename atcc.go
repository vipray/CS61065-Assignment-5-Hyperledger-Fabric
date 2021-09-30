/******************************* Step1 *****************************
Import dependencies and Define smart Contract
*/

package main

import (
  "fmt"
  "log"
  "encoding/json"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an House
type SmartContract struct {
	contractapi.Contract
}

// House describes basic details of what makes up a simple House
type House struct {
	NagarPalikaID    string `json:"NagarPalikaID"`
	Owner          string `json:"Owner"`
	Address          string `json:"Address"`
	Size           int    `json:"Size"`
}


// CreateHouse issues a new House to the world state with given details.
func (s *SmartContract) CreateHouse(ctx contractapi.TransactionContextInterface, id string, address string, size int) error {
	exists, err := s.HouseExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the House %s already exists", id)
	}
	
	owner, ownerErr := ctx.GetClientIdentity().GetMSPID()
	if ownerErr != nil {
		return ownerErr
	}

	House := House{
		NagarPalikaID:  id,
		Address:        address,
		Size:           size,
		Owner:          owner,
	}
	HouseJSON, err := json.Marshal(House)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, HouseJSON)
}


// ReadHouse returns the House stored in the world state with given id.
func (s *SmartContract) ReadHouse(ctx contractapi.TransactionContextInterface, id string) (*House, error) {
	HouseJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if HouseJSON == nil {
		return nil, fmt.Errorf("the House %s does not exist", id)
	}

	var House House
	err = json.Unmarshal(HouseJSON, &House)
	if err != nil {
		return nil, err
	}

	return &House, nil
}


// HouseExists returns true when House with given ID exists in world state
func (s *SmartContract) HouseExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	HouseJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return HouseJSON != nil, nil
}



// TransferHouse updates the owner field of House with given id in world state.
func (s *SmartContract) TransferHouse(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	House, err := s.ReadHouse(ctx, id)
	if err != nil {
		return err
	}
	
	owner, ownerErr := ctx.GetClientIdentity().GetMSPID()
	if ownerErr != nil {
		return ownerErr
	}
	
	if(House.Owner != owner){
		return fmt.Errorf("failed to Transfer Ownership: You are not the owner of this house.");
	}
	
	House.Owner = newOwner
	HouseJSON, err := json.Marshal(House)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, HouseJSON)
}

// GetAllHouses returns all Houses found in world state
func (s *SmartContract) GetAllHouses(ctx contractapi.TransactionContextInterface) ([]*House, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all Houses in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var Houses []*House
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var House House
		err = json.Unmarshal(queryResponse.Value, &House)
		if err != nil {
			return nil, err
		}
		Houses = append(Houses, &House)
	}

	return Houses, nil
}

// The main function which will create the chaincode and start it
func main(){
	// NewChaincode creates a new chaincode using contracts passed.
    HouseChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
      log.Panicf("Error creating House-transfer-basic chaincode: %v", err)
    }

	// Start starts the chaincode in the fabric
    if err := HouseChaincode.Start(); err != nil {
      log.Panicf("Error starting House-transfer-basic chaincode: %v", err)
    }
}



