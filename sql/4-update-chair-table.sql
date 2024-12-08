USE isuride;

ALTER TABLE chairs
    ADD COLUMN total_distance INTEGER NOT NULL DEFAULT 0 COMMENT '移動距離';
ALTER TABLE chairs
    ADD COLUMN total_distance_updated_at DATETIME(6) NULL COMMENT '移動距離更新日時';
