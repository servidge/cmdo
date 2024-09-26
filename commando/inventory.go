package commando

import (
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func (app *appCfg) loadInventoryFromYAML(i *inventory) error {
	yamlFile, err := os.ReadFile(app.inventory)
	if err != nil {
		return err
	}

	err = yaml.UnmarshalStrict(yamlFile, i)
	if err != nil {
		log.Fatal(err)
	}

	filterDevices(i, app.devFilter)

	if len(i.Devices) == 0 {
		return errNoDevices
	}

	app.credentials = i.Credentials
	app.transports = i.Transports

	// user-provided commands (via cli flag) take precedence over inventory
	if app.commands != "" {
		cmds := strings.Split(app.commands, "::")

		for _, device := range i.Devices {
			device.SendCommands = cmds
		}
	}
	// user-provided username & passwords (via cli flag) take precedence over inventory for the default credentials
	app.credentials = map[string]*credentials{
		defaultName: {
			Username:          app.username,
			Password:          app.password,
			SecondaryPassword: app.password,
		},
	}

	return nil
}

func (app *appCfg) loadInventoryFromFlags(i *inventory) error {
	if app.platform == "" {
		return errNoPlatformDefined
	}

	if app.username == "" {
		return errNoUsernameDefined
	}

	if app.password == "" {
		return errNoPasswordDefined
	}

	if app.commands == "" {
		return errNoCommandsDefined
	}

	app.credentials = map[string]*credentials{
		defaultName: {
			Username:          app.username,
			Password:          app.password,
			SecondaryPassword: app.password,
		},
	}

	cmds := strings.Split(app.commands, "::")

	i.Devices = map[string]*device{}

	i.Devices[app.address] = &device{
		Platform:     app.platform,
		Address:      app.address,
		SendCommands: cmds,
	}

	return nil
}

// filterDevices will remove the devices which names do not match the passed filter.
func filterDevices(i *inventory, f string) {
	if f == "" {
		return
	}

	fRe := regexp.MustCompile(f)

	for n := range i.Devices {
		if !fRe.MatchString(n) {
			delete(i.Devices, n)
		}
	}
}
