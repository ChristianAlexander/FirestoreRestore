bin:
	CGO_ENABLED=0 go build .

clean:
	rm FirestoreRestore

.PHONY: bin
