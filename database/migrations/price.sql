CREATE TABLE prices (
                        id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,  -- کلید اصلی با مقدار خودکار
                        sell_price DECIMAL(10, 2) NOT NULL,              -- ستون فروش با دقت 2 رقم اعشار
                        by_price DECIMAL(10, 2) NOT NULL,                -- ستون خرید با دقت 2 رقم اعشار
                        status VARCHAR(255) NOT NULL,                    -- وضعیت به صورت رشته‌ای
                        base_18 DECIMAL(10, 2) NOT NULL,                 -- ستون base_18 با دقت 2 رقم اعشار
                        base_24 DECIMAL(10, 2) NOT NULL,                 -- ستون base_24 با دقت 2 رقم اعشار
                        ojrat DECIMAL(10, 2) NOT NULL,                   -- ستون اجرت با دقت 2 رقم اعشار
                        maliat DECIMAL(10, 2) NOT NULL,                  -- ستون مالیات با دقت 2 رقم اعشار
                        sood DECIMAL(10, 2) NOT NULL,                    -- ستون سود با دقت 2 رقم اعشار
                        updated_at DATETIME(3) DEFAULT NULL,             -- زمان بروزرسانی با دقت 3 میلی‌ثانیه
                        PRIMARY KEY (id),                                -- تعریف کلید اصلی
                        KEY idx_prices_updated_at (updated_at)           -- ایندکس برای ستون updated_at
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
