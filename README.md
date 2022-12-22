# Terminology
* product : 상품 고유 속성을 관리하기 위한 엔티티입니다
    - 기본적으로 상품의 제목, 상세설명, 상품상태, 상품 기본 가격, Author 등을 저장합니다
    - 상품에 부수적인 리뷰, 파일정보, 저작권 정보 등을 추가로 가집니다
    - (중요) 상품은 상품의 고유의 속성으로써, 마켓에서 거래될 수 없는 단위입니다.
    

* deal : 마켓에서 판매를 위한 엔티티입니다
    - editor의 review가 끝난 해당상품에 판매를 위한 가격, 할인정보, 프로모션 정보 등을 추가하여 생성합니다
    - product를 바로 거래하지 않고 deal을 쓰는 이유는 아래와 같습니다
           - 상품 속성과 판매속성의 분리하여 데이터를 효율적으로 관리할 수 있습니다. 사용자의 데이터에 대한 실수를 방지 할 수 있습니다.
           - 하나의 상품을 판매 채널, 프로모션에 따라 다른 가격으로 판매할 수 있습니다.
           - 상품을 주문에서 바로 활용할 경우, 판매정보 변경에 따른 내역을 추적함에 어려움이 있습니다. 딜을 사용하면 해당 문제가 발생하지 않습니다. 

 
# Project Source Structure
* main.go : program의 시작
    - 마켓에서 판매가능한 형태입니다
* config : 앱 구동을 위한 config를 처리
* framework/
    - db : data의 physical interface를 담당(이 프로젝트에서는 간단히 file을 db로 사용)
    - gin : 이 proj.는 golang의 gin web framework로 구현되어 있음. Middleware처리를 담당
    - struct : 이 앱의 입력과 출력을 표준화한 구조체를 정의
    - rbac : 이 앱에 role base access control matrix를 가짐
* function : 공통 기능을 모아서, 함수단위로 handler에게 제공
* data : 간단한 fileDB로써, 데이터를 저장
* model : 데이터를 manage하기 위한 모듈
* handler : 비지니스를 처리하기 위한 모듈
* test : postman을 활용하여 테스트 진행
    - test_case.toml : test코드를 toml형태로 구현
    - upload_to_postman.go : test case를 읽어 postman에 업로드하는 스크립트

# How to run this app
* 사용자 환경에 맞게 go환경을 설치합니다.
https://judo0179.tistory.com/81
* 실행
```
cd {project_root folder}
go mod tidy
go run main.go
```
- test : postman의 테스트케이스를 사용하여 테스트 하실 수 있습니다.
https://www.postman.com/alankim/workspace/acon3d-assessment/overview
