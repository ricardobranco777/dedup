FROM	golang as builder

WORKDIR	/go/src/dedup
COPY	. .

RUN	make

FROM	scratch
COPY	--from=builder /go/src/dedup/dedup /usr/local/bin/dedup

ENTRYPOINT ["/usr/local/bin/dedup"]
