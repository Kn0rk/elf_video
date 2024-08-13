package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("minimal_elf")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write ELF header
	var elfHeader = make([]byte, 64)
	binary.BigEndian.PutUint32(elfHeader[0:], 0x7F454C46)
	elfHeader[4] = 2                                         // 64-bit
	elfHeader[5] = 1                                         // Little-endian
	elfHeader[6] = 1                                         // ELF version
	elfHeader[7] = 0                                         // ABI
	elfHeader[8] = 0                                         // ABI version
	binary.LittleEndian.PutUint16(elfHeader[16:], 2)         // Executable file
	binary.LittleEndian.PutUint16(elfHeader[18:], 0x3e)      // Intel 80386
	binary.LittleEndian.PutUint32(elfHeader[20:], 1)         // ELF version
	binary.LittleEndian.PutUint64(elfHeader[24:], 0x0400078) // Entry point
	binary.LittleEndian.PutUint64(elfHeader[32:], 64)        // Start of program headers
	binary.LittleEndian.PutUint64(elfHeader[40:], 0)         // Start of section headers
	binary.LittleEndian.PutUint32(elfHeader[48:], 0)         // Flags
	binary.LittleEndian.PutUint16(elfHeader[52:], 64)        // Size of this header
	binary.LittleEndian.PutUint16(elfHeader[54:], 56)        // Size of program headers
	binary.LittleEndian.PutUint16(elfHeader[56:], 1)         // Number of program headers
	binary.LittleEndian.PutUint16(elfHeader[58:], 64)        // Size of section headers
	binary.LittleEndian.PutUint16(elfHeader[60:], 0)         // Number of section headers
	binary.LittleEndian.PutUint16(elfHeader[62:], 0)         // Section header string table index

	programHeader := make([]byte, 56)
	binary.LittleEndian.PutUint32(programHeader[0:], 1)          // Type (LOAD)
	binary.LittleEndian.PutUint32(programHeader[4:], 5)          // Flags (R-X)
	binary.LittleEndian.PutUint64(programHeader[8:], 0)          // Offset
	binary.LittleEndian.PutUint64(programHeader[16:], 0x0400000) // Virtual address
	binary.LittleEndian.PutUint64(programHeader[24:], 0x0400000) // Physical address
	binary.LittleEndian.PutUint64(programHeader[32:], 0xa4)      // File size
	binary.LittleEndian.PutUint64(programHeader[40:], 0xa4)      // Memory size
	binary.LittleEndian.PutUint64(programHeader[48:], 0x1000)    // Alignment

	var program = []byte{0xb8, 0x04, 0x00, 0x00, 0x00, 0xbb, 0x01, 0x00, 0x00, 0x00, 0xb9, 0x97, 0x00, 0x40, 0x00, 0xba, 0x0e, 0x00, 0x00, 0x00, 0xcd, 0x80, 0xb8, 0x01, 0x00, 0x00, 0x00, 0x31, 0xdb, 0xcd, 0x80}

	var data_hello_world = append([]byte("Hello, World!"), 0)

	file.Write(elfHeader)
	file.Write(programHeader)
	file.Write(program)
	file.Write(data_hello_world)

}
