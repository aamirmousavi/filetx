package filetx

import (
	"bytes"
	"fmt"
	"os"
)

type FileOperation struct {
	filePath       string
	buffer         *bytes.Buffer
	originalExists bool
	originalData   []byte
}

type FileTransaction struct {
	operations []*FileOperation
	committed  bool
}

func Begin() (*FileTransaction, error) {
	return &FileTransaction{
		operations: []*FileOperation{},
	}, nil
}

func (tx *FileTransaction) Create(filePath string) (*FileOperation, error) {
	op := &FileOperation{
		filePath: filePath,
		buffer:   new(bytes.Buffer),
	}
	if _, err := os.Stat(filePath); err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err == nil {
		op.originalExists = true
		op.originalData, err = os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
	}
	tx.operations = append(tx.operations, op)
	return op, nil
}

func (op *FileOperation) Write(data []byte) (int, error) {
	return op.buffer.Write(data)
}

func (tx *FileTransaction) Commit() error {
	if tx.committed {
		return nil
	}
	for _, op := range tx.operations {
		err := os.WriteFile(op.filePath, op.buffer.Bytes(), 0644)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	tx.committed = true
	return nil
}

func (tx *FileTransaction) Rollback() error {
	if tx.committed {
		return nil
	}
	for _, op := range tx.operations {
		if op.originalExists {
			if err := os.WriteFile(op.filePath, op.originalData, 0644); err != nil {
				return err
			}
		} else {
			if err := os.Remove(op.filePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tx *FileTransaction) Close() error {
	if !tx.committed {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}
	return nil
}
