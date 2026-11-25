package interpreter

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type memoryDump struct {
	XMLName xml.Name     `xml:"memory"`
	Cells   []memoryCell `xml:"cell"`
}

type memoryCell struct {
	Address int   `xml:"address,attr"`
	Value   int32 `xml:",chardata"`
}

func WriteXMLDump(path string, memory []int32, start, end int) error {
	if start < 0 || end < start || end >= len(memory) {
		return fmt.Errorf(
			"invalid memory dump range %d..%d for memory size %d",
			start,
			end,
			len(memory),
		)
	}

	dump := memoryDump{Cells: make([]memoryCell, 0, end-start+1)}
	for address := start; address <= end; address++ {
		dump.Cells = append(dump.Cells, memoryCell{
			Address: address,
			Value:   memory[address],
		})
	}
	encoded, err := xml.MarshalIndent(dump, "", "  ")
	if err != nil {
		return fmt.Errorf("encode XML memory dump: %w", err)
	}
	encoded = append([]byte(xml.Header), encoded...)
	encoded = append(encoded, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create dump directory: %w", err)
	}
	if err := os.WriteFile(path, encoded, 0644); err != nil {
		return fmt.Errorf("write XML memory dump: %w", err)
	}
	return nil
}
