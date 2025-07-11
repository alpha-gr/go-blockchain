package cli

import (
	"bufio"
	"flag"
	"fmt"
	"go-blockchain/blockchain"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - Get balance of an address")
	fmt.Println(" createblockchain -address ADDRESS - Create a new blockchain and the genesis block")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount from one address to another")
	fmt.Println(" mine -address ADDRESS - Mine a new block and receive mining reward")
	fmt.Println(" printchain - Print the blocks in the blockchain")
	fmt.Println(" interactive - Start interactive mode for continuous command execution")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
		// unlike os.Exit, exits by shutting down the goroutine
		// prevents badgerdb form corrupting the database
	}
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iterator := chain.Iterator()

	for {
		block := iterator.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Block Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockchain(address)
	chain.Database.Close()
	fmt.Println("Blockchain created successfully!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOS := chain.FindUTXO(address)

	for _, out := range UTXOS {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Printf("Transaction sent from %s to %s for amount %d successfully\n", from, to, amount)
}

func (cli *CommandLine) mine(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	tx := blockchain.CoinbaseTx(address, "Mining reward")
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Printf("Block mined successfully! Mining reward of 100 coins sent to %s\n", address)
}

func (cli *CommandLine) interactiveMode() {
	fmt.Println("=== Blockchain Interactive Mode ===")
	fmt.Println("Type 'help' for available commands, 'exit' or 'quit' to leave")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("blockchain> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Split input into arguments
		args := strings.Fields(input)
		command := args[0]

		switch command {
		case "exit", "quit":
			fmt.Println("Exiting interactive mode...")
			return
		case "help":
			cli.printInteractiveHelp()
		case "getbalance":
			if len(args) != 2 {
				fmt.Println("Usage: getbalance <address>")
				continue
			}
			cli.getBalance(args[1])
		case "createblockchain":
			if len(args) != 2 {
				fmt.Println("Usage: createblockchain <address>")
				continue
			}
			cli.createBlockChain(args[1])
		case "printchain":
			cli.printChain()
		case "send":
			if len(args) != 4 {
				fmt.Println("Usage: send <from> <to> <amount>")
				continue
			}
			amount, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println("Error: amount must be a valid number")
				continue
			}
			if amount <= 0 {
				fmt.Println("Error: amount must be greater than 0")
				continue
			}
			cli.send(args[1], args[2], amount)
		case "mine":
			if len(args) != 2 {
				fmt.Println("Usage: mine <address>")
				continue
			}
			cli.mine(args[1])
		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Type 'help' for available commands")
		}
		fmt.Println()
	}
}

func (cli *CommandLine) printInteractiveHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  getbalance <address>          - Get balance of an address")
	fmt.Println("  createblockchain <address>    - Create a new blockchain and genesis block")
	fmt.Println("  send <from> <to> <amount>     - Send amount from one address to another")
	fmt.Println("  mine <address>                - Mine a new block and receive mining reward")
	fmt.Println("  printchain                    - Print all blocks in the blockchain")
	fmt.Println("  help                          - Show this help message")
	fmt.Println("  exit, quit                    - Exit interactive mode")
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "Address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "Address to create blockchain for")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	mineAddress := mineCmd.String("address", "", "Address to receive mining reward")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "mine":
		err := mineCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "interactive":
		cli.interactiveMode()
		return
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if mineCmd.Parsed() {
		if *mineAddress == "" {
			mineCmd.Usage()
			runtime.Goexit()
		}
		cli.mine(*mineAddress)
	}

}
