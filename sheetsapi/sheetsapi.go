package sheetsapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/bakyazi/gopass/sheetsapi/util"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const PAGE_SIZE = 100

var (
	ErrNotFound     = errors.New("entry not found")
	ErrAlreadyExist = errors.New("already exists")
)

type PasswordDB interface {
	GetSites() ([]string, error)
	GetPasswords(page, pageSize int) ([]PasswordEntry, error)
	GetPassword(site, username string) (string, error)
	CreatePassword(site, username, password string) error
	UpdatePassword(site, username, password string) error
	DeletePassword(site, username string) error
	Clear() error
}

type PasswordEntry struct {
	Row      int
	Site     string
	Username string
	Password string
}

type passwordDBImpl struct {
	srv     *sheets.Service
	sheetId string
}

func NewPasswordDB(credentialFileName string, sheetId string) (PasswordDB, error) {
	service, err := sheets.NewService(context.Background(), option.WithCredentialsFile(credentialFileName))
	if err != nil {
		return nil, err
	}
	return &passwordDBImpl{srv: service, sheetId: sheetId}, nil
}

func (p *passwordDBImpl) GetSites() ([]string, error) {
	resp, err := p.srv.Spreadsheets.Values.Get(p.sheetId, "A2:A").Do()
	if err != nil {
		return nil, err
	}
	return util.Map(resp.Values, func(t []interface{}, i int) string {
		return t[0].(string)
	}), nil
}

func (p passwordDBImpl) GetPasswords(page, pageSize int) ([]PasswordEntry, error) {
	sht, err := p.srv.Spreadsheets.Get(p.sheetId).Do()
	if err != nil {
		return nil, err
	}
	maxRow := sht.Sheets[0].Properties.GridProperties.RowCount
	start, end := calculateRange(page, pageSize)
	if maxRow < int64(end) {
		end = int(maxRow)
	}
	if start > end {
		return nil, nil
	}
	resp, err := p.srv.Spreadsheets.Values.Get(p.sheetId, fmt.Sprintf("A%d:C%d", start, end)).Do()
	if err != nil {
		return nil, err
	}
	return util.Map(resp.Values, func(t []interface{}, index int) PasswordEntry {
		if len(t) < 3 {
			return PasswordEntry{Row: start + index}
		}
		return PasswordEntry{
			Row:      start + index,
			Site:     t[0].(string),
			Username: t[1].(string),
			Password: t[2].(string),
		}
	}), nil
}

func (p passwordDBImpl) GetPassword(site, username string) (string, error) {
	val, err := p.find(site, username)
	if err != nil {
		return "", errors.New("password not found")
	}
	return val.Password, nil
}

func (p passwordDBImpl) CreatePassword(site, username, password string) error {
	_, err := p.find(site, username)
	if err == nil {
		return ErrAlreadyExist
	}
	if !errors.Is(err, ErrNotFound) {
		return err
	}

	//_, err = p.srv.Spreadsheets.BatchUpdate(p.sheetId, &sheets.BatchUpdateSpreadsheetRequest{
	//	Requests: []*sheets.Request{
	//		{
	//			AppendDimension: &sheets.AppendDimensionRequest{
	//				Dimension: "ROWS",
	//				Length:    1,
	//				SheetId:   0,
	//			},
	//		},
	//	},
	//}).Do()
	//if err != nil {
	//	return err
	//}
	_, err = p.srv.Spreadsheets.Values.Append(p.sheetId, "A:C",
		&sheets.ValueRange{Range: "A:C", Values: [][]interface{}{{site, username, password}}}).
		Do(googleapi.QueryParameter("insertDataOption", "INSERT_ROWS"),
			googleapi.QueryParameter("valueInputOption", "RAW"))
	return err
}

func (p passwordDBImpl) UpdatePassword(site, username, password string) error {
	val, err := p.find(site, username)
	if err != nil {
		return err
	}
	valRange := fmt.Sprintf("C%d:C%d", val.Row, val.Row)
	_, err = p.srv.Spreadsheets.Values.Update(p.sheetId, valRange, &sheets.ValueRange{
		Range: valRange, Values: [][]interface{}{{password}},
	}).Do(googleapi.QueryParameter("valueInputOption", "RAW"))
	return err
}

func (p passwordDBImpl) DeletePassword(site string, username string) error {
	val, err := p.find(site, username)
	if err != nil {
		return err
	}
	_, err = p.srv.Spreadsheets.BatchUpdate(p.sheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteDimension: &sheets.DeleteDimensionRequest{
					Range: &sheets.DimensionRange{
						SheetId:    0,
						Dimension:  "ROWS",
						StartIndex: int64(val.Row - 1),
						EndIndex:   int64(val.Row),
					},
				},
			},
		},
	}).Do()
	return err
}

func (p passwordDBImpl) Clear() error {
	_, err := p.srv.Spreadsheets.BatchUpdate(p.sheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteDimension: &sheets.DeleteDimensionRequest{
					Range: &sheets.DimensionRange{
						SheetId:    0,
						Dimension:  "ROWS",
						StartIndex: 1,
					},
				},
			},
		},
	}).Do()
	return err
}

func calculateRange(page, pageSize int) (int, int) {
	start := 2 + (page-1)*pageSize
	end := start + pageSize
	return start, end - 1
}

func (p *passwordDBImpl) find(site, user string) (PasswordEntry, error) {
	page := 1
	for {
		data, err := p.GetPasswords(page, PAGE_SIZE)
		if err != nil {
			return PasswordEntry{}, err
		}
		if len(data) == 0 {
			break
		}
		for _, v := range data {
			if v.Site == site && v.Username == user {
				return v, nil
			}
		}
		page++

	}
	return PasswordEntry{}, ErrNotFound
}
