package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"jjeth/blockchain"
	"log"
	mrand "math/rand"
	"strings"
	"sync"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	net "github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

var mutex = &sync.Mutex{}

func MakeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader

	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}
	basicHost, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))
	addrs := basicHost.Addrs()
	var addr ma.Multiaddr
	for _, i := range addrs {
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i
			break
		}
	}
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	log.Printf("Now run \"go run main.go -l %d \" on a different terminal \n", listenPort+1)
	return basicHost, nil
}

func handleStream(s net.Stream) {
	log.Println("got a new stream!")
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
READ:

	for {

		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if str == "" {
			return
		}
		if str != "\n" {
			ch := make([]*blockchain.Block, 0)
			if err := json.Unmarshal([]byte(str), &ch); err != nil {
				log.Fatal(err)
			}
			mutex.Lock()
			if len(ch) > len(blockchain.JJchain.Chain) {
				if len(blockchain.JJchain.Chain) == 1 {
					blockchain.JJchain.Chain = make([]*blockchain.Block, 0)
					blockchain.JJchain.Chain = append(blockchain.JJchain.Chain, ch[0])
				}

				for i := len(blockchain.JJchain.Chain); i < len(ch)-1; i++ {
					fmt.Printf("adding block: %v\n", ch[i], blockchain.JJchain)
					err := blockchain.JJchain.AddBlock(ch[i])
					if err != nil {
						fmt.Println(err)
						break READ
					}
				}
				//blockchain.JJchain.Chain = ch
				bytes, err := json.MarshalIndent(blockchain.JJchain.Chain, "", " ")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("\x1b[32m%s\x1b[0m]>", string(bytes))

			}
			mutex.Unlock()
		}
	}

}

func writeData(rw *bufio.ReadWriter) {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(blockchain.JJchain.Chain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()
			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()
		}
	}()

}

func StartServer(target string, port int) {
	ha, err := MakeBasicHost(port, false, 0)
	if err != nil {
		log.Fatal(err)
	}
	if target == "" {
		log.Println("listening for connections")
		ha.SetStreamHandler("/p2p/1.0.0", net.StreamHandler(handleStream))
		select {}
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)
		ipfsaddr, err := ma.NewMultiaddr(target)
		if err != nil {
			log.Fatal(err)
		}
		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatal(err)
		}
		peerid, err := peer.Decode(pid)
		if err != nil {
			log.Fatal(err)
		}
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
		ha.Peerstore().AddAddr(peerid, targetAddr, time.Hour)
		log.Println("opening stream")
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go writeData(rw)
		go readData(rw)
		select {}
	}

}
