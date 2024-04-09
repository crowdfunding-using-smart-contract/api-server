// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package crowdfunding

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ava-labs/coreth/accounts/abi/bind"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = interfaces.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// CrowdfundingProject is an auto generated low-level Go binding around an user-defined struct.
type CrowdfundingProject struct {
	Id             *big.Int
	Title          string
	Description    string
	TargetFunding  *big.Int
	CurrentFunding *big.Int
	StartDate      *big.Int
	EndDate        *big.Int
	Owner          common.Address
	Status         uint8
}

// CrowdfundingMetaData contains all meta data concerning the Crowdfunding contract.
var CrowdfundingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"contributor\",\"type\":\"address\"}],\"name\":\"ContributionMade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"contributor\",\"type\":\"address\"}],\"name\":\"ContributionRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ProjectCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"contribute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"contributions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"targetFunding\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endDate\",\"type\":\"uint256\"}],\"name\":\"createProject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"deleteProject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"targetFunding\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endDate\",\"type\":\"uint256\"}],\"name\":\"editProject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumCrowdfunding.STATUS\",\"name\":\"filterStatus\",\"type\":\"uint8\"}],\"name\":\"getAllProjects\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"getProject\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"targetFunding\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentFunding\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endDate\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"enumCrowdfunding.STATUS\",\"name\":\"status\",\"type\":\"uint8\"}],\"internalType\":\"structCrowdfunding.Project\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"projects\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"targetFunding\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentFunding\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startDate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endDate\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"enumCrowdfunding.STATUS\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b5060405161161738038061161783398101604081905261002e916100f2565b806001600160a01b03811661005c57604051631e4fbdf760e01b81525f600482015260240160405180910390fd5b6100658161008c565b5050600180546001600160a01b0319166001600160a01b039290921691909117905561012a565b5f80546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6001600160a01b03811681146100ef575f80fd5b50565b5f8060408385031215610103575f80fd5b825161010e816100db565b602084015190925061011f816100db565b809150509250929050565b6114e0806101375f395ff3fe608060405234801561000f575f80fd5b50600436106100b1575f3560e01c80638c5909171161006e5780638c5909171461016e5780638da5cb5b14610181578063bd481cf91461019b578063f0f3f2c8146101ae578063f2fde38b146101ce578063f81bf601146101e1575f80fd5b8063107046bd146100b5578063278ecde1146100e657806330e06a01146100fb5780633d891f591461011b5780634cd0cb2814610153578063715018a614610166575b5f80fd5b6100c86100c3366004610eb5565b6101f4565b6040516100dd99989796959493929190610f2e565b60405180910390f35b6100f96100f4366004610eb5565b610357565b005b61010e610109366004610fa0565b61054b565b6040516100dd9190610fc5565b610145610129366004611023565b600460209081525f928352604080842090915290825290205481565b6040519081526020016100dd565b6100f9610161366004610eb5565b6106b0565b6100f961074a565b6100f961017c36600461104d565b61075d565b5f546040516001600160a01b0390911681526020016100dd565b6100f96101a936600461110a565b61099a565b6101c16101bc366004610eb5565b610a7c565b6040516100dd919061118d565b6100f96101dc366004611234565b610c40565b6100f96101ef36600461124d565b610c7d565b60036020525f908152604090208054600182018054919291610215906112c6565b80601f0160208091040260200160405190810160405280929190818152602001828054610241906112c6565b801561028c5780601f106102635761010080835404028352916020019161028c565b820191905f5260205f20905b81548152906001019060200180831161026f57829003601f168201915b5050505050908060020180546102a1906112c6565b80601f01602080910402602001604051908101604052809291908181526020018280546102cd906112c6565b80156103185780601f106102ef57610100808354040283529160200191610318565b820191905f5260205f20905b8154815290600101906020018083116102fb57829003601f168201915b50505060038401546004850154600586015460068701546007909701549596929591945092506001600160a01b0381169060ff600160a01b9091041689565b5f818152600360208190526040909120906007820154600160a01b900460ff16600381111561038857610388610efa565b14806103975750806006015442115b6103de5760405162461bcd60e51b81526020600482015260136024820152721499599d5b991cc81b9bdd08185b1b1bddd959606a1b60448201526064015b60405180910390fd5b5f828152600460209081526040808320338452909152902054806104445760405162461bcd60e51b815260206004820152601a60248201527f4e6f20636f6e747269627574696f6e7320746f20726566756e6400000000000060448201526064016103d5565b5f83815260046020818152604080842033808652925280842093909355600154925163a9059cbb60e01b815291820152602481018390526001600160a01b039091169063a9059cbb906044016020604051808303815f875af11580156104ac573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104d091906112fe565b61050c5760405162461bcd60e51b815260206004820152600d60248201526c1499599d5b990819985a5b1959609a1b60448201526064016103d5565b6040805182815233602082015284917f232bf2213879e79f5b2fa9cf07aad94aa0acac7b3c93bdfe74a545d47a1b0275910160405180910390a2505050565b6002546060905f60015b8281116105c15784600381111561056e5761056e610efa565b5f82815260036020819052604090912060070154600160a01b900460ff169081111561059c5761059c610efa565b036105af57816105ab81611331565b9250505b806105b981611331565b915050610555565b505f8167ffffffffffffffff8111156105dc576105dc61106d565b604051908082528060200260200182016040528015610605578160200160208202803683370190505b5090505f60015b8481116106a55786600381111561062557610625610efa565b5f82815260036020819052604090912060070154600160a01b900460ff169081111561065357610653610efa565b03610693575f81815260036020526040902054835184908490811061067a5761067a611349565b60209081029190910101528161068f81611331565b9250505b8061069d81611331565b91505061060c565b509095945050505050565b5f818152600360205260409020600701546001600160a01b031633146107275760405162461bcd60e51b815260206004820152602660248201527f4f6e6c7920746865206f776e65722063616e2064656c6574652074686520707260448201526537b532b1ba1760d11b60648201526084016103d5565b5f908152600360205260409020600701805460ff60a01b1916600160a01b179055565b610752610ddc565b61075b5f610e08565b565b5f828152600360205260408120906007820154600160a01b900460ff16600381111561078b5761078b610efa565b146107d15760405162461bcd60e51b815260206004820152601660248201527550726f6a656374206d7573742062652061637469766560501b60448201526064016103d5565b806005015442101580156107e9575080600601544211155b6108355760405162461bcd60e51b815260206004820152601d60248201527f50726f6a656374206e6f7420696e2066756e64696e6720706572696f6400000060448201526064016103d5565b6001546040516323b872dd60e01b8152336004820152306024820152604481018490526001600160a01b03909116906323b872dd906064016020604051808303815f875af1158015610889573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906108ad91906112fe565b6108f15760405162461bcd60e51b8152602060048201526015602482015274151bdad95b881d1c985b9cd9995c8819985a5b1959605a1b60448201526064016103d5565b81816004015f828254610904919061135d565b90915550505f8381526004602090815260408083203384529091528120805484929061093190849061135d565b90915550506040805183815233602082015284917f7c50e76f3eff5098fcf915218e219a359f40f4c150e1a6796920587d44923cb4910160405180910390a280600301548160040154106109955760078101805460ff60a01b1916600160a11b1790555b505050565b5f868152600360205260409020600701546001600160a01b03163314610a0e5760405162461bcd60e51b8152602060048201526024808201527f4f6e6c7920746865206f776e65722063616e2065646974207468652070726f6a60448201526332b1ba1760e11b60648201526084016103d5565b5f868152600360205260409020855115610a325760018101610a3087826113c1565b505b845115610a495760028101610a4786826113c1565b505b8315610a5757600381018490555b8215610a6557600581018390555b8115610a7357600681018290555b50505050505050565b610a84610e57565b60035f8381526020019081526020015f20604051806101200160405290815f8201548152602001600182018054610aba906112c6565b80601f0160208091040260200160405190810160405280929190818152602001828054610ae6906112c6565b8015610b315780601f10610b0857610100808354040283529160200191610b31565b820191905f5260205f20905b815481529060010190602001808311610b1457829003601f168201915b50505050508152602001600282018054610b4a906112c6565b80601f0160208091040260200160405190810160405280929190818152602001828054610b76906112c6565b8015610bc15780601f10610b9857610100808354040283529160200191610bc1565b820191905f5260205f20905b815481529060010190602001808311610ba457829003601f168201915b5050509183525050600382810154602083015260048301546040830152600583015460608301526006830154608083015260078301546001600160a01b03811660a084015260c090920191600160a01b900460ff1690811115610c2657610c26610efa565b6003811115610c3757610c37610efa565b90525092915050565b610c48610ddc565b6001600160a01b038116610c7157604051631e4fbdf760e01b81525f60048201526024016103d5565b610c7a81610e08565b50565b60028054905f610c8c83611331565b9091555050600254604080516101208101825282815260208101889052908101869052606081018590525f6080820181905260a0820185905260c082018490523360e08301526101008201525f82815260036020908152604090912082518155908201516001820190610cff90826113c1565b5060408201516002820190610d1490826113c1565b5060608201516003808301919091556080830151600483015560a0830151600583015560c0830151600683015560e08301516007830180546001600160a01b039092166001600160a01b03198316811782556101008601519391926001600160a81b0319161790600160a01b908490811115610d9257610d92610efa565b0217905550905050807fc0c54fed07481d0998e1446b2c13759606bf4f26b78306307413ac4c4309aa828733604051610dcc929190611481565b60405180910390a2505050505050565b5f546001600160a01b0316331461075b5760405163118cdaa760e01b81523360048201526024016103d5565b5f80546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6040518061012001604052805f815260200160608152602001606081526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f6003811115610eb057610eb0610efa565b905290565b5f60208284031215610ec5575f80fd5b5035919050565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b634e487b7160e01b5f52602160045260245ffd5b60048110610f2a57634e487b7160e01b5f52602160045260245ffd5b9052565b5f6101208b8352806020840152610f478184018c610ecc565b90508281036040840152610f5b818b610ecc565b9150508760608301528660808301528560a08301528460c083015260018060a01b03841660e0830152610f92610100830184610f0e565b9a9950505050505050505050565b5f60208284031215610fb0575f80fd5b813560048110610fbe575f80fd5b9392505050565b602080825282518282018190525f9190848201906040850190845b81811015610ffc57835183529284019291840191600101610fe0565b50909695505050505050565b80356001600160a01b038116811461101e575f80fd5b919050565b5f8060408385031215611034575f80fd5b8235915061104460208401611008565b90509250929050565b5f806040838503121561105e575f80fd5b50508035926020909101359150565b634e487b7160e01b5f52604160045260245ffd5b5f82601f830112611090575f80fd5b813567ffffffffffffffff808211156110ab576110ab61106d565b604051601f8301601f19908116603f011681019082821181831017156110d3576110d361106d565b816040528381528660208588010111156110eb575f80fd5b836020870160208301375f602085830101528094505050505092915050565b5f805f805f8060c0878903121561111f575f80fd5b86359550602087013567ffffffffffffffff8082111561113d575f80fd5b6111498a838b01611081565b9650604089013591508082111561115e575f80fd5b5061116b89828a01611081565b945050606087013592506080870135915060a087013590509295509295509295565b60208152815160208201525f60208301516101208060408501526111b5610140850183610ecc565b91506040850151601f198584030160608601526111d28382610ecc565b92505060608501516080850152608085015160a085015260a085015160c085015260c085015160e085015260e0850151610100611219818701836001600160a01b03169052565b860151905061122a85830182610f0e565b5090949350505050565b5f60208284031215611244575f80fd5b610fbe82611008565b5f805f805f60a08688031215611261575f80fd5b853567ffffffffffffffff80821115611278575f80fd5b61128489838a01611081565b96506020880135915080821115611299575f80fd5b506112a688828901611081565b959895975050505060408401359360608101359360809091013592509050565b600181811c908216806112da57607f821691505b6020821081036112f857634e487b7160e01b5f52602260045260245ffd5b50919050565b5f6020828403121561130e575f80fd5b81518015158114610fbe575f80fd5b634e487b7160e01b5f52601160045260245ffd5b5f600182016113425761134261131d565b5060010190565b634e487b7160e01b5f52603260045260245ffd5b808201808211156113705761137061131d565b92915050565b601f82111561099557805f5260205f20601f840160051c8101602085101561139b5750805b601f840160051c820191505b818110156113ba575f81556001016113a7565b5050505050565b815167ffffffffffffffff8111156113db576113db61106d565b6113ef816113e984546112c6565b84611376565b602080601f831160018114611422575f841561140b5750858301515b5f19600386901b1c1916600185901b178555611479565b5f85815260208120601f198616915b8281101561145057888601518255948401946001909101908401611431565b508582101561146d57878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b604081525f6114936040830185610ecc565b905060018060a01b0383166020830152939250505056fea264697066735822122091d824d0b71c716c00c01d52d463c5af4c09d1ca1f9840c3a1026638bc27ab6b64736f6c63430008190033",
}

// CrowdfundingABI is the input ABI used to generate the binding from.
// Deprecated: Use CrowdfundingMetaData.ABI instead.
var CrowdfundingABI = CrowdfundingMetaData.ABI

// CrowdfundingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrowdfundingMetaData.Bin instead.
var CrowdfundingBin = CrowdfundingMetaData.Bin

// DeployCrowdfunding deploys a new Ethereum contract, binding an instance of Crowdfunding to it.
func DeployCrowdfunding(auth *bind.TransactOpts, backend bind.ContractBackend, tokenAddress common.Address, initialOwner common.Address) (common.Address, *types.Transaction, *Crowdfunding, error) {
	parsed, err := CrowdfundingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrowdfundingBin), backend, tokenAddress, initialOwner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Crowdfunding{CrowdfundingCaller: CrowdfundingCaller{contract: contract}, CrowdfundingTransactor: CrowdfundingTransactor{contract: contract}, CrowdfundingFilterer: CrowdfundingFilterer{contract: contract}}, nil
}

// Crowdfunding is an auto generated Go binding around an Ethereum contract.
type Crowdfunding struct {
	CrowdfundingCaller     // Read-only binding to the contract
	CrowdfundingTransactor // Write-only binding to the contract
	CrowdfundingFilterer   // Log filterer for contract events
}

// CrowdfundingCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrowdfundingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdfundingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrowdfundingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdfundingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrowdfundingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdfundingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrowdfundingSession struct {
	Contract     *Crowdfunding     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrowdfundingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrowdfundingCallerSession struct {
	Contract *CrowdfundingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// CrowdfundingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrowdfundingTransactorSession struct {
	Contract     *CrowdfundingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CrowdfundingRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrowdfundingRaw struct {
	Contract *Crowdfunding // Generic contract binding to access the raw methods on
}

// CrowdfundingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrowdfundingCallerRaw struct {
	Contract *CrowdfundingCaller // Generic read-only contract binding to access the raw methods on
}

// CrowdfundingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrowdfundingTransactorRaw struct {
	Contract *CrowdfundingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrowdfunding creates a new instance of Crowdfunding, bound to a specific deployed contract.
func NewCrowdfunding(address common.Address, backend bind.ContractBackend) (*Crowdfunding, error) {
	contract, err := bindCrowdfunding(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Crowdfunding{CrowdfundingCaller: CrowdfundingCaller{contract: contract}, CrowdfundingTransactor: CrowdfundingTransactor{contract: contract}, CrowdfundingFilterer: CrowdfundingFilterer{contract: contract}}, nil
}

// NewCrowdfundingCaller creates a new read-only instance of Crowdfunding, bound to a specific deployed contract.
func NewCrowdfundingCaller(address common.Address, caller bind.ContractCaller) (*CrowdfundingCaller, error) {
	contract, err := bindCrowdfunding(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrowdfundingCaller{contract: contract}, nil
}

// NewCrowdfundingTransactor creates a new write-only instance of Crowdfunding, bound to a specific deployed contract.
func NewCrowdfundingTransactor(address common.Address, transactor bind.ContractTransactor) (*CrowdfundingTransactor, error) {
	contract, err := bindCrowdfunding(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrowdfundingTransactor{contract: contract}, nil
}

// NewCrowdfundingFilterer creates a new log filterer instance of Crowdfunding, bound to a specific deployed contract.
func NewCrowdfundingFilterer(address common.Address, filterer bind.ContractFilterer) (*CrowdfundingFilterer, error) {
	contract, err := bindCrowdfunding(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrowdfundingFilterer{contract: contract}, nil
}

// bindCrowdfunding binds a generic wrapper to an already deployed contract.
func bindCrowdfunding(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrowdfundingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Crowdfunding *CrowdfundingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Crowdfunding.Contract.CrowdfundingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Crowdfunding *CrowdfundingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Crowdfunding.Contract.CrowdfundingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Crowdfunding *CrowdfundingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Crowdfunding.Contract.CrowdfundingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Crowdfunding *CrowdfundingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Crowdfunding.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Crowdfunding *CrowdfundingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Crowdfunding.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Crowdfunding *CrowdfundingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Crowdfunding.Contract.contract.Transact(opts, method, params...)
}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_Crowdfunding *CrowdfundingCaller) Contributions(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Crowdfunding.contract.Call(opts, &out, "contributions", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_Crowdfunding *CrowdfundingSession) Contributions(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _Crowdfunding.Contract.Contributions(&_Crowdfunding.CallOpts, arg0, arg1)
}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_Crowdfunding *CrowdfundingCallerSession) Contributions(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _Crowdfunding.Contract.Contributions(&_Crowdfunding.CallOpts, arg0, arg1)
}

// GetAllProjects is a free data retrieval call binding the contract method 0x30e06a01.
//
// Solidity: function getAllProjects(uint8 filterStatus) view returns(uint256[])
func (_Crowdfunding *CrowdfundingCaller) GetAllProjects(opts *bind.CallOpts, filterStatus uint8) ([]*big.Int, error) {
	var out []interface{}
	err := _Crowdfunding.contract.Call(opts, &out, "getAllProjects", filterStatus)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAllProjects is a free data retrieval call binding the contract method 0x30e06a01.
//
// Solidity: function getAllProjects(uint8 filterStatus) view returns(uint256[])
func (_Crowdfunding *CrowdfundingSession) GetAllProjects(filterStatus uint8) ([]*big.Int, error) {
	return _Crowdfunding.Contract.GetAllProjects(&_Crowdfunding.CallOpts, filterStatus)
}

// GetAllProjects is a free data retrieval call binding the contract method 0x30e06a01.
//
// Solidity: function getAllProjects(uint8 filterStatus) view returns(uint256[])
func (_Crowdfunding *CrowdfundingCallerSession) GetAllProjects(filterStatus uint8) ([]*big.Int, error) {
	return _Crowdfunding.Contract.GetAllProjects(&_Crowdfunding.CallOpts, filterStatus)
}

// GetProject is a free data retrieval call binding the contract method 0xf0f3f2c8.
//
// Solidity: function getProject(uint256 projectId) view returns((uint256,string,string,uint256,uint256,uint256,uint256,address,uint8))
func (_Crowdfunding *CrowdfundingCaller) GetProject(opts *bind.CallOpts, projectId *big.Int) (CrowdfundingProject, error) {
	var out []interface{}
	err := _Crowdfunding.contract.Call(opts, &out, "getProject", projectId)

	if err != nil {
		return *new(CrowdfundingProject), err
	}

	out0 := *abi.ConvertType(out[0], new(CrowdfundingProject)).(*CrowdfundingProject)

	return out0, err

}

// GetProject is a free data retrieval call binding the contract method 0xf0f3f2c8.
//
// Solidity: function getProject(uint256 projectId) view returns((uint256,string,string,uint256,uint256,uint256,uint256,address,uint8))
func (_Crowdfunding *CrowdfundingSession) GetProject(projectId *big.Int) (CrowdfundingProject, error) {
	return _Crowdfunding.Contract.GetProject(&_Crowdfunding.CallOpts, projectId)
}

// GetProject is a free data retrieval call binding the contract method 0xf0f3f2c8.
//
// Solidity: function getProject(uint256 projectId) view returns((uint256,string,string,uint256,uint256,uint256,uint256,address,uint8))
func (_Crowdfunding *CrowdfundingCallerSession) GetProject(projectId *big.Int) (CrowdfundingProject, error) {
	return _Crowdfunding.Contract.GetProject(&_Crowdfunding.CallOpts, projectId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Crowdfunding *CrowdfundingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Crowdfunding.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Crowdfunding *CrowdfundingSession) Owner() (common.Address, error) {
	return _Crowdfunding.Contract.Owner(&_Crowdfunding.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Crowdfunding *CrowdfundingCallerSession) Owner() (common.Address, error) {
	return _Crowdfunding.Contract.Owner(&_Crowdfunding.CallOpts)
}

// Projects is a free data retrieval call binding the contract method 0x107046bd.
//
// Solidity: function projects(uint256 ) view returns(uint256 id, string title, string description, uint256 targetFunding, uint256 currentFunding, uint256 startDate, uint256 endDate, address owner, uint8 status)
func (_Crowdfunding *CrowdfundingCaller) Projects(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id             *big.Int
	Title          string
	Description    string
	TargetFunding  *big.Int
	CurrentFunding *big.Int
	StartDate      *big.Int
	EndDate        *big.Int
	Owner          common.Address
	Status         uint8
}, error) {
	var out []interface{}
	err := _Crowdfunding.contract.Call(opts, &out, "projects", arg0)

	outstruct := new(struct {
		Id             *big.Int
		Title          string
		Description    string
		TargetFunding  *big.Int
		CurrentFunding *big.Int
		StartDate      *big.Int
		EndDate        *big.Int
		Owner          common.Address
		Status         uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Title = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Description = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.TargetFunding = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.CurrentFunding = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.StartDate = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.EndDate = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Owner = *abi.ConvertType(out[7], new(common.Address)).(*common.Address)
	outstruct.Status = *abi.ConvertType(out[8], new(uint8)).(*uint8)

	return *outstruct, err

}

// Projects is a free data retrieval call binding the contract method 0x107046bd.
//
// Solidity: function projects(uint256 ) view returns(uint256 id, string title, string description, uint256 targetFunding, uint256 currentFunding, uint256 startDate, uint256 endDate, address owner, uint8 status)
func (_Crowdfunding *CrowdfundingSession) Projects(arg0 *big.Int) (struct {
	Id             *big.Int
	Title          string
	Description    string
	TargetFunding  *big.Int
	CurrentFunding *big.Int
	StartDate      *big.Int
	EndDate        *big.Int
	Owner          common.Address
	Status         uint8
}, error) {
	return _Crowdfunding.Contract.Projects(&_Crowdfunding.CallOpts, arg0)
}

// Projects is a free data retrieval call binding the contract method 0x107046bd.
//
// Solidity: function projects(uint256 ) view returns(uint256 id, string title, string description, uint256 targetFunding, uint256 currentFunding, uint256 startDate, uint256 endDate, address owner, uint8 status)
func (_Crowdfunding *CrowdfundingCallerSession) Projects(arg0 *big.Int) (struct {
	Id             *big.Int
	Title          string
	Description    string
	TargetFunding  *big.Int
	CurrentFunding *big.Int
	StartDate      *big.Int
	EndDate        *big.Int
	Owner          common.Address
	Status         uint8
}, error) {
	return _Crowdfunding.Contract.Projects(&_Crowdfunding.CallOpts, arg0)
}

// Contribute is a paid mutator transaction binding the contract method 0x8c590917.
//
// Solidity: function contribute(uint256 projectId, uint256 amount) returns()
func (_Crowdfunding *CrowdfundingTransactor) Contribute(opts *bind.TransactOpts, projectId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.contract.Transact(opts, "contribute", projectId, amount)
}

// Contribute is a paid mutator transaction binding the contract method 0x8c590917.
//
// Solidity: function contribute(uint256 projectId, uint256 amount) returns()
func (_Crowdfunding *CrowdfundingSession) Contribute(projectId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.Contribute(&_Crowdfunding.TransactOpts, projectId, amount)
}

// Contribute is a paid mutator transaction binding the contract method 0x8c590917.
//
// Solidity: function contribute(uint256 projectId, uint256 amount) returns()
func (_Crowdfunding *CrowdfundingTransactorSession) Contribute(projectId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.Contribute(&_Crowdfunding.TransactOpts, projectId, amount)
}

// CreateProject is a paid mutator transaction binding the contract method 0xf81bf601.
//
// Solidity: function createProject(string title, string description, uint256 targetFunding, uint256 startDate, uint256 endDate) returns()
func (_Crowdfunding *CrowdfundingTransactor) CreateProject(opts *bind.TransactOpts, title string, description string, targetFunding *big.Int, startDate *big.Int, endDate *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.contract.Transact(opts, "createProject", title, description, targetFunding, startDate, endDate)
}

// CreateProject is a paid mutator transaction binding the contract method 0xf81bf601.
//
// Solidity: function createProject(string title, string description, uint256 targetFunding, uint256 startDate, uint256 endDate) returns()
func (_Crowdfunding *CrowdfundingSession) CreateProject(title string, description string, targetFunding *big.Int, startDate *big.Int, endDate *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.CreateProject(&_Crowdfunding.TransactOpts, title, description, targetFunding, startDate, endDate)
}

// CreateProject is a paid mutator transaction binding the contract method 0xf81bf601.
//
// Solidity: function createProject(string title, string description, uint256 targetFunding, uint256 startDate, uint256 endDate) returns()
func (_Crowdfunding *CrowdfundingTransactorSession) CreateProject(title string, description string, targetFunding *big.Int, startDate *big.Int, endDate *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.CreateProject(&_Crowdfunding.TransactOpts, title, description, targetFunding, startDate, endDate)
}

// DeleteProject is a paid mutator transaction binding the contract method 0x4cd0cb28.
//
// Solidity: function deleteProject(uint256 projectId) returns()
func (_Crowdfunding *CrowdfundingTransactor) DeleteProject(opts *bind.TransactOpts, projectId *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.contract.Transact(opts, "deleteProject", projectId)
}

// DeleteProject is a paid mutator transaction binding the contract method 0x4cd0cb28.
//
// Solidity: function deleteProject(uint256 projectId) returns()
func (_Crowdfunding *CrowdfundingSession) DeleteProject(projectId *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.DeleteProject(&_Crowdfunding.TransactOpts, projectId)
}

// DeleteProject is a paid mutator transaction binding the contract method 0x4cd0cb28.
//
// Solidity: function deleteProject(uint256 projectId) returns()
func (_Crowdfunding *CrowdfundingTransactorSession) DeleteProject(projectId *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.DeleteProject(&_Crowdfunding.TransactOpts, projectId)
}

// EditProject is a paid mutator transaction binding the contract method 0xbd481cf9.
//
// Solidity: function editProject(uint256 projectId, string title, string description, uint256 targetFunding, uint256 startDate, uint256 endDate) returns()
func (_Crowdfunding *CrowdfundingTransactor) EditProject(opts *bind.TransactOpts, projectId *big.Int, title string, description string, targetFunding *big.Int, startDate *big.Int, endDate *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.contract.Transact(opts, "editProject", projectId, title, description, targetFunding, startDate, endDate)
}

// EditProject is a paid mutator transaction binding the contract method 0xbd481cf9.
//
// Solidity: function editProject(uint256 projectId, string title, string description, uint256 targetFunding, uint256 startDate, uint256 endDate) returns()
func (_Crowdfunding *CrowdfundingSession) EditProject(projectId *big.Int, title string, description string, targetFunding *big.Int, startDate *big.Int, endDate *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.EditProject(&_Crowdfunding.TransactOpts, projectId, title, description, targetFunding, startDate, endDate)
}

// EditProject is a paid mutator transaction binding the contract method 0xbd481cf9.
//
// Solidity: function editProject(uint256 projectId, string title, string description, uint256 targetFunding, uint256 startDate, uint256 endDate) returns()
func (_Crowdfunding *CrowdfundingTransactorSession) EditProject(projectId *big.Int, title string, description string, targetFunding *big.Int, startDate *big.Int, endDate *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.EditProject(&_Crowdfunding.TransactOpts, projectId, title, description, targetFunding, startDate, endDate)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 projectId) returns()
func (_Crowdfunding *CrowdfundingTransactor) Refund(opts *bind.TransactOpts, projectId *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.contract.Transact(opts, "refund", projectId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 projectId) returns()
func (_Crowdfunding *CrowdfundingSession) Refund(projectId *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.Refund(&_Crowdfunding.TransactOpts, projectId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 projectId) returns()
func (_Crowdfunding *CrowdfundingTransactorSession) Refund(projectId *big.Int) (*types.Transaction, error) {
	return _Crowdfunding.Contract.Refund(&_Crowdfunding.TransactOpts, projectId)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Crowdfunding *CrowdfundingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Crowdfunding.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Crowdfunding *CrowdfundingSession) RenounceOwnership() (*types.Transaction, error) {
	return _Crowdfunding.Contract.RenounceOwnership(&_Crowdfunding.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Crowdfunding *CrowdfundingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Crowdfunding.Contract.RenounceOwnership(&_Crowdfunding.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Crowdfunding *CrowdfundingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Crowdfunding.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Crowdfunding *CrowdfundingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Crowdfunding.Contract.TransferOwnership(&_Crowdfunding.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Crowdfunding *CrowdfundingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Crowdfunding.Contract.TransferOwnership(&_Crowdfunding.TransactOpts, newOwner)
}

// CrowdfundingContributionMadeIterator is returned from FilterContributionMade and is used to iterate over the raw logs and unpacked data for ContributionMade events raised by the Crowdfunding contract.
type CrowdfundingContributionMadeIterator struct {
	Event *CrowdfundingContributionMade // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdfundingContributionMadeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdfundingContributionMade)
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
		it.Event = new(CrowdfundingContributionMade)
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
func (it *CrowdfundingContributionMadeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdfundingContributionMadeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdfundingContributionMade represents a ContributionMade event raised by the Crowdfunding contract.
type CrowdfundingContributionMade struct {
	ProjectId   *big.Int
	Amount      *big.Int
	Contributor common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContributionMade is a free log retrieval operation binding the contract event 0x7c50e76f3eff5098fcf915218e219a359f40f4c150e1a6796920587d44923cb4.
//
// Solidity: event ContributionMade(uint256 indexed projectId, uint256 amount, address contributor)
func (_Crowdfunding *CrowdfundingFilterer) FilterContributionMade(opts *bind.FilterOpts, projectId []*big.Int) (*CrowdfundingContributionMadeIterator, error) {

	var projectIdRule []interface{}
	for _, projectIdItem := range projectId {
		projectIdRule = append(projectIdRule, projectIdItem)
	}

	logs, sub, err := _Crowdfunding.contract.FilterLogs(opts, "ContributionMade", projectIdRule)
	if err != nil {
		return nil, err
	}
	return &CrowdfundingContributionMadeIterator{contract: _Crowdfunding.contract, event: "ContributionMade", logs: logs, sub: sub}, nil
}

// WatchContributionMade is a free log subscription operation binding the contract event 0x7c50e76f3eff5098fcf915218e219a359f40f4c150e1a6796920587d44923cb4.
//
// Solidity: event ContributionMade(uint256 indexed projectId, uint256 amount, address contributor)
func (_Crowdfunding *CrowdfundingFilterer) WatchContributionMade(opts *bind.WatchOpts, sink chan<- *CrowdfundingContributionMade, projectId []*big.Int) (event.Subscription, error) {

	var projectIdRule []interface{}
	for _, projectIdItem := range projectId {
		projectIdRule = append(projectIdRule, projectIdItem)
	}

	logs, sub, err := _Crowdfunding.contract.WatchLogs(opts, "ContributionMade", projectIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdfundingContributionMade)
				if err := _Crowdfunding.contract.UnpackLog(event, "ContributionMade", log); err != nil {
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

// ParseContributionMade is a log parse operation binding the contract event 0x7c50e76f3eff5098fcf915218e219a359f40f4c150e1a6796920587d44923cb4.
//
// Solidity: event ContributionMade(uint256 indexed projectId, uint256 amount, address contributor)
func (_Crowdfunding *CrowdfundingFilterer) ParseContributionMade(log types.Log) (*CrowdfundingContributionMade, error) {
	event := new(CrowdfundingContributionMade)
	if err := _Crowdfunding.contract.UnpackLog(event, "ContributionMade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdfundingContributionRefundedIterator is returned from FilterContributionRefunded and is used to iterate over the raw logs and unpacked data for ContributionRefunded events raised by the Crowdfunding contract.
type CrowdfundingContributionRefundedIterator struct {
	Event *CrowdfundingContributionRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdfundingContributionRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdfundingContributionRefunded)
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
		it.Event = new(CrowdfundingContributionRefunded)
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
func (it *CrowdfundingContributionRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdfundingContributionRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdfundingContributionRefunded represents a ContributionRefunded event raised by the Crowdfunding contract.
type CrowdfundingContributionRefunded struct {
	ProjectId   *big.Int
	Amount      *big.Int
	Contributor common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContributionRefunded is a free log retrieval operation binding the contract event 0x232bf2213879e79f5b2fa9cf07aad94aa0acac7b3c93bdfe74a545d47a1b0275.
//
// Solidity: event ContributionRefunded(uint256 indexed projectId, uint256 amount, address contributor)
func (_Crowdfunding *CrowdfundingFilterer) FilterContributionRefunded(opts *bind.FilterOpts, projectId []*big.Int) (*CrowdfundingContributionRefundedIterator, error) {

	var projectIdRule []interface{}
	for _, projectIdItem := range projectId {
		projectIdRule = append(projectIdRule, projectIdItem)
	}

	logs, sub, err := _Crowdfunding.contract.FilterLogs(opts, "ContributionRefunded", projectIdRule)
	if err != nil {
		return nil, err
	}
	return &CrowdfundingContributionRefundedIterator{contract: _Crowdfunding.contract, event: "ContributionRefunded", logs: logs, sub: sub}, nil
}

// WatchContributionRefunded is a free log subscription operation binding the contract event 0x232bf2213879e79f5b2fa9cf07aad94aa0acac7b3c93bdfe74a545d47a1b0275.
//
// Solidity: event ContributionRefunded(uint256 indexed projectId, uint256 amount, address contributor)
func (_Crowdfunding *CrowdfundingFilterer) WatchContributionRefunded(opts *bind.WatchOpts, sink chan<- *CrowdfundingContributionRefunded, projectId []*big.Int) (event.Subscription, error) {

	var projectIdRule []interface{}
	for _, projectIdItem := range projectId {
		projectIdRule = append(projectIdRule, projectIdItem)
	}

	logs, sub, err := _Crowdfunding.contract.WatchLogs(opts, "ContributionRefunded", projectIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdfundingContributionRefunded)
				if err := _Crowdfunding.contract.UnpackLog(event, "ContributionRefunded", log); err != nil {
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

// ParseContributionRefunded is a log parse operation binding the contract event 0x232bf2213879e79f5b2fa9cf07aad94aa0acac7b3c93bdfe74a545d47a1b0275.
//
// Solidity: event ContributionRefunded(uint256 indexed projectId, uint256 amount, address contributor)
func (_Crowdfunding *CrowdfundingFilterer) ParseContributionRefunded(log types.Log) (*CrowdfundingContributionRefunded, error) {
	event := new(CrowdfundingContributionRefunded)
	if err := _Crowdfunding.contract.UnpackLog(event, "ContributionRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdfundingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Crowdfunding contract.
type CrowdfundingOwnershipTransferredIterator struct {
	Event *CrowdfundingOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdfundingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdfundingOwnershipTransferred)
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
		it.Event = new(CrowdfundingOwnershipTransferred)
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
func (it *CrowdfundingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdfundingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdfundingOwnershipTransferred represents a OwnershipTransferred event raised by the Crowdfunding contract.
type CrowdfundingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Crowdfunding *CrowdfundingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CrowdfundingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Crowdfunding.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CrowdfundingOwnershipTransferredIterator{contract: _Crowdfunding.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Crowdfunding *CrowdfundingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CrowdfundingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Crowdfunding.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdfundingOwnershipTransferred)
				if err := _Crowdfunding.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Crowdfunding *CrowdfundingFilterer) ParseOwnershipTransferred(log types.Log) (*CrowdfundingOwnershipTransferred, error) {
	event := new(CrowdfundingOwnershipTransferred)
	if err := _Crowdfunding.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdfundingProjectCreatedIterator is returned from FilterProjectCreated and is used to iterate over the raw logs and unpacked data for ProjectCreated events raised by the Crowdfunding contract.
type CrowdfundingProjectCreatedIterator struct {
	Event *CrowdfundingProjectCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log          // Log channel receiving the found contract events
	sub  interfaces.Subscription // Subscription for errors, completion and termination
	done bool                    // Whether the subscription completed delivering logs
	fail error                   // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdfundingProjectCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdfundingProjectCreated)
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
		it.Event = new(CrowdfundingProjectCreated)
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
func (it *CrowdfundingProjectCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdfundingProjectCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdfundingProjectCreated represents a ProjectCreated event raised by the Crowdfunding contract.
type CrowdfundingProjectCreated struct {
	ProjectId *big.Int
	Title     string
	Owner     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterProjectCreated is a free log retrieval operation binding the contract event 0xc0c54fed07481d0998e1446b2c13759606bf4f26b78306307413ac4c4309aa82.
//
// Solidity: event ProjectCreated(uint256 indexed projectId, string title, address owner)
func (_Crowdfunding *CrowdfundingFilterer) FilterProjectCreated(opts *bind.FilterOpts, projectId []*big.Int) (*CrowdfundingProjectCreatedIterator, error) {

	var projectIdRule []interface{}
	for _, projectIdItem := range projectId {
		projectIdRule = append(projectIdRule, projectIdItem)
	}

	logs, sub, err := _Crowdfunding.contract.FilterLogs(opts, "ProjectCreated", projectIdRule)
	if err != nil {
		return nil, err
	}
	return &CrowdfundingProjectCreatedIterator{contract: _Crowdfunding.contract, event: "ProjectCreated", logs: logs, sub: sub}, nil
}

// WatchProjectCreated is a free log subscription operation binding the contract event 0xc0c54fed07481d0998e1446b2c13759606bf4f26b78306307413ac4c4309aa82.
//
// Solidity: event ProjectCreated(uint256 indexed projectId, string title, address owner)
func (_Crowdfunding *CrowdfundingFilterer) WatchProjectCreated(opts *bind.WatchOpts, sink chan<- *CrowdfundingProjectCreated, projectId []*big.Int) (event.Subscription, error) {

	var projectIdRule []interface{}
	for _, projectIdItem := range projectId {
		projectIdRule = append(projectIdRule, projectIdItem)
	}

	logs, sub, err := _Crowdfunding.contract.WatchLogs(opts, "ProjectCreated", projectIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdfundingProjectCreated)
				if err := _Crowdfunding.contract.UnpackLog(event, "ProjectCreated", log); err != nil {
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

// ParseProjectCreated is a log parse operation binding the contract event 0xc0c54fed07481d0998e1446b2c13759606bf4f26b78306307413ac4c4309aa82.
//
// Solidity: event ProjectCreated(uint256 indexed projectId, string title, address owner)
func (_Crowdfunding *CrowdfundingFilterer) ParseProjectCreated(log types.Log) (*CrowdfundingProjectCreated, error) {
	event := new(CrowdfundingProjectCreated)
	if err := _Crowdfunding.contract.UnpackLog(event, "ProjectCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
