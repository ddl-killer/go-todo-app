package gotodoapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type TodoManager interface {
	Add(task string)
	Complete(i int) error
	Delete(i int) error
	Save(file string) error
	Load(file string) error
}

type todo struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

func (t *todo) Status() string {
	if t.Done {
		return green("done")
	}
	return blue("todo")
}

func (t *todo) MarkDone() {
	t.Done = true
	t.CompletedAt = time.Now()
}

func (t *todo) HasBeenCompleted() bool {
	return t.Done
}

type Todos []todo

func (ts *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{
				Align: simpletable.AlignCenter,
				Text:  "#",
			},
			{
				Align: simpletable.AlignCenter,
				Text:  "Task",
			},
			{
				Align: simpletable.AlignCenter,
				Text:  "Status",
			},
			{
				Align: simpletable.AlignLeft,
				Text:  "CreatedAt",
			},
			{
				Align: simpletable.AlignLeft,
				Text:  "CompletedAt",
			},
		},
	}

	var cells [][]*simpletable.Cell

	for i, t := range *ts {
		i++
		task := blue(t.Task)
		status := t.Status()
		cells = append(cells, *&[]*simpletable.Cell{
			{
				Text: fmt.Sprintf("%d", i),
			},
			{
				Text: task,
			},
			{
				Text: status,
			},
			{
				Text: t.CreatedAt.Format(time.RFC822),
			},
			{
				Text: t.CreatedAt.Format(time.RFC822),
			},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("You have %d pending todos", ts.CountPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (ts *Todos) CountPending() int {
	res := 0
	for _, t := range *ts {
		if !t.Done {
			res++
		}
	}
	return res
}


func (ts *Todos) Add(task string) {
	t := todo{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*ts = append(*ts, t)
}

func (ts *Todos) Complete(i int) error {
	t := *ts
	if i <= 0 || i > len(t) {
		return errors.New("invalid index")
	}
	t[i-1].MarkDone()
	return nil
}

func (ts *Todos) Delete(i int) error {
	t := *ts
	if i <= 0 || i > len(t) {
		return errors.New("invalid index")
	}
	*ts = append(t[:i-1], t[i:]...)
	return nil
}

func (ts *Todos) Save(fileName string) error {
	bytes, err := json.Marshal(ts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, bytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (ts *Todos) Load(fileName string) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	err = json.Unmarshal(file, ts)
	if err != nil {
		return err
	}
	return nil
}
