USE isuride;

ALTER TABLE chairs
    ADD COLUMN total_distance INTEGER NOT NULL DEFAULT 0 COMMENT '移動距離';
ALTER TABLE chairs
    ADD COLUMN total_distance_updated_at DATETIME(6) NULL COMMENT '移動距離更新日時';

ALTER TABLE chairs ADD COLUMN is_free TINYINT(1) NOT NULL DEFAULT 1 COMMENT '椅子が空いているかどうか';

ALTER TABLE rides ADD COLUMN matched_at DATETIME(6) NULL COMMENT '椅子とのマッチング日時';
ALTER TABLE rides ADD COLUMN arrived_first_at DATETIME(6) NULL COMMENT '椅子到着日時';
ALTER TABLE rides ADD COLUMN picked_up_at DATETIME(6) NULL COMMENT '椅子乗車日時';
ALTER TABLE rides ADD COLUMN arrived_at DATETIME(6) NULL COMMENT '目的地到着日時';
