package tests

import "github.com/OpenSlides/openslides-permission-service/internal/types"

// HandlerStoreMock implements the types.HandlerStore interface.
type HandlerStoreMock struct {
	WriteHandler map[string]types.Writer
	ReadHandler  map[string]types.Reader
}

// RegisterReadHandler registers a read handler.
func (m *HandlerStoreMock) RegisterReadHandler(name string, reader types.Reader) {
	if m.ReadHandler == nil {
		m.ReadHandler = make(map[string]types.Reader)
	}
	m.ReadHandler[name] = reader
}

// RegisterWriteHandler registers a write handler.
func (m *HandlerStoreMock) RegisterWriteHandler(name string, writer types.Writer) {
	if m.WriteHandler == nil {
		m.WriteHandler = make(map[string]types.Writer)
	}
	m.WriteHandler[name] = writer

}
