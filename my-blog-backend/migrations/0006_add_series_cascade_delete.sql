-- 为 series_sections 表添加级联删除约束
-- 这样当删除系列时，会自动删除关联的章节、子章节和文章关联

-- 删除旧的约束（如果存在）
ALTER TABLE series_sections DROP FOREIGN KEY fk_series_sections_series;

-- 添加新的级联删除约束
ALTER TABLE series_sections
ADD CONSTRAINT fk_series_sections_series
FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE;
