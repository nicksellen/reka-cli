oses = darwin linux windows solaris
arches = amd64

clean:
	@rm -rf build

build:
	@for os in $(oses); do \
		for arch in $(arches); do \
			echo "building $$os-$$arch"; \
			mkdir -p build/$$os-$$arch; \
			GOOS=$$os GOARCH=$$arch go build -o build/$$os-$$arch/reka; \
		done \
	done

s3: build
	@for os in $(oses); do \
		for arch in $(arches); do \
			echo "uploading $$os-$$arch"; \
			aws s3 \
				cp build/$$os-$$arch/reka s3://reka/$$os-$$arch/reka \
				--grants \
		      read=uri=http://acs.amazonaws.com/groups/global/AllUsers; \
		done \
	done
