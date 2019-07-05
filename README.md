# castflow

castflow는 Lazypic 에서
캐릭터 Intellectual property rights (IPR) 을 관리하는 소프트웨어이다.

### 사용법
캐릭터 등록

```bash
$ castflow -add -id chocoala -email chocoala@lazypic.org -regnum C-2018-019061 -manager [담당자] -foa [활동범위] -concept [콘셉]
```

캐릭터 수정
```bash
$ castflow -set -id [id] -concept [수정내용]
```


캐릭터 삭제

```bash
$ sudo castflow -rm -id [id]
```

캐릭터 검색

```bash
$ castflow -search [검색어]
```
### 참고사항
대한민국에서 보호하는 지적재산관은 총 3가지이다.

- 산업재산권: 특허법, 실용신안법, 상표법, 디자인 보호법
- 저작권법: 문화,예술 부문
- 신지식재산권: 사회, 기술 변화에 따른 새로운 형태의 지식

### AWS DB권한 설정
AWS DB접근 권한을 설정할 계정에 아래 권한을 부여합니다.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ListAndDescribe",
            "Effect": "Allow",
            "Action": [
                "dynamodb:List*",
                "dynamodb:DescribeReservedCapacity*",
                "dynamodb:DescribeLimits",
                "dynamodb:DescribeTimeToLive"
            ],
            "Resource": "*"
        },
        {
            "Sid": "SpecificTable",
            "Effect": "Allow",
            "Action": [
                "dynamodb:BatchGet*",
                "dynamodb:DescribeStream",
                "dynamodb:DescribeTable",
                "dynamodb:Get*",
                "dynamodb:Query",
                "dynamodb:Scan",
                "dynamodb:BatchWrite*",
                "dynamodb:CreateTable",
                "dynamodb:Delete*",
                "dynamodb:Update*",
                "dynamodb:PutItem"
            ],
            "Resource": "arn:aws:dynamodb:*:*:table/castflow"
        }
    ]
}
```
