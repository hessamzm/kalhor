create table wallet_gold
(
    melli_number   String,
    balance_in     Float64,
    balance_out    Float64,
    feebalance_in  Float64,
    feebalance_out Float64,
    freez_bl_in    Float64,
    freez_bl_out   Float64,
    ban_bl_in      Float64,
    ban_bl_out     Float64,
    event_time     DateTime
)
    engine = MergeTree ORDER BY (melli_number, event_time)
        SETTINGS index_granularity = 8192;

