package usecase

func (uc *BlockChainUseCase) GetBlockChain() int {
	uc.Logger.Infoln("[get_block_chain] Info: GetBlockChain::", uc.BlockChainRepo.GetBlockChainLength())
	return uc.BlockChainRepo.GetBlockChainLength()
}
