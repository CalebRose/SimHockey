package dbprovider

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Provider struct {
}

var db *gorm.DB
var once sync.Once
var instance *Provider

func GetInstance() *Provider {
	once.Do(func() {
		instance = &Provider{}
	})
	return instance
}

func (p *Provider) InitDatabase() bool {
	fmt.Println("Database initializing...")

	sshConfig := GetSSHConfig()
	localPort, localErr := setupSSHTunnel(&sshConfig)
	if localErr != nil {
		log.Fatalf("Failed to establish SSH tunnel: %v", localErr)
	}

	var err error
	c := Config(localPort) // c["cs"]
	db, err = gorm.Open(mysql.Open(c["cs"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return false
	}

	// AutoMigrations -- uncomment when needing to update a table
	//
	// General
	// College
	// db.AutoMigrate(&structs.CollegePlayer{})
	// db.AutoMigrate(&structs.HistoricCollegePlayer{})
	// db.AutoMigrate(&structs.CollegePlayerGameStats{})
	// db.AutoMigrate(&structs.CollegePlayerSeasonStats{})
	// db.AutoMigrate(&structs.CollegeLineup{})
	// db.AutoMigrate(&structs.CollegeShootoutLineup{})
	// db.AutoMigrate(&structs.CollegeTeam{})
	// db.AutoMigrate(&structs.CollegeTeamSeasonStats{})
	// db.AutoMigrate(&structs.CollegeTeamGameStats{})
	// db.AutoMigrate(&structs.CollegeGame{})

	// Professional
	// db.AutoMigrate(&structs.ProfessionalTeam{})
	// db.AutoMigrate(&structs.ProfessionalTeamFranchise{})
	// db.AutoMigrate(&structs.ProfessionalPlayer{})
	// db.AutoMigrate(&structs.RetiredPlayer{})
	// db.AutoMigrate(&structs.ProfessionalLineup{})
	// db.AutoMigrate(&structs.ProfessionalShootoutLineup{})
	// db.AutoMigrate(&structs.ProfessionalGame{})

	// Recruiting
	// db.AutoMigrate(&structs.Recruit{})
	// Administrative
	// db.AutoMigrate(&structs.Arena{})
	// db.AutoMigrate(&structs.GlobalPlayer{})
	// db.AutoMigrate(&structs.Timestamp{})
	return true
}

func (p *Provider) GetDB() *gorm.DB {
	return db
}

// setupSSHTunnel establishes an SSH tunnel and forwards a local port to the remote database port.
// Returns the local port and any error encountered.
func setupSSHTunnel(config *SshTunnelConfig) (string, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.SshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SshPassword),
		},
		// CAUTION: In production, you should use a more secure HostKeyCallback.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(config.SshHost, config.SshPort), sshConfig)
	if err != nil {
		return "", err
	}

	// Setup local port forwarding
	localListener, err := net.Listen("tcp", "localhost:"+config.LocalPort)
	if err != nil {
		return "", err
	}

	go func() {
		defer localListener.Close()
		for {
			localConn, err := localListener.Accept()
			if err != nil {
				log.Printf("Failed to accept local connection: %s", err)
				continue
			}

			// Handle the connection in a new goroutine
			go func() {
				defer localConn.Close()

				// Connect to the remote database server through the SSH tunnel
				remoteConn, err := sshClient.Dial("tcp", net.JoinHostPort(config.DbHost, config.DbPort))
				if err != nil {
					log.Printf("Failed to dial remote server: %s", err)
					return
				}
				defer remoteConn.Close()

				// Copy data between the local connection and the remote connection
				copyConn(localConn, remoteConn)
			}()
		}
	}()

	return localListener.Addr().String(), nil
}

// copyConn copies data between two io.ReadWriteCloser objects (e.g., network connections)
func copyConn(localConn, remoteConn io.ReadWriteCloser) {
	// Start goroutine to copy data from local to remote
	go func() {
		_, err := io.Copy(remoteConn, localConn)
		if err != nil {
			log.Printf("Error copying from local to remote: %v", err)
		}
		localConn.Close()
		remoteConn.Close()
	}()

	// Copy data from remote to local in the main goroutine (or vice versa)
	_, err := io.Copy(localConn, remoteConn)
	if err != nil {
		log.Printf("Error copying from remote to local: %v", err)
	}
	// Ensure connections are closed when copying is done or an error occurs
	localConn.Close()
	remoteConn.Close()
}
