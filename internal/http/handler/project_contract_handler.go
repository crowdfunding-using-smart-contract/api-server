package handler

import (
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"

	"math/big"

	"github.com/ava-labs/coreth/ethclient"
	"github.com/ethereum/go-ethereum/common"

	"fund-o/api-server/pkg/contracts/crowdfunding"
)

type ProjectContractHandler struct {
	ethClient           *ethclient.Client
	crowdfundingAddress common.Address
}

func NewProjectContractHandler(client *ethclient.Client, address string) *ProjectContractHandler {
	return &ProjectContractHandler{
		ethClient:           client,
		crowdfundingAddress: common.HexToAddress(address),
	}
}

func (h *ProjectContractHandler) GetAllCrowdfundingProjects(c *gin.Context) {
	const activeStatus uint8 = 0

	contract, err := crowdfunding.NewCrowdfunding(h.crowdfundingAddress, h.ethClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contract"})
		return
	}

	projects, err := contract.GetAllProjects(nil, activeStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}

	projectIDs := make([]uint64, len(projects))
	for i, id := range projects {
		projectIDs[i] = id.Uint64()
	}

	c.JSON(http.StatusOK, gin.H{"projects": projectIDs})
}

func (h *ProjectContractHandler) CreateCrowdfundingProject(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create crowdfunding project mock response"})
}

func (h *ProjectContractHandler) ContributeToProject(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Contribute to project mock response"})
}

func (h *ProjectContractHandler) RefundProject(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Refund project contribution mock response"})
}

func (h *ProjectContractHandler) DeleteProjectContract(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete project contract mock response"})
}

func (h *ProjectContractHandler) GetProjectContractByID(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	contract, err := crowdfunding.NewCrowdfunding(h.crowdfundingAddress, h.ethClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contract"})
		return
	}

	project, err := contract.GetProject(nil, new(big.Int).SetUint64(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

func (h *ProjectContractHandler) GetUserProjectContracts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user project contracts mock response"})
}

func (h *ProjectContractHandler) EditProjectContract(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Edit project contract mock response"})
}
