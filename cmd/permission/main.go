package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/OpenSlides/openslides-permission-service/internal/definitions"
	"github.com/OpenSlides/openslides-permission-service/pkg/permission"
)

type fakeDataProvider struct {
	ctx context.Context
}

func (dp fakeDataProvider) Get(fqfields ...definitions.Fqfield) ([]json.RawMessage, error) {
	fmt.Printf("Access context: %p\n", &dp.ctx)
	log.Println("TEST")
	m := make([]json.RawMessage, len(fqfields))
	for i := range fqfields {
		m[i] = json.RawMessage(strconv.Itoa(i))
	}
	return m, nil
}

func (dp fakeDataProvider) setContext(ctx context.Context) {
	dp.ctx = ctx
}

func main() {
	myDataProvider := fakeDataProvider{}
	ps := permission.New(myDataProvider)

	myDataProvider.setContext(context.TODO())
	result, addition, err := ps.IsAllowed("", 0, nil)
	fmt.Println(result, addition, err)

	myDataProvider.setContext(context.TODO())
	result, addition, err = ps.IsAllowed("topic.create", 0, nil)
	fmt.Println(result, addition, err)

	myDataProvider.setContext(context.TODO())
	data := definitions.FqfieldData{
		"meeting_id": "1",
	}
	result, addition, err = ps.IsAllowed("topic.create", 0, data)
	fmt.Println(result, addition, err)

	myDataProvider.setContext(context.TODO())
	result, addition, err = ps.IsAllowed("topic.update", 0, nil)
	fmt.Println(result, addition, err)

	myDataProvider.setContext(context.TODO())
	result, addition, err = ps.IsAllowed("topic.delete", 0, nil)
	fmt.Println(result, addition, err)
}
