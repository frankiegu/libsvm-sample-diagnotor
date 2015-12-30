all: diagnotor


diagnotor:
	go build -o diagnotor core/main.go core/diagnose.go core/config.go


clean:
	rm -rf diagnotor


