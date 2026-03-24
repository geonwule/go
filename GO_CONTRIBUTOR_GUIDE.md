# Go 오픈소스 분석 & 컨트리뷰터 가이드

> 현재 레포: Go 1.27 (개발 중) — `src/go.mod` 기준

---

## 1. 레포 전체 구조

```
go/
├── src/          # 핵심: 표준 라이브러리 + 툴체인 소스
├── cmd/          # (src/cmd/) 컴파일러, 링커, go 툴 등
├── api/          # 각 버전별 공개 API 목록 (호환성 검증용)
├── doc/          # 릴리즈 노트, 언어 스펙, 어셈블리 가이드
├── lib/          # 타임존 데이터, WASM 실행기, FIPS 모듈
├── misc/         # CGo 예제, Chrome 확장, Android/iOS 실행기
└── test/         # (src/testdata/) 컴파일러 동작 검증 테스트
```

---

## 2. src/ 핵심 디렉토리 분류

### 표준 라이브러리 (입문하기 좋은 영역)
| 패키지 | 설명 |
|--------|------|
| `src/fmt` | 포맷 출력/파싱 |
| `src/strings`, `src/bytes` | 문자열/바이트 조작 |
| `src/strconv` | 타입 변환 |
| `src/bufio` | 버퍼 I/O |
| `src/io` | I/O 인터페이스 |
| `src/errors` | 에러 처리 |
| `src/slices`, `src/maps` | 제네릭 컬렉션 유틸 (비교적 최신, 코드 적음) |
| `src/cmp` | 비교 유틸리티 |

### 네트워크 / 시스템
| 패키지 | 설명 |
|--------|------|
| `src/net` | TCP/UDP/HTTP 기반 |
| `src/os` | OS 인터페이스 |
| `src/syscall` | 시스템 콜 래퍼 |
| `src/sync` | 동기화 프리미티브 |
| `src/context` | 컨텍스트 전파 |

### 툴체인 (고급 영역)
| 디렉토리 | 설명 |
|----------|------|
| `src/cmd/compile` | Go 컴파일러 |
| `src/cmd/link` | 링커 |
| `src/cmd/go` | `go` CLI 툴 전체 |
| `src/cmd/asm` | 어셈블러 |
| `src/cmd/gofmt` | 코드 포매터 |
| `src/cmd/vet` | 정적 분석기 |
| `src/cmd/cgo` | C 인터롭 |

### 런타임 (가장 복잡한 영역)
| 파일 패턴 | 설명 |
|-----------|------|
| `src/runtime/proc.go` | 고루틴 스케줄러 (G/M/P 모델) |
| `src/runtime/malloc.go` | 메모리 할당기 |
| `src/runtime/mgc*.go` | GC (가비지 컬렉터) |
| `src/runtime/chan.go` | 채널 구현 |
| `src/runtime/map*.go` | 맵 구현 |
| `src/runtime/asm_*.s` | 아키텍처별 어셈블리 |

---

## 3. 빌드 & 테스트 방법

### 전체 빌드
```bash
cd src
./make.bash          # Go 툴체인 빌드
./all.bash           # 빌드 + 전체 테스트 (시간 오래 걸림)
```

### 특정 패키지 테스트
```bash
# 이미 설치된 go 바이너리가 있다면
go test ./src/strings/...
go test ./src/fmt/...

# 레이스 컨디션 검사 포함
go test -race ./src/sync/...
```

### 빌드 환경 변수 (make.bash 참고)
```bash
GOOS=linux GOARCH=amd64   # 크로스 컴파일 타겟
CGO_ENABLED=0              # CGo 비활성화
GO_GCFLAGS="-N -l"         # 최적화 끄기 (디버깅용)
```

---

## 4. 파일 네이밍 컨벤션

Go 레포는 빌드 태그 대신 파일명으로 플랫폼을 구분합니다.

```
foo.go                  # 모든 플랫폼
foo_linux.go            # Linux 전용
foo_linux_amd64.go      # Linux + amd64 전용
foo_test.go             # 테스트 파일
foo_amd64.s             # amd64 어셈블리
export_test.go          # 테스트에서 내부 심볼 노출용
```

---

## 5. 컨트리뷰션 프로세스

Go는 GitHub PR이 아닌 **Gerrit 코드 리뷰 시스템**을 사용합니다.

### 단계별 흐름
```
1. 이슈 확인/등록
   → https://github.com/golang/go/issues
   → "help wanted", "good first issue" 라벨 필터링

2. CLA 서명
   → https://cla.developers.google.com/

3. Gerrit 계정 설정
   → https://go-review.googlesource.com/

4. git-codereview 설치
   → go install golang.org/x/review/git-codereview@latest

5. 변경 작업
   → git checkout -b my-fix
   → 코드 수정
   → go test ./affected/package/...

6. 변경 제출
   → git codereview change    # 커밋 생성
   → git codereview mail      # Gerrit에 업로드

7. 리뷰 & 머지
   → 리뷰어가 LGTM + Approved 주면 자동 머지
```

### 커밋 메시지 형식
```
패키지명: 변경 내용 요약 (소문자, 명령형)

더 자세한 설명 (필요시)

Fixes #이슈번호
```
예시: `strings: fix IndexByte returning wrong offset for empty string`

---

## 6. 입문 전략 (추천 순서)

### 1단계: 작은 패키지부터 읽기
`src/slices`, `src/maps`, `src/cmp` — 코드가 적고 제네릭 활용 패턴을 볼 수 있음

### 2단계: 테스트 파일로 동작 이해
각 패키지의 `*_test.go`를 먼저 읽으면 API 사용법과 엣지 케이스를 빠르게 파악 가능

### 3단계: 이슈 트래커 탐색
```
https://github.com/golang/go/issues?q=label%3A"help+wanted"
https://github.com/golang/go/issues?q=label%3A"good+first+issue"
```

### 4단계: 기여 유형 선택
- 문서/주석 개선 (가장 쉬운 시작점)
- 테스트 케이스 추가
- 버그 수정
- 성능 개선
- 새 기능 (proposal 프로세스 필요)

---

## 7. 유용한 링크

| 리소스 | URL |
|--------|-----|
| 공식 기여 가이드 | https://go.dev/doc/contribute |
| Gerrit 코드 리뷰 | https://go-review.googlesource.com |
| 이슈 트래커 | https://github.com/golang/go/issues |
| Go 언어 스펙 | `doc/go_spec.html` |
| 어셈블리 가이드 | `doc/asm.html` |
| 런타임 해킹 가이드 | `src/runtime/HACKING.md` |
| 메모리 모델 | `doc/go_mem.html` |
| Go 위키 | https://go.dev/wiki |
| golang-dev 메일링 | https://groups.google.com/g/golang-dev |

---

## 8. 분석 시작점 추천

처음 코드를 읽는다면 이 순서를 추천합니다:

```
src/builtin/builtin.go     → 내장 함수/타입 문서 (실제 구현 아님, 문서용)
src/errors/errors.go       → 에러 인터페이스 구현 (매우 짧음)
src/cmp/cmp.go             → 제네릭 비교 (Go 1.21+, 코드 적음)
src/slices/slices.go       → 제네릭 슬라이스 유틸
src/strings/strings.go     → 문자열 처리 (실용적, 읽기 좋음)
src/sync/mutex.go          → 뮤텍스 구현 (동기화 이해)
src/runtime/chan.go        → 채널 내부 구현 (중급)
src/runtime/proc.go        → 스케줄러 (고급)
```
