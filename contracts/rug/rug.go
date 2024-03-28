// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rug

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// RugMetaData contains all meta data concerning the Rug contract.
var RugMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"MinterRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"MINTED_AMOUNT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isMinter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"mintTokens\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceMinter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052674563918244f400006007553480156200001d57600080fd5b50604051620014df380380620014df833981810160405260408110156200004357600080fd5b81019080805160405193929190846401000000008211156200006457600080fd5b9083019060208201858111156200007a57600080fd5b82516401000000008111828201881017156200009557600080fd5b82525081516020918201929091019080838360005b83811015620000c4578181015183820152602001620000aa565b50505050905090810190601f168015620000f25780820380516001836020036101000a031916815260200191505b50604052602001805160405193929190846401000000008211156200011657600080fd5b9083019060208201858111156200012c57600080fd5b82516401000000008111828201881017156200014757600080fd5b82525081516020918201929091019080838360005b83811015620001765781810151838201526020016200015c565b50505050905090810190601f168015620001a45780820380516001836020036101000a031916815260200191505b50604052505050818160128260009080519060200190620001c7929190620004ed565b508151620001dd906001906020850190620004ed565b506002805460ff191660ff92909216919091179055506200021290506200020362000239565b6001600160e01b036200023e16565b620002313369d3c21bcecceda10000006001600160e01b036200029016565b50506200058f565b335b90565b620002598160066200039560201b62000c8d1790919060201c565b6040516001600160a01b038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b6001600160a01b038216620002ec576040805162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b62000308816005546200042260201b62000a431790919060201c565b6005556001600160a01b0382166000908152600360209081526040909120546200033d91839062000a4362000422821b17901c565b6001600160a01b03831660008181526003602090815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b620003aa82826001600160e01b036200048416565b15620003fd576040805162461bcd60e51b815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c6500604482015290519081900360640190fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b6000828201838110156200047d576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b60006001600160a01b038216620004cd5760405162461bcd60e51b8152600401808060200182810382526022815260200180620014bd6022913960400191505060405180910390fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200053057805160ff191683800117855562000560565b8280016001018555821562000560579182015b828111156200056057825182559160200191906001019062000543565b506200056e92915062000572565b5090565b6200023b91905b808211156200056e576000815560010162000579565b610f1e806200059f6000396000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c806395d89b41116100a2578063a9059cbb11610071578063a9059cbb1461031d578063aa271e1a14610349578063abccf2301461036f578063dd62ed3e14610377578063eeb9635c146103a55761010b565b806395d89b41146102b9578063983b2d56146102c157806398650275146102e9578063a457c2d7146102f15761010b565b8063313ce567116100de578063313ce5671461021d578063395093511461023b57806340c10f191461026757806370a08231146102935761010b565b806306fdde0314610110578063095ea7b31461018d57806318160ddd146101cd57806323b872dd146101e7575b600080fd5b6101186103ad565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561015257818101518382015260200161013a565b50505050905090810190601f16801561017f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101b9600480360360408110156101a357600080fd5b506001600160a01b038135169060200135610443565b604080519115158252519081900360200190f35b6101d5610460565b60408051918252519081900360200190f35b6101b9600480360360608110156101fd57600080fd5b506001600160a01b03813581169160208101359091169060400135610466565b6102256104f3565b6040805160ff9092168252519081900360200190f35b6101b96004803603604081101561025157600080fd5b506001600160a01b0381351690602001356104fc565b6101b96004803603604081101561027d57600080fd5b506001600160a01b038135169060200135610550565b6101d5600480360360208110156102a957600080fd5b50356001600160a01b03166105a7565b6101186105c2565b6102e7600480360360208110156102d757600080fd5b50356001600160a01b0316610622565b005b6102e7610674565b6101b96004803603604081101561030757600080fd5b506001600160a01b038135169060200135610686565b6101b96004803603604081101561033357600080fd5b506001600160a01b0381351690602001356106f4565b6101b96004803603602081101561035f57600080fd5b50356001600160a01b0316610708565b6101d5610721565b6101d56004803603604081101561038d57600080fd5b506001600160a01b0381358116916020013516610727565b6102e7610752565b60008054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156104395780601f1061040e57610100808354040283529160200191610439565b820191906000526020600020905b81548152906001019060200180831161041c57829003601f168201915b5050505050905090565b600061045761045061075e565b8484610762565b50600192915050565b60055490565b600061047384848461084e565b6104e98461047f61075e565b6104e485604051806060016040528060288152602001610e32602891396001600160a01b038a166000908152600460205260408120906104bd61075e565b6001600160a01b03168152602081019190915260400160002054919063ffffffff6109ac16565b610762565b5060019392505050565b60025460ff1690565b600061045761050961075e565b846104e4856004600061051a61075e565b6001600160a01b03908116825260208083019390935260409182016000908120918c16815292529020549063ffffffff610a4316565b600061056261055d61075e565b610708565b61059d5760405162461bcd60e51b8152600401808060200182810382526030815260200180610de16030913960400191505060405180910390fd5b6104578383610aa4565b6001600160a01b031660009081526003602052604090205490565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156104395780601f1061040e57610100808354040283529160200191610439565b61062d61055d61075e565b6106685760405162461bcd60e51b8152600401808060200182810382526030815260200180610de16030913960400191505060405180910390fd5b61067181610b96565b50565b61068461067f61075e565b610bde565b565b600061045761069361075e565b846104e485604051806060016040528060258152602001610ec560259139600460006106bd61075e565b6001600160a01b03908116825260208083019390935260409182016000908120918d1681529252902054919063ffffffff6109ac16565b600061045761070161075e565b848461084e565b600061071b60068363ffffffff610c2616565b92915050565b60075481565b6001600160a01b03918216600090815260046020908152604080832093909416825291909152205490565b61068433600754610aa4565b3390565b6001600160a01b0383166107a75760405162461bcd60e51b8152600401808060200182810382526024815260200180610ea16024913960400191505060405180910390fd5b6001600160a01b0382166107ec5760405162461bcd60e51b8152600401808060200182810382526022815260200180610d996022913960400191505060405180910390fd5b6001600160a01b03808416600081815260046020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b0383166108935760405162461bcd60e51b8152600401808060200182810382526025815260200180610e7c6025913960400191505060405180910390fd5b6001600160a01b0382166108d85760405162461bcd60e51b8152600401808060200182810382526023815260200180610d766023913960400191505060405180910390fd5b61091b81604051806060016040528060268152602001610dbb602691396001600160a01b038616600090815260036020526040902054919063ffffffff6109ac16565b6001600160a01b038085166000908152600360205260408082209390935590841681522054610950908263ffffffff610a4316565b6001600160a01b0380841660008181526003602090815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008184841115610a3b5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610a005781810151838201526020016109e8565b50505050905090810190601f168015610a2d5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b505050900390565b600082820183811015610a9d576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b6001600160a01b038216610aff576040805162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b600554610b12908263ffffffff610a4316565b6005556001600160a01b038216600090815260036020526040902054610b3e908263ffffffff610a4316565b6001600160a01b03831660008181526003602090815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b610ba760068263ffffffff610c8d16565b6040516001600160a01b038216907f6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f690600090a250565b610bef60068263ffffffff610d0e16565b6040516001600160a01b038216907fe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb6669290600090a250565b60006001600160a01b038216610c6d5760405162461bcd60e51b8152600401808060200182810382526022815260200180610e5a6022913960400191505060405180910390fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b610c978282610c26565b15610ce9576040805162461bcd60e51b815260206004820152601f60248201527f526f6c65733a206163636f756e7420616c72656164792068617320726f6c6500604482015290519081900360640190fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b610d188282610c26565b610d535760405162461bcd60e51b8152600401808060200182810382526021815260200180610e116021913960400191505060405180910390fd5b6001600160a01b0316600090815260209190915260409020805460ff1916905556fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e63654d696e746572526f6c653a2063616c6c657220646f6573206e6f74206861766520746865204d696e74657220726f6c65526f6c65733a206163636f756e7420646f6573206e6f74206861766520726f6c6545524332303a207472616e7366657220616d6f756e74206578636565647320616c6c6f77616e6365526f6c65733a206163636f756e7420697320746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f206164647265737345524332303a2064656372656173656420616c6c6f77616e63652062656c6f77207a65726fa265627a7a72315820377b7b359996a017d76f769f3db958f50e1dd364f239d2d8ad04eed4e75f9fd464736f6c63430005100032526f6c65733a206163636f756e7420697320746865207a65726f2061646472657373",
}

// RugABI is the input ABI used to generate the binding from.
// Deprecated: Use RugMetaData.ABI instead.
var RugABI = RugMetaData.ABI

// RugBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RugMetaData.Bin instead.
var RugBin = RugMetaData.Bin

// DeployRug deploys a new Ethereum contract, binding an instance of Rug to it.
func DeployRug(auth *bind.TransactOpts, backend bind.ContractBackend, _name string, _symbol string) (common.Address, *types.Transaction, *Rug, error) {
	parsed, err := RugMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RugBin), backend, _name, _symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Rug{RugCaller: RugCaller{contract: contract}, RugTransactor: RugTransactor{contract: contract}, RugFilterer: RugFilterer{contract: contract}}, nil
}

// Rug is an auto generated Go binding around an Ethereum contract.
type Rug struct {
	RugCaller     // Read-only binding to the contract
	RugTransactor // Write-only binding to the contract
	RugFilterer   // Log filterer for contract events
}

// RugCaller is an auto generated read-only Go binding around an Ethereum contract.
type RugCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RugTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RugTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RugFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RugFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RugSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RugSession struct {
	Contract     *Rug              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RugCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RugCallerSession struct {
	Contract *RugCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RugTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RugTransactorSession struct {
	Contract     *RugTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RugRaw is an auto generated low-level Go binding around an Ethereum contract.
type RugRaw struct {
	Contract *Rug // Generic contract binding to access the raw methods on
}

// RugCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RugCallerRaw struct {
	Contract *RugCaller // Generic read-only contract binding to access the raw methods on
}

// RugTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RugTransactorRaw struct {
	Contract *RugTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRug creates a new instance of Rug, bound to a specific deployed contract.
func NewRug(address common.Address, backend bind.ContractBackend) (*Rug, error) {
	contract, err := bindRug(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rug{RugCaller: RugCaller{contract: contract}, RugTransactor: RugTransactor{contract: contract}, RugFilterer: RugFilterer{contract: contract}}, nil
}

// NewRugCaller creates a new read-only instance of Rug, bound to a specific deployed contract.
func NewRugCaller(address common.Address, caller bind.ContractCaller) (*RugCaller, error) {
	contract, err := bindRug(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RugCaller{contract: contract}, nil
}

// NewRugTransactor creates a new write-only instance of Rug, bound to a specific deployed contract.
func NewRugTransactor(address common.Address, transactor bind.ContractTransactor) (*RugTransactor, error) {
	contract, err := bindRug(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RugTransactor{contract: contract}, nil
}

// NewRugFilterer creates a new log filterer instance of Rug, bound to a specific deployed contract.
func NewRugFilterer(address common.Address, filterer bind.ContractFilterer) (*RugFilterer, error) {
	contract, err := bindRug(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RugFilterer{contract: contract}, nil
}

// bindRug binds a generic wrapper to an already deployed contract.
func bindRug(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RugMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rug *RugRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rug.Contract.RugCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rug *RugRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rug.Contract.RugTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rug *RugRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rug.Contract.RugTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rug *RugCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rug.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rug *RugTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rug.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rug *RugTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rug.Contract.contract.Transact(opts, method, params...)
}

// MINTEDAMOUNT is a free data retrieval call binding the contract method 0xabccf230.
//
// Solidity: function MINTED_AMOUNT() view returns(uint256)
func (_Rug *RugCaller) MINTEDAMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rug.contract.Call(opts, &out, "MINTED_AMOUNT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINTEDAMOUNT is a free data retrieval call binding the contract method 0xabccf230.
//
// Solidity: function MINTED_AMOUNT() view returns(uint256)
func (_Rug *RugSession) MINTEDAMOUNT() (*big.Int, error) {
	return _Rug.Contract.MINTEDAMOUNT(&_Rug.CallOpts)
}

// MINTEDAMOUNT is a free data retrieval call binding the contract method 0xabccf230.
//
// Solidity: function MINTED_AMOUNT() view returns(uint256)
func (_Rug *RugCallerSession) MINTEDAMOUNT() (*big.Int, error) {
	return _Rug.Contract.MINTEDAMOUNT(&_Rug.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Rug *RugCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Rug.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Rug *RugSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Rug.Contract.Allowance(&_Rug.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Rug *RugCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Rug.Contract.Allowance(&_Rug.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Rug *RugCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Rug.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Rug *RugSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Rug.Contract.BalanceOf(&_Rug.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Rug *RugCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Rug.Contract.BalanceOf(&_Rug.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Rug *RugCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Rug.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Rug *RugSession) Decimals() (uint8, error) {
	return _Rug.Contract.Decimals(&_Rug.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Rug *RugCallerSession) Decimals() (uint8, error) {
	return _Rug.Contract.Decimals(&_Rug.CallOpts)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_Rug *RugCaller) IsMinter(opts *bind.CallOpts, account common.Address) (bool, error) {
	var out []interface{}
	err := _Rug.contract.Call(opts, &out, "isMinter", account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_Rug *RugSession) IsMinter(account common.Address) (bool, error) {
	return _Rug.Contract.IsMinter(&_Rug.CallOpts, account)
}

// IsMinter is a free data retrieval call binding the contract method 0xaa271e1a.
//
// Solidity: function isMinter(address account) view returns(bool)
func (_Rug *RugCallerSession) IsMinter(account common.Address) (bool, error) {
	return _Rug.Contract.IsMinter(&_Rug.CallOpts, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Rug *RugCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Rug.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Rug *RugSession) Name() (string, error) {
	return _Rug.Contract.Name(&_Rug.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Rug *RugCallerSession) Name() (string, error) {
	return _Rug.Contract.Name(&_Rug.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Rug *RugCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Rug.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Rug *RugSession) Symbol() (string, error) {
	return _Rug.Contract.Symbol(&_Rug.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Rug *RugCallerSession) Symbol() (string, error) {
	return _Rug.Contract.Symbol(&_Rug.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Rug *RugCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rug.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Rug *RugSession) TotalSupply() (*big.Int, error) {
	return _Rug.Contract.TotalSupply(&_Rug.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Rug *RugCallerSession) TotalSupply() (*big.Int, error) {
	return _Rug.Contract.TotalSupply(&_Rug.CallOpts)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_Rug *RugTransactor) AddMinter(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "addMinter", account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_Rug *RugSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _Rug.Contract.AddMinter(&_Rug.TransactOpts, account)
}

// AddMinter is a paid mutator transaction binding the contract method 0x983b2d56.
//
// Solidity: function addMinter(address account) returns()
func (_Rug *RugTransactorSession) AddMinter(account common.Address) (*types.Transaction, error) {
	return _Rug.Contract.AddMinter(&_Rug.TransactOpts, account)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Rug *RugTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Rug *RugSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.Approve(&_Rug.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Rug *RugTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.Approve(&_Rug.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Rug *RugTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Rug *RugSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.DecreaseAllowance(&_Rug.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Rug *RugTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.DecreaseAllowance(&_Rug.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Rug *RugTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Rug *RugSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.IncreaseAllowance(&_Rug.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Rug *RugTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.IncreaseAllowance(&_Rug.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_Rug *RugTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_Rug *RugSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.Mint(&_Rug.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_Rug *RugTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.Mint(&_Rug.TransactOpts, account, amount)
}

// MintTokens is a paid mutator transaction binding the contract method 0xeeb9635c.
//
// Solidity: function mintTokens() returns()
func (_Rug *RugTransactor) MintTokens(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "mintTokens")
}

// MintTokens is a paid mutator transaction binding the contract method 0xeeb9635c.
//
// Solidity: function mintTokens() returns()
func (_Rug *RugSession) MintTokens() (*types.Transaction, error) {
	return _Rug.Contract.MintTokens(&_Rug.TransactOpts)
}

// MintTokens is a paid mutator transaction binding the contract method 0xeeb9635c.
//
// Solidity: function mintTokens() returns()
func (_Rug *RugTransactorSession) MintTokens() (*types.Transaction, error) {
	return _Rug.Contract.MintTokens(&_Rug.TransactOpts)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_Rug *RugTransactor) RenounceMinter(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "renounceMinter")
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_Rug *RugSession) RenounceMinter() (*types.Transaction, error) {
	return _Rug.Contract.RenounceMinter(&_Rug.TransactOpts)
}

// RenounceMinter is a paid mutator transaction binding the contract method 0x98650275.
//
// Solidity: function renounceMinter() returns()
func (_Rug *RugTransactorSession) RenounceMinter() (*types.Transaction, error) {
	return _Rug.Contract.RenounceMinter(&_Rug.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Rug *RugTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Rug *RugSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.Transfer(&_Rug.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Rug *RugTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.Transfer(&_Rug.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Rug *RugTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Rug *RugSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.TransferFrom(&_Rug.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Rug *RugTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rug.Contract.TransferFrom(&_Rug.TransactOpts, sender, recipient, amount)
}

// RugApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Rug contract.
type RugApprovalIterator struct {
	Event *RugApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RugApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RugApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RugApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RugApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RugApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RugApproval represents a Approval event raised by the Rug contract.
type RugApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Rug *RugFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*RugApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Rug.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &RugApprovalIterator{contract: _Rug.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Rug *RugFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *RugApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Rug.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RugApproval)
				if err := _Rug.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Rug *RugFilterer) ParseApproval(log types.Log) (*RugApproval, error) {
	event := new(RugApproval)
	if err := _Rug.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RugMinterAddedIterator is returned from FilterMinterAdded and is used to iterate over the raw logs and unpacked data for MinterAdded events raised by the Rug contract.
type RugMinterAddedIterator struct {
	Event *RugMinterAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RugMinterAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RugMinterAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RugMinterAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RugMinterAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RugMinterAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RugMinterAdded represents a MinterAdded event raised by the Rug contract.
type RugMinterAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterAdded is a free log retrieval operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_Rug *RugFilterer) FilterMinterAdded(opts *bind.FilterOpts, account []common.Address) (*RugMinterAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Rug.contract.FilterLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &RugMinterAddedIterator{contract: _Rug.contract, event: "MinterAdded", logs: logs, sub: sub}, nil
}

// WatchMinterAdded is a free log subscription operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_Rug *RugFilterer) WatchMinterAdded(opts *bind.WatchOpts, sink chan<- *RugMinterAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Rug.contract.WatchLogs(opts, "MinterAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RugMinterAdded)
				if err := _Rug.contract.UnpackLog(event, "MinterAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMinterAdded is a log parse operation binding the contract event 0x6ae172837ea30b801fbfcdd4108aa1d5bf8ff775444fd70256b44e6bf3dfc3f6.
//
// Solidity: event MinterAdded(address indexed account)
func (_Rug *RugFilterer) ParseMinterAdded(log types.Log) (*RugMinterAdded, error) {
	event := new(RugMinterAdded)
	if err := _Rug.contract.UnpackLog(event, "MinterAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RugMinterRemovedIterator is returned from FilterMinterRemoved and is used to iterate over the raw logs and unpacked data for MinterRemoved events raised by the Rug contract.
type RugMinterRemovedIterator struct {
	Event *RugMinterRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RugMinterRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RugMinterRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RugMinterRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RugMinterRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RugMinterRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RugMinterRemoved represents a MinterRemoved event raised by the Rug contract.
type RugMinterRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMinterRemoved is a free log retrieval operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_Rug *RugFilterer) FilterMinterRemoved(opts *bind.FilterOpts, account []common.Address) (*RugMinterRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Rug.contract.FilterLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &RugMinterRemovedIterator{contract: _Rug.contract, event: "MinterRemoved", logs: logs, sub: sub}, nil
}

// WatchMinterRemoved is a free log subscription operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_Rug *RugFilterer) WatchMinterRemoved(opts *bind.WatchOpts, sink chan<- *RugMinterRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Rug.contract.WatchLogs(opts, "MinterRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RugMinterRemoved)
				if err := _Rug.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMinterRemoved is a log parse operation binding the contract event 0xe94479a9f7e1952cc78f2d6baab678adc1b772d936c6583def489e524cb66692.
//
// Solidity: event MinterRemoved(address indexed account)
func (_Rug *RugFilterer) ParseMinterRemoved(log types.Log) (*RugMinterRemoved, error) {
	event := new(RugMinterRemoved)
	if err := _Rug.contract.UnpackLog(event, "MinterRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RugTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Rug contract.
type RugTransferIterator struct {
	Event *RugTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RugTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RugTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(RugTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *RugTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RugTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RugTransfer represents a Transfer event raised by the Rug contract.
type RugTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Rug *RugFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RugTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Rug.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RugTransferIterator{contract: _Rug.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Rug *RugFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *RugTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Rug.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RugTransfer)
				if err := _Rug.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Rug *RugFilterer) ParseTransfer(log types.Log) (*RugTransfer, error) {
	event := new(RugTransfer)
	if err := _Rug.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
