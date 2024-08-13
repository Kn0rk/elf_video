; nasm -f elf64 -o min.o min.asm && ld -o min_asm min.o && ./min_asm

section .data 
    hello db 'Hello, World!'; 

section .text 
    global _start

_start:

    ; Write "Hello, World!" to stdout
    mov eax, 4          ; System call number for sys_write (4)
    mov ebx, 1          ; File descriptor 1 is stdout
    mov ecx, hello      ; Pointer to the string to write
    mov edx, 14         ; Length of the string (13 characters + newline)
    int 0x80            ; Call the kernel
    ; Exit the program
    mov eax, 1          ; System call number for sys_exit (1)
    xor ebx, ebx        ; Exit code 0
    int 0x80            ; Call the kernel


