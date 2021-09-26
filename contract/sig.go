package contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/asm"
)

var (
	ERC721Methods map[string]string = map[string]string{
		"approve(address,uint256)":                        "095ea7b3",
		"balanceOf(address)":                              "70a08231",
		"getApproved(uint256)":                            "081812fc",
		"isApprovedForAll(address,address)":               "e985e9c5",
		"name()":                                          "06fdde03",
		"ownerOf(uint256)":                                "6352211e",
		"safeTransferFrom(address,address,uint256)":       "42842e0e",
		"safeTransferFrom(address,address,uint256,bytes)": "b88d4fde",
		"setApprovalForAll(address,bool)":                 "a22cb465",
		"symbol()":                                        "95d89b41",
		"tokenURI(uint256)":                               "c87b56dd",
		"transferFrom(address,address,uint256)":           "23b872dd",
	}
)

//CheckERC721 checks for an ERC721 token by searching through the disassembled bytecode
//If we find that a find that this bytecode is sufficient with
func CheckERC721(bytecode []byte, tolerance int) error {
	disassembly, err := asm.Disassemble(bytecode)

	if err != nil {
		return err
	}

	fmt.Println(disassembly)
	return nil
}
