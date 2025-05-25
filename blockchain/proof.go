package blockchain

// ProofOfStake contains the data required to validate block
type ProofOfStake struct {
	Block        Block
	Requirements *Requirements
}

type Requirements struct {
	OnOffenderList      bool
	CarInsurance        bool
	ConfirmRequirements bool
}

func (pos *ProofOfStake) Run() {
	pos.Requirements = &Requirements{
		OnOffenderList:      isOnOffenderList(),
		CarInsurance:        hasCarInsurance(),
		ConfirmRequirements: validatorHasConfirmedRequirements(),
	}
}

func NewProof(b Block) *ProofOfStake {
	return &ProofOfStake{
		Block: b,
	}
}

func isOnOffenderList() bool {
	return false
}

func hasCarInsurance() bool {
	return true
}

func validatorHasConfirmedRequirements() bool {
	return true
}
