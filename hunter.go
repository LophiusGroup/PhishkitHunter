package main

import (
	"flag"
	"log"
	"os"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"io/ioutil"
	"errors"
	"bufio"
)

func main() {
		
	var endPoint = flag.String("endPoint", "", "target url to brute")
   	flag.StringVar(endPoint, "e", "", "target url to brute")

   	var outDir = flag.String("outDir", "", "directory to write files to")
   	flag.StringVar(outDir, "o", "", "directory to write files to")

	var wordList = flag.String("wordList", "", "wordlist to use")
	flag.StringVar(wordList, "w", "", "Input file to triage")

	var logFile = flag.String("logFile", "", "send stdout to a log file")
	flag.StringVar(logFile, "l", "", "send stdout to a log file")

	flag.Parse()

  	//Setup logfile stuff
	if *logFile != "" {
		logTown, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println("error opening file: %v", err)
		}
		defer logTown.Close()
		log.SetOutput(logTown)
		log.Println("Log file started!")
	}

	if *wordList != "" {
		err := ScanListTor(*endPoint, *wordList, *outDir)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Done!")
		}
	}
}

func ScanListTor(endpoint, wordlist, outDir string) error {
    	//Open wordlist
	words := ReadLines(wordlist)
	for _, word := range words {
		err := ReqThroughTor(endpoint, word, outDir)
		if err != nil {
			log.Printf("Failed to request with wordlist target: %v, error: %v\n", word, err)
			return err
		}
	}
	return nil
}


// Tor HTTP Request, largely from: https://gist.github.com/Yawning/bac58e08a05fc378a8cc
// ReqThroughTor
func ReqThroughTor(endpoint, target, outDir string) error {
	// Create a transport that uses Tor Browser's SocksPort.  If
	// talking to a system tor, this may be an AF_UNIX socket, or
	// 127.0.0.1:9050 instead.
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		log.Printf("Failed to parse proxy URL: %v\n", err)
		return err
	}

	// Get a proxy Dialer that will create the connection on our
	// behalf via the SOCKS5 proxy.  Specify the authentication
	// and re-create the dialer/transport/client if tor's
	// IsolateSOCKSAuth is needed.
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		log.Printf("Failed to obtain proxy dialer: %v\n", err)
		return err
	}

	// Make a http.Transport that uses the proxy dialer, and a
	// http.Client that uses the transport.
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}

	// Example: Fetch something.  Real code will probably want to use
	// client.Do() so they can change the User-Agent.
	resp, err := client.Get(endpoint+target)
	if err != nil {
		log.Printf("Failed to issue GET request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("GET returned: %v\n", resp.StatusCode)
    	if resp.StatusCode == 200 {
    		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		    	log.Printf("Failed to read the body: %v\n", err)
			return err
	    	}
	    	//log.Printf("----- Body -----\n%s\n----- Body -----", body)
		//Need to save the kit in a better way
    		err = CreateFile(body, outDir+"/"+target)
    		if err != nil {
	   		log.Printf("Failed to save file: %v\n", err)
			return err
		} else {
			return nil
		}
    	}
	log.Printf("Did not find kit")
	return nil
}

//CreateFile takes a byte array and a file path and writes the bytes to that location. 
//It uses Exists to make sure the file path is available before writing
func CreateFile(bytes []byte, path string) error {
    	// Check if the file already exists
    	if Exists(path) {
        	return errors.New("The file to create already exists so we won't overwite it")
    	}
    	// write the lines to the file
    	err := ioutil.WriteFile(path, bytes, 0700)
    	if err != nil {
        	return err
    	}
    	return nil
}

//Exists checks a path and returns a bool if there is a file there
func Exists(path string) bool {
    	// Run stat on a file
   	 _, err := os.Stat(path)
    	// If it runs fine the file exists
    	if err == nil {
        	return true
    	}
    	// If stat fails then the file does not exist
    	return false
}

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
