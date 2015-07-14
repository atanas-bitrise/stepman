package cli

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-pathutil/pathutil"
	"github.com/bitrise-io/stepman/stepman"
	"github.com/codegangsta/cli"
)

func activate(c *cli.Context) {
	log.Info("[STEPMAN] - Activate")

	// Input validation
	id := c.String(IDKey)
	if id == "" {
		log.Fatal("[STEPMAN] - Missing step id")
	}

	version := c.String(VersionKey)
	if version == "" {
		log.Fatal("[STEPMAN] - Missing step version")
	}

	path := c.String(PathKey)
	if path == "" {
		log.Fatal("[STEPMAN] - Missing destination path")
	}

	// Get step
	collection, err := stepman.ReadStepSpec()
	if err != nil {
		log.Fatalln("[STEPMAN] - Failed to read steps spec")
	}

	exist, step := collection.GetStep(id, version)
	if exist == false {
		log.Fatalf("[STEPMAN] - Failed to activate Step: %s (v%s), does not exist in local cache.", id, version)
	}

	pth := stepman.GetStepPath(step)
	if exist, err := pathutil.IsPathExists(pth); err != nil {
		log.Fatal("[STEPMAN] - Failed to check path:", err)
	} else if exist == false {
		log.Info("[STEPMAN] - Step does not exist, download it")
		if err := stepman.DownloadStep(collection, step); err != nil {
			log.Fatal("[STEPMAN] - Failed to download step:", err)
		}
	}

	// Copy to specified path
	srcFolder := pth
	destFolder := path

	if exist, err = pathutil.IsPathExists(destFolder); err != nil {
		log.Fatal("[STEPMAN] - Failed to check path:", err)
	} else if exist == false {
		if err := os.MkdirAll(destFolder, 0777); err != nil {
			log.Fatal("[STEPMAN] - Failed to create path:", err)
		}
	}

	if err = stepman.RunCommand("cp", []string{"-rf", srcFolder, destFolder}...); err != nil {
		log.Fatal("[STEPMAN] - Failed to copy step:", err)
	}
}
