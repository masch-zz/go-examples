package main

//Source: https://talks.golang.org/2013/bestpractices.slide#3

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

type Gopher struct {
	Name     string
	AgeYears int
}

//#########################################First method##########################################
func (g *Gopher) WriteTo(w io.Writer) (size int64, err error) {
	err = binary.Write(w, binary.LittleEndian, int32(len(g.Name)))
	if err == nil {
		size += 4
		var n int
		n, err = w.Write([]byte(g.Name))
		size += int64(n)
		if err == nil {
			err = binary.Write(w, binary.LittleEndian, int64(g.AgeYears))
			if err == nil {
				size += 8
			}
			return
		}
		return
	}
	return
}

//#############################Avoid Nesting By Handling Error First#############################
func (g *Gopher) WriteToAvoidNestingByHandlingErrorFirst(w io.Writer) (size int64, err error) {
	err = binary.Write(w, binary.LittleEndian, int32(len(g.Name)))
	if err != nil {
		return
	}
	size += 4

	n, err := w.Write([]byte(g.Name))
	size += int64(n)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.LittleEndian, int64(g.AgeYears))
	if err != nil {
		return
	}
	size += 8

	return
}

//################################Avoid Repetition When Possible#################################
type binWriter struct {
	w    io.Writer
	size int64
	err  error
}

// WriteAvoidRepetitionWhenPossible writes a value to the provided write in little endian form.
func (w *binWriter) WriteAvoidRepetitionWhenPossible(v interface{}) {
	if w.err != nil {
		return
	}
	if w.err = binary.Write(w.w, binary.LittleEndian, v); w.err == nil {
		w.size += int64(binary.Size(v))
	}
}

func (g *Gopher) WriteToAvoidRepetitionWhenPossible(w io.Writer) (int64, error) {
	bw := &binWriter{w: w}
	bw.WriteAvoidRepetitionWhenPossible(int32(len(g.Name)))
	bw.WriteAvoidRepetitionWhenPossible([]byte(g.Name))
	bw.WriteAvoidRepetitionWhenPossible(int64(g.AgeYears))
	return bw.size, bw.err
}

//##########################Type switch with short variable declaration##########################
func (w *binWriter) WriteTypeSwitchWithShortVariableDeclaration(v interface{}) {
	if w.err != nil {
		return
	}

	switch x := v.(type) {
	case string:
		w.WriteTypeSwitchWithShortVariableDeclaration(int32(len(x)))
		w.WriteTypeSwitchWithShortVariableDeclaration([]byte(x))
	case int:
		w.WriteTypeSwitchWithShortVariableDeclaration(int64(x))
	default:
		if w.err = binary.Write(w.w, binary.LittleEndian, v); w.err == nil {
			w.size += int64(binary.Size(v))
		}
	}
}

func (g *Gopher) WriteTypeSwitchWithShortVariableDeclaration(w io.Writer) (int64, error) {
	bw := &binWriter{w: w}
	bw.WriteTypeSwitchWithShortVariableDeclaration(g.Name)
	bw.WriteTypeSwitchWithShortVariableDeclaration(g.AgeYears)
	return bw.size, bw.err
}

//#################################Writing everything or nothing#################################

type binWriterBuffer struct {
	w   io.Writer
	buf bytes.Buffer
	err error
}

// Write writes a value to the provided writer in little endian form
func (w *binWriterBuffer) Write(v interface{}) {
	if w.err != nil {
		return
	}

	switch x := v.(type) {
	case string:
		w.Write(int32(len(x)))
		w.Write([]byte(x))
	case int:
		w.Write(int64(x))
	default:
		w.err = binary.Write(&w.buf, binary.LittleEndian, v)
	}
}

// Flush writes any pending values into the writer if no error has occurred.
// If an error has occurred, earlier or with a write by Flush, the error is
// returned.
func (w *binWriterBuffer) Flush() (int64, error) {
	if w.err != nil {
		return 0, w.err
	}

	return w.buf.WriteTo(w.w)
}

func (g *Gopher) WriteToWritingEverythingOrNothing(w io.Writer) (int64, error) {
	bw := &binWriterBuffer{w: w}
	bw.Write(g.Name)
	bw.Write(g.AgeYears)
	return bw.Flush()
}

func TestBuff(t *testing.T) {

	var expectedGopherSize int64 = 17
	gopherBuffer := new(bytes.Buffer)
	gopher := Gopher{"Bauer", 45}

	calculatedGopherSize, err := gopher.WriteTo(gopherBuffer)
	validateGopherBuffer(err, t, expectedGopherSize, calculatedGopherSize)

	calculatedGopherSize, err = gopher.WriteToAvoidNestingByHandlingErrorFirst(gopherBuffer)
	validateGopherBuffer(err, t, expectedGopherSize, calculatedGopherSize)

	calculatedGopherSize, err = gopher.WriteToAvoidRepetitionWhenPossible(gopherBuffer)
	validateGopherBuffer(err, t, expectedGopherSize, calculatedGopherSize)

	calculatedGopherSize, err = gopher.WriteTypeSwitchWithShortVariableDeclaration(gopherBuffer)
	validateGopherBuffer(err, t, expectedGopherSize, calculatedGopherSize)

	calculatedGopherSize, err = gopher.WriteToWritingEverythingOrNothing(gopherBuffer)
	validateGopherBuffer(err, t, expectedGopherSize, calculatedGopherSize)
}

func validateGopherBuffer(err error, t *testing.T, expectedGopherSize, actualGopherSize int64) {
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if expectedGopherSize != actualGopherSize {
		t.Fatalf("Buffers binaries sizes are not equals. Expected: %d - Actual: %d", expectedGopherSize, actualGopherSize)
	}
}
