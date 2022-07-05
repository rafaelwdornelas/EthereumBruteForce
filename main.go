package main

import (
	"fmt"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"log"
	"context"
	"math"
	"math/big"
	"math/rand"
        "time"
        "os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/aherve/gopool"
)

var tokens = []string{
    "a1743f084f8a46bfb3696389eeb9f217",
	"483c1730b99b46729c7f82f49302bbf8",
	"93ca33aa55d147f08666ac82d7cc69fd",
	"7aef3f0cd1f64408b163814b22cc643c",
	"3258e142b54447a89b8c002ee7465a6d",
	"c4898a6d00c44e82a3bbbe9e188f6717",
	"9db80c814ba04b51a0b79d439a3420e0",
	"e2da43769d694713b61702628ef124da",
	"e3e81da39a694d27b8b864800d70e5db",
	"ad5b634f40b349be8e5d77186c181531",
	"89dbeeaefe1847c085a3eabc7794fa47",
	"5b34306b94414c55a0779ae253f80b6f",
	"5f9c7a4d6f0e4c9aa48ae4044daee46b",
	"84842078b09946638c03157f83405213",
	"9aa3d95b3bc440fa88ea12eaa4456161",
	"bd72ec702e1d4b2689eb63012a860eb8",
	"800ff156fefb4424b2c7ba09e03a7ab0",
	"c69dc12364b34877806ae76b87c27fce",
	"ad612c8d545c4cc1a3d72812769fa1f5",
	"6a48654961d14572921bdc5be83c23a9",
}
var contagem = 0
func main() {
	pool := gopool.NewPool(10) 
    for {
        pool.Add(1)
        go func( pool *gopool.GoPool) {
            defer pool.Done()
			contagem++
			fmt.Println("Testado:",contagem)
            gerar()
        }(pool)
    }
    pool.Wait()
}



func gerar() {
    
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, "")	// Here you can choose to pass in the specified password or empty string , Different passwords generate different mnemonics 

	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
    
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0") // The last digit is the address of the same mnemonic word id, from 0 Start , The same mnemonic can produce unlimited addresses 
	account, err := wallet.Derive(path, false)
	if err != nil {
    
		log.Fatal(err)
	}

	
	address := account.Address.Hex()
	privateKey, _ := wallet.PrivateKeyHex(account)
	publicKey, _ := wallet.PublicKeyHex(account)  	

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(20 - 1)
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + tokens[n])
	if err != nil {
		//fmt.Println(err)
	}
	
	account2 := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account2, nil)
	if err != nil {
		//fmt.Println(err)
	}
	

	 if len(balance.Bits())== 0 {
		fmt.Println("mnemonic:", mnemonic, "Balance:", balance)
		
	} else {
		fbalance := new(big.Float)
		fbalance.SetString(balance.String())
		ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
		fmt.Println("mnemonic:", mnemonic, "Balance:", ethValue)
		salvaLog(mnemonic)
		salvaLog(address)
		salvaLog(privateKey)
		salvaLog(publicKey)
		salvaLog("-----------------------------------------------------")

	} 
}


func salvaLog(texto string) {
    f, err := os.OpenFile("retorno.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    if _, err := f.Write([]byte(texto + "\n")); err != nil {
        log.Fatal(err)
    }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}
