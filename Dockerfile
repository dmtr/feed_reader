FROM golang
ARG PROJ_NAME=github.com/dmtr/fbreader
RUN go get github.com/huandu/facebook
ADD fbutil /go/src/$PROJ_NAME/fbutil
ADD web /go/src/$PROJ_NAME/web
RUN go build $PROJ_NAME/fbutil
RUN go install $PROJ_NAME/fbutil
RUN go build $PROJ_NAME/web
RUN go install $PROJ_NAME/web
