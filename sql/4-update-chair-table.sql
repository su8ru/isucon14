ALTER TABLE chairs
    ADD COLUMN total_distance INTEGER NOT NULL DEFAULT 0 COMMENT '移動距離';
ALTER TABLE chairs
    ADD COLUMN total_distance_updated_at DATETIME(6) NULL COMMENT '移動距離更新日時';

CREATE TRIGGER IF NOT EXISTS update_chair_total_distance
    BEFORE INSERT
    ON chair_locations
    FOR EACH ROW
BEGIN
    DECLARE distance INTEGER;

    SET distance = (SELECT ABS(NEW.latitude - latitude) + ABS(NEW.longitude - longitude)
                    FROM chair_locations
                    WHERE chair_id = NEW.chair_id
                    ORDER BY created_at DESC
                    LIMIT 1);

    UPDATE chairs
    SET total_distance            = IFNULL(chairs.total_distance + distance, 0),
        total_distance_updated_at = NEW.created_at
    WHERE id = NEW.chair_id;
END;
