package block

import (
	"github.com/aayampokharel/fyp/models"
	"github.com/aayampokharel/fyp/utils/global"
)

func InsertBlock(block models.CertificateData) error {
	// global array check garyo
	// if array is empty insert directly there , recalcualte the headers thing .
	// else call evaluete_new_header, and do stuffs . and insert the data there .
	lastInsertedBlock := global.BlockChain[len(global.BlockChain)-1]
	return nil
}
