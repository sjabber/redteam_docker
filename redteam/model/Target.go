package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
)

// Target(훈련대상)을 관리하기 위한 json 구조체
type Target struct {
	FakeNo           int       `json:"fake_no"`
	TargetNo         int       `json:"tg_no"`
	TargetName       string    `json:"tg_name"`
	TargetEmail      string    `json:"tg_email"`
	TargetPhone      string    `json:"tg_phone"`
	TargetOrganize   string    `json:"tg_organize"` //소속
	TargetPosition   string    `json:"tg_position"` //직급
	TargetTag        [3]string `json:"tg_tag"`      //태그 내보낼 때 사용
	TargetCreateTime string    `json:"created_t"`
	TagArray         []string  `json:"tag_no"` // 태그 입력받을 때 사용
}

// 삭제할 Target(훈련대상)의 시퀀스 넘버를 프론트엔드로 부터 받아오기 위한 변수
type TargetNumber struct {
	TargetNumber []string `json:"target_list"` //front javascript 와 이름을 일치시켜야함.
}

type Tag struct {
	TagNo         int    `json:"tag_no"`
	TagName       string `json:"tag_name"`
	TagCreateTime string `json:"created_t"`
}

// 해시테이블에 해당하는 태그 값이 들어가 있는지 점검할때 사용하는 메서드
func isValueIn(value string, list map[int]string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func (t *Target) CreateTarget(conn *sql.DB, num int) (int, error) {

	t.TargetName = strings.Trim(t.TargetName, " ")
	t.TargetEmail = strings.Trim(t.TargetEmail, " ")
	t.TargetPhone = strings.Trim(t.TargetPhone, " ")
	t.TargetOrganize = strings.Trim(t.TargetOrganize, " ")
	t.TargetPosition = strings.Trim(t.TargetPosition, " ")

	if len(t.TargetEmail) < 1 {
		return 400, fmt.Errorf("Target's E-mail is empty ")
	} else if len(t.TargetName) < 1 {
		return 400, fmt.Errorf("Target's name is empty ")
	}

	//else if len(t.TargetPhone) < 1 {
	//	errcode = 400
	//	return errcode, fmt.Errorf(" Target's Phone number is empty ")
	//} else if len(t.TargetOrganize) < 1 {
	//	errcode = 400
	//	return errcode, fmt.Errorf(" Target's Organize is empty ")
	//} else if len(t.TargetPosition) < 1 {
	//	errcode = 400
	//	return errcode, fmt.Errorf(" Target's Position is empty ")
	//} else if len(t.TagArray) < 1 {
	//	return errcode, fmt.Errorf(" Target's Tag is empty ")
	//}

	// 이메일 형식검사
	var validEmail, _ = regexp.MatchString(
		"^[_A-Za-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", t.TargetEmail)

	if validEmail != true {
		return 402, fmt.Errorf("Email format is incorrect. ")
	}

	// 이름 형식검사 (한글, 영어 이름만 허용)
	var validName, _ = regexp.MatchString("^[가-힣A-Za-z0-9\\s]{1,30}$", t.TargetName)

	if validName != true {
		return 402, fmt.Errorf("Name format is incorrect. ")
	}

	// 핸드폰 형식검사, 안넣거나 형식에 맞게 넣거나.
	if len(t.TargetPhone) > 0 {
		var phoneNumber, _ = regexp.MatchString(
			"^[0-9]{9,11}$", t.TargetPhone)

		if phoneNumber != true {
			return 402, fmt.Errorf("Phone number format is incorrect. ")
		}
	}

	// 등록된 대상자 수를 조회한다.
	row := conn.QueryRow(`SELECT count(target_no)
								FROM target_info
								WHERE user_no = $1;`, num)
	err := row.Scan(&t.TargetNo)
	if err != nil {
		return 500, fmt.Errorf("%v", err)
	}

	// 등록된 대상자 수 검사 (405에러)
	if t.TargetNo >= 300 {
		return 405, fmt.Errorf(" The target is already full. ")
	}

	// 태그 중복제거
	keys := make(map[string]bool)
	ue := []string{}

	for _, value := range t.TagArray {
		if _, saveValue := keys[value]; !saveValue { // 중복제거 핵심포인트

			keys[value] = true
			ue = append(ue, value)
		}
	}

	t.TagArray = nil
	t.TagArray = ue

	// t.TagArray 값이 비어있으면 에러나는 관계로 값을 채워준다.
	for i := 1; i <= 3; i++ {
		if len(t.TagArray) < i {
			t.TagArray = append(t.TagArray, "0")
		}
	}

	// 엑셀파일의 중간에 값이 없는 경우, 잘못된 형식이 들어가 있을경우 이를 검사할 필요가 있음.

	query1 := "INSERT INTO target_info (target_name, target_email, target_phone, target_organize, target_position," +
		"tag1, tag2, tag3, user_no) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	_, err = conn.Exec(query1, t.TargetName, t.TargetEmail, t.TargetPhone, t.TargetOrganize, t.TargetPosition,
		t.TagArray[0], t.TagArray[1], t.TagArray[2], num)
	if err != nil {
		fmt.Println(err)
		return 500, fmt.Errorf("Target create error. ")
	}

	return 200, nil
}

func ReadTarget(conn *sql.DB, num int, page int) ([]Target, int, int, error) {
	var pageNum int // 몇번째 페이지부터 가져올지 결정하는 변수
	var pages int   // 총 페이지 수
	var total int   // 총 훈련대상자들의 수를 담을 변수

	// ex) 1페이지 -> 1~10, 2페이지 -> 11~20
	// 페이지번호에 따라 가져올 목록이 달라진다.
	pageNum = (page - 1) * 20

	// 대상목록들을 20개씩만 잘라서 반하여 페이징처리한다.
	query := `
    SELECT
       row_num,
       target_name,
       target_email,
       target_phone,
       target_organize,
       target_position,
       COALESCE(tag_name1, '') as tag_name1, -- Null 일 경우 공백으로 대체
       COALESCE(tag_name2, '') as tag_name2, -- Null 일 경우 공백으로 대체
       COALESCE(tag_name3, '') as tag_name3, -- Null 일 경우 공백으로 대체
       to_char(modified_time, 'YYYY-MM-DD'),
       target_no
    FROM (SELECT ROW_NUMBER() over (ORDER BY target_no) AS row_num,
             target_no,
             target_name,
             target_email,
             target_phone,
             target_organize,
             target_position,
             tag1,
             tag2,
             tag3,
             modified_time,
             user_no
          FROM target_info
          WHERE user_no = $1
         ) AS T
        left join (select tag_name as tag_name1, user_no, tag_no
                    from tag_info
                    where user_no = $1) ti1 on ti1.tag_no = T.tag1
        left join (select tag_name as tag_name2, user_no, tag_no
                    from tag_info
                    where user_no = $1) ti2 on ti2.tag_no = T.tag2
        left join (select tag_name as tag_name3, user_no, tag_no
                    from tag_info
                    where user_no = $1) ti3 on ti3.tag_no = T.tag3
    WHERE row_num > $2
    ORDER BY target_no asc
    LIMIT 20; -- 개수 20개 제한
`

	// 조건에 맞는 데이터를 조회한다.
	rows, err := conn.Query(query, num, pageNum)
	if err != nil {
		return nil, 0, 0, fmt.Errorf(err.Error())
	}

	defer conn.Close()

	var targets []Target
	tg := Target{}

	// 목록들을 하나하나 읽어들여온다.
	for rows.Next() {
		err = rows.Scan(&tg.FakeNo, &tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tg.TargetTag[0], &tg.TargetTag[1], &tg.TargetTag[2], &tg.TargetCreateTime, &tg.TargetNo)
		if err != nil {
			SugarLogger.Error(err.Error())
			continue
		}

		// 프론트단에서 처리하도록 수정 완료.
		//var sub [3]string
		//phone := []rune(tg.TargetPhone)
		//
		//if len(tg.TargetPhone) < 10 {
		//	sub[0] = string(phone[0:2])
		//	sub[1] = string(phone[2:5])
		//	sub[2] = string(phone[5:9])
		//	tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		//} else if string(phone[1:2]) == "2" && len(tg.TargetPhone) == 10 {
		//	sub[0] = string(phone[0:2])
		//	sub[1] = string(phone[2:6])
		//	sub[2] = string(phone[6:10])
		//	tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		//} else if len(tg.TargetPhone) == 10 {
		//	sub[0] = string(phone[0:3])
		//	sub[1] = string(phone[3:6])
		//	sub[2] = string(phone[6:10])
		//	tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		//} else if len(tg.TargetPhone) == 11 {
		//	sub[0] = string(phone[0:3])
		//	sub[1] = string(phone[3:7])
		//	sub[2] = string(phone[7:11])
		//	tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		//}

		targets = append(targets, tg)

		// 태그 이름 비워주기
		// slice 로 변경되면 다른 방식으로 값을 비운다.
		//tg.TargetTag[0] = ""
		//tg.TargetTag[1] = ""
		//tg.TargetTag[2] = ""
	}

	// 전체 타겟(훈련대상)의 수를 반환한다.
	query = `
    select count(target_no) 
    from target_info 
    where user_no = $1`

	pageCount := conn.QueryRow(query, num)
	_ = pageCount.Scan(&total) // 훈련 대상자들의 전체 수를 pages 에 바인딩.

	pages = (total / 20) + 1 // 전체훈련 대상자들을 토대로 전체 페이지수를 계산한다.

	// 각각 표시할 대상 20개, 대상의 총 갯수, 총 페이지 수, 에러를 반환한다.
	return targets, total, pages, nil
}

func (t *TargetNumber) DeleteTarget(conn *sql.DB, num int) error {

	for i := 0; i < len(t.TargetNumber); i++ {
		number, _ := strconv.Atoi(t.TargetNumber[i])

		if t.TargetNumber == nil {
			return fmt.Errorf("Please enter the number of the object to be deleted. ")
		}

		// target_info 테이블에서 대상을 지운다.
		_, err := conn.Exec("DELETE FROM target_info WHERE user_no = $1 AND target_no = $2", num, number)
		if err != nil {
			SugarLogger.Error(err.Error())
			return fmt.Errorf("Error deleting target. ")
		}
	}

	defer conn.Close()

	return nil
}

// 반복해서 읽고 값을 넣는것을 메서드로 구현하고 API는 이걸 그냥 사용하기만 하면됨.
// Excel 파일로부터 대상의 정보를 일괄적으로 읽어 DB에 등록한다.
func (t *Target) ImportTargets(conn *sql.DB, num int, file multipart.File) (int, error) {

	// 훈련대상자 수를 300명으로 제한하기 위해 조회한다.
	var count int // 등록된 훈련대상자 수
	row := conn.QueryRow(`SELECT count(target_no) 
								FROM target_info
								WHERE user_no = $1;`, num)
	err := row.Scan(&count)
	if count >= 300 {
		return 405, fmt.Errorf("The target is already full. ")
	}

	f, err := excelize.OpenReader(file)
	if err != nil {
		SugarLogger.Info(err.Error())
		return 500, nil
	}

	user1 := strconv.Itoa(num) //int -> string

	// Bulk insert 하기 위해 값들을 쌓아놓을 변수
	var BigString string

	// 태그를 미리 담아놓는 변수
	var list = make(map[int]string)

	rows, err := conn.Query(`SELECT tag_no, tag_name
								  FROM tag_info
								  WHERE user_no = $1
								  ORDER BY tag_no ASC;`, num)
	if err != nil {
		SugarLogger.Error(err.Error())
		return 500, nil
	}

	// list 변수에 DB로 부터 조회한 태그 정보를 담아놓는다.
	for rows.Next() {
		var key1 int
		var value string

		err = rows.Scan(&key1, &value)
		if err != nil {
			SugarLogger.Error(err.Error())
		}

		list[key1] = value
	}

	i := 2 // 2행부터 값을 읽어온다.
	for i <= 301 {
		if count >= 301 {
			return 405, fmt.Errorf("exceeded. ")
		}

		str := strconv.Itoa(i)

		t.TargetName = f.GetCellValue("Sheet1", "A"+str)
		t.TargetEmail = f.GetCellValue("Sheet1", "B"+str)

		// 이메일 형식검사
		var validEmail, _ = regexp.MatchString(
			"^[_A-Za-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", t.TargetEmail)

		// 이름 형식검사 (한글, 영어 이름만 허용)
		var validName, _ = regexp.MatchString("^[가-힣A-Za-z0-9\\s]{2,30}$", t.TargetName)

		// 필수적인 정보가 누락됐거나 형식이 잘못된 경우 그 즉시 입력을 중단한다.
		if validName != true || t.TargetName == "" {
			break
		} else if validEmail != true || t.TargetEmail == "" {
			// todo 추후 중복으로 입력된 이메일도 검사해주면 좋을 것 같다.
			break
		}

		t.TargetPhone = f.GetCellValue("Sheet1", "C"+str)

		// 핸드폰 형식검사
		var validPhone, _ = regexp.MatchString(
			"^[0-9]{9,11}$", t.TargetPhone)

		// 핸드폰 번호 형식이 올바르지 않을 경우에는 공백처리한다.
		if validPhone != true {
			t.TargetPhone = ""
		}

		// slice 변수
		var sub []string

		// 여기부터는 선택정보.
		t.TargetOrganize = f.GetCellValue("Sheet1", "D"+str)
		t.TargetPosition = f.GetCellValue("Sheet1", "E"+str)
		sub = append(sub, f.GetCellValue("Sheet1", "F"+str))
		sub = append(sub, f.GetCellValue("Sheet1", "G"+str))
		sub = append(sub, f.GetCellValue("Sheet1", "H"+str))

		// 태그 중복제거
		keys := make(map[string]bool)
		var ue []string

		for _, value := range sub {
			if _, saveValue := keys[value]; !saveValue { // 중복제거 핵심포인트

				keys[value] = true
				ue = append(ue, value)
			}
		}

		// 중복 제거완료
		sub = ue

		// 비어있는 경우 0을 채워넣어 크기 3을 맞춘다.
		for j := 0; j < 3; j++ {
			if len(sub) < j+1 {
				sub = append(sub, "0")
			}

			if sub[j] == "" {
				sub[j] = "0"
			}
		}

	Loop1:
		for k := 0; k < len(sub); k++ {
			// 엑셀 파일의 태그정보와 DB의 정보가 일치할 경우
			if isValueIn(sub[k], list) {
			Loop2:
				for key2, val := range list {
					if val == sub[k] {
						sub[k] = strconv.Itoa(key2)
						break Loop2
					}
				}
				// 그렇지 않을 경우
			} else {
				sub[k] = "0"
				continue Loop1
			}
		}

		//Bulk insert로 삽입할 내용들을 텍스트로 만든다.
		//Note xlsx 파일은 psql 에서 인코딩문제로 bulk insert 불가, csv, txt 등은 가능함.
		BigString += "('" + t.TargetName + "', '" + t.TargetEmail + "', '" +
			t.TargetPhone + "', '" + t.TargetOrganize + "', '" + t.TargetPosition + "', " +
			user1 + "," + sub[0] + "," + sub[1] + "," + sub[2] + ")," + "\n"

		i++
		count++
	}

	BigString = BigString[:len(BigString)-2]

	query := "INSERT INTO target_info (target_name, target_email, target_phone," +
		"target_organize, target_position, user_no, tag1, tag2, tag3) VALUES" +
		BigString

	_, err = conn.Exec(query)
	if err != nil {
		SugarLogger.Error(err.Error())
	}

	defer conn.Close()

	//bulkFile.Close()

	return 200, nil
}

// DB에 저장된 값들을 읽어 엘셀파일에 일괄적으로 작성하여 저장한다.
func ExportTargets(conn *sql.DB, num int, tagNumber int) (bytes.Buffer, error) {

	var buffer bytes.Buffer

	// tagNumber 가 0인 경우 (전체 선택)
	if tagNumber == 0 {
		query := `SELECT target_no,
						   target_name,
						   target_email,
						   target_phone,
						   target_organize,
						   target_position,
						   to_char(modified_time, 'YYYY-MM-DD HH24:MI'),
						   COALESCE(tag_name1, '') as tag_name1,
						   COALESCE(tag_name2, '') as tag_name2,
						   COALESCE(tag_name3, '') as tag_name3
					from target_info as ta
							 LEFT JOIN (SELECT tag_name as tag_name1, user_no, tag_no
										FROM tag_info
										WHERE user_no = $1) ti1 on ti1.tag_no = ta.tag1
							 LEFT JOIN (SELECT tag_name as tag_name2, user_no, tag_no
										FROM tag_info
										WHERE user_no = $1) ti2 on ti2.tag_no = ta.tag2
							 LEFT JOIN (SELECT tag_name as tag_name3, user_no, tag_no
										FROM tag_info
										WHERE user_no = $1) ti3 on ti3.tag_no = ta.tag3
					WHERE ta.user_no = $1
					ORDER BY target_no;`

		rows, err := conn.Query(query, num)
		if err != nil {
			SugarLogger.Error(err.Error())
			return buffer, fmt.Errorf("%v", err)
		}

		// todo 1 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 1은 전부 같은 경로로 수정, api_Target.go 파일의 todo 1 참고)
		// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
		// 서버에 있는 sample 파일에 내용을 작성한 다음 다른 이름의 파일로 클라이언트에게 전송한다.

		//f, err := excelize.OpenFile("./Spreadsheet/sample.xlsx")
		f, err := excelize.OpenFile("../../root/redteam/Spreadsheet/sample.xlsx")
		if err != nil {
			SugarLogger.Error(err.Error())
			return buffer, fmt.Errorf(err.Error())
		}
		//index := f.NewSheet("Sheet1")

		i := 2
		for rows.Next() {
			tg := Target{}
			err = rows.Scan(&tg.TargetNo, &tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
				&tg.TargetPosition, &tg.TargetCreateTime, &tg.TargetTag[0], &tg.TargetTag[1], &tg.TargetTag[2])
			if err != nil {
				SugarLogger.Error(err.Error())
				return buffer, fmt.Errorf(err.Error())
			}

			str := strconv.Itoa(i)
			f.SetCellValue("Sheet1", "A"+str, tg.TargetName)
			f.SetCellValue("Sheet1", "B"+str, tg.TargetEmail)
			f.SetCellValue("Sheet1", "C"+str, tg.TargetPhone)
			f.SetCellValue("Sheet1", "D"+str, tg.TargetOrganize)
			f.SetCellValue("Sheet1", "E"+str, tg.TargetPosition)
			f.SetCellValue("Sheet1", "F"+str, tg.TargetTag[0])
			f.SetCellValue("Sheet1", "G"+str, tg.TargetTag[1])
			f.SetCellValue("Sheet1", "H"+str, tg.TargetTag[2])
			f.SetCellValue("Sheet1", "I"+str, tg.TargetCreateTime)

			i++

			// 태그의 값을 마지막엔 비워준다.
			tg.TargetTag[0] = ""
			tg.TargetTag[1] = ""
			tg.TargetTag[2] = "" // slice 로 변경되면 다른 방식으로 값을 비운다.
		}

		// 메모리에 엑셀 파일을 작성한다.
		if err = f.Write(&buffer); err != nil {
			SugarLogger.Error(err.Error())
			return buffer, err
		}

		return buffer, nil

		// todo -------------------아래부턴 특정 태그만 골라서 내보낼 경우에 해당함.-----------------------------------
	} else {

		query := `SELECT target_name,
						   target_email,
						   target_phone,
						   target_organize,
						   target_position,
						   to_char(modified_time, 'YYYY-MM-DD HH24:MI'),
						   COALESCE(tag_name1, '') as tag_name1,
						   COALESCE(tag_name2, '') as tag_name2,
						   COALESCE(tag_name3, '') as tag_name3
					FROM (SELECT target_name,
								 target_email,
								 target_phone,
								 target_organize,
								 target_position,
								 modified_time,
								 tag1,
								 tag2,
								 tag3
						  FROM target_info as ta
						  WHERE user_no = $1) as T
							 LEFT JOIN (SELECT tag_name as tag_name1, user_no, tag_no
										FROM tag_info
										WHERE user_no = $1) ti1 on ti1.tag_no = T.tag1
							 LEFT JOIN (SELECT tag_name as tag_name2, user_no, tag_no
										FROM tag_info
										WHERE user_no = $1) ti2 on ti2.tag_no = T.tag2
							 LEFT JOIN (SELECT tag_name as tag_name3, user_no, tag_no
										FROM tag_info
										WHERE user_no = $1) ti3 on ti3.tag_no = T.tag3
					WHERE tag1 = $2
					   OR tag2 = $2
					   OR tag3 = $2;`

		result, err := conn.Query(query, num, tagNumber)
		if err != nil {
			SugarLogger.Error(err.Error())
			return buffer, fmt.Errorf(err.Error())
		}

		i := 2
		// todo 1 : 추후 서버에 업로드할 때 경로를 바꿔주어야 한다. (todo 1은 전부 같은 경로로 수정, api_Target.go 파일의 todo 1 참고)
		// 현재는 프로젝트파일의 Spreadsheet 파일에 보관해둔다.
		// 서버에 있는 sample 파일에 내용을 작성한 다음 다른 이름의 파일로 클라이언트에게 전송한다.
		//f, err := excelize.OpenFile("./Spreadsheet/sample.xlsx")
		f, err := excelize.OpenFile("../../root/redteam/Spreadsheet/sample.xlsx")
		if err != nil {
			SugarLogger.Error(err.Error())
			return buffer, fmt.Errorf(err.Error())
		}

		//index := f.NewSheet("Sheet1")

		for result.Next() {
			tg := Target{}

			// 해당 태그에 속하는 대상들을 하나하나 가져온다.
			err = result.Scan(&tg.TargetName, &tg.TargetEmail, &tg.TargetPhone,
				&tg.TargetOrganize, &tg.TargetPosition, &tg.TargetCreateTime,
				&tg.TargetTag[0], &tg.TargetTag[1], &tg.TargetTag[2]) //조회한 값들을 하나하나 바인딩
			if err != nil {
				SugarLogger.Error(err.Error())
				return buffer, fmt.Errorf(err.Error())
			}

			str := strconv.Itoa(i)
			f.SetCellValue("Sheet1", "A"+str, tg.TargetName)
			f.SetCellValue("Sheet1", "B"+str, tg.TargetEmail)
			f.SetCellValue("Sheet1", "C"+str, tg.TargetPhone)
			f.SetCellValue("Sheet1", "D"+str, tg.TargetOrganize)
			f.SetCellValue("Sheet1", "E"+str, tg.TargetPosition)
			f.SetCellValue("Sheet1", "F"+str, tg.TargetTag[0])
			f.SetCellValue("Sheet1", "G"+str, tg.TargetTag[1])
			f.SetCellValue("Sheet1", "H"+str, tg.TargetTag[2])
			f.SetCellValue("Sheet1", "I"+str, tg.TargetCreateTime)

			i++
		}

		// 메모리에 엑셀 파일을 작성한다.
		if err = f.Write(&buffer); err != nil {
			SugarLogger.Error(err.Error())
			return buffer, err
		}
	}

	defer conn.Close()

	return buffer, nil
}

func (t *Tag) CreateTag(conn *sql.DB, num int) (error, int) {

	// 태그 이름 검사 (400 에러)
	var validName, _ = regexp.MatchString("^[가-힣A-Za-z0-9\\s]{1,20}$", t.TagName)
	if validName != true {
		return fmt.Errorf(" Tag Name is not correct. "), 400
	}

	// 태그 중복여부와 개수를 검사한다.
	rows, err := conn.Query(`SELECT tag_name
									FROM tag_info
									WHERE user_no = $1
									GROUP BY tag_name;`, num)
	if err != nil {
		SugarLogger.Error(err.Error())
		return fmt.Errorf("%v", err), 500
	}

	var tagName1 []string
	var tagName2 string
	var count int

	for rows.Next() {
		err = rows.Scan(&tagName2)
		count += 1
		tagName1 = append(tagName1, tagName2)
	}

	// 태그 개수 검사 (405 에러)
	if count >= 5 {
		return fmt.Errorf(" The tag is already full. "), 405
	}

	// 태그 이름 중복검사 (400 에러)
	for i := 0; i < len(tagName1); i++ {
		if t.TagName == tagName1[i] {
			return fmt.Errorf(" That tag name already exists. "), 400
		}
	}

	// 위 조건들 전부 충족할 경우 태그 등록
	_, err = conn.Exec("INSERT INTO tag_info(tag_name, user_no) VALUES ($1, $2)", t.TagName, num)
	if err != nil {
		SugarLogger.Error(err.Error())
		return fmt.Errorf("%v", err), 500
	}

	defer conn.Close()

	return nil, 200
}

func (t *Tag) DeleteTag(conn *sql.DB, num int) error {

	_, err := conn.Exec("DELETE FROM tag_info WHERE tag_no = $1 AND user_no = $2", t.TagNo, num)
	if err != nil {
		SugarLogger.Error(err.Error())
		return fmt.Errorf("Error deleting tag on tag_info ")
	}

	// 해당 태그를 사용하는 프로젝트가 아무런 태그를 가지지 않게 될 경우 삭제한다.
	_, err = conn.Exec("DELETE FROM project_info WHERE tag1 = 0 AND tag2 = 0 AND tag3 = 0")
	if err != nil {
		SugarLogger.Error(err.Error())
		return fmt.Errorf("Error deleting tag on tag_info ")
	}

	defer conn.Close()

	return nil
}

// todo 4 : tag_info 에서 사용자 번호로 태그정보를 가져온다.
func GetTag(conn *sql.DB, num int) ([]Tag, error) {

	defer conn.Close()

	var query string

	query = `SELECT tag_no, tag_name, to_char(modified_time, 'YYYY-MM-DD')
			  FROM tag_info
			  WHERE user_no = $1
			  ORDER BY tag_no asc
`
	tags, err := conn.Query(query, num)
	if err != nil {
		SugarLogger.Error(err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	var tag []Tag
	tg := Tag{}

	for tags.Next() {
		err = tags.Scan(&tg.TagNo, &tg.TagName, &tg.TagCreateTime)
		if err != nil {
			SugarLogger.Error(err.Error())
			return nil, fmt.Errorf("%v", err)
		}

		tag = append(tag, tg)
	}

	return tag, nil
}

func SearchTarget(conn *sql.DB, num int, page int, searchDivision string, searchText string) ([]Target, int, int, error) {
	var pageNum int // 몇번째 페이지부터 가져올지 결정하는 변수
	var pages int   // 총 페이지 수
	var total int   // 총 훈련대상자들의 수를 담을 변수

	// ex) 1페이지 -> 1~10, 2페이지 -> 11~20
	// 페이지번호에 따라 가져올 목록이 달라진다.
	pageNum = (page - 1) * 20

	// 대상목록들을 20개씩만 잘라서 반하여 페이징처리한다.
	query := "SELECT row_num, " +
		"target_name, " +
		"target_email, " +
		"target_phone, " +
		"target_organize, " +
		"target_position, " +
		"COALESCE(tag_name1, '') as tag_name1, " +
		"COALESCE(tag_name2, '') as tag_name2, " +
		"COALESCE(tag_name3, '') as tag_name3, " +
		"to_char(modified_time, 'YYYY-MM-DD')," +
		"target_no " +
		"FROM (SELECT ROW_NUMBER() over (ORDER BY target_no) AS row_num, " +
		"target_no, " +
		"target_name, " +
		"target_email, " +
		"target_phone, " +
		"target_organize, " +
		"target_position, " +
		"tag1, " +
		"tag2, " +
		"tag3, " +
		"modified_time " +
		"FROM target_info " +
		"WHERE user_no = $1 AND " + searchDivision + " LIKE $2 " +
		") AS T " +
		"LEFT JOIN (SELECT tag_name as tag_name1, user_no, tag_no " +
		"FROM tag_info " +
		"WHERE user_no = $1) ti1 on ti1.tag_no = T.tag1 " +
		"LEFT JOIN (SELECT tag_name as tag_name2, user_no, tag_no " +
		"FROM tag_info " +
		"WHERE user_no = $1) ti2 on ti2.tag_no = T.tag2 " +
		"LEFT JOIN (SELECT tag_name as tag_name3, user_no, tag_no " +
		"FROM tag_info " +
		"WHERE user_no = $1) ti3 on ti3.tag_no = T.tag3 " +
		"WHERE row_num > $3 " +
		"ORDER BY target_no asc " +
		"LIMIT 20;"

	searchText = "%" + searchText + "%"
	rows, err := conn.Query(query, num, searchText, pageNum)
	if err != nil {
		SugarLogger.Error(err.Error())
		return nil, 0, 0, fmt.Errorf("%v", err)
	}

	var targets []Target
	tg := Target{}

	for rows.Next() { // 목록들을 하나하나 읽어들여온다.
		err = rows.Scan(&tg.FakeNo, &tg.TargetName, &tg.TargetEmail, &tg.TargetPhone, &tg.TargetOrganize,
			&tg.TargetPosition, &tg.TargetTag[0], &tg.TargetTag[1], &tg.TargetTag[2], &tg.TargetCreateTime, &tg.TargetNo)
		if err != nil {
			SugarLogger.Error(err.Error())
			fmt.Printf("%v", err)
			continue
		}

		// 추후 프론트에서 처리하도록 수정함.
		//var sub [3]string
		//phone := []rune(tg.TargetPhone)
		//
		//if len(tg.TargetPhone) < 10 {
		//	sub[0] = string(phone[0:2])
		//	sub[1] = string(phone[2:5])
		//	sub[2] = string(phone[5:9])
		//	tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		//} else if string(phone[1:2]) == "2" && len(tg.TargetPhone) == 10 {
		//	sub[0] = string(phone[0:2])
		//	sub[1] = string(phone[2:6])
		//	sub[2] = string(phone[6:10])
		//	tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		//} else if len(tg.TargetPhone) == 10 {
		//	sub[0] = string(phone[0:3])
		//	sub[1] = string(phone[3:6])
		//	sub[2] = string(phone[6:10])
		//	tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		//} else if len(tg.TargetPhone) == 11 {
		//	sub[0] = string(phone[0:3])
		//	sub[1] = string(phone[3:7])
		//	sub[2] = string(phone[7:11])
		//	tg.TargetPhone = sub[0] + "-" + sub[1] + "-" + sub[2]
		//}

		targets = append(targets, tg)

		// slice 로 변경되면 다른 방식으로 값을 비운다.
		//tg.TargetTag[0] = ""
		//tg.TargetTag[1] = ""
		//tg.TargetTag[2] = ""
	}

	// 전체 타겟(훈련대상)의 수를 반환한다.
	query = "SELECT count(target_no) " +
		"FROM target_info " +
		"WHERE user_no = $1 AND " + searchDivision + " LIKE $2"

	pageCount := conn.QueryRow(query, num, searchText)
	_ = pageCount.Scan(&total) // 훈련 대상자들의 전체 수를 pages 에 바인딩.

	pages = (total / 20) + 1 // 전체훈련 대상자들을 토대로 전체 페이지수를 계산한다.

	defer conn.Close()

	// 각각 표시할 대상 20개, 대상의 총 갯수, 총 페이지 수, 에러를 반환한다.
	return targets, total, pages, nil
}
