all: diagnotor


diagnotor:
	go build -o diagnotor core/main.go core/diagnose.go


clean:
	rm -rf diagnotor


