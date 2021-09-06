# PhishkitHunter
A project to take known phishing sites and brute force them to find phishing kits


Requires a tor service to be running on 9050


## Example usage:
```
go build
./PhishkitHunter -w wordlist.txt -e https://tradeswarehouse.com/ -o outfiles/
./PhishkitHunter -u -e https://tradeswarehouse.com/1drvme/ -o outfiles/
```
