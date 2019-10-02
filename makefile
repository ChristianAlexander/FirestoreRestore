bin:
	CGO_ENABLED=0 go build -o firestorerestore .

clean:
	rm FirestoreRestore

.PHONY: bin
