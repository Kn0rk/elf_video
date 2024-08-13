# Creating ELF executable from scratch

## Prerequisites

- Go (Golang):  Installation instructions: https://golang.org/doc/install
- NASM (Netwide Assembler): sudo apt-get install nasm



## File Descriptions

### min.asm

This creates our reference executable so we can explore a binary file that is properly structured.

To compile and run:
```
nasm -f elf64 -o min.o min.asm
ld -o min_asm min.o
./min_asm
```

### minimal_elf.go

This Go program creates a minimal ELF file. It will only create the parts were absolutely necessary to run a executable:
- ELF header
- Program header
- Instructions/Data

To run:
```
go run min.go
chmod +x minimal_elf
./minimal_elf
```


### valid_elf.go

This is  creates an (almost) valid elf file by copying the structure from the min.asm nasm output.  Therefore it will work with most utilities like objdump or readelf.

To run:
```
go run valid_elf.go
```
This will overwrite `minimal_elf`.

## Utility programs

- objdump
- readelf
- vbindiff


## Video

[Video explanation](https://youtu.be/HyKyn8Zbj24)

## Note

This project is for educational purposes only. The created ELF files are minimal examples and should not be used as templates for production executables.