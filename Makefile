all: midimap midimap.map.5.pdf

midimap: midimap.go
	go build midimap.go

midimap.map.5.pdf: midimap.map.5
	mandoc -Tpdf midimap.map.5 > midimap.map.5.pdf
