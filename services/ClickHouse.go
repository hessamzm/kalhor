package services

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"kalhor/app/models"
	"kalhor/utils"
	"log"
	"net"
	"os"
	"time"
)

type WalletService struct {
	db driver.Conn
}
type WalletServiceRial struct {
	db driver.Conn
}
type MellatService struct {
	db driver.Conn
}

// back pointer
func NewWalletService() (*WalletService, error) {
	conn, err := Gonn()
	if err != nil {
		return nil, err
	}
	return &WalletService{
		db: conn,
	}, nil
}
func NewWalletServiceRial() (*WalletServiceRial, error) {
	conn, err := Gonn()
	if err != nil {
		return nil, err
	}
	return &WalletServiceRial{
		db: conn,
	}, nil
}
func NewMellatService() (*MellatService, error) {
	conn, err := Gonn()
	if err != nil {
		return nil, err
	}
	return &MellatService{
		db: conn,
	}, nil
}

// insert to

func (s *MellatService) InsertPaymentGateway(id int64, refId, encPan, enc, phoneNumber, body, saleOrderID, saleReference, resCode string, amount float64) error {
	if utils.KlDebug {
		fmt.Printf("Inserting into payment_gateway: %d, %s, %s, %s, %s, %s, %s, %s, %s, %.2f\n", id, refId, encPan, enc, phoneNumber, body, saleOrderID, saleReference, resCode, amount)
	}
	query := `
		INSERT INTO payment_gateway (ID, RefID, EncPan, Enc, PhoneNumber, Body, SaleOrderID, SaleReference, ResCode, Amount, CreatedAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	return s.db.Exec(context.Background(), query, id, refId, encPan, enc, phoneNumber, body, saleOrderID, saleReference, resCode, amount, now, now)
}
func (s *MellatService) InsertMellatForm(refId, phoneNumber, body, encPan, encMelliNumber string) error {
	if utils.KlDebug {
		fmt.Println("insert mellatdb", refId, phoneNumber, body, encPan, encMelliNumber)
	}
	query := `INSERT INTO mellatform (RefID, PhoneNumber, Body, EncPan, Enc) VALUES (?, ?, ?, ?, ?)`
	return s.db.Exec(context.Background(), query, refId, phoneNumber, body, encPan, encMelliNumber)
}

func (s *WalletService) InsertWalletGold(gw *models.WalletGold) error {
	if utils.KlDebug {
		fmt.Println("insert wallet gold:", gw)
	}
	gw.EventTime = time.Now()
	// حذف 'id' از کوئری
	query := `INSERT INTO wallet_gold (melli_number, balance_in, feebalance_in, event_time) VALUES (?, ?, ?, ?)`
	return s.db.Exec(context.Background(), query, "kal"+gw.MelliNumber, gw.BalanceIn, gw.FeebalanceIn, gw.EventTime)
}
func (s *WalletService) TakeOutWalletGold(gw *models.WalletGold) error {
	if utils.KlDebug {
		fmt.Println("insert wallet gold:", gw)
	}
	gw.EventTime = time.Now()
	// حذف 'id' از کوئری
	query := `INSERT INTO wallet_gold (melli_number, balance_out, feebalance_out, event_time) VALUES (?, ?, ?, ?)`
	return s.db.Exec(context.Background(), query, "kal"+gw.MelliNumber, gw.BalanceOut, gw.FeebalanceOut, gw.EventTime)
}
func (s *WalletServiceRial) OutWalletGoldRial(gr *models.WalletRial) error {
	if utils.KlDebug {
		fmt.Println("insert wallet gold:", gr)
	}
	gr.EventTime = time.Now()
	// حذف 'id' از کوئری
	query := `INSERT INTO wallet_rial (trakonesh_id , melli_number , balance_out ,event_time) VALUES (?, ?, ?, ?)`
	return s.db.Exec(context.Background(), query, gr.TrakoneshId, "kal"+gr.MelliNumber, gr.BalanceOut, gr.EventTime)
}
func (s *WalletServiceRial) InsertWalletGoldRial(gr *models.WalletRial) error {
	if utils.KlDebug {
		fmt.Println("insert wallet gold:", gr)
	}
	gr.EventTime = time.Now()
	// حذف 'id' از کوئری
	query := `INSERT INTO wallet_rial (trakonesh_id , melli_number , balance_out ,event_time) VALUES (?, ?, ?, ?)`
	return s.db.Exec(context.Background(), query, gr.TrakoneshId, "kal"+gr.MelliNumber, gr.BalanceOut, gr.EventTime)
}
func (s *WalletService) ShowWalletGold(m string) ([]*models.WalletGold, error) {
	m = "kal" + m
	query := fmt.Sprintf("SELECT * FROM wallet_gold WHERE melli_number = '%s';", m)

	if utils.KlDebug {
		fmt.Println("query:", query)
	}

	// اجرای کوئری و دریافت چندین سطر
	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*models.WalletGold
	for rows.Next() {
		wallet := &models.WalletGold{}
		e := rows.Scan(
			&wallet.MelliNumber,
			&wallet.BalanceIn,
			&wallet.BalanceOut,
			&wallet.FreezBlIn,
			&wallet.FreezBlOut,
			&wallet.BanBlIn,
			&wallet.BanBlOut,
			&wallet.EventTime,
		)
		if e != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	// بررسی اینکه آیا سطری دریافت شده است یا خیر
	if len(wallets) == 0 {
		fmt.Println("No user found with the given MelliNumber.")
		return nil, nil
	}
	if utils.KlDebug {
		fmt.Println("wallets:", wallets)
	}

	return wallets, nil
}

func (s *WalletService) Queryforgold(q string) ([]*models.WalletGold, error) {
	fmt.Println(q)
	rows, err := s.db.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println(rows)
	var wallets []*models.WalletGold
	fmt.Printf("wallet 2", wallets)
	for rows.Next() {
		fmt.Println("test for")
		wallet := &models.WalletGold{}
		if err := rows.Scan(
			&wallet.MelliNumber,
			&wallet.BalanceIn,
			&wallet.BalanceOut,
			&wallet.FeebalanceIn,
			&wallet.FeebalanceOut,
			&wallet.FreezBlIn,
			&wallet.FreezBlOut,
			&wallet.BanBlIn,
			&wallet.BanBlOut,
			&wallet.EventTime,
		); err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
		fmt.Printf("Scanned wallet: %v\n", wallet) // چاپ نتایج اسکن شده
	}

	fmt.Printf("wallet 3", wallets)
	// بررسی اینکه آیا سطری دریافت شده است یا خیر
	if len(wallets) == 0 {
		fmt.Println("No user found with the given MelliNumber.")
		return nil, nil
	}
	if utils.KlDebug {
		fmt.Println("wallets:", wallets)
	}
	return wallets, err
}
func (r *WalletServiceRial) Queryforrial(q string) ([]*models.WalletRial, error) {
	fmt.Println(q)
	rows, err := r.db.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*models.WalletRial
	fmt.Printf("wallet 2", wallets)
	for rows.Next() {
		wallet := &models.WalletRial{}
		fmt.Println("test for")
		err := rows.Scan(
			&wallet.MelliNumber,
			&wallet.BalanceIn,
			&wallet.BalanceOut,
			&wallet.FreezBlIn,
			&wallet.FreezBlOut,
			&wallet.BanBlIn,
			&wallet.BanBlOut,
			&wallet.EventTime,
			&wallet.TrakoneshId,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
		fmt.Printf("Scanned wallet: %v\n", wallet)
	}
	fmt.Printf("wallet 3", wallets)
	// بررسی اینکه آیا سطری دریافت شده است یا خیر
	if len(wallets) == 0 {
		fmt.Println("No user found with the given MelliNumber.")
		return nil, nil
	}
	if utils.KlDebug {
		fmt.Println("wallets:", wallets)
	}
	return wallets, err
}

func (s *WalletService) GetBalanceDifference(m string) (float64, error) {
	m = "kal" + m
	query := fmt.Sprintf("SELECT SUM(balance_in) AS total_balance_in, SUM(balance_out) AS total_balance_out FROM wallet_gold WHERE melli_number = '%s';", m)

	if utils.KlDebug {
		fmt.Println("query:", query)
	}

	// اجرای کوئری
	row := s.db.QueryRow(context.Background(), query)

	var totalBalanceIn, totalBalanceOut float64
	err := row.Scan(&totalBalanceIn, &totalBalanceOut)
	if err != nil {
		return 0, err
	}

	// محاسبه اختلاف
	balanceDifference := totalBalanceIn - totalBalanceOut

	return balanceDifference, nil
}
func (s *WalletService) GetBalanceDifferenceRial(m string) (float64, error) {
	m = "kal" + m
	query := fmt.Sprintf("SELECT SUM(balance_in) AS total_balance_in, SUM(balance_out) AS total_balance_out FROM wallet_rial WHERE melli_number = '%s';", m)

	if utils.KlDebug {
		fmt.Println("query:", query)
	}

	// اجرای کوئری
	row := s.db.QueryRow(context.Background(), query)

	var totalBalanceIn, totalBalanceOut float64
	err := row.Scan(&totalBalanceIn, &totalBalanceOut)
	if err != nil {
		return 0, err
	}

	// محاسبه اختلاف
	balanceDifference := totalBalanceIn - totalBalanceOut

	return balanceDifference, nil
}
func (s *WalletService) GetTotalFeeBy(m string) (float64, error) {
	m = "kal" + m
	query := fmt.Sprintf("SELECT SUM(feebalance_in) AS total_fee_by FROM wallet_gold WHERE melli_number = '%s';", m)

	if utils.KlDebug {
		fmt.Println("query:", query)
	}

	// اجرای کوئری
	row := s.db.QueryRow(context.Background(), query)

	var total_fee_by float64
	err := row.Scan(&total_fee_by)
	if err != nil {
		return 0, err
	}

	// محاسبه اختلاف
	//balanceDifference := totalBalanceIn - totalBalanceOut

	return total_fee_by, nil
}

// base connection
func Gonn() (driver.Conn, error) {
	// مسیر فایل لاگ در پوشه database
	logFilePath := "storage/logs/clickhouse"

	// ایجاد یا باز کردن فایل برای ذخیره لاگ‌ها
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("خطا در باز کردن فایل لاگ:", err)
		return nil, err
	}

	// تنظیم لاگر برای نوشتن در فایل به جای ترمینال
	logger := log.New(file, "", log.LstdFlags)

	// اضافه کردن جداکننده بین اجراهای مختلف
	logger.Println("===================================")

	dialCount := 0
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "kalhoregold",
			Username: "",
			Password: "",
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			logger.Printf(format, v...) // لاگ‌ها را به فایل می‌نویسد
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "k.d-app", Version: "0.0.1"},
			},
		},
	})

	if err != nil {
		logger.Println("خطا در اتصال به ClickHouse:", err)
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		logger.Println("خطا در پینگ:", err)
		return nil, err
	}

	logger.Println("اتصال به ClickHouse برقرار شد.")
	return conn, err
}
