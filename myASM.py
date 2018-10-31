#! /usr/bin/env python3

import sys

# A simple assembler for a CPU that we made in software.
# We where ment to do all of the machine code by hand.
# But I wrote this up do I didn't have to

if len(sys.argv) != 3:
	print("Usage\t file_IN file_OUT")
	raise ValueError

def toBin(number):
	out = bin(int(number, 16))[2:]
	return ((4-len(out))*"0") + out


# changes a register to a bin
def deref(value):
	out = bin(int(value[1:]))[2:]
	return ((4-len(out))*"0") + out
	

def NOP():
	return "0000"*4

def LD(cmd):
	out = '0001'
	out += deref(cmd[1])
	out += '0000'
	out += toBin(cmd[2])
	return out

def MOV(cmd):
	out = '0010'
	out += deref(cmd[1])
	out += deref(cmd[2])
	out += '0000'
	return out

def DISP(cmd):
	out = '0011'
	out += '0000'
	if len(cmd) == 1:
		out += '0000'
		out += '0001'
	else:
		out += deref(cmd[1])
		out += deref(cmd[2])
	return out

def SHOW(cmd, outPut):
	for letter in cmd[1]:
		letter = str(hex(ord(letter))[2:])
		push(LD(["", "r0", letter[:1]]), outPut)
		push(LD(["", "r1", letter[1:]]), outPut)
		push(DISP(["", "r0", "r1"]), outPut)


def ALU_OP(cmd):
	out = deref(cmd[1])
	out += deref(cmd[2])
	out += deref(cmd[3])
	return out

def push(out, fOut):
	if out != '':
		fOut.write(hex(int(out, 2))[2:] + '\n')


instructions = {"NOP" :'0000',
				"LD"  :'0001',
				"MOV" :'0010',
				"DISP":'0011', 
				"XOR" :'0100',
				"AND" :'0101',
				"OR"  :'0110',
				"ADD" :'0111',
				"XXXX":'1000', # JMP //Unconditinal jump
				"XXXX":'1001', # JEZ //Jumps when the last operation was 0
				"XXXX":'1010', # JGZ //Jumps when the last operation was >0
				"XXXX":'1011', # JLZ //Jumps when the last operation was <0
				"XXXX":'1100', # TEST //subtract without the WE bit set
				"XXXX":'1101',
				"XXXX":'1110',
				"SUB" :'1111'}

"""
4bit
SHOW A
=
LD r5 
LD r6
DISP r5, r6
"""

fileIn = sys.argv[1]
fileOut = sys.argv[2]

#fIn = open(fileIn, 'r')
#fOut = open(fileOut, 'w')
with open(fileIn, 'r') as fIn:
	with open(fileOut, 'w') as fOut:
		fOut.write("v2.0 raw\n")
		for line in fIn:
			out = ""
			cmd = (line.replace(',','').split())
			try:
				op = cmd[0].upper()
			except IndexError:
				op = ""

			if op == "NOP":
				out = NOP()
				
			elif op == "LD":
				out = LD(cmd)

			elif op == "DISP":
				out = DISP(cmd)

			elif op == "SHOW":
				SHOW(cmd, fOut)
				
			elif op == "XOR":
				out = '0100'+ALU_OP(cmd)
				
			elif op == "AND":
				out = '0101'+ALU_OP(cmd)
				
			elif op == "OR":
				out = '0110'+ALU_OP(cmd)
				
			elif op == "ADD":
				out = '0111'+ALU_OP(cmd)
			
			elif op == "SUB":
				out = '1111'+ALU_OP(cmd)

			#print(line)
			#print(out)

			if out != '':
				#print("'"+str(int(out,2))+"'")
				fOut.write(hex(int(out, 2))[2:]+'\n')








