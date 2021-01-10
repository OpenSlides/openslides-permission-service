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

	UserID     *int `yaml:"user_id"`
	userID     int
	MeetingID  int `yaml:"meeting_id"`
	Permission string

	Payload map[string]interface{}
	Action  string

	IsAllowed *bool    `yaml:"is_allowed"`
	CanSee    []string `yaml:"can_see"`

	Cases []*Case
}

func (c *Case) walk(f func(*Case)) {
	f(c)
	for _, s := range c.Cases {
		s.walk(f)
	}
}

func (c *Case) test(t *testing.T) {
	if onlyTest := os.Getenv("TEST_CASE"); onlyTest != "" {
		onlyTest = strings.TrimPrefix(onlyTest, "TestCases/")
		if c.Name != onlyTest {
			return
		}
	}
	if c.IsAllowed != nil {
		c.testWrite(t)
	}
	if c.CanSee != nil {
		c.testRead(t)
	}
}

func (c *Case) loadDB() (map[string]json.RawMessage, error) {
	data := make(map[string]json.RawMessage)
	for dbKey, dbValue := range c.DB {
		parts := strings.Split(dbKey, "/")
		switch len(parts) {
		case 1:
			map1, ok := dbValue.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid type in db key %s: %T", dbKey, dbValue)
			}
			for rawID, rawObject := range map1 {
				id, ok := rawID.(int)
				if !ok {
					return nil, fmt.Errorf("invalid id type: got %T expected int", rawID)
				}
				field, ok := rawObject.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("invalid object type: got %T, expected map[string]interface{}", rawObject)
				}

				for fieldName, fieldValue := range field {
					fqfield := fmt.Sprintf("%s/%d/%s", dbKey, id, fieldName)
					bs, err := json.Marshal(fieldValue)
					if err != nil {
						return nil, fmt.Errorf("creating test db. Key %s: %w", fqfield, err)
					}
					data[fqfield] = bs
				}

				idField := fmt.Sprintf("%s/%d/id", dbKey, id)
				data[idField] = json.RawMessage(strconv.Itoa(id))
			}

		case 2:
			field, ok := dbValue.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid object type: got %T, expected map[string]interface{}", dbValue)
			}

			for fieldName, fieldValue := range field {
				fqfield := fmt.Sprintf("%s/%s/%s", parts[0], parts[1], fieldName)
				bs, err := json.Marshal(fieldValue)
				if err != nil {
					return nil, fmt.Errorf("creating test db. Key %s: %w", fqfield, err)
				}
				data[fqfield] = bs
			}

			idField := fmt.Sprintf("%s/%s/id", parts[0], parts[1])
			data[idField] = []byte(parts[1])

		case 3:
			bs, err := json.Marshal(dbValue)
			if err != nil {
				return nil, fmt.Errorf("creating test db. Key %s: %w", dbKey, err)
			}
			data[dbKey] = bs

			idField := fmt.Sprintf("%s/%s/id", parts[0], parts[1])
			data[idField] = []byte(parts[1])
		default:
			return nil, fmt.Errorf("invalid db key %s", dbKey)
		}

	}

	return data, nil
}

func (c *Case) service() (*permission.Permission, error) {
	data, err := c.loadDB()
	if err != nil {
		return nil, fmt.Errorf("loading database: %w", err)
	}

	// Make sure the user does exists.
	userFQID := fmt.Sprintf("user/%d", c.userID)
	if data[userFQID+"/id"] == nil {
		data[userFQID+"/id"] = []byte(strconv.Itoa(c.userID))
	}

	// Make sure, the user is in the meeting.
	meetingFQID := fmt.Sprintf("meeting/%d", c.MeetingID)
	data[meetingFQID+"/user_ids"] = jsonAddInt(data[meetingFQID+"/user_ids"], c.userID)

	// Create group with the user and the given permissions.
	data["group/1337/id"] = []byte("1337")
	data[meetingFQID+"/group_ids"] = []byte("[1337]")
	data["group/1337/user_ids"] = []byte(fmt.Sprintf("[%d]", c.userID))
	f := fmt.Sprintf("user/%d/group_$%d_ids", c.userID, c.MeetingID)
	data[f] = jsonAddInt(data[f], 1337)
	data["group/1337/meeting_id"] = []byte(strconv.Itoa(c.MeetingID))
	if c.Permission != "" {
		data["group/1337/permissions"] = []byte(fmt.Sprintf(`["%s"]`, c.Permission))
	}

	return permission.New(&dataProvider{data}), nil
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

	got, err := p.IsAllowed(context.Background(), c.Action, c.userID, dataList)
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

	got, err := p.RestrictFQFields(context.Background(), c.userID, c.FQFields)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	if len(got) != len(c.CanSee) {
		var gotFields []string
		for k, v := range got {
			if v {
				gotFields = append(gotFields, k)
			}
		}
		t.Errorf("Got %v, expected %v", gotFields, c.CanSee)
	}

	for _, f := range c.CanSee {
		if !got[f] {
			t.Errorf("Did not allow %s", f)
		}
	}
}

func (c *Case) initSub() {
	for i, s := range c.Cases {
		name := s.Name
		if name == "" {
			name = fmt.Sprintf("case_%d", i)
		}
		name = strings.ReplaceAll(name, " ", "_")
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

		s.userID = c.userID
		if s.UserID != nil {
			s.userID = *s.UserID
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

	name := strings.TrimPrefix(path, "../../tests/")
	if c.Name != "" {
		name += ":" + c.Name
	}
	c.Name = name

	if c.MeetingID == 0 {
		c.MeetingID = 1
	}
	if c.userID == 0 {
		c.userID = 1
	}

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
