package wecom

import (
	"encoding/xml"
)

func ParseEvent(body []byte) (*CallbackEvent, error) {
	var evt CallbackEvent
	err := xml.Unmarshal(body, &evt)
	if err != nil {
		return nil, err
	}
	return &evt, nil
}

func ParseBatchJob(b []byte) ([]BatchJob, error) {
	var a = make(map[string][]byte)
	err := xml.Unmarshal(b, &a)
	if err != nil {
		return nil, err
	}
	var r []BatchJob
	if v, ok := a["BatchJob"]; ok {
		err = xml.Unmarshal(v, &r)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}
