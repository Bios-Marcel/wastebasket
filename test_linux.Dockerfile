# nanoserver doesn't have a shell, so we can't run the test, as we require
# shell systemcalls such as SHFileOperationW.
FROM docker.io/debian:trixie-slim 

WORKDIR /test
copy wastebasket.test ./test
RUN ./test

