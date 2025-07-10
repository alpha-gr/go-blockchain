package main

import (
	"flag"
	"fmt"
	"go-blockchain/blockchain"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	Blockchain *blockchain.Blockchain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  add -data BLOCK_DATA - Add a block to the blockchain")
	fmt.Println("  print - Print the blocks in the blockchain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
		// unlike os.Exit, exits by shutting down the goroutine
		// prevents badgerdb form corrupting the database
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.Blockchain.AddBlock(data)
	fmt.Println("Block added successfully!")
}

func (cli *CommandLine) printChain() {
	iterator := cli.Blockchain.Iterator()
	for {
		block := iterator.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func main() {
	chain := blockchain.InitBlockchain()
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()
}
