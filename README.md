# Channel Page Collector

유튜브 채널 목록 페이지를 스크래핑하여 키워드별로 채널에 접근할 수 있는 **subscriber name(ex. @mango1135)** 을 추출합니다.<br/>
페이지 접근 후 스크롤 다운하여 새로운 데이터를 계속 반환받을 수 있으며, **gRPC stream response** 기능을 이용하여 스크롤 다운 할 때 마다 바인딩되는 채널 subscriber name 목록을 클라이언트가 값을 반환받을 수 있습니다. <br/>

만약 "economics" 키워드로 채널 데이터들을 수집하려 한다면 아래와 같은 URL에 접근하게 됩니다.<br/>
https://www.youtube.com/results?search_query=economics&sp=EgIQAg%253D%253D

![Screenshot 2025-02-23 at 8.59.51 PM.png](public/image/Screenshot%202025-02-23%20at%208.59.51%E2%80%AFPM.png)

## gRPC
- [Proto 호출 명세](https://github.com/Sujin1135/channel-page-collector-interface/blob/main/protobuf/channel_page.proto)
![Screenshot 2025-02-23 at 9.04.36 PM.png](public/image/Screenshot%202025-02-23%20at%209.04.36%E2%80%AFPM.png)
## Requirements
|    | Main version  |
|----|---------------|
| Go | 1.24+         |

## Run

**Execute**
```bash
## build
go build -o channel-page-collector cmd/main.go

## execute
./channel-page-collector
```

**Generate pb.go files from protobuf**
```bash
protoc \
  --proto_path=api/grpc/interface/channel-page-collector-interface/protobuf \
  --go_out=paths=source_relative:api/grpc/interface/channel-page-collector-interface/protobuf \
  --go-grpc_out=paths=source_relative:api/grpc/interface/channel-page-collector-interface/protobuf \
  api/grpc/interface/channel-page-collector-interface/protobuf/*.proto
```
