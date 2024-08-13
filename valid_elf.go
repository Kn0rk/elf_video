package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func elfHeader() []byte {
	// ELF Header
	elfHeaderArr := make([]byte, 64)
	// 5B7x 2H5I 6H
	//4s
	binary.BigEndian.PutUint32(elfHeaderArr[0:], 0x7F454C46) // Magic number ascii == elf
	// 5B
	elfHeaderArr[4] = 2 // 64-bit
	elfHeaderArr[5] = 1 // Little-endian
	elfHeaderArr[6] = 1 // ELF version
	elfHeaderArr[7] = 0 // ABI
	elfHeaderArr[8] = 0 // ABI version
	// 16
	binary.LittleEndian.PutUint16(elfHeaderArr[16:], 2)    // Executable file
	binary.LittleEndian.PutUint16(elfHeaderArr[18:], 0x3e) // Intel 80386
	binary.LittleEndian.PutUint32(elfHeaderArr[20:], 1)
	binary.LittleEndian.PutUint64(elfHeaderArr[24:], 0x00401000) // Entry point
	binary.LittleEndian.PutUint64(elfHeaderArr[32:], 64)         // Start of program headers
	binary.LittleEndian.PutUint64(elfHeaderArr[40:], 0x2108)     // Start of section headers
	binary.LittleEndian.PutUint32(elfHeaderArr[48:], 0)          // Flags
	binary.LittleEndian.PutUint16(elfHeaderArr[52:], 64)         // Size of this header
	binary.LittleEndian.PutUint16(elfHeaderArr[54:], 56)         // Size of program headers
	binary.LittleEndian.PutUint16(elfHeaderArr[56:], 3)          // Number of program headers
	binary.LittleEndian.PutUint16(elfHeaderArr[58:], 64)         // Size of section headers
	binary.LittleEndian.PutUint16(elfHeaderArr[60:], 6)          // Number of section headers
	binary.LittleEndian.PutUint16(elfHeaderArr[62:], 5)          // Section header string table index
	return elfHeaderArr
}

func program() []byte {
	return []byte{0xb8, 0x04, 0x00, 0x00, 0x00, 0xbb, 0x01, 0x00, 0x00, 0x00, 0xb9, 0x00, 0x20, 0x40, 0x00, 0xba, 0x0e, 0x00, 0x00, 0x00, 0xcd, 0x80, 0xb8, 0x01, 0x00, 0x00, 0x00, 0x31, 0xdb, 0xcd, 0x80}
}

type ProgramHeader struct {
	load                uint32
	flag                uint32
	file_offset         uint64
	virtual_address     uint64
	physical_address    uint64
	segment_byte_length uint64
	memory_size         uint64
	alignment           uint64
}

func (header ProgramHeader) toBytes() []byte {
	programHeader := make([]byte, 56)
	binary.LittleEndian.PutUint32(programHeader[0:], header.load)                 // Type (LOAD)
	binary.LittleEndian.PutUint32(programHeader[4:], header.flag)                 // Flags (R-X)
	binary.LittleEndian.PutUint64(programHeader[8:], header.file_offset)          // Offset
	binary.LittleEndian.PutUint64(programHeader[16:], header.virtual_address)     // Virtual address
	binary.LittleEndian.PutUint64(programHeader[24:], header.virtual_address)     // Physical address
	binary.LittleEndian.PutUint64(programHeader[32:], header.segment_byte_length) // File size
	// Memory size
	if header.memory_size == 0 {
		binary.LittleEndian.PutUint64(programHeader[40:], header.segment_byte_length)
	} else {
		binary.LittleEndian.PutUint64(programHeader[40:], header.memory_size)
	}
	// Alignment
	if header.alignment == 0 {
		binary.LittleEndian.PutUint64(programHeader[48:], 0x1000)
	} else {
		binary.LittleEndian.PutUint64(programHeader[48:], header.alignment)
	}

	return programHeader
}

type SectionHeader struct {
	name_offset     uint32
	section_type    uint32
	flags           uint64
	virtual_address uint64
	file_offset     uint64
	length          uint64
	link            uint32
	info            uint32
	align           uint64
	ent_size        uint64
}

func (sec SectionHeader) toBytes() []byte {
	programSection := make([]byte, 64)
	binary.LittleEndian.PutUint32(programSection[0:], sec.name_offset)      // Offset into string table
	binary.LittleEndian.PutUint32(programSection[4:], sec.section_type)     // Type e.g. symtab, strtab,...
	binary.LittleEndian.PutUint64(programSection[8:], sec.flags)            // Flags writable,...
	binary.LittleEndian.PutUint64(programSection[16:], sec.virtual_address) // Virtual address
	binary.LittleEndian.PutUint64(programSection[24:], sec.file_offset)     // file offset
	binary.LittleEndian.PutUint64(programSection[32:], sec.length)          // section size in file
	binary.LittleEndian.PutUint32(programSection[40:], sec.link)            // sh link
	binary.LittleEndian.PutUint32(programSection[44:], sec.info)            // sh info
	binary.LittleEndian.PutUint64(programSection[48:], sec.align)           // Alignment
	binary.LittleEndian.PutUint64(programSection[56:], sec.ent_size)        // ent size if s
	return programSection
}

func main() {
	file, err := os.Create("minimal_elf")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write ELF header

	var programHeaders = make([]ProgramHeader, 0)
	programHeaders = append(programHeaders,
		ProgramHeader{
			load:                1,
			flag:                4,
			file_offset:         0,
			virtual_address:     0x00400000,
			segment_byte_length: 0xe8,
		})
	programHeaders = append(programHeaders, ProgramHeader{
		load:                1,
		flag:                5,
		file_offset:         0x1000,
		virtual_address:     0x00401000,
		segment_byte_length: 0x1f,
	})
	programHeaders = append(programHeaders, ProgramHeader{
		load:                1,
		flag:                6,
		file_offset:         0x2000,
		virtual_address:     0x00402000,
		segment_byte_length: 0xd,
	})

	var data_hello_world = append([]byte("Hello, World!"), 0)

	var strtab = []byte{}
	strtab = append(strtab, append([]byte("min.asm"), 0)...)
	strtab = append(strtab, append([]byte("hello"), 0)...)
	strtab = append(strtab, append([]byte("__bss_start"), 0)...)
	strtab = append(strtab, append([]byte("_edata"), 0)...)
	strtab = append(strtab, append([]byte("_end"), []byte{0, 0}...)...)
	strtab = append(strtab, append([]byte(".symtab"), 0)...)
	strtab = append(strtab, append([]byte(".strtab"), 0)...)
	strtab = append(strtab, append([]byte(".shstrtab"), 0)...)
	strtab = append(strtab, append([]byte(".text"), 0)...)
	strtab = append(strtab, append([]byte(".data"), 0)...)

	var sections = make([]SectionHeader, 0)

	sections = append(sections,
		SectionHeader{
			// All 0
		})

	sections = append(sections, SectionHeader{
		name_offset:     0x1B,
		section_type:    0x1,
		flags:           6,
		virtual_address: 0x00401000,
		file_offset:     0x01000,
		align:           0x10,
		length:          uint64(len(program())),
	})

	sections = append(sections, SectionHeader{
		name_offset:     0x21,
		section_type:    0x1,
		flags:           3,
		virtual_address: 0x00402000,
		file_offset:     0x02000,
		align:           4,
		length:          uint64(len(data_hello_world)) - 1,
	})

	sections = append(sections, SectionHeader{
		name_offset:     0x01, // symtab
		section_type:    0x2,  // symtab
		flags:           0,
		virtual_address: 0x00,
		file_offset:     0x02010,
		align:           8,
		length:          0xa8,
		ent_size:        0x18,
		info:            3,
		link:            4,
	})

	sections = append(sections, SectionHeader{
		name_offset:     0x9, // strtab
		section_type:    0x3, //strtab
		flags:           0,
		virtual_address: 0x0,
		file_offset:     0x020b8,
		align:           1,
		length:          0x27,
	})

	sections = append(sections, SectionHeader{
		name_offset:     0x11, // shstrtab
		section_type:    0x3,  //strtab
		flags:           0,
		virtual_address: 0x0,
		file_offset:     0x020df,
		align:           1,
		length:          0x27,
	})

	var bytes = make([]byte, 0)
	bytes = append(bytes, elfHeader()...)
	for _, seg := range programHeaders {
		bytes = append(bytes, seg.toBytes()...)
	}
	bytes = append(bytes, make([]byte, 0x1000-len(bytes))...)

	bytes = append(bytes, program()...)
	bytes = append(bytes, make([]byte, 0x2000-len(bytes))...)
	bytes = append(bytes, data_hello_world...)
	bytes = append(bytes, make([]byte, 0x2010-len(bytes))...)

	bytes = append(bytes, make([]byte, 0x20b9-len(bytes))...)
	bytes = append(bytes, strtab...)
	// // skip symtab
	bytes = append(bytes, make([]byte, 0x2108-len(bytes))...)

	for _, sec := range sections {
		bytes = append(bytes, sec.toBytes()...)
	}

	file.Write(bytes)

}
