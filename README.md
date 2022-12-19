# requireement
- 각 상품은 ( 제목 / 본문 / 판매 가격 / 수수료 )을 포함합니다.
    - 제목과 본문은 언어별로 여러 쌍이 존재합니다. (한국어, 영어, 중국어)
    - 판매 가격은 각 나라의 환율을 적용하여 계산됩니다.
    - 수수료는 에디터가 결정합니다.
- 3가지 종류의 사용자가 있습니다.
    - 작가: Acon3d에서 판매를 희망하고, 상품을 작성하여 올립니다.
    - 에디터: 작가가 업로드한 상품을 검토합니다. 에디터에게 승인된 작품만이 Acon3d 쇼핑몰에서 노출됩니다.
    - 고객: Acon3d에서 상품들을 둘러보고 구매합니다.
- 상품은 다음과 같은 단계를 거쳐서 쇼핑몰에 등록됩니다.
    1. 작가가 새로운 상품을 "한국어로" **작성**하고, **검토를 요청**합니다.
    2. 에디터가 검토 요청이 들어온 상품들을 확인하여, 읽어보고 수정합니다.
    3. 에디터가 **검토를 완료**한 이후에는 쇼핑몰에 노출이 됩니다. 이 때, 쇼핑몰에 설정된 언어에 맞게 상품 정보가 표시되어야 합니다.
- 가끔은 상품 설명을 작성하다가 실수를 할 때가 있습니다.
    - 에디터는 쇼핑몰에 노출되고 있는 상품을 수정할 수 있는 방법이 필요합니다.
    - 작가들이 상품을 수정하고 싶을 때에는, 에디터에게 수정 요청을 할 수 있어야 합니다.

# requirements
- "작가"가 상품을 최초로 등록
     - POST /product
- "에디터"가 검토가 필요한 상품들을 조회
     - GET /product
- "작가" 및 "에디터"가 상품을 수정
     - PUT /product
- 검토를 완료하여 쇼핑몰에 노출
     - POST /deal
- 쇼핑몰에 노출되고 있는 상품들을 조회
     - GET /deal?lang="KR|JP|CN"