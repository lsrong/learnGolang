package main

// Sample program demonstrating decoupling with interface composition.
// 演示与接口组合解耦的示例程序。

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Data 数据实体
type Data struct {
	Line string
}

// Puller 拉取行为抽象。
type Puller interface {
	Pull(d *Data) error
}

// Storer 存储行为抽象。
type Storer interface {
	Store(d *Data) error
}

// PullStorer 定义同时具有拉取和存储的行为接口。通过接口组合完成
type PullStorer interface {
	Puller
	Storer
}

// =============================================================================

// Xenia 拉取数据操作体
type Xenia struct {
	Host    string
	Timeout time.Duration
}

func (*Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF
	case 5:
		return errors.New("error reading data from Xenia")
	default:
		d.Line = "data"
		fmt.Println("In: ", d.Line)
		return nil
	}
}

// Pillar 保存数据操作体。
type Pillar struct {
	Host    string
	Timeout time.Duration
}

func (*Pillar) Store(d *Data) error {
	fmt.Println("Out: ", d.Line)
	return nil
}

// System 接口组合的方式实现解耦。######### 通过接口来解耦，不依赖于具体的实现。
type System struct {
	Puller
	Storer
}

// pull 具体实例×Xenia替换成抽象接口Puller.
func pull(p Puller, data []Data) (int, error) {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store 具体实例×Pillar替换成抽象接口Storer.
func store(t Storer, data []Data) (int, error) {
	for i := range data {
		if err := t.Store(&data[i]); err != nil {
			return i, err
		}
	}
	return len(data), nil
}

// Copy 拷贝操作，先拉取数据，然后保存数据
// 将*System替换成接口PullStorer，达到抽象操作。
func Copy(ps PullStorer, batch int) error {
	data := make([]Data, batch)
	for {
		i, err := pull(ps, data)
		if i > 0 {
			if _, err := store(ps, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

func main() {
	// System 组合Xenia，Pillar分别实现行为接口Puller, Storer.
	sys := System{
		Puller: &Xenia{Host: "localhost:3000"},
		Storer: &Pillar{Host: "localhost:4000"},
	}
	batch := 3
	if err := Copy(&sys, batch); err != io.EOF {
		fmt.Println(err)
	}
}
