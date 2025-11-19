package dbprovider

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/CalebRose/SimHockey/structs"
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
	// db.AutoMigrate(&structs.CollegePlayByPlay{})
	// db.AutoMigrate(&structs.CollegePlayerGameStats{})
	// db.AutoMigrate(&structs.CollegePlayerSeasonStats{})
	db.AutoMigrate(&structs.CollegeGameplan{})
	// db.AutoMigrate(&structs.CollegeLineup{})
	// db.AutoMigrate(&structs.CollegeShootoutLineup{})
	// db.AutoMigrate(&structs.CollegeTeam{})
	// db.AutoMigrate(&structs.CollegeTeamRequest{})
	// db.AutoMigrate(&structs.CollegeTeamSeasonStats{})
	// db.AutoMigrate(&structs.CollegeTeamGameStats{})
	// db.AutoMigrate(&structs.CollegeGame{})
	// db.AutoMigrate(&structs.CollegeSeries{})
	// db.AutoMigrate(&structs.CollegeStandings{})
	// db.AutoMigrate(&structs.CollegePollOfficial{})
	// db.AutoMigrate(&structs.CollegePollSubmission{})

	// Professional
	// db.AutoMigrate(&structs.DraftablePlayer{})
	// db.AutoMigrate(&structs.DraftPick{})
	// db.AutoMigrate(&structs.ExtensionOffer{})
	// db.AutoMigrate(&structs.FreeAgencyOffer{})
	// db.AutoMigrate(&structs.ProCapsheet{})
	// db.AutoMigrate(&structs.ProContract{})
	// db.AutoMigrate(&structs.ProPlayByPlay{})
	// db.AutoMigrate(&structs.ProfessionalPlayerGameStats{})
	// db.AutoMigrate(&structs.ProfessionalPlayerSeasonStats{})
	// db.AutoMigrate(&structs.ProfessionalTeam{})
	// db.AutoMigrate(&structs.ProfessionalTeamFranchise{})
	// db.AutoMigrate(&structs.ProfessionalTeamGameStats{})
	// db.AutoMigrate(&structs.ProfessionalTeamSeasonStats{})
	// db.AutoMigrate(&structs.ProfessionalPlayer{})
	// db.AutoMigrate(&structs.RetiredPlayer{})
	db.AutoMigrate(&structs.ProGameplan{})
	// db.AutoMigrate(&structs.ProfessionalLineup{})
	// db.AutoMigrate(&structs.ProfessionalShootoutLineup{})
	// db.AutoMigrate(&structs.ProfessionalGame{})
	// db.AutoMigrate(&structs.ProfessionalStandings{})
	// db.AutoMigrate(&structs.PlayoffSeries{})
	// db.AutoMigrate(&structs.ProSeries{})
	// db.AutoMigrate(&structs.ProTeamRequest{})
	// db.AutoMigrate(&structs.TradeProposal{})
	// db.AutoMigrate(&structs.TradeOption{})
	// db.AutoMigrate(&structs.TradePreferences{})
	// db.AutoMigrate(&structs.WaiverOffer{})

	// Recruiting
	// db.AutoMigrate(&structs.Recruit{})
	// db.AutoMigrate(&structs.RecruitPlayerProfile{})
	// db.AutoMigrate(&structs.RecruitingTeamProfile{})
	// db.AutoMigrate(&structs.RecruitPointAllocation{})

	// Portal
	// db.AutoMigrate(&structs.CollegePromise{})
	// db.AutoMigrate(&structs.TransferPortalProfile{})

	// Administrative
	// db.AutoMigrate(&structs.Arena{})
	// db.AutoMigrate(&structs.GlobalPlayer{})
	// db.AutoMigrate(&structs.FaceData{})
	// db.AutoMigrate(&structs.NewsLog{})
	// db.AutoMigrate(&structs.Notification{})
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
