package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/OpenSlides/openslides-permission-service/pkg/permission"
	"gopkg.in/yaml.v3"
)

// Case object for testing.
type Case struct {
	Name     string
	DB       map[string]interface{}
	FQFields []string

	UserID     int
	MeetingID  int
	Permission string

	Payload map[string]interface{}
	Action  string

	IsAllowed  *bool `yaml:"is_allowed"`
	Restricted []string

	Cases []*Case
}

func (c *Case) walk(f func(*Case)) {
	f(c)
	for _, s := range c.Cases {
		s.walk(f)
	}
}

func (c *Case) test(t *testing.T) {
	if c.IsAllowed != nil {
		c.testWrite(t)
	}
	if c.Restricted != nil {
		c.testRead(t)
	}
}

func (c *Case) service() (*permission.Permission, error) {
	data := make(map[string]json.RawMessage)
	for k, v := range c.DB {
		bs, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("creating test db. Key %s: %w", k, err)
		}
		data[k] = bs

		parts := strings.Split(k, "/")
		idField := fmt.Sprintf("%s/%s/id", parts[0], parts[1])
		data[idField] = []byte(parts[1])
	}

	// Make sure the user does exists.
	c.UserID = defaultInt(c.UserID, 1337)
	meetingID := defaultInt(c.MeetingID, 1)
	userFQID := fmt.Sprintf("user/%d", c.UserID)
	if data[userFQID+"/id"] == nil {
		data[userFQID+"/id"] = []byte(strconv.Itoa(c.UserID))
	}

	// Make sure, the user is in the meeting.
	meetingFQID := fmt.Sprintf("meeting/%d", meetingID)
	data[meetingFQID+"/user_ids"] = jsonAddInt(data[meetingFQID+"/user_ids"], c.UserID)

	// Create group with the user and the given permissions.
	data["group/1337/id"] = []byte("1337")
	data[meetingFQID+"/group_ids"] = []byte("[1337]")
	data["group/1337/user_ids"] = []byte(fmt.Sprintf("[%d]", c.UserID))
	f := fmt.Sprintf("user/%d/group_$%d_ids", c.UserID, meetingID)
	data[f] = jsonAddInt(data[f], 1337)
	data["group/1337/meeting_id"] = []byte(strconv.Itoa(meetingID))
	if c.Permission != "" {
		data["group/1337/permissions"] = []byte(fmt.Sprintf(`["%s"]`, c.Permission))
	}

	return permission.New(&TestDataProvider{data}), nil
}

func (c *Case) testWrite(t *testing.T) {
	p, err := c.service()
	if err != nil {
		t.Fatalf("Can not create permission service: %v", err)
	}

	payload := make(map[string]json.RawMessage, len(c.Payload))
	for k, v := range c.Payload {
		bs, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("Invalid Payload: %v", err)
		}
		payload[k] = bs

	}
	dataList := []map[string]json.RawMessage{payload}

	got, err := p.IsAllowed(context.Background(), c.Action, c.UserID, dataList)
	if err != nil {
		t.Fatalf("IsAllowed retuend unexpected error: %v", err)
	}

	if got != *c.IsAllowed {
		t.Errorf("Got %t, expected %t", got, *c.IsAllowed)
	}
}

func (c *Case) testRead(t *testing.T) {
	p, err := c.service()
	if err != nil {
		t.Fatalf("Can not create permission service: %v", err)
	}

	got, err := p.RestrictFQFields(context.Background(), c.UserID, c.FQFields)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	if len(got) != len(c.Restricted) {
		var gotFields []string
		for k, v := range got {
			if v {
				gotFields = append(gotFields, k)
			}
		}
		t.Errorf("Got %v, expected %v", gotFields, c.Restricted)
	}

	for _, f := range c.Restricted {
		if !got[f] {
			t.Errorf("Did not allow %s", f)
		}
	}
}

func (c *Case) initSub() {
	for i, s := range c.Cases {
		name := s.Name
		if name == "" {
			name = fmt.Sprintf("case %d", i)
		}
		s.Name = c.Name + ":" + name

		db := make(map[string]interface{})
		for k, v := range c.DB {
			db[k] = v
		}
		for k, v := range s.DB {
			db[k] = v
		}
		s.DB = db

		fields := append([]string{}, c.FQFields...)
		fields = append(fields, s.FQFields...)
		s.FQFields = fields

		if s.UserID == 0 {
			s.UserID = c.UserID
		}
		if s.MeetingID == 0 {
			s.MeetingID = c.MeetingID
		}
		if s.Permission == "" {
			s.Permission = c.Permission
		}
		if s.Payload == nil {
			s.Payload = c.Payload
		}
		if s.Action == "" {
			s.Action = c.Action
		}

		s.initSub()
	}
}

func walk(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() == false && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			files = append(files, path)
		}
		return nil

	})
	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", path, err)
	}
	return files, nil
}

func loadFile(path string) (*Case, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	var c Case
	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}

	name := path
	if c.Name != "" {
		name += ":" + c.Name
	}
	c.Name = name

	c.initSub()

	return &c, nil
}

// defaultInt returns returns the given value or d, if value == 0
func defaultInt(value int, d int) int {
	if value == 0 {
		return d
	}
	return value
}

// jsonAddInt adds the given int to the encoded json list.
//
// If the value exists in the list, the list is returned unchanged.
func jsonAddInt(list json.RawMessage, value int) json.RawMessage {
	var decoded []int
	if list != nil {
		json.Unmarshal(list, &decoded)
	}

	for _, i := range decoded {
		if i == value {
			return list
		}
	}

	decoded = append(decoded, value)
	list, _ = json.Marshal(decoded)
	return list
}