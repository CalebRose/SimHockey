package dbprovider

import (
	"fmt"
	"os"
)

type SshTunnelConfig struct {
	SshHost     string // SSH server host
	SshPort     string // SSH server port
	SshUser     string // SSH username
	SshPassword string // SSH password
	DbHost      string // Remote database host (from the perspective of the SSH server)
	DbPort      string // Remote database port
	LocalPort   string // Local port to forward to dbPort over the SSH tunnel
}

func GetSSHConfig() SshTunnelConfig {
	hostName, hnExists := os.LookupEnv("SSHHN")
	sshPort, sshPoExists := os.LookupEnv("SSHPO")
	sshUser, sshUExists := os.LookupEnv("SSHU")
	sshPassword, sshPExists := os.LookupEnv("SSHP")
	dbHost, dbHExists := os.LookupEnv("DBH")
	dbPort, dbPExists := os.LookupEnv("DBP")
	localPort, localExists := os.LookupEnv("LCP")
	if hnExists && sshPoExists && sshUExists && sshPExists && dbHExists && dbPExists && localExists {
		return SshTunnelConfig{
			SshHost:     hostName,
			SshPort:     sshPort,
			SshUser:     sshUser,
			SshPassword: sshPassword,
			DbHost:      dbHost,
			DbPort:      dbPort,
			LocalPort:   localPort,
		}
	}
	fmt.Println("WARNING! COULD NOT GET ENV VARIABLES. TRYING ALT METHOD")
	hostName = os.Getenv("SSHHN")
	sshPort = os.Getenv("SSHPO")
	sshUser = os.Getenv("SSHU")
	sshPassword = os.Getenv("SSHP")
	dbHost = os.Getenv("DBH")
	dbPort = os.Getenv("DBP")
	localPort = os.Getenv("LCP")
	return SshTunnelConfig{
		SshHost:     hostName,
		SshPort:     sshPort,
		SshUser:     sshUser,
		SshPassword: sshPassword,
		DbHost:      dbHost,
		DbPort:      dbPort,
		LocalPort:   localPort,
	}
}

func Config(local string) map[string]string {
	dbName, exists := os.LookupEnv("DB")
	csUserName, csUNExists := os.LookupEnv("CSUSERNAME")
	csPW, csPWExists := os.LookupEnv("CSPW")
	hostDB, dbExists := os.LookupEnv("DBNAME")
	lcp, lcpExists := os.LookupEnv("LCP")

	if exists && csPWExists && csUNExists && dbExists && lcpExists {
		connstring := csUserName + ":" + csPW + "@tcp(localhost:" + lcp + ")/" + hostDB + "?parseTime=true"
		return map[string]string{
			"db": dbName,
			"cs": connstring,
		}
	}
	dbName = os.Getenv("DB")
	csUserName = os.Getenv("CSUSERNAME")
	csPW = os.Getenv("CSPW")
	hostDB = os.Getenv("DBNAME")
	lcp = os.Getenv("LCP")
	connstring := csUserName + ":" + csPW + "@tcp(localhost:" + lcp + ")/" + hostDB + "?parseTime=true"
	return map[string]string{
		"db": dbName,
		"cs": connstring,
	}
}

func SFAConfig() map[string]string {
	sfaKey, exists := os.LookupEnv("SFAKEY")
	sfaUser, userExists := os.LookupEnv("SFAUSER")
	if exists && userExists {
		return map[string]string{
			"sfaKey":  sfaKey,
			"sfaUser": sfaUser,
		}
	}
	sfaKey = os.Getenv("SFAKEY")
	sfaUser = os.Getenv("SFAUSER")
	return map[string]string{
		"sfaKey":  sfaKey,
		"sfaUser": sfaUser,
	}
}
