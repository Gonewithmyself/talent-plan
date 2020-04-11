.PHONY: all

all: test_example test_homework cleanup gendata

profile_example:
	go test -v -run=TestExampleURLTop -cpuprofile=excpu.out -memprofile=exmem.out

test_example:
	go test -v -run=TestExampleURLTop

test_homework:
	go test -v -run=TestURLTop 

profile_homework:
	go test -v -run=TestURLTop -cpuprofile=cpu.out -memprofile=mem.out

cleanup:
	go test -v -run=TestCleanData

gendata:
	go test -v -run=TestGenData
