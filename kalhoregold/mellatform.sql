create table mellatform
(
    ID          Int64,
    RefID       String,
    EncPan      String,
    Enc         String,
    PhoneNumber String,
    Body        String,
    CreatedAt   DateTime
)
    engine = MergeTree ORDER BY ID
        SETTINGS index_granularity = 8192;

