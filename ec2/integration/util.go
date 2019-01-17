package integration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/shawncatz/automagical/ec2"
)

func loadEvent(name string) ec2.Event {
	evt := ec2.Event{}

	file, err := ioutil.ReadFile("../fixtures/" + name + ".json")
	if err != nil {
		//t.Error("could not read json file: ", err)
		return evt
	}

	err = json.Unmarshal(file, &evt)
	if err != nil {
		//t.Error("could not unmarshal json: ", err)
		return evt
	}

	return evt
}
