package container

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func MountVolume(rootURL string, mntURL string, volumeURLs []string) {
	parentUrl := volumeURLs[0]
	handleParentUrl(parentUrl)

	containerUrl := volumeURLs[1]
	containerVolumeURL := mntURL + containerUrl
	handleContainerVolumeUrl(containerUrl)
	
	dirs := "dirs=" + parentUrl
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount volume failed. %v", err)
	}

}

func DeleteMountPointWithVolume(rootURL string, mntURL string, volumeURLs []string) {
	containerUrl := mntURL + volumeURLs[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Umount volume failed. %v", err)
	}

	cmd = exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Umount mountpoint failed. %v", err)
	}

	if err := os.RemoveAll(mntURL); err != nil {
		log.Infof("Remove mountpoint dir %s error %v", mntURL, err)
	}
}

func handleParentUrl(parentUrl string) {
	if exist, err := PathExists(parentUrl); !exist {
		if err == nil {
			if err := os.Mkdir(parentUrl, 0777); err != nil {
				log.Infof("Mkdir parent dir %s error. %v", parentUrl, err)
			}
		} else {
			log.Infof("parent dir %s error. %v", parentUrl, err)
		}
	}
}

func handleContainerVolumeUrl(containerVolumeURL string) {
	if exist, err := PathExists(containerVolumeURL); !exist {
		if err == nil {
			if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
				log.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
			}
		}
	} else {
		log.Infof("container dir %s error. %v", containerVolumeURL, err)
	}

}
