package main

import (
	"fmt"
	"time"
)

type Type struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskString string
}

type resultStruct struct {
	taskRESULT []byte
}

func taskCreature(superChan chan Type, a Type) {
	for {
		ct := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков.
			ct = "Some error occured"
		}

		time.Sleep(time.Millisecond * 1000)                                             // это для того, чтобы id не повторялись
		superChan <- Type{cT: ct, id: int(time.Now().Unix()), taskString: a.taskString} // передаем таск на выполнение
	}
}

func taskWorker(a Type, r resultStruct) (Type, resultStruct) {
	for {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			r.taskRESULT = []byte("task has been successes")
		} else {
			r.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a, r
	}
}

func taskSorter(superChan chan Type, doneTasks chan Type, undoneTasks chan error, t Type, r resultStruct) {
	for t := range superChan {

		l, r := taskWorker(t, r)

		if string(r.taskRESULT[14:]) == "successes" {

			t.taskString = string(r.taskRESULT)
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id: %d \nCreation time: %s, \nExecution time: %s, \nError: %s", t.id, t.cT, l.fT, r.taskRESULT)
		}
	}
}

func Result(doneTasks chan Type, undoneTasks chan error) {
	go func() {
		for d := range doneTasks {
			fmt.Println("\nDone tasks:")
			fmt.Println(d)
		}
	}()

	for u := range undoneTasks {
		fmt.Println("\nErrors!")
		fmt.Println(u)
	}
}

func main() {
	superChan := make(chan Type, 1)
	doneTasks := make(chan Type, 1)
	undoneTasks := make(chan error, 1)

	go taskCreature(superChan, Type{})
	go taskWorker(Type{}, resultStruct{})
	go taskSorter(superChan, doneTasks, undoneTasks, Type{}, resultStruct{})
	Result(doneTasks, undoneTasks)
}
