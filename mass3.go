package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/miekg/dns"
	"github.com/remeh/sizedwaitgroup"
)

const (
	s3HostPrefix = "s3.amazonaws.com"
)

var (
	threads     int
	baseName    string
	wordList    string
	resolvers   string
	output      string
	results     []string
	nameservers []string
)

func main() {
	flag.StringVar(&wordList, "w", "", "path to the word list")
	flag.StringVar(&resolvers, "r", "lists/resolvers.txt", "path to the word list")
	flag.IntVar(&threads, "t", 10, "number of threads")
	flag.StringVar(&output, "o", "out.csv", "path to the output file")
	flag.Parse()

	if wordListExists() == false {
		return
	}
	prepareNameservers()
	wordListProcess()

	err := writeResultsToCsv(results, output)
	if err != nil {
		fmt.Printf("Couldn't write the results to the output file: %v", err)
	}
}

func writeResultsToCsv(results []string, outputFilePath string) error {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Couldn't create the output file: %v", err)
		return err
	}
	defer outputFile.Close()

	if len(results) != 0 {
		for _, str := range results {
			outputFile.WriteString(str + "\n")
		}
	} else {
		outputFile.WriteString("NA")
	}
	return nil
}

func prepareNameservers() {
	file, err := os.Open(resolvers)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nameservers = append(nameservers, scanner.Text())
	}
}

func randomNameserver() string {
	n := rand.Int() % len(nameservers)
	return nameservers[n]
}

// processing CNAME
func resolveCNAME(name string) []dns.RR {
	msg := dns.Msg{}
	msg.SetQuestion(name, dns.TypeCNAME)
	client := &dns.Client{Net: "tcp"}
	nameserver := randomNameserver()
	r, _, err := client.Exchange(&msg, nameserver+":53")
	if err != nil {
		return nil
	}
	return r.Answer
}

// wordlist processing

func wordListResolve(hosts []string) {
	swg := sizedwaitgroup.New(threads)

	for j := 0; j < len(hosts); j++ {
		swg.Add()
		host := fmt.Sprintf("%s.%s.", hosts[j], s3HostPrefix)
		go func(host string) {
			defer swg.Done()
			result := resolveCNAME(host)
			if len(result) == 0 {
				// No results
				return
			}

			if v, ok := result[0].(*dns.CNAME); ok {
				if !strings.Contains(v.Target, "s3-directional") {
					print(host + "\n")
					results = append(results, host)
				}
			}
		}(host)
	}
	swg.Wait()
}

func wordListProcess() {
	var hosts []string
	file, err := os.Open(wordList)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hosts = append(hosts, scanner.Text())
	}

	wordListResolve(hosts)
}

func wordListExists() bool {
	if _, err := os.Stat(wordList); os.IsNotExist(err) {
		print("Wordlist doens't exist.")
		return false
	}
	return true
}
